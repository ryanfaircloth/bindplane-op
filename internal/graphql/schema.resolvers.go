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

package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/bindplane-op/internal/eventbus"
	"github.com/observiq/bindplane-op/internal/graphql/generated"
	model1 "github.com/observiq/bindplane-op/internal/graphql/model"
	"github.com/observiq/bindplane-op/internal/server"
	"github.com/observiq/bindplane-op/internal/server/report"
	"github.com/observiq/bindplane-op/internal/store"
	"github.com/observiq/bindplane-op/internal/util/semver"
	"github.com/observiq/bindplane-op/model"
	"github.com/observiq/bindplane-op/model/graph"
	"github.com/observiq/bindplane-op/model/otel"
	"go.uber.org/zap"
)

// Labels is the resolver for the labels field.
func (r *agentResolver) Labels(ctx context.Context, obj *model.Agent) (map[string]interface{}, error) {
	labels := map[string]interface{}{}
	for k := range obj.Labels.Set {
		labels[k] = obj.Labels.Get(k)
	}
	return labels, nil
}

// Status is the resolver for the status field.
func (r *agentResolver) Status(ctx context.Context, obj *model.Agent) (int, error) {
	return int(obj.Status), nil
}

// Configuration is the resolver for the configuration field.
func (r *agentResolver) Configuration(ctx context.Context, obj *model.Agent) (*model1.AgentConfiguration, error) {
	ac := &model1.AgentConfiguration{}
	if err := mapstructure.Decode(obj.Configuration, ac); err != nil {
		return &model1.AgentConfiguration{}, err
	}

	return ac, nil
}

// ConfigurationResource is the resolver for the configurationResource field.
func (r *agentResolver) ConfigurationResource(ctx context.Context, obj *model.Agent) (*model.Configuration, error) {
	return r.bindplane.Store().AgentConfiguration(ctx, obj.ID)
}

// UpgradeAvailable is the resolver for the upgradeAvailable field.
func (r *agentResolver) UpgradeAvailable(ctx context.Context, obj *model.Agent) (*string, error) {
	if !obj.SupportsUpgrade() {
		return nil, nil
	}

	latestVersion, err := r.bindplane.Versions().LatestVersion(ctx)
	if err != nil {
		return nil, nil
	}

	if latestVersion.SemanticVersion().IsNewer(semver.Parse(obj.Version)) {
		latestVersionString := latestVersion.Version()
		return &latestVersionString, nil
	}
	return nil, nil
}

// Features is the resolver for the features field.
func (r *agentResolver) Features(ctx context.Context, obj *model.Agent) (int, error) {
	return int(obj.Features()), nil
}

// MatchLabels is the resolver for the matchLabels field.
func (r *agentSelectorResolver) MatchLabels(ctx context.Context, obj *model.AgentSelector) (map[string]interface{}, error) {
	labels := map[string]interface{}{}
	for k := range obj.MatchLabels {
		labels[k] = obj.MatchLabels[k]
	}
	return labels, nil
}

// Status is the resolver for the status field.
func (r *agentUpgradeResolver) Status(ctx context.Context, obj *model.AgentUpgrade) (int, error) {
	return int(obj.Status), nil
}

// Kind is the resolver for the kind field.
func (r *configurationResolver) Kind(ctx context.Context, obj *model.Configuration) (string, error) {
	return string(obj.GetKind()), nil
}

// AgentCount is the resolver for the agentCount field.
func (r *configurationResolver) AgentCount(ctx context.Context, obj *model.Configuration) (*int, error) {
	ids, err := r.bindplane.Store().AgentsIDsMatchingConfiguration(ctx, obj)
	if err != nil {
		return nil, err
	}
	count := len(ids)
	return &count, nil
}

// Graph is the resolver for the graph field.
func (r *configurationResolver) Graph(ctx context.Context, obj *model.Configuration) (*graph.Graph, error) {
	return obj.Graph(ctx, r.bindplane.Store())
}

// Rendered is the resolver for the rendered field.
func (r *configurationResolver) Rendered(ctx context.Context, obj *model.Configuration) (*string, error) {
	rendered, err := obj.Render(ctx, nil, r.bindplane.Config(), r.bindplane.Store())
	if err != nil {
		return nil, err
	}
	return &rendered, nil
}

// Kind is the resolver for the kind field.
func (r *destinationResolver) Kind(ctx context.Context, obj *model.Destination) (string, error) {
	return string(obj.GetKind()), nil
}

// Kind is the resolver for the kind field.
func (r *destinationTypeResolver) Kind(ctx context.Context, obj *model.DestinationType) (string, error) {
	return string(obj.GetKind()), nil
}

// Labels is the resolver for the labels field.
func (r *metadataResolver) Labels(ctx context.Context, obj *model.Metadata) (map[string]interface{}, error) {
	labels := map[string]interface{}{}
	for k := range obj.Labels.Set {
		labels[k] = obj.Labels.Get(k)
	}
	return labels, nil
}

// UpdateProcessors is the resolver for the updateProcessors field.
func (r *mutationResolver) UpdateProcessors(ctx context.Context, input model1.UpdateProcessorsInput) (*bool, error) {
	config, err := r.bindplane.Store().Configuration(ctx, input.Configuration)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("configuration not found")
	}

	processors := make([]model.ResourceConfiguration, len(input.Processors))
	for ix, p := range input.Processors {
		processors[ix] = *p
	}

	switch input.ResourceType {
	case model1.ResourceTypeKindDestination:
		config.Spec.Destinations[input.ResourceIndex].Processors = processors
	case model1.ResourceTypeKindSource:
		config.Spec.Sources[input.ResourceIndex].Processors = processors
	default:
		return nil, fmt.Errorf("invalid resource type, should be source or destination")
	}

	statuses, err := r.bindplane.Store().ApplyResources(ctx, []model.Resource{config})
	if err != nil {
		return nil, err
	}

	if statuses[0].Status == model.StatusError {
		return nil, errors.New(statuses[0].Reason)
	}

	return nil, nil
}

// Type is the resolver for the type field.
func (r *parameterDefinitionResolver) Type(ctx context.Context, obj *model.ParameterDefinition) (model1.ParameterType, error) {
	switch obj.Type {
	case "strings":
		return model1.ParameterTypeStrings, nil

	case "string":
		return model1.ParameterTypeString, nil
	case "enum":
		return model1.ParameterTypeEnum, nil

	case "bool":
		return model1.ParameterTypeBool, nil

	case "int":
		return model1.ParameterTypeInt, nil

	case "map":
		return model1.ParameterTypeMap, nil

	case "yaml":
		return model1.ParameterTypeYaml, nil

	case "enums":
		return model1.ParameterTypeEnums, nil

	case "timezone":
		return model1.ParameterTypeTimezone, nil

	case "metrics":
		return model1.ParameterTypeMetrics, nil

	case "awsCloudwatchNamedField":
		return model1.ParameterTypeAwsCloudwatchNamedField, nil

	default:
		return "", errors.New("unknown parameter type")
	}
}

// Kind is the resolver for the kind field.
func (r *processorResolver) Kind(ctx context.Context, obj *model.Processor) (string, error) {
	return string(obj.GetKind()), nil
}

// Kind is the resolver for the kind field.
func (r *processorTypeResolver) Kind(ctx context.Context, obj *model.ProcessorType) (string, error) {
	return string(obj.GetKind()), nil
}

// OverviewPage is the resolver for the overviewPage field.
func (r *queryResolver) OverviewPage(ctx context.Context, configIDs []string, destinationIDs []string, period string, telemetryType string) (*model1.OverviewPage, error) {
	graph, err := overviewGraph(ctx, r.bindplane, configIDs, destinationIDs, period, telemetryType)
	if err != nil {
		return nil, err
	}

	return &model1.OverviewPage{
		Graph: graph,
	}, nil
}

// Agents is the resolver for the agents field.
func (r *queryResolver) Agents(ctx context.Context, selector *string, query *string) (*model1.Agents, error) {
	ctx, span := tracer.Start(ctx, "graphql/Agents")
	defer span.End()

	options, suggestions, err := r.queryOptionsAndSuggestions(selector, query, r.bindplane.Store().AgentIndex(ctx))
	if err != nil {
		r.bindplane.Logger().Error("error getting query options and suggestion", zap.Error(err))
		return nil, err
	}

	agents, err := r.bindplane.Store().Agents(ctx, options...)
	if err != nil {
		r.bindplane.Logger().Error("error in graphql Agents", zap.Error(err))
		return nil, err
	}

	return &model1.Agents{
		Agents:        agents,
		Query:         query,
		Suggestions:   suggestions,
		LatestVersion: r.bindplane.Versions().LatestVersionString(ctx),
	}, nil
}

// Agent is the resolver for the agent field.
func (r *queryResolver) Agent(ctx context.Context, id string) (*model.Agent, error) {
	return r.bindplane.Store().Agent(ctx, id)
}

// Configurations is the resolver for the configurations field.
func (r *queryResolver) Configurations(ctx context.Context, selector *string, query *string, onlyDeployedConfigurations *bool) (*model1.Configurations, error) {
	options, suggestions, err := r.queryOptionsAndSuggestions(selector, query, r.bindplane.Store().ConfigurationIndex(ctx))
	if err != nil {
		r.bindplane.Logger().Error("error getting query options and suggestion", zap.Error(err))
		return nil, err
	}

	configurations, err := r.bindplane.Store().Configurations(ctx, options...)
	if err != nil {
		return nil, err
	}
	// filter out configurations that are not deployed
	if onlyDeployedConfigurations != nil && *onlyDeployedConfigurations {
		filteredConfigurations := []*model.Configuration{}
		for _, configuration := range configurations {
			ids, err := r.bindplane.Store().AgentsIDsMatchingConfiguration(ctx, configuration)
			if err != nil {
				return nil, err
			}
			if len(ids) > 0 {
				filteredConfigurations = append(filteredConfigurations, configuration)
			}
		}
		configurations = filteredConfigurations
	}
	return &model1.Configurations{
		Configurations: configurations,
		Query:          query,
		Suggestions:    suggestions,
	}, nil
}

// Configuration is the resolver for the configuration field.
func (r *queryResolver) Configuration(ctx context.Context, name string) (*model.Configuration, error) {
	return r.bindplane.Store().Configuration(ctx, name)
}

// Sources is the resolver for the sources field.
func (r *queryResolver) Sources(ctx context.Context) ([]*model.Source, error) {
	return r.bindplane.Store().Sources(ctx)
}

// Source is the resolver for the source field.
func (r *queryResolver) Source(ctx context.Context, name string) (*model.Source, error) {
	return r.bindplane.Store().Source(ctx, name)
}

// SourceTypes is the resolver for the sourceTypes field.
func (r *queryResolver) SourceTypes(ctx context.Context) ([]*model.SourceType, error) {
	return r.bindplane.Store().SourceTypes(ctx)
}

// SourceType is the resolver for the sourceType field.
func (r *queryResolver) SourceType(ctx context.Context, name string) (*model.SourceType, error) {
	return r.bindplane.Store().SourceType(ctx, name)
}

// Processors is the resolver for the processors field.
func (r *queryResolver) Processors(ctx context.Context) ([]*model.Processor, error) {
	return r.bindplane.Store().Processors(ctx)
}

// Processor is the resolver for the processor field.
func (r *queryResolver) Processor(ctx context.Context, name string) (*model.Processor, error) {
	return r.bindplane.Store().Processor(ctx, name)
}

// ProcessorTypes is the resolver for the processorTypes field.
func (r *queryResolver) ProcessorTypes(ctx context.Context) ([]*model.ProcessorType, error) {
	return r.bindplane.Store().ProcessorTypes(ctx)
}

// ProcessorType is the resolver for the processorType field.
func (r *queryResolver) ProcessorType(ctx context.Context, name string) (*model.ProcessorType, error) {
	return r.bindplane.Store().ProcessorType(ctx, name)
}

// Destinations is the resolver for the destinations field.
func (r *queryResolver) Destinations(ctx context.Context) ([]*model.Destination, error) {
	return r.bindplane.Store().Destinations(ctx)
}

// Destination is the resolver for the destination field.
func (r *queryResolver) Destination(ctx context.Context, name string) (*model.Destination, error) {
	return r.bindplane.Store().Destination(ctx, name)
}

// DestinationWithType is the resolver for the destinationWithType field.
func (r *queryResolver) DestinationWithType(ctx context.Context, name string) (*model1.DestinationWithType, error) {
	resp := &model1.DestinationWithType{}

	dest, err := r.bindplane.Store().Destination(ctx, name)
	if err != nil {
		return resp, err
	}

	if dest == nil {
		return resp, nil
	}

	destinationType, err := r.bindplane.Store().DestinationType(ctx, dest.Spec.Type)
	if err != nil {
		return resp, err
	}

	return &model1.DestinationWithType{
		Destination:     dest,
		DestinationType: destinationType,
	}, nil
}

// DestinationsInConfigs is the resolver for the destinationsInConfigs field.
func (r *queryResolver) DestinationsInConfigs(ctx context.Context) ([]*model.Destination, error) {
	// returns only destinations that are in non-raw (managed?) configs and deployed to agents
	configs, err := r.bindplane.Store().Configurations(ctx)
	if err != nil {
		return nil, err
	}

	// create a map from destination name to destination
	destinationsMap := make(map[string]string)
	destinations := make([]*model.Destination, 0, len(destinationsMap))

	// loop through configs, collect all destinations from configs that aren't raw
	for _, config := range configs {
		if config.Spec.Raw != "" {
			continue
		}
		ids, err := r.bindplane.Store().AgentsIDsMatchingConfiguration(ctx, config)
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			for _, destination := range config.Spec.Destinations {
				_, ok := destinationsMap[destination.Name]
				if !ok {
					dest, err := r.bindplane.Store().Destination(ctx, destination.Name)
					if err != nil {
						return destinations, err
					}
					destinationsMap[destination.Name] = "Remember that we've already seen this destination!"
					destinations = append(destinations, dest)
				}
			}
		}
	}
	return destinations, nil
}

// DestinationTypes is the resolver for the destinationTypes field.
func (r *queryResolver) DestinationTypes(ctx context.Context) ([]*model.DestinationType, error) {
	return r.bindplane.Store().DestinationTypes(ctx)
}

// DestinationType is the resolver for the destinationType field.
func (r *queryResolver) DestinationType(ctx context.Context, name string) (*model.DestinationType, error) {
	return r.bindplane.Store().DestinationType(ctx, name)
}

// Snapshot is the resolver for the snapshot field.
func (r *queryResolver) Snapshot(ctx context.Context, agentID string, pipelineType otel.PipelineType) (*model1.Snapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	signals := &model1.Snapshot{}

	// construct a reporting config for this agent
	config, err := r.bindplane.Store().AgentConfiguration(ctx, agentID)
	if err != nil {
		return signals, err
	}
	if config == nil {
		return signals, fmt.Errorf("no configuration available for agent %s", agentID)
	}

	reportRequest := func(id string) report.Configuration {
		rc := report.Configuration{
			Snapshot: report.Snapshot{
				Processor:    string(otel.SnapshotProcessorName),
				PipelineType: pipelineType,
				Endpoint: report.Endpoint{
					URL: fmt.Sprintf("%s/v1/otlphttp/v1/%s", r.bindplane.Config().BindPlaneURL(), pipelineType),
					Header: http.Header{
						server.HeaderSessionID: []string{id},
					},
				},
			},
		}
		r.bindplane.Logger().Info("Requesting report", zap.Any("config", rc))
		return rc
	}

	// all three cases follow the same pattern:
	//
	// 1) receive a channel to await results from relayer
	//
	// 2) send a message to the agent with the specified session id
	//
	// 3) wait for results or timeout

	switch pipelineType {
	case otel.Logs:
		id, result, cancel := r.bindplane.Relayers().Logs.AwaitResult()
		defer cancel()

		if err := r.bindplane.Manager().RequestReport(ctx, agentID, reportRequest(id)); err != nil {
			return signals, err
		}

		select {
		case <-ctx.Done():
		case signals.Logs = <-result:
		}
	case otel.Metrics:
		id, result, cancel := r.bindplane.Relayers().Metrics.AwaitResult()
		defer cancel()

		if err := r.bindplane.Manager().RequestReport(ctx, agentID, reportRequest(id)); err != nil {
			return signals, err
		}

		select {
		case <-ctx.Done():
		case signals.Metrics = <-result:
		}
	case otel.Traces:
		id, result, cancel := r.bindplane.Relayers().Traces.AwaitResult()
		defer cancel()

		if err := r.bindplane.Manager().RequestReport(ctx, agentID, reportRequest(id)); err != nil {
			return signals, err
		}

		select {
		case <-ctx.Done():
		case signals.Traces = <-result:
		}
	}

	return signals, nil
}

// AgentMetrics is the resolver for the agentMetrics field.
func (r *queryResolver) AgentMetrics(ctx context.Context, period string, ids []string) (*model1.GraphMetrics, error) {
	return agentMetrics(ctx, r.bindplane, period, ids)
}

// ConfigurationMetrics is the resolver for the configurationMetrics field.
func (r *queryResolver) ConfigurationMetrics(ctx context.Context, period string, name *string) (*model1.GraphMetrics, error) {
	return configurationMetrics(ctx, r.bindplane, period, name)
}

// OverviewMetrics is the resolver for the overviewMetrics field.
func (r *queryResolver) OverviewMetrics(ctx context.Context, period string, configIDs []string, destinationIDs []string) (*model1.GraphMetrics, error) {
	return overviewMetrics(ctx, r.bindplane, period, configIDs, destinationIDs)
}

// Operator is the resolver for the operator field.
func (r *relevantIfConditionResolver) Operator(ctx context.Context, obj *model.RelevantIfCondition) (model1.RelevantIfOperatorType, error) {
	return model1.RelevantIfOperatorType(obj.Operator), nil
}

// Kind is the resolver for the kind field.
func (r *sourceResolver) Kind(ctx context.Context, obj *model.Source) (string, error) {
	return string(obj.GetKind()), nil
}

// Kind is the resolver for the kind field.
func (r *sourceTypeResolver) Kind(ctx context.Context, obj *model.SourceType) (string, error) {
	return string(obj.GetKind()), nil
}

// AgentChanges is the resolver for the agentChanges field.
func (r *subscriptionResolver) AgentChanges(ctx context.Context, selector *string, query *string) (<-chan []*model1.AgentChange, error) {
	parsedSelector, parsedQuery, err := r.parseSelectorAndQuery(selector, query)
	if err != nil {
		return nil, err
	}

	// we can ignore the unsubscribe function because this will automatically unsubscribe when the context is done. we
	// could subscribe directly to store.AgentChanges, but the resolver is setup to relay events and the filter and
	// dispatch will happen in a separate goroutine.
	channel, _ := eventbus.SubscribeWithFilterUntilDone(ctx, r.updates, func(updates *store.Updates) (result []*model1.AgentChange, accept bool) {
		// if the observer is using a selector or query, we want to change Update to Remove if it no longer matches the
		// selector or query
		events := applySelectorToChanges(parsedSelector, updates.Agents)
		events = applyQueryToChanges(parsedQuery, r.bindplane.Store().AgentIndex(ctx), events)

		return model1.ToAgentChangeArray(events), !events.Empty()
	})

	return channel, nil
}

// ConfigurationChanges is the resolver for the configurationChanges field.
func (r *subscriptionResolver) ConfigurationChanges(ctx context.Context, selector *string, query *string) (<-chan []*model1.ConfigurationChange, error) {
	parsedSelector, parsedQuery, err := r.parseSelectorAndQuery(selector, query)
	if err != nil {
		return nil, err
	}

	// we can ignore the unsubscribe function because this will automatically unsubscribe when the context is done.
	channel, _ := eventbus.SubscribeWithFilterUntilDone(ctx, r.updates, func(updates *store.Updates) (result []*model1.ConfigurationChange, accept bool) {
		// if the observer is using a selector or query, we want to change Update to Remove if it no longer matches the
		// selector or query

		configUpdates := updates.Configurations
		if r.hasAgentConfigurationChanges(updates) {
			configUpdates = configUpdates.Clone()
			// add all configurations here as updates since we don't know what agent counts could be affected
			if configs, err := r.bindplane.Store().Configurations(ctx); err == nil {
				for _, config := range configs {
					// don't add configuration pseudo-updates that already have updates associated with them
					if _, ok := configUpdates[config.UniqueKey()]; !ok {
						configUpdates.Include(config, store.EventTypeUpdate)
					}
				}
			} else {
				r.bindplane.Logger().Error("unable to get configurations to include in agent changes", zap.Error(err))
			}
		}

		events := applySelectorToEvents(parsedSelector, configUpdates)
		events = applyQueryToEvents(parsedQuery, r.bindplane.Store().ConfigurationIndex(ctx), events)

		return model1.ToConfigurationChanges(events), len(events) > 0
	})

	return channel, nil
}

// AgentMetricsSubscription is the resolver for the agentMetricsSubscription field.
func (r *subscriptionResolver) AgentMetrics(ctx context.Context, period string, ids []string) (<-chan *model1.GraphMetrics, error) {
	channel := make(chan *model1.GraphMetrics)

	updateTicker := time.NewTicker(agentMetricsUpdateInterval)

	sendMetrics := func() {
		if metrics, err := agentMetrics(ctx, r.bindplane, period, ids); err != nil {
			r.bindplane.Logger().Error("failed to get agentMetrics", zap.Error(err))
		} else {
			channel <- metrics
		}
	}

	go metricSubscriber(ctx, sendMetrics, updateTicker)

	return channel, nil
}

// ConfigurationMetricsSubscription is the resolver for the configurationMetricsSubscription field.
func (r *subscriptionResolver) ConfigurationMetrics(ctx context.Context, period string, name *string, agent *string) (<-chan *model1.GraphMetrics, error) {
	channel := make(chan *model1.GraphMetrics)

	updateTicker := time.NewTicker(configurationMetricsUpdateInterval)

	sendMetrics := func() {
		if agent != nil && *agent != "" {
			ids := []string{*agent}
			if metrics, err := agentMetrics(ctx, r.bindplane, period, ids); err != nil {
				r.bindplane.Logger().Error("failed to get agentMetrics", zap.Error(err))
			} else {
				channel <- metrics
			}
		} else {
			if metrics, err := configurationMetrics(ctx, r.bindplane, period, name); err != nil {
				r.bindplane.Logger().Error("failed to get configurationMetrics", zap.Error(err))
			} else {
				channel <- metrics
			}
		}
	}

	go metricSubscriber(ctx, sendMetrics, updateTicker)

	return channel, nil
}

// OverviewMetricsSubscription is the resolver for the overviewMetricsSubscription field.
func (r *subscriptionResolver) OverviewMetrics(ctx context.Context, period string, configIDs []string, destinationIDs []string) (<-chan *model1.GraphMetrics, error) {
	channel := make(chan *model1.GraphMetrics)

	updateTicker := time.NewTicker(overviewMetricsUpdateInterval)

	sendMetrics := func() {
		if metrics, err := overviewMetrics(ctx, r.bindplane, period, configIDs, destinationIDs); err != nil {
			r.bindplane.Logger().Error("failed to get overviewMetrics", zap.Error(err))
		} else {
			channel <- metrics
		}
	}

	go metricSubscriber(ctx, sendMetrics, updateTicker)

	return channel, nil
}

// Agent returns generated.AgentResolver implementation.
func (r *Resolver) Agent() generated.AgentResolver { return &agentResolver{r} }

// AgentSelector returns generated.AgentSelectorResolver implementation.
func (r *Resolver) AgentSelector() generated.AgentSelectorResolver { return &agentSelectorResolver{r} }

// AgentUpgrade returns generated.AgentUpgradeResolver implementation.
func (r *Resolver) AgentUpgrade() generated.AgentUpgradeResolver { return &agentUpgradeResolver{r} }

// Configuration returns generated.ConfigurationResolver implementation.
func (r *Resolver) Configuration() generated.ConfigurationResolver { return &configurationResolver{r} }

// Destination returns generated.DestinationResolver implementation.
func (r *Resolver) Destination() generated.DestinationResolver { return &destinationResolver{r} }

// DestinationType returns generated.DestinationTypeResolver implementation.
func (r *Resolver) DestinationType() generated.DestinationTypeResolver {
	return &destinationTypeResolver{r}
}

// Metadata returns generated.MetadataResolver implementation.
func (r *Resolver) Metadata() generated.MetadataResolver { return &metadataResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// ParameterDefinition returns generated.ParameterDefinitionResolver implementation.
func (r *Resolver) ParameterDefinition() generated.ParameterDefinitionResolver {
	return &parameterDefinitionResolver{r}
}

// Processor returns generated.ProcessorResolver implementation.
func (r *Resolver) Processor() generated.ProcessorResolver { return &processorResolver{r} }

// ProcessorType returns generated.ProcessorTypeResolver implementation.
func (r *Resolver) ProcessorType() generated.ProcessorTypeResolver { return &processorTypeResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// RelevantIfCondition returns generated.RelevantIfConditionResolver implementation.
func (r *Resolver) RelevantIfCondition() generated.RelevantIfConditionResolver {
	return &relevantIfConditionResolver{r}
}

// Source returns generated.SourceResolver implementation.
func (r *Resolver) Source() generated.SourceResolver { return &sourceResolver{r} }

// SourceType returns generated.SourceTypeResolver implementation.
func (r *Resolver) SourceType() generated.SourceTypeResolver { return &sourceTypeResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type agentResolver struct{ *Resolver }
type agentSelectorResolver struct{ *Resolver }
type agentUpgradeResolver struct{ *Resolver }
type configurationResolver struct{ *Resolver }
type destinationResolver struct{ *Resolver }
type destinationTypeResolver struct{ *Resolver }
type metadataResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type parameterDefinitionResolver struct{ *Resolver }
type processorResolver struct{ *Resolver }
type processorTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type relevantIfConditionResolver struct{ *Resolver }
type sourceResolver struct{ *Resolver }
type sourceTypeResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
