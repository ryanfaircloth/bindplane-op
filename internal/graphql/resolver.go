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

package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/observiq/bindplane-op/internal/eventbus"
	model1 "github.com/observiq/bindplane-op/internal/graphql/model"
	"github.com/observiq/bindplane-op/internal/otlp/record"
	"github.com/observiq/bindplane-op/internal/server"
	"github.com/observiq/bindplane-op/internal/store"
	"github.com/observiq/bindplane-op/internal/store/search"
	"github.com/observiq/bindplane-op/internal/store/stats"
	"github.com/observiq/bindplane-op/model"
	bpotel "github.com/observiq/bindplane-op/model/otel"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

var tracer = otel.Tracer("graphql")

//go:generate gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver TODO(doc)
type Resolver struct {
	bindplane server.BindPlane
	updates   eventbus.Source[*store.Updates]
}

// NewResolver returns a new Resolver and starts a go routine
// that sends agent updates to observers.
func NewResolver(bindplane server.BindPlane) *Resolver {
	resolver := &Resolver{
		bindplane: bindplane,
		updates:   eventbus.NewSource[*store.Updates](),
	}

	// relay events from the store to the resolver where they will be dispatched to individual graphql subscriptions
	eventbus.Relay(context.Background(), bindplane.Store().Updates(), resolver.updates)

	return resolver
}

func applySelectorToChanges(selector *model.Selector, changes store.Events[*model.Agent]) store.Events[*model.Agent] {
	if selector == nil {
		return changes
	}
	result := store.NewEvents[*model.Agent]()
	for _, change := range changes {
		if change.Type != store.EventTypeRemove && !selector.Matches(change.Item.Labels) {
			result.Include(change.Item, store.EventTypeRemove)
		} else {
			result.Include(change.Item, change.Type)
		}
	}
	return result
}

func applyQueryToChanges(query *search.Query, index search.Index, changes store.Events[*model.Agent]) store.Events[*model.Agent] {
	if query == nil {
		return changes
	}
	result := store.NewEvents[*model.Agent]()
	for _, change := range changes {
		if change.Type != store.EventTypeRemove && !index.Matches(query, change.Item.ID) {
			result.Include(change.Item, store.EventTypeRemove)
		} else {
			result.Include(change.Item, change.Type)
		}
	}
	return result
}

func applySelectorToEvents[T model.Resource](selector *model.Selector, events store.Events[T]) store.Events[T] {
	if selector == nil {
		return events
	}
	result := store.NewEvents[T]()
	for _, event := range events {
		if event.Type != store.EventTypeRemove && !selector.Matches(event.Item.GetLabels()) {
			result.Include(event.Item, store.EventTypeRemove)
		} else {
			result.Include(event.Item, event.Type)
		}
	}
	return result
}

func applyQueryToEvents[T model.Resource](query *search.Query, index search.Index, events store.Events[T]) store.Events[T] {
	if query == nil || index == nil {
		return events
	}
	result := store.NewEvents[T]()
	for _, event := range events {
		if event.Type != store.EventTypeRemove && !index.Matches(query, event.Item.Name()) {
			result.Include(event.Item, store.EventTypeRemove)
		} else {
			result.Include(event.Item, event.Type)
		}
	}
	return result
}

func (r *Resolver) parseSelectorAndQuery(selector *string, query *string) (*model.Selector, *search.Query, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var parsedSelector *model.Selector
	if selector != nil {
		sel, err := model.SelectorFromString(*selector)
		if err != nil {
			return nil, nil, err
		}
		parsedSelector = &sel
	}

	// parse the parsedQuery, if any
	var parsedQuery *search.Query
	if query != nil && *query != "" {
		q := search.ParseQuery(*query)
		q.ReplaceVersionLatest(ctx, r.bindplane.Versions())
		parsedQuery = q
	}

	return parsedSelector, parsedQuery, nil
}

func (r *Resolver) queryOptionsAndSuggestions(selector *string, query *string, index search.Index) ([]store.QueryOption, []*search.Suggestion, error) {
	parsedSelector, parsedQuery, err := r.parseSelectorAndQuery(selector, query)
	if err != nil {
		return nil, nil, err
	}

	options := []store.QueryOption{}
	if parsedSelector != nil {
		options = append(options, store.WithSelector(*parsedSelector))
	}

	var suggestions []*search.Suggestion
	if parsedQuery != nil {
		options = append(options, store.WithQuery(parsedQuery))

		s, err := index.Suggestions(parsedQuery)
		if err != nil {
			return nil, nil, err
		}

		suggestions = s
	}
	return options, suggestions, nil
}

// hasAgentConfigurationChanges determines if there is an agent update
// in updates that would affect the list of configurations
func (r *Resolver) hasAgentConfigurationChanges(updates *store.Updates) bool {
	for _, change := range updates.Agents {
		// Only events type Remove, Label, and Insert could affect
		// the agentCount field.
		if change.Type != store.EventTypeUpdate {
			return true
		}
	}
	return false
}



func configurationNodeIDResolver(metric *record.Metric, position model.MeasurementPosition, pipelineType bpotel.PipelineType, resourceName string) string {
	switch position {
	case model.MeasurementPositionSourceBeforeProcessors:
		return fmt.Sprintf("source/%s", resourceName)
	case model.MeasurementPositionSourceAfterProcessors:
		return fmt.Sprintf("source/%s/processors", resourceName)
	case model.MeasurementPositionDestinationBeforeProcessors:
		return fmt.Sprintf("destination/%s/processors", resourceName)
	case model.MeasurementPositionDestinationAfterProcessors:
		return fmt.Sprintf("destination/%s", resourceName)
	}
	return resourceName
}

func overviewMetrics(ctx context.Context, bindplane server.BindPlane, period string) (*model1.GraphMetrics, error) {
	if period == "" {
		period = "1m"
	}
	d, err := time.ParseDuration(period)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period %s", period)
	}
	metrics, err := bindplane.Store().Measurements().OverviewMetrics(ctx, stats.WithPeriod(d))
	if err != nil {
		return nil, err
	}

	includeMetric := func(metricMap map[string]*model1.GraphMetric, pipelineType string, nodeID string, metric *record.Metric) {
		// separate metric per pipelineType
		key := fmt.Sprintf("%s_%s", pipelineType, nodeID)
		if cur, ok := metricMap[key]; ok {
			// already exists, include in sum
			if value, ok := stats.Value(metric); ok {
				cur.Value += value
			} else {
				bindplane.Logger().Debug("unable to parse value as float", zap.Any("value", metric.Value))
			}
		} else {
			// doesn't exist, create a metric
			m, err := model1.ToGraphMetric(metric)
			if err != nil {
				bindplane.Logger().Debug("unable to convert record.Metric to GraphMetric", zap.Error(err))
				return
			}
			m.NodeID = nodeID
			metricMap[key] = m
		}
	}

	// map of processor (includes type and name) => metric
	destinationMetrics := map[string]*model1.GraphMetric{}
	includeDestination := func(metric *record.Metric, pipelineType, destinationName string) {
		nodeID := fmt.Sprintf("destination/%s", destinationName)
		includeMetric(destinationMetrics, pipelineType, nodeID, metric)
	}

	// map of configuration name => metric
	configurationMetrics := map[string]*model1.GraphMetric{}
	includeConfiguration := func(metric *record.Metric, pipelineType string) {
		configurationName := stats.Configuration(metric)
		nodeID := fmt.Sprintf("configuration/%s", configurationName)
		includeMetric(configurationMetrics, pipelineType, nodeID, metric)
	}

	for _, metric := range metrics {
		position, pipelineType, resourceName := stats.ProcessorParsed(metric)
		if position != string(model.MeasurementPositionDestinationAfterProcessors) {
			continue
		}
		includeDestination(metric, pipelineType, resourceName)
		includeConfiguration(metric, pipelineType)
	}

	var graphMetrics []*model1.GraphMetric

	// add all of the totals for destinations and configurations
	graphMetrics = append(graphMetrics, maps.Values(destinationMetrics)...)
	graphMetrics = append(graphMetrics, maps.Values(configurationMetrics)...)

	return &model1.GraphMetrics{
		Metrics: graphMetrics,
	}, nil
}

func configurationMetrics(ctx context.Context, bindplane server.BindPlane, period string, name *string) (*model1.GraphMetrics, error) {
	if period == "" {
		period = "1m"
	}
	configurationName := ""
	if name != nil {
		configurationName = *name
	}
	d, err := time.ParseDuration(period)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period %s with duration %s", period, d)
	}
	metrics, err := bindplane.Store().Measurements().ConfigurationMetrics(ctx, configurationName, stats.WithPeriod(d))
	if err != nil {
		return nil, err
	}

	return assignMetricsToGraph(metrics, configurationNodeIDResolver, bindplane), nil
}

func agentMetrics(ctx context.Context, bindplane server.BindPlane, period string, ids []string) (*model1.GraphMetrics, error) {
	if period == "" {
		period = "1m"
	}

	d, err := time.ParseDuration(period)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period %s with duration %s", period, d)
	}
	if len(ids) == 0 {
		agents, _ := bindplane.Store().Agents(ctx)
		for _, a := range agents {
			ids = append(ids, a.ID)
		}
	}

	var graphMetrics []*model1.GraphMetric

	returnMetrics := &model1.GraphMetrics{
		Metrics: graphMetrics,
	}

	for _, id := range ids {
		metrics, err := bindplane.Store().Measurements().AgentMetrics(ctx, []string{id}, stats.WithPeriod(d))
		if err != nil {
			return nil, err
		}

		singleIDMetrics := assignMetricsToGraph(metrics, configurationNodeIDResolver, bindplane)

		for _, m := range singleIDMetrics.Metrics {
			// This prevents each AgentID being pointed at the same string
			aID := id
			m.AgentID = &aID
		}
		returnMetrics.Metrics = append(returnMetrics.Metrics, singleIDMetrics.Metrics...)
	}

	return returnMetrics, nil
}

// NodeIDResolver is a function that assigns the appropriate NodeID to a GraphMetric based on its position,
// pipelineType, and resourceName parsed out of the metric processor name. If an empty string is returned, this metric
// will be ignored.
type NodeIDResolver func(metric *record.Metric, position model.MeasurementPosition, pipelineType bpotel.PipelineType, resourceName string) string

func assignMetricsToGraph(metrics []*record.Metric, resolver NodeIDResolver, bindplane server.BindPlane) *model1.GraphMetrics {
	var graphMetrics []*model1.GraphMetric

	for _, m := range metrics {
		graphMetric, err := model1.ToGraphMetric(m)
		if err != nil {
			bindplane.Logger().Debug("unable to convert record.Metric to GraphMetric", zap.Error(err))
			continue
		}

		// figure out what node this is. this must be sync'd with model.Configuration::Graph()
		position, pipelineType, resourceName := stats.ProcessorParsed(m)
		graphMetric.NodeID = resolver(m, model.MeasurementPosition(position), bpotel.PipelineType(pipelineType), resourceName)

		graphMetric.PipelineType = pipelineType
		graphMetrics = append(graphMetrics, graphMetric)
	}

	return &model1.GraphMetrics{
		Metrics: graphMetrics,
	}
}

const configurationMetricsUpdateInterval = 10 * time.Second
const overviewMetricsUpdateInterval = 10 * time.Second
const agentMetricsUpdateInterval = 10 * time.Second

func metricSubscriber(ctx context.Context, sendMetrics func(), updateTicker *time.Ticker) {
	defer updateTicker.Stop()

	sendMetrics()
	for {
		select {
		case <-updateTicker.C:
			// tick, send data
			sendMetrics()

		case <-ctx.Done():
			return
		}
	}
}
