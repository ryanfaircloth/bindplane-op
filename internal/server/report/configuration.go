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

package report

import (
	"net/http"

	"github.com/observiq/bindplane-op/model/otel"
	"gopkg.in/yaml.v2"
)

// ConfigurationName is the name of the configuration file name
const ConfigurationName = "report.yaml"

// Configuration represents the "report.yaml" config sent to the agent via opamp
type Configuration struct {
	Snapshot Snapshot `json:"snapshot" yaml:"snapshot"`
}

// Snapshot contains the portion of the configuration specific to the snapshot capability
type Snapshot struct {
	// Processor is the full ComponentID of the snapshot processor
	Processor string `json:"processor" yaml:"processor"`

	// PipelineType will be "logs", "metrics", or "traces"
	PipelineType otel.PipelineType `json:"pipeline_type" yaml:"pipeline_type"`

	// Endpoint indicates where OTLP telemetry should be sent
	Endpoint Endpoint `json:"endpoint" yaml:"endpoint,omitempty"`
}

// Endpoint contains the headers and url where OTLP data should be sent
type Endpoint struct {
	// Header should be added as HTTP headers with the payload sent to the endpoint
	Header http.Header `json:"headers" yaml:"headers,omitempty"`

	// URL indicates the target for the HTTP POST with OTLP data
	URL string `json:"url" yaml:"url"`
}

// YAML returns the encoded Configuration
func (c *Configuration) YAML() ([]byte, error) {
	bytes, err := yaml.Marshal(c)
	return bytes, err
}
