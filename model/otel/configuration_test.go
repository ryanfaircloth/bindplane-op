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

package otel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueComponentID(t *testing.T) {
	tests := []struct {
		original     string
		typeName     string
		resourceName string
		expect       string
	}{
		{
			original:     "plugin",
			typeName:     "macos",
			resourceName: "1",
			expect:       "plugin/1",
		},
		{
			original:     "plugin",
			typeName:     "macos",
			resourceName: "name",
			expect:       "plugin/name",
		},
		{
			original:     "plugin/foo",
			typeName:     "macos",
			resourceName: "name",
			expect:       "plugin/name__foo",
		},
		{
			// This is malformed, but uniqueName isn't responsible
			original:     "plugin/foo/bar",
			typeName:     "macos",
			resourceName: "name",
			expect:       "plugin/name__foo/bar",
		},
	}
	for _, test := range tests {
		t.Run(test.original, func(t *testing.T) {
			result := UniqueComponentID(test.original, test.typeName, test.resourceName)
			require.Equal(t, test.expect, string(result))

			// round-trip original and expect
			for _, str := range []string{test.original, test.expect} {
				t.Run(str, func(t *testing.T) {
					pipelineType, name := ParseComponentID(ComponentID(str))
					id := NewComponentID(pipelineType, name)
					require.Equal(t, str, string(id))
				})
			}
		})
	}
}

func TestEmptyConfiguration(t *testing.T) {
	c := NewConfiguration()
	yaml, err := c.YAML()
	require.NoError(t, err)
	require.Equal(t, NoopConfig, yaml)
}

func TestNilConfiguration(t *testing.T) {
	var c *Configuration
	yaml, err := c.YAML()
	require.NoError(t, err)
	require.Equal(t, NoopConfig, yaml)
}

func TestPipelineTypeFlags_Set(t *testing.T) {
	tests := []struct {
		name   string
		flags  PipelineTypeFlags
		add    PipelineTypeFlags
		expect PipelineTypeFlags
	}{
		{
			name:   "zeros",
			flags:  0,
			add:    0,
			expect: 0,
		},
		{
			name:   "metrics",
			flags:  0,
			add:    MetricsFlag,
			expect: MetricsFlag,
		},
		{
			name:   "logs|metrics",
			flags:  LogsFlag,
			add:    MetricsFlag,
			expect: LogsFlag | MetricsFlag,
		},
		{
			name:   "logs|metrics + metrics",
			flags:  LogsFlag | MetricsFlag,
			add:    MetricsFlag,
			expect: LogsFlag | MetricsFlag,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.flags.Set(test.add)
			require.Equal(t, test.expect, test.flags)
		})
	}
}
