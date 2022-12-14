// Copyright  observIQ, Inc
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
	"sort"
	"time"

	"github.com/observiq/bindplane-op/internal/store/search"
	"github.com/observiq/bindplane-op/internal/util/semver"
)

// AgentTypeName is the name of a type of agent
type AgentTypeName string

const (
	// AgentTypeNameObservIQOtelCollector is the name of the observIQ Distro for OpenTelemetry Collector
	AgentTypeNameObservIQOtelCollector AgentTypeName = "observiq-otel-collector"
)

// AgentStatus TODO(doc)
type AgentStatus uint8

const (
	// Disconnected is the state of an agent that was formerly Connected to the management platform but is no longer
	// connected. This could mean that the agent has stopped running, the network connection has been interrupted, or that
	// the server terminated the connection.
	Disconnected AgentStatus = 0

	// Connected is the normal state of a healthy agent that is Connected to the management platform.
	Connected AgentStatus = 1

	// Error occurs if there is an error running or Configuring the agent.
	Error AgentStatus = 2

	// ComponentFailed is deprecated.
	ComponentFailed AgentStatus = 4

	// Deleted is set on a deleted Agent before notifying observers of the change.
	Deleted AgentStatus = 5

	// Configuring is set on an Agent when it is sent a new configuration that has not been applied. After successful
	// Configuring, it will transition back to Connected. If there is an error Configuring, it will transition to Error.
	Configuring AgentStatus = 6

	// Upgrading is set on an Agent when it has been sent a new package that is being applied. After Upgrading, it will
	// transition back to Connected or Error unless it already has the Configuring status.
	Upgrading AgentStatus = 7
)

// AgentUpgradeStatus is the status of the AgentUpgrade
type AgentUpgradeStatus uint8

const (
	// UpgradePending is set when the upgrade is initially started
	UpgradePending AgentUpgradeStatus = 0
	// UpgradeStarted is set when the upgrade has been sent to the agent
	UpgradeStarted AgentUpgradeStatus = 1
	// UpgradeFailed is set when the upgrade is complete but there was an error. If the upgrade is successful, the Agent
	// Upgrade field will be set to nil and there is no corresponding status.
	UpgradeFailed AgentUpgradeStatus = 2
)

// AgentUpgrade stores information on an Agent about the upgrade process.
type AgentUpgrade struct {
	// Status indicates the progress of the agent upgrade
	Status AgentUpgradeStatus `json:"status" yaml:"status"`

	// Version is used to indicate that an agent should be or is being upgraded. The agent status will be set to Upgrading
	// when the upgrade begins.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// AllPackagesHash is the hash of the packages sent to the agent to upgrade
	AllPackagesHash []byte `json:"allPackagesHash,omitempty" yaml:"allPackagesHash,omitempty"`

	// Error is set if there were errors upgrading the agent
	Error string `json:"error,omitempty" yaml:"error,omitempty"`
}

// AgentFeatures is a bitmask of features supported by the Agent, usually based on its version.
type AgentFeatures uint32

const (
	// AgentSupportsUpgrade will be set if this Agent can be upgraded remotely.
	AgentSupportsUpgrade AgentFeatures = 1 << iota

	// AgentSupportsSnapshots will be set if this Agent can send snapshots of recent telemetry signals.
	AgentSupportsSnapshots

	// AgentSupportsMeasurements will be set if this Agent supports the throughputmeasurement processor for measuring throughput.
	AgentSupportsMeasurements

	// AgentSupportsLogBasedMetrics will be set if this agent supports the logcount processor and route receiver for creating metrics from logs.
	AgentSupportsLogBasedMetrics
)

// AgentFeaturesDefault is the default bitmask of features supported by Agents
const AgentFeaturesDefault = AgentSupportsUpgrade

// Agent TODO(doc)
type Agent struct {
	ID              string `json:"id" yaml:"id"`
	Name            string `json:"name" yaml:"name"`
	Type            string `json:"type" yaml:"type"`
	Architecture    string `json:"arch" yaml:"arch"`
	HostName        string `json:"hostname" yaml:"hostname"`
	Labels          Labels `json:"labels,omitempty" yaml:"labels"`
	Version         string `json:"version" yaml:"version"`
	Home            string `json:"home" yaml:"home"`
	Platform        string `json:"platform" yaml:"platform"`
	OperatingSystem string `json:"operatingSystem" yaml:"operatingSystem"`
	MacAddress      string `json:"macAddress" yaml:"macAddress"`
	RemoteAddress   string `json:"remoteAddress,omitempty" yaml:"remoteAddress,omitempty"`

	// SecretKey is provided by the agent to authenticate
	SecretKey string `json:"-" yaml:"-"`

	// Upgrade stores information about an agent upgrade
	Upgrade *AgentUpgrade `json:"upgrade,omitempty" yaml:"upgrade,omitempty"`

	// reported by Status messages
	Status       AgentStatus `json:"status"`
	ErrorMessage string      `json:"errorMessage,omitempty" yaml:"errorMessage,omitempty"`

	// tracked by BindPlane
	Configuration  interface{} `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	ConnectedAt    *time.Time  `json:"connectedAt,omitempty" yaml:"connectedAt,omitempty"`
	DisconnectedAt *time.Time  `json:"disconnectedAt,omitempty" yaml:"disconnectedAt,omitempty"`

	// used by the agent management protocol
	Protocol string      `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	State    interface{} `json:"state,omitempty" yaml:"state,omitempty"`
}

var _ search.Indexed = (*Agent)(nil)
var _ HasUniqueKey = (*Agent)(nil)
var _ Labeled = (*Agent)(nil)

// UniqueKey returns the agent ID to uniquely identify an Agent
func (a *Agent) UniqueKey() string {
	return a.ID
}

// StatusDisplayText returns the string representation of the agent's status.
func (a *Agent) StatusDisplayText() string {
	switch a.Status {
	case Disconnected:
		return "Disconnected"
	case Connected:
		return "Connected"
	case Error:
		return "Error"
	case ComponentFailed:
		return "Component Failed"
	case Deleted:
		return "Deleted"
	case Configuring:
		return "Configuring"
	case Upgrading:
		return "Upgrading"
	default:
		return "Unknown"
	}
}

// GetLabels implements the Labeled interface for Agents
func (a *Agent) GetLabels() Labels {
	return a.Labels
}

// ConnectedDurationDisplayText TODO(doc)
func (a *Agent) ConnectedDurationDisplayText() string {
	if a.Status == Disconnected {
		return "-"
	}
	return durationDisplay(a.ConnectedAt)
}

// DisconnectedDurationDisplayText TODO(doc) What RFC?
func (a *Agent) DisconnectedDurationDisplayText() string {
	return durationDisplay(a.DisconnectedAt)
}

// MatchesSelector returns true if the given selector matches the agent's labels.
func (a *Agent) MatchesSelector(selector Selector) bool {
	return selector.Matches(a.Labels)
}

// DisconnectedSince returns true if the agent has been disconnected since a given time.
func (a *Agent) DisconnectedSince(since time.Time) bool {
	return a.DisconnectedAt != nil || a.DisconnectedAt.Before(since)
}

// Connect updates the ConnectedAt and DisconnectedAt fields of the agent and should be called when the
// agent connects.
func (a *Agent) Connect(newAgentVersion string) {
	// only update ConnectedAt if this is a new version or never connected
	if a.Version != newAgentVersion || a.ConnectedAt == nil {
		now := time.Now()
		a.ConnectedAt = &now
	}
	a.DisconnectedAt = nil
}

// Disconnect updates the DisconnectedAt and Status fields of the agent and should be called when the agent disconnects.
func (a *Agent) Disconnect() {
	now := time.Now()
	a.DisconnectedAt = &now
	a.Status = Disconnected
}

func durationDisplay(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "-"
	}
	return time.Since(*t).Round(time.Second).String()
}

// ----------------------------------------------------------------------
// features

var v1_6_0 = semver.Parse("1.6.0")
var v1_8_0 = semver.Parse("1.8.0")
var v1_9_2 = semver.Parse("1.9.2")
var v1_14_0 = semver.Parse("1.14.0")

// Features returns a bitmask of the features supported by this Agent
func (a *Agent) Features() AgentFeatures {
	agentVersion := semver.Parse(a.Version)

	// arrange version checks newest first
	switch {
	// 1.14.0 introduced the logcount processor and route receiver
	case !agentVersion.IsOlder(v1_14_0):
		return AgentSupportsUpgrade | AgentSupportsSnapshots | AgentSupportsMeasurements | AgentSupportsLogBasedMetrics
	// 1.9.2 introduced the throughputmeasurement processor
	case !agentVersion.IsOlder(v1_9_2):
		return AgentSupportsUpgrade | AgentSupportsSnapshots | AgentSupportsMeasurements

	// 1.8.0 introduced snapshots
	case !agentVersion.IsOlder(v1_8_0):
		return AgentSupportsUpgrade | AgentSupportsSnapshots

	// 1.6.0 introduced upgrade
	case !agentVersion.IsOlder(v1_6_0):
		return AgentSupportsUpgrade

	}
	return 0
}

// HasFeatures returns true if this Agent supports the specified feature or features.
func (a *Agent) HasFeatures(feature AgentFeatures) bool {
	return a.Features().Has(feature)
}

// Has returns true if the AgentFeatures bitmask has the specified features enabled. This is simply implemented as
// agentFeatures&feature != 0 but I prefer to avoid using bitmasks directly when possible.
func (agentFeatures AgentFeatures) Has(feature AgentFeatures) bool {
	return agentFeatures&feature != 0
}

// ----------------------------------------------------------------------
// upgrading

// SupportsUpgrade returns true if this agent supports upgrade
func (a *Agent) SupportsUpgrade() bool {
	// Ideally this would be based on the opamp flag AgentCapabilities_AcceptsPackages but agent capabilities aren't
	// currently available on the Agent model. That should change but for now the version will be checked.
	return a.HasFeatures(AgentSupportsUpgrade)
}

// UpgradeTo begins an upgrade by setting the status to Upgrading and setting the Upgrade field to the specified
// version.
func (a *Agent) UpgradeTo(version string) {
	if !a.SupportsUpgrade() {
		return
	}
	a.Upgrade = &AgentUpgrade{
		Version: version,
		Status:  UpgradePending,
	}
	a.Status = Upgrading
}

// UpgradeStarted it set when the upgrade instructions have actually been sent to the Agent.
func (a *Agent) UpgradeStarted(version string, allPackagesHash []byte) {
	a.Upgrade = &AgentUpgrade{
		Version:         version,
		Status:          UpgradeStarted,
		AllPackagesHash: allPackagesHash,
	}
	a.Status = Upgrading
}

// UpgradeComplete completes an upgrade by setting the status back to either Connected or Error (depending on
// ErrorMessage) and either removing the AgentUpgrade field or setting the Error on it if the specified errorMessage is
// not empty.
func (a *Agent) UpgradeComplete(version, errorMessage string) {
	if errorMessage != "" {
		// set the errorMessage on the AgentUpgrade
		if a.Upgrade == nil {
			a.Upgrade = &AgentUpgrade{}
		}
		a.Upgrade.Status = UpgradeFailed
		a.Upgrade.Error = errorMessage
		if version != "" {
			a.Upgrade.Version = version
		}
	} else {
		// clear the upgrade and reset the Status. if the Status is Configuring, allow the configuring process to continue.
		a.Upgrade = nil
	}
	if a.Status != Configuring {
		// preserve Error state if a configuration error exists after update.
		if a.ErrorMessage != "" {
			a.Status = Error
		} else {
			a.Status = Connected
		}
	}
}

// ----------------------------------------------------------------------
// sorting

type byName []*Agent

func (s byName) Len() int {
	return len(s)
}
func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// SortAgentsByName TODO(doc)
func SortAgentsByName(agents []*Agent) {
	sort.Sort(byName(agents))
}

// ----------------------------------------------------------------------
// indexing

// IndexID returns an ID used to identify the resource that is indexed
func (a *Agent) IndexID() string {
	return a.ID
}

// IndexFields returns a map of field name to field value to be stored in the index
func (a *Agent) IndexFields(index search.Indexer) {
	index("id", a.ID)
	index("arch", a.Architecture)
	index("hostname", a.HostName)
	index("platform", a.Platform)
	index("version", a.Version)
	index("name", a.Name)
	index("home", a.Home)
	index("os", a.OperatingSystem)
	index("macAddress", a.MacAddress)
	index("type", a.Type)
	index("status", a.StatusDisplayText())
}

// IndexLabels returns a map of label name to label value to be stored in the index
func (a *Agent) IndexLabels(index search.Indexer) {
	for n, v := range a.Labels.Set {
		index(n, v)
	}
}

// ----------------------------------------------------------------------
// Printable

// PrintableKindSingular returns the singular form of the Kind, e.g. "Agent"
func (a *Agent) PrintableKindSingular() string {
	return "Agent"
}

// PrintableKindPlural returns the singular form of the Kind, e.g. "Agents"
func (a *Agent) PrintableKindPlural() string {
	return "Agents"
}

// PrintableFieldTitles returns the list of field titles, used for printing a table of resources
func (a *Agent) PrintableFieldTitles() []string {
	return []string{"ID", "Name", "Version", "Status", "Connected", "Disconnected", "Labels"}
}

// PrintableFieldValue returns the field value for a title, used for printing a table of resources
func (a *Agent) PrintableFieldValue(title string) string {
	switch title {
	case "ID":
		return a.ID
	case "Name":
		return a.Name
	case "Version":
		return a.Version
	case "Status":
		return a.StatusDisplayText()
	case "Connected":
		return a.ConnectedDurationDisplayText()
	case "Disconnected":
		return a.DisconnectedDurationDisplayText()
	case "Labels":
		return a.Labels.Custom().String()
	}
	return ""
}
