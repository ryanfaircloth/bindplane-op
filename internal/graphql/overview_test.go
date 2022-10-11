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
	"testing"

	"github.com/observiq/bindplane-op/internal/store"
	"github.com/observiq/bindplane-op/model"
	"github.com/observiq/bindplane-op/model/graph"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
	destinations   []*model.Destination
	configurations []*model.Configuration
	agentIds       map[string][]string
	store.Store
}

func (m *mockStore) Configurations(ctx context.Context, options ...store.QueryOption) ([]*model.Configuration, error) {
	return m.configurations, nil
}

func (m *mockStore) Destinations(ctx context.Context) ([]*model.Destination, error) {
	return m.destinations, nil
}

func (m *mockStore) AgentsIDsMatchingConfiguration(ctx context.Context, c *model.Configuration) ([]string, error) {
	return m.agentIds[c.Name()], nil
}

func newMockStore() *mockStore {
	return &mockStore{
		destinations:   make([]*model.Destination, 0),
		configurations: make([]*model.Configuration, 0),
		agentIds:       map[string][]string{},
	}
}

func Test_overviewGraph(t *testing.T) {
	tests := []struct {
		description    string
		destinations   []*model.Destination
		configurations []*model.Configuration
		agentIds       map[string][]string
		rawConfig      bool
		want           *graph.Graph
		wantErr        bool
	}{
		{
			"No configs, no destinations",
			make([]*model.Destination, 0),
			make([]*model.Configuration, 0),
			make(map[string][]string),
			false,
			&graph.Graph{
				Sources: make([]*graph.Node, 0),
				Targets: make([]*graph.Node, 0),
				Edges:   make([]*graph.Edge, 0),
			},
			false,
		},
		{
			"one config, one destination, one edge",
			[]*model.Destination{
				testDestination("d-1"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 1),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
				},
			},
			false,
		},
		{
			"1 config to 2 destinations",
			[]*model.Destination{
				testDestination("d-1"),
				testDestination("d-2"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1", "d-2"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 1),
					testNode("d-2", "destinationNode", 1),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
					graph.NewEdge("configuration/c-1", "destination/d-2"),
				},
			},
			false,
		},
		{
			"2 configs to 1 destination,",
			[]*model.Destination{
				testDestination("d-1"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1"}),
				testConfiguration("c-2", false, []string{"d-1"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
				"c-2": {"agent2"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
					testNode("c-2", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 2),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
					graph.NewEdge("configuration/c-2", "destination/d-1"),
				},
			},
			false,
		},
		{
			"2 configs to 2 to destinations, no overlap",
			[]*model.Destination{
				testDestination("d-1"),
				testDestination("d-2"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1"}),
				testConfiguration("c-2", false, []string{"d-2"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
				"c-2": {"agent2"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
					testNode("c-2", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 1),
					testNode("d-2", "destinationNode", 1),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
					graph.NewEdge("configuration/c-2", "destination/d-2"),
				},
			},
			false,
		},
		{
			"2 configs to 2 to destinations, one config uses both",
			[]*model.Destination{
				testDestination("d-1"),
				testDestination("d-2"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1", "d-2"}),
				testConfiguration("c-2", false, []string{"d-2"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
				"c-2": {"agent2"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
					testNode("c-2", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 1),
					testNode("d-2", "destinationNode", 2),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
					graph.NewEdge("configuration/c-1", "destination/d-2"),
					graph.NewEdge("configuration/c-2", "destination/d-2"),
				},
			},
			false,
		},
		{
			"does not show configs with no destination",
			[]*model.Destination{
				testDestination("d-1"),
			},
			[]*model.Configuration{
				testConfiguration("c-1", false, []string{"d-1"}),
				testConfiguration("c-2", false, []string{}),
			},
			map[string][]string{
				"c-1": {"agent1"},
				"c-2": {"agent2"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("c-1", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("d-1", "destinationNode", 1),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/c-1", "destination/d-1"),
				},
			},
			false,
		},
		{
			"works as expected with destination and configuration with same name",
			[]*model.Destination{
				testDestination("GCP"),
			},
			[]*model.Configuration{
				testConfiguration("GCP", false, []string{"GCP"}),
				testConfiguration("c-1", false, []string{"GCP"}),
			},
			map[string][]string{
				"c-1": {"agent1"},
				"GCP": {"agent2"},
			},
			false,
			&graph.Graph{
				Sources: []*graph.Node{
					testNode("GCP", "configurationNode", 1),
					testNode("c-1", "configurationNode", 1),
				},
				Targets: []*graph.Node{
					testNode("GCP", "destinationNode", 2),
				},
				Edges: []*graph.Edge{
					graph.NewEdge("configuration/GCP", "destination/GCP"),
					graph.NewEdge("configuration/c-1", "destination/GCP"),
				},
			},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mock := newMockStore()
			mock.destinations = test.destinations
			mock.configurations = test.configurations
			mock.agentIds = test.agentIds

			got, err := overviewGraph(context.Background(), mock)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.want, got)
		})
	}
}

func testDestination(name string) *model.Destination {
	return &model.Destination{
		ResourceMeta: model.ResourceMeta{
			APIVersion: "v1",
			Kind:       "destinationNode",
			Metadata: model.Metadata{
				Name: name,
			},
		},
		Spec: model.ParameterizedSpec{},
	}
}

func testConfiguration(name string, raw bool, destinationNames []string) *model.Configuration {
	meta := model.ResourceMeta{
		APIVersion: "v1",
		Kind:       "configurationNode",
		Metadata: model.Metadata{
			Name: name,
		},
	}

	var spec model.ConfigurationSpec
	if raw {
		spec = model.ConfigurationSpec{
			Raw: "Set",
		}
	} else {
		spec = model.ConfigurationSpec{
			Destinations: destinationResourceConfigurations(destinationNames),
		}
	}

	return &model.Configuration{
		ResourceMeta: meta,
		Spec:         spec,
	}
}

func destinationResourceConfigurations(names []string) []model.ResourceConfiguration {
	if len(names) == 0 {
		return nil
	}

	resourceConfigurations := make([]model.ResourceConfiguration, 0)
	for _, name := range names {
		resourceConfigurations = append(
			resourceConfigurations,
			model.ResourceConfiguration{
				Name: name,
			},
		)
	}
	return resourceConfigurations
}

func testNode(name, nodeType string, agentCount int) *graph.Node {
	var id string
	var kind string

	if nodeType == "configurationNode" {
		id = fmt.Sprintf("configuration/%s", name)
		kind = string(model.KindConfiguration)
	} else {
		id = fmt.Sprintf("destination/%s", name)
		kind = string(model.KindDestination)
	}
	return &graph.Node{
		ID:    id,
		Label: name,
		Type:  nodeType,
		Attributes: map[string]any{
			"kind":       kind,
			"resourceId": name,
			"agentCount": agentCount,
		},
	}
}
