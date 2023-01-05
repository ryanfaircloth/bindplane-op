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

import (
	"context"
	"fmt"
	"sort"
	"strings"

	model1 "github.com/observiq/bindplane-op/internal/graphql/model"
	"github.com/observiq/bindplane-op/internal/server"
	"github.com/observiq/bindplane-op/model"
	"github.com/observiq/bindplane-op/model/graph"
	"github.com/observiq/bindplane-op/model/otel"
)

func overviewGraph(ctx context.Context, b server.BindPlane, configsIDs []string, destinationIDs []string, period string, telemetryType string) (*graph.Graph, error) {
	if configsIDs == nil {
		configsIDs = []string{}
	}
	if destinationIDs == nil {
		destinationIDs = []string{}
	}

	g := graph.NewGraph()
	var activeFlags otel.PipelineTypeFlags

	store := b.Store()

	configs, err := store.Configurations(ctx)
	if err != nil {
		return nil, err
	}

	everythingOrSelected := func(resourceKey string, kind model.Kind) (string, bool) {
		resourcesSelected := []string{}
		switch kind {
		case model.KindConfiguration:
			resourcesSelected = configsIDs
		case model.KindDestination:
			resourcesSelected = destinationIDs
		}
		for _, resourceID := range resourcesSelected {
			if strings.HasSuffix(resourceID, resourceKey) {
				// resources is selected
				return fmt.Sprintf("%s/%s", strings.ToLower(string(kind)), resourceKey), false
			}
		}
		// resource is in the everything node
		return fmt.Sprintf("everything/%s", strings.ToLower(string(kind))), true
	}

	// keep track of destinations and reuse them
	destNodes := make(map[string]*graph.Node)
	// Also keep an ordered slice of the nodes for consistent ordering
	destNodesSlice := make([]*graph.Node, 0)

	edges := make([]*graph.Edge, 0)

	// track which edges we've drawn so far
	edgesSlice := make(map[string]*graph.Edge, 0)

	// keep track of configs and reuse them
	configNodes := make(map[string]*graph.Node)

	// Also keep an ordered slice of the nodes for consistent ordering
	configNodesSlice := make([]*graph.Node, 0)

	for _, c := range configs {
		// don't include raw configurations because their destinations aren't configured
		if c.Spec.Raw != "" {
			continue
		}

		// don't include configurations with no destinations
		if len(c.Spec.Destinations) == 0 {
			continue
		}

		// don't include configurations that aren't deployed to agents
		agentIDs, err := store.AgentsIDsMatchingConfiguration(ctx, c)
		if err != nil || len(agentIDs) == 0 {
			continue
		}

		configNodeID, isEverything := everythingOrSelected(c.Name(), model.KindConfiguration)
		configKey := ""
		if isEverything {
			configKey = configNodeID
		} else {
			configKey = c.Name()
		}

		configUsage := c.Usage(ctx, store)

		if configNode, ok := configNodes[configKey]; !ok {
			configAttrs := graph.MakeAttributes(string(model.KindConfiguration), configKey)
			configAttrs.AddAttribute("agentCount", len(agentIDs))
			configAttrs.AddAttribute("activeTypeFlags", configUsage.ActiveFlags())

			configNode := &graph.Node{
				ID:         configNodeID,
				Label:      configKey,
				Type:       "configurationNode",
				Attributes: configAttrs,
			}

			configNodes[configNodeID] = configNode
			configNodesSlice = append(configNodesSlice, configNode)
		} else {
			configNode.Attributes["activeTypeFlags"] = configNode.Attributes["activeTypeFlags"].(otel.PipelineTypeFlags) | configUsage.ActiveFlags()
			configNode.Attributes["agentCount"] = configNode.Attributes["agentCount"].(int) + len(agentIDs)
		}
		activeFlags = activeFlags | configUsage.ActiveFlags()

		for _, d := range c.Spec.Destinations {
			// ignore inline destinations which are not supported in the UI
			if d.Name == "" {
				continue
			}

			destNodeID, isEverything := everythingOrSelected(d.Name, model.KindDestination)

			destinationKey := ""
			if isEverything {
				destinationKey = destNodeID
			} else {
				destinationKey = d.Name
			}

			if destNode, ok := destNodes[destinationKey]; !ok {
				destAttrs := graph.MakeAttributes(string(model.KindDestination), destinationKey)
				destAttrs.AddAttribute("agentCount", len(agentIDs))
				destAttrs.AddAttribute("activeTypeFlags", configUsage.ActiveFlagsForDestination(d.Name))

				destinationNode := &graph.Node{
					ID:         destNodeID,
					Label:      destinationKey,
					Type:       "destinationNode",
					Attributes: destAttrs,
				}
				destNodes[destinationKey] = destinationNode
				destNodesSlice = append(destNodesSlice, destinationNode)
			} else {
				// increment the agent count
				destNode.Attributes["activeTypeFlags"] = destNode.Attributes["activeTypeFlags"].(otel.PipelineTypeFlags) | configUsage.ActiveFlagsForDestination(d.Name)
				destNode.Attributes["agentCount"] = destNode.Attributes["agentCount"].(int) + len(agentIDs)
			}

			edgeKey := fmt.Sprintf("%s|%s", configNodeID, destNodeID)
			if _, ok := edgesSlice[edgeKey]; !ok {
				newEdge := graph.NewEdge(configNodeID, destNodeID)
				edgesSlice[edgeKey] = newEdge
				edges = append(edges, newEdge)
			}
		}
	}

	// get the metrics for each config
	if period == "" {
		period = "10s"
	}

	metrics, err := overviewMetrics(ctx, b, period, configsIDs, destinationIDs)

	if err != nil {
		return nil, err
	}

	metricMap := make(map[string]*model1.GraphMetric)

	for _, m := range metrics.Metrics {
		if strings.Contains(m.Name, telemetryType) {
			metricMap[m.NodeID] = m
		}
	}

	// sort the configs by metrics
	sort.SliceStable(configNodesSlice, func(i, j int) bool {
		// always put the everything node last
		if configNodesSlice[i].ID == "everything/configuration" {
			return false
		}
		if configNodesSlice[j].ID == "everything/configuration" {
			return true
		}
		// find metrics for each config
		mi, ok := metricMap[configNodesSlice[i].ID]
		if !ok {
			mi = &model1.GraphMetric{}
		}
		mj, ok := metricMap[configNodesSlice[j].ID]
		if !ok {
			mj = &model1.GraphMetric{}
		}

		return mi.Value > mj.Value
	})

	// sort the destinations by metrics,
	sort.SliceStable(destNodesSlice, func(i, j int) bool {
		// always put the everything node last
		if destNodesSlice[i].ID == "everything/destination" {
			return false
		}
		if destNodesSlice[j].ID == "everything/destination" {
			return true
		}
		// find metrics for each config
		mi, ok := metricMap[destNodesSlice[i].ID]
		if !ok {
			mi = &model1.GraphMetric{}
		}
		mj, ok := metricMap[destNodesSlice[j].ID]
		if !ok {
			mj = &model1.GraphMetric{}
		}
		return mi.Value > mj.Value
	})

	g.Sources = configNodesSlice
	g.Targets = destNodesSlice
	g.Edges = edges
	g.Attributes["activeTypeFlags"] = activeFlags

	return g, nil
}
