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
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/observiq/bindplane-op/internal/store/search"
	"github.com/observiq/bindplane-op/model/graph"
	"github.com/observiq/bindplane-op/model/otel"
	"github.com/observiq/bindplane-op/model/validation"
	otelExt "go.opentelemetry.io/otel"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var tracer = otelExt.Tracer("model/configuration")

// ConfigurationType indicates the kind of configuration. It is based on the presence of the Raw, Sources, and
// Destinations fields.
type ConfigurationType string

const (
	// ConfigurationTypeRaw configurations have a configuration in the Raw field that is passed directly to the agent.
	ConfigurationTypeRaw ConfigurationType = "raw"

	// ConfigurationTypeModular configurations have Sources and Destinations that are used to generate the configuration to pass to an agent.
	ConfigurationTypeModular ConfigurationType = "modular"
	// TODO(andy): Do we like Modular for configurations with Sources/Destinations?
)

// Configuration is the resource for the entire agent configuration
type Configuration struct {
	// ResourceMeta TODO(doc)
	ResourceMeta `yaml:",inline" json:",inline" mapstructure:",squash"`
	// Spec TODO(doc)
	Spec ConfigurationSpec `json:"spec" yaml:"spec" mapstructure:"spec"`
}

var _ HasAgentSelector = (*Configuration)(nil)

// NewConfiguration creates a new configuration with the specified name
func NewConfiguration(name string) *Configuration {
	return NewConfigurationWithSpec(name, ConfigurationSpec{})
}

// NewRawConfiguration creates a new configuration with the specified name and raw configuration
func NewRawConfiguration(name string, raw string) *Configuration {
	return NewConfigurationWithSpec(name, ConfigurationSpec{
		Raw: raw,
	})
}

// NewConfigurationWithSpec creates a new configuration with the specified name and spec
func NewConfigurationWithSpec(name string, spec ConfigurationSpec) *Configuration {
	return &Configuration{
		ResourceMeta: ResourceMeta{
			APIVersion: V1,
			Kind:       KindConfiguration,
			Metadata: Metadata{
				Name:   name,
				Labels: MakeLabels(),
			},
		},
		Spec: spec,
	}
}

// GetKind returns "Configuration"
func (c *Configuration) GetKind() Kind {
	return KindConfiguration
}

// ConfigurationSpec is the spec for a configuration resource
type ConfigurationSpec struct {
	ContentType  string                  `json:"contentType" yaml:"contentType" mapstructure:"contentType"`
	Raw          string                  `json:"raw,omitempty" yaml:"raw,omitempty" mapstructure:"raw"`
	Sources      []ResourceConfiguration `json:"sources,omitempty" yaml:"sources,omitempty" mapstructure:"sources"`
	Destinations []ResourceConfiguration `json:"destinations,omitempty" yaml:"destinations,omitempty" mapstructure:"destinations"`
	Selector     AgentSelector           `json:"selector" yaml:"selector" mapstructure:"selector"`
}

// ResourceConfiguration represents a source or destination configuration
type ResourceConfiguration struct {
	Name              string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name"`
	ParameterizedSpec `yaml:",inline" json:",inline" mapstructure:",squash"`
}

// Validate validates most of the configuration, but if a store is available, ValidateWithStore should be used to
// validate the sources and destinations.
func (c *Configuration) Validate() (warnings string, errors error) {
	errs := validation.NewErrors()
	c.validate(errs)
	return errs.Warnings(), errs.Result()
}

func (c *Configuration) validate(errs validation.Errors) {
	c.ResourceMeta.validate(errs)
	c.Spec.validate(errs)
}

// ValidateWithStore checks that the configuration is valid, returning an error if it is not. It uses the store to
// retrieve source types and destination types so that parameter values can be validated against the parameter
// definitions.
func (c *Configuration) ValidateWithStore(ctx context.Context, store ResourceStore) (warnings string, errors error) {
	errs := validation.NewErrors()

	c.validate(errs)
	c.Spec.validateSourcesAndDestinations(ctx, errs, store)

	return errs.Warnings(), errs.Result()
}

// Type returns the ConfigurationType. It is based on the presence of the Raw, Sources, and Destinations fields.
func (c *Configuration) Type() ConfigurationType {
	if c.Spec.Raw != "" {
		// we always prefer raw
		return ConfigurationTypeRaw
	}
	return ConfigurationTypeModular
}

// AgentSelector returns the Selector for this configuration that can be used to match this resource to agents.
func (c *Configuration) AgentSelector() Selector {
	return c.Spec.Selector.Selector()
}

// IsForAgent returns true if this configuration matches a given agent's labels.
func (c *Configuration) IsForAgent(agent *Agent) bool {
	return isResourceForAgent(c, agent)
}

// ResourceStore provides access to resources required to render configurations that use Sources and Destinations.
type ResourceStore interface {
	Source(ctx context.Context, name string) (*Source, error)
	SourceType(ctx context.Context, name string) (*SourceType, error)
	Processor(ctx context.Context, name string) (*Processor, error)
	ProcessorType(ctx context.Context, name string) (*ProcessorType, error)
	Destination(ctx context.Context, name string) (*Destination, error)
	DestinationType(ctx context.Context, name string) (*DestinationType, error)
}

// BindPlaneConfiguration includes configuration information needed to render configurations
type BindPlaneConfiguration interface {
	BindPlaneURL() string
}

// Render converts the Configuration model to a configuration yaml that can be sent to an agent. The specified Agent can
// be nil if this configuration is not being rendered for a specific agent.
func (c *Configuration) Render(ctx context.Context, agent *Agent, config BindPlaneConfiguration, store ResourceStore) (string, error) {
	ctx, span := tracer.Start(ctx, "model/Configuration/Render")
	defer span.End()

	if c.Spec.Raw != "" {
		// we always prefer raw
		return c.Spec.Raw, nil
	}
	return c.renderComponents(ctx, agent, config, store)
}

func (c *Configuration) renderComponents(ctx context.Context, agent *Agent, config BindPlaneConfiguration, store ResourceStore) (string, error) {
	configuration, err := c.otelConfiguration(ctx, agent, config, store)
	if err != nil {
		return "", err
	}
	return configuration.YAML()
}

type renderContext struct {
	*otel.RenderContext
	pipelineTypeUsage *PipelineTypeUsage
}

func (c *Configuration) otelConfiguration(ctx context.Context, agent *Agent, config BindPlaneConfiguration, store ResourceStore) (*otel.Configuration, error) {
	if len(c.Spec.Sources) == 0 || len(c.Spec.Destinations) == 0 {
		return nil, nil
	}

	agentID := ""
	agentFeatures := AgentFeaturesDefault
	if agent != nil {
		agentID = agent.ID
		agentFeatures = agent.Features()
	}

	rc := &renderContext{
		RenderContext:     otel.NewRenderContext(agentID, c.Name(), config.BindPlaneURL()),
		pipelineTypeUsage: newPipelineTypeUsage(),
	}
	rc.IncludeSnapshotProcessor = agentFeatures.Has(AgentSupportsSnapshots)
	rc.IncludeMeasurements = agentFeatures.Has(AgentSupportsMeasurements)

	return c.otelConfigurationWithRenderContext(ctx, rc, store)
}

func (c *Configuration) otelConfigurationWithRenderContext(ctx context.Context, rc *renderContext, store ResourceStore) (*otel.Configuration, error) {
	configuration := otel.NewConfiguration()

	// match each source with each destination to produce a pipeline
	sources, destinations, err := c.evalComponents(ctx, store, rc)
	if err != nil {
		return nil, err
	}

	// to keep configurations consistent, iterate over the sorted keys instead of just iterating over the map directly.
	sourceNames := maps.Keys(sources)
	destinationNames := maps.Keys(destinations)
	sort.Strings(sourceNames)
	sort.Strings(destinationNames)

	for _, sourceName := range sourceNames {
		source := sources[sourceName]
		for _, destinationName := range destinationNames {
			destination := destinations[destinationName]

			name := fmt.Sprintf("%s__%s", sourceName, destinationName)
			configuration.AddPipeline(name, otel.Logs, sourceName, source, destinationName, destination, rc.RenderContext)
			configuration.AddPipeline(name, otel.Metrics, sourceName, source, destinationName, destination, rc.RenderContext)
			configuration.AddPipeline(name, otel.Traces, sourceName, source, destinationName, destination, rc.RenderContext)
		}
	}

	configuration.AddAgentMetricsPipeline(rc.RenderContext)

	return configuration, nil
}

func (c *Configuration) evalComponents(ctx context.Context, store ResourceStore, rc *renderContext) (sources map[string]otel.Partials, destinations map[string]otel.Partials, err error) {
	errorHandler := func(e error) {
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	sources = map[string]otel.Partials{}
	destinations = map[string]otel.Partials{}

	for i, source := range c.Spec.Sources {
		source := source // copy to local variable to securely pass a reference to a loop variable
		sourceName, srcParts := evalSource(ctx, &source, fmt.Sprintf("source%d", i), store, rc, errorHandler)
		sources[sourceName] = srcParts
	}

	for i, destination := range c.Spec.Destinations {
		destination := destination // copy to local variable to securely pass a reference to a loop variable
		destName, destParts := evalDestination(ctx, &destination, fmt.Sprintf("destination%d", i), store, rc, errorHandler)
		destinations[destName] = destParts
	}

	return sources, destinations, err
}

func evalSource(ctx context.Context, source *ResourceConfiguration, defaultName string, store ResourceStore, rc *renderContext, errorHandler TemplateErrorHandler) (string, otel.Partials) {
	src, srcType, err := findSourceAndType(ctx, source, defaultName, store)
	if err != nil {
		errorHandler(err)
		return "", nil
	}

	srcName := src.Name()
	if src.Spec.Disabled {
		return srcName, otel.NewPartials()
	}
	partials := srcType.eval(src, errorHandler)

	if rc.pipelineTypeUsage != nil {
		rc.pipelineTypeUsage.sources.setSupported(srcName, partials)
	}

	addMeasureProcessors(partials, MeasurementPositionSourceBeforeProcessors, src.Name(), rc)

	// evaluate the processors associated with the source
	for i, processor := range source.Processors {
		processor := processor
		_, processorParts := evalProcessor(ctx, &processor, fmt.Sprintf("%s__processor%d", srcName, i), store, rc, errorHandler)
		if processorParts == nil {
			continue
		}
		partials.Append(processorParts)
	}

	addMeasureProcessors(partials, MeasurementPositionSourceAfterProcessors, src.Name(), rc)

	return srcName, partials
}

func evalProcessor(ctx context.Context, processor *ResourceConfiguration, defaultName string, store ResourceStore, rc *renderContext, errorHandler TemplateErrorHandler) (string, otel.Partials) {
	prc, prcType, err := findProcessorAndType(ctx, processor, defaultName, store)
	if err != nil {
		errorHandler(err)
		return "", nil
	}

	return prc.Name(), prcType.eval(prc, errorHandler)
}

func evalDestination(ctx context.Context, destination *ResourceConfiguration, defaultName string, store ResourceStore, rc *renderContext, errorHandler TemplateErrorHandler) (string, otel.Partials) {
	dest, destType, err := findDestinationAndType(ctx, destination, defaultName, store)
	if err != nil {
		errorHandler(err)
		return "", nil
	}

	destName := dest.Name()
	if dest.Spec.Disabled || destination.Disabled {
		return destName, otel.NewPartials()
	}
	partials := destType.eval(dest, errorHandler)

	if rc.pipelineTypeUsage != nil {
		rc.pipelineTypeUsage.sources.setSupported(destName, partials)
	}

	d0partials := otel.NewPartials()
	addMeasureProcessors(d0partials, MeasurementPositionDestinationBeforeProcessors, destName, rc)

	destProcessors := otel.NewPartials()
	// evaluate the processors associated with the destination
	for i, processor := range destination.Processors {
		processor := processor
		_, processorParts := evalProcessor(ctx, &processor, fmt.Sprintf("%s__processor%d", destName, i), store, rc, errorHandler)
		if processorParts == nil {
			continue
		}
		destProcessors.Append(processorParts)
	}

	d1partials := otel.NewPartials()
	addMeasureProcessors(d1partials, MeasurementPositionDestinationAfterProcessors, destName, rc)

	// destination processors are prepended to the destination
	partials.Prepend(d1partials)
	if !destProcessors.Empty() {
		partials.Prepend(destProcessors)
		partials.Prepend(d0partials)
	}

	return destName, partials
}

func findSourceAndType(ctx context.Context, source *ResourceConfiguration, defaultName string, store ResourceStore) (*Source, *SourceType, error) {
	src, err := FindSource(ctx, source, defaultName, store)
	if err != nil {
		return nil, nil, err
	}

	srcType, err := store.SourceType(ctx, src.Spec.Type)
	if err == nil && srcType == nil {
		err = fmt.Errorf("unknown %s: %s", KindSourceType, src.Spec.Type)
	}
	if err != nil {
		return src, nil, err
	}

	return src, srcType, nil
}

func findProcessorAndType(ctx context.Context, source *ResourceConfiguration, defaultName string, store ResourceStore) (*Processor, *ProcessorType, error) {
	prc, err := FindProcessor(ctx, source, defaultName, store)
	if err != nil {
		return nil, nil, err
	}

	prcType, err := store.ProcessorType(ctx, prc.Spec.Type)
	if err == nil && prcType == nil {
		err = fmt.Errorf("unknown %s: %s", KindProcessorType, prc.Spec.Type)
	}
	if err != nil {
		return prc, nil, err
	}

	return prc, prcType, nil
}

func findDestinationAndType(ctx context.Context, destination *ResourceConfiguration, defaultName string, store ResourceStore) (*Destination, *DestinationType, error) {
	dest, err := FindDestination(ctx, destination, defaultName, store)
	if err != nil {
		return nil, nil, err
	}

	destType, err := store.DestinationType(ctx, dest.Spec.Type)
	if err == nil && destType == nil {
		err = fmt.Errorf("unknown %s: %s", KindDestinationType, dest.Spec.Type)
	}
	if err != nil {
		return dest, nil, err
	}

	return dest, destType, nil
}

func findResourceAndType(ctx context.Context, resourceKind Kind, resource *ResourceConfiguration, defaultName string, store ResourceStore) (Resource, *ResourceType, error) {
	switch resourceKind {
	case KindSource:
		src, srcType, err := findSourceAndType(ctx, resource, defaultName, store)
		if srcType == nil {
			return src, nil, err
		}
		return src, &srcType.ResourceType, err
	case KindProcessor:
		prc, prcType, err := findProcessorAndType(ctx, resource, defaultName, store)
		if prcType == nil {
			return prc, nil, err
		}
		return prc, &prcType.ResourceType, err
	case KindDestination:
		dest, destType, err := findDestinationAndType(ctx, resource, defaultName, store)
		if destType == nil {
			return dest, nil, err
		}
		return dest, &destType.ResourceType, err
	}
	return nil, nil, nil
}

// ----------------------------------------------------------------------

func (cs *ConfigurationSpec) validate(errors validation.Errors) {
	cs.validateSpecFields(errors)
	cs.validateRaw(errors)
	cs.Selector.validate(errors)
}

func (cs *ConfigurationSpec) validateSpecFields(errors validation.Errors) {
	if cs.Raw != "" {
		if len(cs.Destinations) > 0 || len(cs.Sources) > 0 {
			errors.Add(fmt.Errorf("configuration must specify raw or sources and destinations"))
		}
	}
}

func (cs *ConfigurationSpec) validateRaw(errors validation.Errors) {
	if cs.Raw == "" {
		return
	}
	parsed := map[string]any{}
	err := yaml.Unmarshal([]byte(cs.Raw), parsed)
	if err != nil {
		errors.Add(fmt.Errorf("unable to parse spec.raw as yaml: %w", err))
	}
}

func (cs *ConfigurationSpec) validateSourcesAndDestinations(ctx context.Context, errors validation.Errors, store ResourceStore) {
	for _, source := range cs.Sources {
		source.validate(ctx, KindSource, errors, store)
	}
	for _, destination := range cs.Destinations {
		destination.validate(ctx, KindDestination, errors, store)
	}
}

func (rc *ResourceConfiguration) localName(kind Kind, index int) string {
	if rc.Name != "" {
		return rc.Name
	}
	return fmt.Sprintf("%s%d", strings.ToLower(string(kind)), index)
}

func (rc *ResourceConfiguration) validate(ctx context.Context, resourceKind Kind, errors validation.Errors, store ResourceStore) {
	if rc.validateHasNameOrType(resourceKind, errors) {
		rc.validateParameters(ctx, resourceKind, errors, store)
	}
	rc.validateProcessors(ctx, resourceKind, errors, store)
}

func (rc *ResourceConfiguration) validateHasNameOrType(resourceKind Kind, errors validation.Errors) bool {
	// must have name or type
	if rc.Name == "" && rc.Type == "" {
		errors.Add(fmt.Errorf("all %s must have either a name or type", resourceKind))
		return false
	}
	return true
}

func (rc *ResourceConfiguration) validateParameters(ctx context.Context, resourceKind Kind, errors validation.Errors, store ResourceStore) {
	// must have a name
	for _, parameter := range rc.Parameters {
		if parameter.Name == "" {
			errors.Add(fmt.Errorf("all %s parameters must have a name", resourceKind))
		}
	}
	_, resourceType, err := findResourceAndType(ctx, resourceKind, rc, string(resourceKind), store)
	if err != nil {
		errors.Add(err)
		return
	}
	// ensure parameters are valid
	for _, parameter := range rc.Parameters {
		if parameter.Name == "" {
			continue
		}
		def := resourceType.Spec.ParameterDefinition(parameter.Name)
		if def == nil {
			errors.Warn(fmt.Errorf("ignoring parameter %s not defined in type %s", parameter.Name, resourceType.Name()))
			continue
		}
		err := def.validateValue(parameter.Value)
		if err != nil {
			errors.Add(err)
		}
	}
}

func (rc *ResourceConfiguration) validateProcessors(ctx context.Context, resourceKind Kind, errors validation.Errors, store ResourceStore) {
	for _, processor := range rc.Processors {
		processor.validate(ctx, KindProcessor, errors, store)
	}
}

// ----------------------------------------------------------------------
// Printable

// PrintableFieldTitles returns the list of field titles, used for printing a table of resources
func (c *Configuration) PrintableFieldTitles() []string {
	return []string{"Name", "Match"}
}

// PrintableFieldValue returns the field value for a title, used for printing a table of resources
func (c *Configuration) PrintableFieldValue(title string) string {
	switch title {
	case "Name":
		return c.Name()
	case "Match":
		return c.AgentSelector().String()
	default:
		return "-"
	}
}

// ----------------------------------------------------------------------
// Indexed

// IndexFields returns a map of field name to field value to be stored in the index
func (c *Configuration) IndexFields(index search.Indexer) {
	c.ResourceMeta.IndexFields(index)

	// add the type of configuration
	index("type", string(c.Type()))

	// add source, sourceType fields
	for _, source := range c.Spec.Sources {
		source.indexFields("source", "sourceType", index)
	}

	// add destination, destinationType fields
	for _, destination := range c.Spec.Destinations {
		destination.indexFields("destination", "destinationType", index)
	}

	// add pipeline fields
	//
	// TODO(andy): I was going to add pipeline:traces, pipeline:logs, and pipeline:metrics because I thought it would be a
	// useful way to filter configurations. However, we need a ResourceStore implementation to call otelConfiguration and
	// we don't have that here, even though indexing is actually done in the store. I think the best solution is to cache
	// the output on the Spec and keep that up to date as any dependent sourceTypes and destinationTypes change. This will
	// improve performance when comparing configurations and displaying the configuration in UI.
}

func (rc *ResourceConfiguration) indexFields(resourceName string, resourceTypeName string, index search.Indexer) {
	index(resourceName, rc.Name)
	index(resourceTypeName, rc.Type)
}

// Duplicate copies the value of the current configuration and returns
// a duplicate with the new name.  It should be identical except for the
// Metadata.Name, Metadata.ID, and Spec.Selector fields.
func (c *Configuration) Duplicate(name string) *Configuration {
	copy := *c

	// Change the metadata values
	copy.Metadata.Name = name
	copy.Metadata.ID = uuid.NewString()

	// replace the configuration matchLabel
	matchLabels := copy.Spec.Selector.MatchLabels
	matchLabels["configuration"] = name
	return &copy
}

// ----------------------------------------------------------------------
// topology

// Graph returns a graph representing the topology of a configuration

// Graph returns a graph representing the topology of a configuration
func (c *Configuration) Graph(ctx context.Context, store ResourceStore) (*graph.Graph, error) {
	g := graph.NewGraph()

	// lastNodes is a list of the last node for each source that will be connected to the destinations
	lastNodes := make([]*graph.Node, 0, len(c.Spec.Sources))

	pipelineUsage := c.determinePipelineTypeUsage(ctx, store)
	g.Attributes["activeTypeFlags"] = pipelineUsage.ActiveFlags()

	for i, source := range c.Spec.Sources {
		sourceName := source.localName(KindSource, i)
		usage := pipelineUsage.sources.usage(sourceName)

		attributes := graph.MakeAttributes(string(KindSource), sourceName)
		attributes["activeTypeFlags"] = usage.active
		attributes["supportedTypeFlags"] = usage.supported
		s := &graph.Node{
			ID:         fmt.Sprintf("source/%s", sourceName),
			Type:       "sourceNode",
			Label:      source.Type,
			Attributes: attributes,
		}
		g.AddSource(s)

		// For now only add one intermediate node for each
		// source which represents all the processors on the source.
		p := &graph.Node{
			ID:    fmt.Sprintf("source/%s/processors", sourceName),
			Type:  "processorNode",
			Label: "Processors",
			Attributes: map[string]any{
				"activeTypeFlags":    usage.active,
				"supportedTypeFlags": usage.supported,
			},
		}
		g.AddIntermediate(p)
		g.Connect(s, p)

		lastNodes = append(lastNodes, p)
	}

	for i, destination := range c.Spec.Destinations {
		destinationName := destination.localName(KindDestination, i)
		usage := pipelineUsage.destinations.usage(destinationName)

		// For now only add one intermediate node for each
		// destination which represents all the processors on the destination.
		p := &graph.Node{
			ID:    fmt.Sprintf("destination/%s/processors", destinationName),
			Type:  "processorNode",
			Label: "Processors",
			Attributes: map[string]any{
				"activeTypeFlags":    usage.active,
				"supportedTypeFlags": usage.supported,
			},
		}
		g.AddIntermediate(p)
		for _, l := range lastNodes {
			g.Connect(l, p)
		}

		attributes := graph.MakeAttributes(string(KindDestination), destination.Name)
		// Needed to determine which destination card to display.
		attributes["isInline"] = destination.Name == ""
		attributes["activeTypeFlags"] = usage.active
		attributes["supportedTypeFlags"] = usage.supported
		d := &graph.Node{
			ID:         fmt.Sprintf("destination/%s", destinationName),
			Type:       "destinationNode",
			Label:      destination.Name,
			Attributes: attributes,
		}
		g.AddTarget(d)
		g.Connect(p, d)
	}

	return g, nil
}

// MeasurementPosition is a position within the graph of the measurements processor
type MeasurementPosition string

const (
	// MeasurementPositionSourceBeforeProcessors is the initial throughput of the source
	MeasurementPositionSourceBeforeProcessors MeasurementPosition = "s0"

	// MeasurementPositionSourceAfterProcessors is the throughput after source processors
	MeasurementPositionSourceAfterProcessors MeasurementPosition = "s1"

	// MeasurementPositionDestinationBeforeProcessors is the throughput to the destination (from all sources) before
	// destination processors
	MeasurementPositionDestinationBeforeProcessors MeasurementPosition = "d0"

	// MeasurementPositionDestinationAfterProcessors is the throughput to the destination (from all sources) after
	// destination processors
	MeasurementPositionDestinationAfterProcessors MeasurementPosition = "d1"
)

func addMeasureProcessors(partials otel.Partials, position MeasurementPosition, resourceName string, rc *renderContext) {
	if !rc.IncludeMeasurements {
		return
	}
	for _, pipelineType := range []otel.PipelineType{otel.Logs, otel.Metrics, otel.Traces} {
		processorName := otel.ComponentID(fmt.Sprintf("%s/_%s_%s_%s", otel.MeasureProcessorName, position, pipelineType, resourceName))
		addMeasureProcessor(partials[pipelineType], processorName)
	}
}

func addMeasureProcessor(partial *otel.Partial, processorName otel.ComponentID) {
	partial.Processors = append(partial.Processors, otel.ComponentMap{
		processorName: map[string]any{
			"enabled":        true,
			"sampling_ratio": 1,
		},
	})
}

type pipelineUsage struct {
	active    otel.PipelineTypeFlags
	supported otel.PipelineTypeFlags
}

type pipelineTypeUsageMap map[string]*pipelineUsage

func (p pipelineTypeUsageMap) usage(name string) *pipelineUsage {
	if u, ok := p[name]; ok {
		return u
	}
	result := &pipelineUsage{}
	p[name] = result
	return result
}

func (p pipelineTypeUsageMap) setActive(name string, pipelineType otel.PipelineType) {
	p.usage(name).active.Set(pipelineType.Flag())
}

func (p pipelineTypeUsageMap) setSupported(name string, partials otel.Partials) {
	p.usage(name).supported.Set(partials.PipelineTypes())
}

// PipelineTypeUsage contains information about active telemetry on the Configuration
// and its sources and destinations.
type PipelineTypeUsage struct {
	sources      pipelineTypeUsageMap
	destinations pipelineTypeUsageMap

	// active refers to the top level configuration
	active otel.PipelineTypeFlags
}

// ActiveFlagsForDestination returns the PipelineTypeFlags that are active for a destination with given name.
func (p *PipelineTypeUsage) ActiveFlagsForDestination(name string) otel.PipelineTypeFlags {
	return p.destinations.usage(name).active
}

// ActiveFlags returns the pipeline type flags that are in use by the configuration.
func (p *PipelineTypeUsage) ActiveFlags() otel.PipelineTypeFlags {
	return p.active
}

// setActive sets the top level activeFlags for a pipelineTypeUsage
func (p *PipelineTypeUsage) setActive(t otel.PipelineType) {
	p.active.Set(t.Flag())
}

func newPipelineTypeUsage() *PipelineTypeUsage {
	return &PipelineTypeUsage{
		sources:      pipelineTypeUsageMap{},
		destinations: pipelineTypeUsageMap{},
	}
}

func (c *Configuration) determinePipelineTypeUsage(ctx context.Context, store ResourceStore) *PipelineTypeUsage {
	p := newPipelineTypeUsage()

	// the agent ID and URL values aren't important
	rc := &renderContext{
		RenderContext:     otel.NewRenderContext("AGENT_ID", c.Name(), "BINDPLANE_URL"),
		pipelineTypeUsage: p,
	}
	config, err := c.otelConfigurationWithRenderContext(ctx, rc, store)
	if err != nil {
		// pipeline type usage won't be available if there is an error rendering the configuration, but that's ok.
		return p
	}

	for _, pipeline := range config.Service.Pipelines {
		p.sources.setActive(pipeline.SourceName(), pipeline.Type())
		p.destinations.setActive(pipeline.DestinationName(), pipeline.Type())

		p.setActive(pipeline.Type())
	}

	return p
}

// Usage returns a PipelineTypeUsage struct which contains information about the active and
// supported telemetry types on a configuration.
func (c *Configuration) Usage(ctx context.Context, store ResourceStore) *PipelineTypeUsage {
	return c.determinePipelineTypeUsage(ctx, store)
}
