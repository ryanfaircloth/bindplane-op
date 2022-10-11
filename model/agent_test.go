// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgentApplyLabels(t *testing.T) {
	agent := Agent{}

	tests := []struct {
		selector string
		success  bool
		expect   Labels
	}{
		{
			selector: "app=mindplane",
			success:  true,
			expect: LabelsFromValidatedMap(map[string]string{
				"app": "mindplane",
			}),
		},
		{
			selector: "app=mindplane,env=production",
			success:  true,
			expect: LabelsFromValidatedMap(map[string]string{
				"app": "mindplane",
				"env": "production",
			}),
		},
		{
			selector: "app=mindplane, env = production",
			success:  true,
			expect: LabelsFromValidatedMap(map[string]string{
				"app": "mindplane",
				"env": "production",
			}),
		},
		{
			selector: "app=====",
			success:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.selector, func(t *testing.T) {
			labels, err := LabelsFromSelector(test.selector)
			agent.Labels = labels
			if test.success {
				require.NoError(t, err)
				require.Equal(t, test.expect, agent.Labels)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAgentMatchesSelector(t *testing.T) {
	tests := []struct {
		labels   map[string]string
		selector string
		matches  bool
	}{
		{
			labels: map[string]string{
				"app":     "mindplane",
				"os":      "Darwin",
				"version": "2.0.6",
			},
			selector: "app=mindplane",
			matches:  true,
		},
		{
			labels: map[string]string{
				"app":     "mindplane",
				"os":      "Darwin",
				"version": "2.0.6",
			},
			selector: "app=mindplane,version=2",
			matches:  false,
		},
		{
			labels: map[string]string{
				"app":     "mindplane",
				"os":      "Darwin",
				"version": "2.0.6",
			},
			selector: "os=Darwin,app=mindplane",
			matches:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.selector, func(t *testing.T) {
			selector, err := SelectorFromString(test.selector)
			require.NoError(t, err)
			require.Equal(t, test.matches, selector.Matches(LabelsFromValidatedMap(test.labels)))
		})
	}
}

func TestAgentSupportsUpgrade(t *testing.T) {
	tests := []struct {
		version string
		expect  bool
	}{
		{
			version: "v1.5.0",
			expect:  false,
		},
		{
			version: "v1.6.0",
			expect:  true,
		},
		{
			version: "v1.6.1",
			expect:  true,
		},
		{
			version: "v2.0.0",
			expect:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			agent := &Agent{
				Version: test.version,
			}
			require.Equal(t, test.expect, agent.SupportsUpgrade())
		})
	}
}

func TestAgentUpgradeComplete(t *testing.T) {
	tests := []struct {
		name          string
		prepareAgent  func(a *Agent)
		errorMessage  string
		expectStatus  AgentStatus
		expectUpgrade *AgentUpgrade
	}{
		{
			name: "success",
			prepareAgent: func(a *Agent) {
				a.UpgradeTo("v1.1")
			},
			errorMessage:  "",
			expectStatus:  Connected,
			expectUpgrade: nil,
		},
		{
			name: "success with existing error",
			prepareAgent: func(a *Agent) {
				a.UpgradeStarted("v1.1", []byte{1})
				a.ErrorMessage = "error"
			},
			errorMessage:  "",
			expectStatus:  Error,
			expectUpgrade: nil,
		},
		{
			name: "fail",
			prepareAgent: func(a *Agent) {
				a.Status = Connected
			},
			errorMessage: "upgrade error",
			expectStatus: Connected,
			expectUpgrade: &AgentUpgrade{
				Status:  UpgradeFailed,
				Version: "v1.2",
				Error:   "upgrade error",
			},
		},
		{
			name: "fail with upgrade",
			prepareAgent: func(a *Agent) {
				a.UpgradeTo("v1.1")
				a.Status = Connected
			},
			errorMessage: "upgrade error",
			expectStatus: Connected,
			expectUpgrade: &AgentUpgrade{
				Status:  UpgradeFailed,
				Version: "v1.2",
				Error:   "upgrade error",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			agent := &Agent{}
			test.prepareAgent(agent)
			agent.UpgradeComplete("v1.2", test.errorMessage)
			require.Equal(t, test.expectStatus, agent.Status)
			require.Equal(t, test.expectUpgrade, agent.Upgrade)
		})
	}
}

func TestFeatures(t *testing.T) {
	tests := []struct {
		version        string
		expectFeatures AgentFeatures
	}{
		{
			version:        "",
			expectFeatures: 0,
		},
		{
			version:        "1.5.0",
			expectFeatures: 0,
		},
		{
			version:        "1.6.0",
			expectFeatures: AgentSupportsUpgrade,
		},
		{
			version:        "1.7.0",
			expectFeatures: AgentSupportsUpgrade,
		},
		{
			version:        "1.8.0",
			expectFeatures: AgentSupportsUpgrade | AgentSupportsSnapshots,
		},
		{
			version:        "1.9.2",
			expectFeatures: AgentSupportsUpgrade | AgentSupportsSnapshots | AgentSupportsMeasurements,
		},
		{
			version:        "2.0.0",
			expectFeatures: AgentSupportsUpgrade | AgentSupportsSnapshots | AgentSupportsMeasurements,
		},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			agent := Agent{
				Version: test.version,
			}
			require.Equal(t, test.expectFeatures, agent.Features())
		})
	}
}
