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
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

// PipelineType is the telemetry type specified in pipeline names, e.g. metrics in metrics/redis.
type PipelineType string

// OpenTelemetry currently supports "metrics", "logs", and "traces"
const (
	Metrics PipelineType = "metrics"
	Logs    PipelineType = "logs"
	Traces  PipelineType = "traces"
)

// Flag returns the PipelineTypeFlags corresponding to this PipelineType
func (pt PipelineType) Flag() PipelineTypeFlags {
	switch pt {
	case Logs:
		return LogsFlag
	case Metrics:
		return MetricsFlag
	case Traces:
		return TracesFlag
	}
	return 0
}

// PipelineTypeFlags expresses a set of PipelineTypes
type PipelineTypeFlags byte

// OpenTelemetry currently supports "metrics", "logs", and "traces"
const (
	LogsFlag PipelineTypeFlags = 1 << iota
	MetricsFlag
	TracesFlag
)

// IncludesType returns true if the specified type is included
func (p PipelineTypeFlags) IncludesType(pipelineType PipelineType) bool {
	return p.Includes(pipelineType.Flag())
}

// Includes returns true if the specified flags are included
func (p PipelineTypeFlags) Includes(flags PipelineTypeFlags) bool {
	return flags&p != 0
}

// Set will return a new PipelineTypeFlags including the specified flags
func (p *PipelineTypeFlags) Set(flags PipelineTypeFlags) {
	*p = *p | flags
}

// ComponentID is a the name of an individual receiver, processor, exporter, or extension.
type ComponentID string

// ComponentMap is a map of individual receivers, processors, etc.
type ComponentMap map[ComponentID]any

// ComponentList is an ordered list of individual receivers, processors, etc. The order is especially important for
// processors.
type ComponentList []map[ComponentID]any

// NewComponentID creates a new ComponentID from its type and name
func NewComponentID(pipelineType, name string) ComponentID {
	return ComponentID(fmtComponentID(pipelineType, name))
}

func fmtComponentID(pipelineType, name string) string {
	if name == "" {
		return pipelineType
	}
	return fmt.Sprintf("%s/%s", pipelineType, name)
}

// ParseComponentID returns the typeName and name from a ComponentID
func ParseComponentID(id ComponentID) (pipelineType, name string) {
	parts := strings.SplitN(string(id), "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// ----------------------------------------------------------------------

var present = struct{}{}

// RenderContext keeps track of state needed to render configurations
type RenderContext struct {
	IncludeSnapshotProcessor bool
	IncludeMeasurements      bool
	AgentID                  string
	ConfigurationName        string
	BindPlaneURL             string
	measurementProcessors    map[ComponentID]struct{}
}

// NewRenderContext creates a new render context used to render configurations
func NewRenderContext(agentID string, configurationName string, bindplaneURL string) *RenderContext {
	return &RenderContext{
		AgentID:               agentID,
		ConfigurationName:     configurationName,
		BindPlaneURL:          bindplaneURL,
		measurementProcessors: map[ComponentID]struct{}{},
	}
}

func (rc *RenderContext) addMeasurementProcessor(processor ComponentID) {
	rc.measurementProcessors[processor] = present
}

func (rc *RenderContext) hasMeasurementProcessor(processor ComponentID) bool {
	_, ok := rc.measurementProcessors[processor]
	return ok
}

// ----------------------------------------------------------------------

// Configuration is a rough approximation of an OpenTelemetry configuration. It is used to help assemble a configuration
// and marshal it to a string to send to an agent.
type Configuration struct {
	Receivers  ComponentMap `yaml:"receivers,omitempty"`
	Processors ComponentMap `yaml:"processors,omitempty"`
	Exporters  ComponentMap `yaml:"exporters,omitempty"`
	Extensions ComponentMap `yaml:"extensions,omitempty"`
	Service    Service      `yaml:"service"`
}

// NewConfiguration creates a new configuration with initialized fields
func NewConfiguration() *Configuration {
	c := &Configuration{
		Receivers:  ComponentMap{},
		Processors: ComponentMap{},
		Exporters:  ComponentMap{},
		Extensions: ComponentMap{},
		Service: Service{
			Pipelines: Pipelines{},
		},
	}
	return c
}

// SnapshotProcessorName is the name of the snapshot processor that is inserted into each pipeline before the
// exporter
const SnapshotProcessorName ComponentID = "snapshotprocessor"

// MeasureProcessorName is the name of the measurement processor this is inserted into the pipeline to measure
// throughput
const MeasureProcessorName ComponentID = "throughputmeasurement"

// YAML marshals the configuration to yaml
func (c *Configuration) YAML() (string, error) {
	if c == nil || !c.HasPipelines() {
		return NoopConfig, nil
	}
	bytes, err := yaml.Marshal(c)
	return string(bytes), err
}

// HasPipelines returns true if there are pipelines
func (c *Configuration) HasPipelines() bool {
	return len(c.Service.Pipelines) > 0
}

// Service is the part of the configuration that defines the pipelines which consist of references to the components in
// the Configuration.
type Service struct {
	Extensions []ComponentID `yaml:"extensions,omitempty"`
	Pipelines  Pipelines     `yaml:"pipelines"`
}

// AddPipeline adds a pipeline to the map of pipelines for the service
func (s *Service) AddPipeline(p Pipeline) {
	s.Pipelines[p.ComponentID()] = p
}

// Pipelines are identified by a pipeline type and name in the form type/name where type is "metrics", "logs", or
// "traces"
type Pipelines map[ComponentID]Pipeline

// Pipeline is an ordered list of receivers, processors, and exporters.
type Pipeline struct {
	Receivers       []ComponentID `yaml:"receivers"`
	Processors      []ComponentID `yaml:"processors"`
	Exporters       []ComponentID `yaml:"exporters"`
	name            string
	pipelineType    PipelineType
	sourceName      string
	destinationName string
}

// NewPipeline creates a new pipeline with the specified type
func NewPipeline(pipelineType PipelineType, name string, sourceName string, destinationName string) Pipeline {
	return Pipeline{
		pipelineType:    pipelineType,
		name:            name,
		sourceName:      sourceName,
		destinationName: destinationName,
	}
}

// ComponentID returns the ComponentID of the
func (p *Pipeline) ComponentID() ComponentID {
	return ComponentID(fmt.Sprintf("%s/%s", p.pipelineType, p.name))
}

// Type returns the PipelineType of the pipeline
func (p *Pipeline) Type() PipelineType {
	return p.pipelineType
}

// SourceName returns the name of the Source responsible for the receiver(s) of this pipeline.
func (p *Pipeline) SourceName() string {
	return p.sourceName
}

// DestinationName returns the name of the Destination responsible for the exporters(s) of this pipeline.
func (p *Pipeline) DestinationName() string {
	return p.destinationName
}

// Incomplete returns true if there are zero Receivers or zero Exporters
func (p *Pipeline) Incomplete() bool {
	return len(p.Receivers) == 0 || len(p.Exporters) == 0
}

// Partial represents a fragment of configuration produced by an individual resource.
type Partial struct {
	Receivers  ComponentList
	Processors ComponentList
	Exporters  ComponentList
	Extensions ComponentList
}

// Size returns the number of components in the partial configuration
func (p *Partial) Size() int {
	return len(p.Receivers) + len(p.Processors) + len(p.Exporters) + len(p.Extensions)
}

// HasNoReceiversOrExporters returns true if this Partial doesn't have any receivers or exporters
func (p *Partial) HasNoReceiversOrExporters() bool {
	return len(p.Receivers)+len(p.Exporters) == 0
}

// Append adds components from another partial by appending each of the component lists together
func (p *Partial) Append(o *Partial) {
	p.Receivers = append(p.Receivers, o.Receivers...)
	p.Processors = append(p.Processors, o.Processors...)
	p.Exporters = append(p.Exporters, o.Exporters...)
	p.Extensions = append(p.Extensions, o.Extensions...)
}

func prepend[T any](list []T, others ...T) []T {
	var result []T
	result = append(result, others...)
	result = append(result, list...)
	return result
}

// Prepend adds components from another partial by prepending each of the component lists together.
func (p *Partial) Prepend(o *Partial) {
	p.Receivers = prepend(p.Receivers, o.Receivers...)
	p.Processors = prepend(p.Processors, o.Processors...)
	p.Exporters = prepend(p.Exporters, o.Exporters...)
	p.Extensions = prepend(p.Extensions, o.Extensions...)
}

// Partials represents a fragments of configuration for each type of telemetry.
type Partials map[PipelineType]*Partial

// NewPartials initializes a new Partials with empty Logs, Metrics, and Traces Partial configurations
func NewPartials() Partials {
	return map[PipelineType]*Partial{
		Logs:    {},
		Metrics: {},
		Traces:  {},
	}
}

// Append combines the individual Logs, Metrics, and Traces Partial configurations
func (p Partials) Append(o Partials) {
	p[Logs].Append(o[Logs])
	p[Metrics].Append(o[Metrics])
	p[Traces].Append(o[Traces])
}

// Prepend combines the individual Logs, Metrics, and Traces Partial configurations
func (p Partials) Prepend(o Partials) {
	p[Logs].Prepend(o[Logs])
	p[Metrics].Prepend(o[Metrics])
	p[Traces].Prepend(o[Traces])
}

// Empty returns true if there are no components for any of the individual partial configurations for Logs, Metrics, and Traces
func (p Partials) Empty() bool {
	return p[Logs].Size()+p[Metrics].Size()+p[Traces].Size() == 0
}

// PipelineTypes returns the list of PipelineTypes with components in these Partials.
func (p Partials) PipelineTypes() PipelineTypeFlags {
	var types PipelineTypeFlags
	for _, pipelineType := range []PipelineType{Logs, Metrics, Traces} {
		if p[pipelineType].Size() > 0 {
			types.Set(pipelineType.Flag())
		}
	}
	return types
}

// ComponentIDProvider can provide ComponentIDs for component names
type ComponentIDProvider interface {
	ComponentID(componentName string) ComponentID
}

// UniqueComponentID ensures that each ComponentID is unique by including the type and resource name. To make them easy
// to find in a completed configuration, we preserve the part before the / and then insert the type and resource name
// separated by 2 underscores.
func UniqueComponentID(original, typeName, resourceName string) ComponentID {
	// replace type/name with type/resourceType__resourceName__name
	pipelineType, name := ParseComponentID(ComponentID(original))

	var newName string
	if name != "" {
		newName = fmt.Sprintf("%s__%s", resourceName, name)
	} else {
		newName = resourceName
	}
	return NewComponentID(pipelineType, newName)
}

// AddExtensions adds all of the extensions to the configuration, replacing any extensions with the same id
func (c *Configuration) AddExtensions(extensions ComponentList) {
	for _, extension := range extensions {
		for n, v := range extension {
			c.AddExtension(n, v)
		}
	}
}

// AddExtension adds the specified extension with the specified id, replace any extension with the same id
func (c *Configuration) AddExtension(name ComponentID, extension any) {
	c.Extensions[name] = extension
	if !slices.Contains(c.Service.Extensions, name) {
		c.Service.Extensions = append(c.Service.Extensions, name)
	}
}

// AddAgentMetricsPipeline adds the measurements pipeline to the configuration
func (c *Configuration) AddAgentMetricsPipeline(rc *RenderContext) {
	if !rc.IncludeMeasurements {
		return
	}

	endpoint := fmt.Sprintf("%s/v1/otlphttp", rc.BindPlaneURL)
	otlphttp := map[string]any{
		"endpoint": endpoint,
	}
	if strings.HasPrefix(endpoint, "https") {
		otlphttp["tls"] = map[string]any{
			"insecure": true,
		}
	}

	parts := Partial{
		Receivers: ComponentList{{
			"prometheus/_agent_metrics": map[string]any{
				"config": map[string]any{
					"scrape_configs": []map[string]any{
						{
							"job_name":        "observiq-otel-collector",
							"scrape_interval": "10s",
							"static_configs": []map[string]any{
								{
									"targets": []string{"0.0.0.0:8888"},
									"labels": map[string]string{
										"configuration": rc.ConfigurationName,
										"agent":         rc.AgentID,
									},
								},
							},
							"metric_relabel_configs": []map[string]any{
								{
									"source_labels": []string{"__name__"},
									"action":        "keep",
									"regex":         "otelcol_processor_throughputmeasurement_.*",
								},
							},
						},
					},
				},
			},
		}},
		Processors: ComponentList{{
			"batch/_agent_metrics": nil,
		}},
		Exporters: ComponentList{{
			"otlphttp/_agent_metrics": otlphttp,
		}},
	}

	p := NewPipeline(Metrics, "_agent_metrics", "_agent_metrics", "_agent_metrics")
	p.AddReceivers(rc, c.Receivers.addComponents(parts.Receivers))
	p.AddProcessors(rc, c.Processors.addComponents(parts.Processors))
	p.AddExporters(rc, c.Exporters.addComponents(parts.Exporters))

	c.Service.AddPipeline(p)
}

// AddPipeline adds a pipeline and all of the corresponding components to the configuration
func (c *Configuration) AddPipeline(name string, pipelineType PipelineType, sourceName string, source Partials, destinationName string, destination Partials, rc *RenderContext) {
	s := source[pipelineType]
	d := destination[pipelineType]
	if s.HasNoReceiversOrExporters() || d.HasNoReceiversOrExporters() {
		// not all pipelineType will have components, ignore these
		return
	}

	p := NewPipeline(pipelineType, name, sourceName, destinationName)

	// add any receivers specified
	p.AddReceivers(rc, c.Receivers.addComponents(s.Receivers))
	p.AddReceivers(rc, c.Receivers.addComponents(d.Receivers))

	// add any processors specified
	p.AddProcessors(rc, c.Processors.addComponents(s.Processors))
	p.AddProcessors(rc, c.Processors.addComponents(d.Processors))

	if rc.IncludeSnapshotProcessor {
		// add the snapshot processor before the exporter
		p.AddProcessors(rc, c.Processors.addComponents(ComponentList{{SnapshotProcessorName: nil}}))
	}

	// add any exporters specified
	p.AddExporters(rc, c.Exporters.addComponents(s.Exporters))
	p.AddExporters(rc, c.Exporters.addComponents(d.Exporters))

	// skip any incomplete pipelines
	if p.Incomplete() {
		return
	}

	// any extensions are added and shared by all pipelines
	c.AddExtensions(s.Extensions)
	c.AddExtensions(d.Extensions)

	c.Service.AddPipeline(p)
}

// hasComponent returns true if the ComponentMap already has a component with the specified name
func (c ComponentMap) hasComponent(id ComponentID) bool {
	_, ok := c[id]
	return ok
}

// addComponents adds the components to the map and returns their ids as a convenience to build the pipeline
func (c ComponentMap) addComponents(componentList ComponentList) []ComponentID {
	ids := []ComponentID{}
	for _, components := range componentList {
		for id, component := range components {
			c[id] = component
			ids = append(ids, id)
		}
	}
	return ids
}

// AddReceivers adds receivers to the pipeline
func (p *Pipeline) AddReceivers(rc *RenderContext, id []ComponentID) {
	p.Receivers = append(p.Receivers, id...)
}

// AddProcessors adds processors to the pipeline
func (p *Pipeline) AddProcessors(rc *RenderContext, ids []ComponentID) {
	// when adding processors, only add measurement processors once
	for _, id := range ids {
		if isMeasurementProcessor(id) {
			if rc.hasMeasurementProcessor(id) {
				continue
			}
			rc.addMeasurementProcessor(id)
		}
		p.Processors = append(p.Processors, id)
	}
}

// AddExporters adds exporters to the pipeline
func (p *Pipeline) AddExporters(rc *RenderContext, id []ComponentID) {
	p.Exporters = append(p.Exporters, id...)
}

func isMeasurementProcessor(id ComponentID) bool {
	componentType, _ := ParseComponentID(id)
	return componentType == string(MeasureProcessorName)
}
