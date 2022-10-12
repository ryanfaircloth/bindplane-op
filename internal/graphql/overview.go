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

	"github.com/observiq/bindplane-op/internal/store"
	"github.com/observiq/bindplane-op/model"
	"github.com/observiq/bindplane-op/model/graph"
)

func overviewGraph(ctx context.Context, store store.Store) (*graph.Graph, error) {
	configs, err := store.Configurations(ctx)
	if err != nil {
		return nil, err
	}

	// keep track of destinations and reuse them
	destNodes := make(map[string]*graph.Node)
	// Also keep an ordered slice of the nodes for consistent ordering
	destNodesSlice := make([]*graph.Node, 0)

	edges := make([]*graph.Edge, 0)
	configNodes := make([]*graph.Node, 0)
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

		configNodeID := fmt.Sprintf("configuration/%s", c.Name())
		configAttrs := graph.MakeAttributes(string(model.KindConfiguration), c.Name())

		configUsage := c.Usage(ctx, store)
		configAttrs.AddAttribute("agentCount", len(agentIDs))
		configAttrs.AddAttribute("activeTypeFlags", configUsage.ActiveFlags())

		configNodes = append(configNodes, &graph.Node{
			ID:         configNodeID,
			Label:      c.Name(),
			Type:       "configurationNode",
			Attributes: configAttrs,
		})

		for _, d := range c.Spec.Destinations {
			// ignore inline destinations which are not supported in the UI
			if d.Name == "" {
				continue
			}

			destNodeID := fmt.Sprintf("destination/%s", d.Name)
			destNode, ok := destNodes[destNodeID]
			if !ok {
				destAttrs := graph.MakeAttributes(string(model.KindDestination), d.Name)
				destAttrs.AddAttribute("agentCount", len(agentIDs))
				destAttrs.AddAttribute("activeTypeFlags", configUsage.ActiveFlagsForDestination(d.Name))

				node := &graph.Node{
					ID:         destNodeID,
					Label:      d.Name,
					Type:       "destinationNode",
					Attributes: destAttrs,
				}
				destNodes[destNodeID] = node
				destNodesSlice = append(destNodesSlice, node)
			} else {
				// increment the agent count
				destNode.Attributes["agentCount"] = destNode.Attributes["agentCount"].(int) + len(agentIDs)
			}
			edges = append(edges, graph.NewEdge(configNodeID, destNodeID))
		}
	}

	return &graph.Graph{
		Sources: configNodes,
		Targets: destNodesSlice,
		Edges:   edges,
	}, nil
}
