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

	"github.com/observiq/bindplane-op/internal/store/mocks"
	"github.com/observiq/bindplane-op/model"
	"github.com/observiq/bindplane-op/model/graph"
	"github.com/observiq/bindplane-op/model/otel"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type overviewGraphTest struct {
	destinations   []*model.Destination
	configurations []*model.Configuration
	agentIds       map[string][]string
	rawConfig      bool
	want           *graph.Graph
	wantErr        bool	
}

func Test_overviewGraph(t *testing.T) {
	tests := []struct {
		description    string
		testData overviewGraphTest 
	}{
		{
			"No configs, no destinations",
			overviewGraphTest {
				make([]*model.Destination, 0),
				make([]*model.Configuration, 0),
				make(map[string][]string),
				false,
				&graph.Graph{
					Sources: 		make([]*graph.Node, 0),
					Targets: 		make([]*graph.Node, 0),
					Intermediates: 	make([]*graph.Node, 0),
					Edges:   		make([]*graph.Edge, 0),
					Attributes:    	testAttributes(),
				},
				false,
			},
		},
		{
			"one config, one destination, one edge",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 1),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"1 config to 2 destinations",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
						graph.NewEdge("configuration/c-1", "destination/d-2"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"2 configs to 1 destination,",
			overviewGraphTest{
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
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 2),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
						graph.NewEdge("configuration/c-2", "destination/d-1"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"2 configs to 2 to destinations, no overlap",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),
					
					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 1),
						testNode("d-2", "destinationNode", 1),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
						graph.NewEdge("configuration/c-2", "destination/d-2"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"2 configs to 2 to destinations, one config uses both",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),

					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 1),
						testNode("d-2", "destinationNode", 2),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
						graph.NewEdge("configuration/c-1", "destination/d-2"),
						graph.NewEdge("configuration/c-2", "destination/d-2"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"3 configs to 3 to destinations, one config uses two, second config shares one dest with the first, third config uses it's own destination.",
			overviewGraphTest {
				[]*model.Destination{
					testDestination("d-1"),
					testDestination("d-2"),
					testDestination("d-3"),
				},
				[]*model.Configuration{
					testConfiguration("c-1", false, []string{"d-1", "d-2"}),
					testConfiguration("c-2", false, []string{"d-2"}),
					testConfiguration("c-3", false, []string{"d-3"}),
				},
				map[string][]string{
					"c-1": {"agent1"},
					"c-2": {"agent2"},
					"c-3": {"agent3"},
				},
				false,
				&graph.Graph{
					Sources: []*graph.Node{
						testNode("c-1", "configurationNode", 1),
						testNode("c-2", "configurationNode", 1),
						testNode("c-3", "configurationNode", 1),
					},
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 1),
						testNode("d-2", "destinationNode", 2),
						testNode("d-3", "destinationNode", 1),						
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
						graph.NewEdge("configuration/c-1", "destination/d-2"),
						graph.NewEdge("configuration/c-2", "destination/d-2"),
						graph.NewEdge("configuration/c-3", "destination/d-3"),						
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"does not show configs with no destination",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("d-1", "destinationNode", 1),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/c-1", "destination/d-1"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"works as expected with 2 configurations and 1 destination, one configuration with same name as destination",
			overviewGraphTest {
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
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("GCP", "destinationNode", 2),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/GCP", "destination/GCP"),
						graph.NewEdge("configuration/c-1", "destination/GCP"),
					},
					Attributes: testAttributes(),
				},
				false,
			},
		},
		{
			"works as expected with 2 configurations and 2 destinations, one configuration with same name as destination",
			overviewGraphTest {
				[]*model.Destination{
					testDestination("GCP"),
					testDestination("Kafka"),					
				},
				[]*model.Configuration{
					testConfiguration("GCP", false, []string{"GCP"}),
					testConfiguration("c-1", false, []string{"Kafka"}),
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
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("GCP", "destinationNode", 1),
						testNode("Kafka", "destinationNode", 1),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/GCP", "destination/GCP"),
						graph.NewEdge("configuration/c-1", "destination/Kafka"),
					},
					Attributes: testAttributes(),					
				},
				false,
			},
		},
		{
			"works as expected with 2 configurations and 2 destinations, both configurations with same name as destinations",
			overviewGraphTest {
				[]*model.Destination{
					testDestination("GCP"),
					testDestination("Kafka"),					

				},
				[]*model.Configuration{
					testConfiguration("GCP", false, []string{"GCP"}),
					testConfiguration("Kafka", false, []string{"Kafka"}),
				},
				map[string][]string{
					"Kafka": {"agent1"},
					"GCP": {"agent2"},
				},
				false,
				&graph.Graph{
					Sources: []*graph.Node{
						testNode("GCP", "configurationNode", 1),
						testNode("Kafka", "configurationNode", 1),
					},
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("GCP", "destinationNode", 1),
						testNode("Kafka", "destinationNode", 1),
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/GCP", "destination/GCP"),
						graph.NewEdge("configuration/Kafka", "destination/Kafka"),
					},
					Attributes: testAttributes(),					
				},
				false,
			},
		},
		{
			"works as expected with 3 configurations and 3 destinations, all configurations with same name as destinations",
			overviewGraphTest {
				[]*model.Destination{
					testDestination("GCP"),
					testDestination("Kafka"),					
					testDestination("Logging"),										

				},
				[]*model.Configuration{
					testConfiguration("GCP", false, []string{"GCP"}),
					testConfiguration("Kafka", false, []string{"Kafka"}),
					testConfiguration("Logging", false, []string{"Logging"}),					
				},
				map[string][]string{
					"Kafka": {"agent1"},
					"GCP": {"agent2"},
					"Logging": {"agent3"},						
				},
				false,
				&graph.Graph{
					Sources: []*graph.Node{
						testNode("GCP", "configurationNode", 1),
						testNode("Kafka", "configurationNode", 1),
						testNode("Logging", "configurationNode", 1),						
					},
					Intermediates: make([]*graph.Node, 0),
					Targets: []*graph.Node{
						testNode("GCP", "destinationNode", 1),
						testNode("Kafka", "destinationNode", 1),
						testNode("Logging", "destinationNode", 1),						
					},
					Edges: []*graph.Edge{
						graph.NewEdge("configuration/GCP", "destination/GCP"),
						graph.NewEdge("configuration/Kafka", "destination/Kafka"),
						graph.NewEdge("configuration/Logging", "destination/Logging"),						
					},
					Attributes: testAttributes(),					
				},
				false,
			},
		},
		{
			"Raw configuration",
			overviewGraphTest {
				[]*model.Destination {
					testDestination("d1"),
				},
				[]*model.Configuration{
					testConfiguration("c1", true, []string{"d1"}),
				},
				map[string][]string{
					"c1": {"agent"},						
				},				
				true,
				&graph.Graph{
					Sources: make([]*graph.Node, 0),
					Intermediates: make([]*graph.Node, 0),
					Targets: make([]*graph.Node, 0),
					Edges:   make([]*graph.Edge, 0),
					Attributes: testAttributes(),					
				},
				false,
			},
		},
		{
			"Undeployed configuration",
			overviewGraphTest {
				[]*model.Destination {
					testDestination("d1"),
				},
				[]*model.Configuration{
					testConfiguration("c1", false, []string{"d1"}),
				},
				map[string][]string{
									
				},				
				false,
				&graph.Graph{
					Sources: make([]*graph.Node, 0),
					Intermediates: make([]*graph.Node, 0),
					Targets: make([]*graph.Node, 0),
					Edges:   make([]*graph.Edge, 0),
					Attributes: testAttributes(),					
				},
				false,
			},
		},
		{
			"Inline destination",
			overviewGraphTest {
				[]*model.Destination {
					testDestination(""),
				},
				[]*model.Configuration{
					testConfiguration("c1", false, []string{""}),
				},
				map[string][]string{
					"c1": {"agent"},						
				},				
				false,
				&graph.Graph{
					Sources: []*graph.Node{
						testNode("c1", "configurationNode", 1),					
					},
					Intermediates: make([]*graph.Node, 0),
					Targets: make([]*graph.Node, 0),
					Edges:   make([]*graph.Edge, 0),
					Attributes: testAttributes(),					
				},
				false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			m := mocks.NewStore(t)
			m.On("Configurations", mock.Anything, mock.Anything).Return(test.testData.configurations, nil)
			m.On("AgentsIDsMatchingConfiguration", mock.Anything, mock.Anything).Maybe().Return(
				func(ctx context.Context, config *model.Configuration)  []string {
					return test.testData.agentIds[config.Name()]
				},
				func(ctx context.Context, config *model.Configuration)  error {
					return nil
				},
			)
			m.On("Destinations", mock.Anything).Maybe().Return(
				func(ctx context.Context) []*model.Destination {
					return test.testData.destinations
				},
				func(ctx context.Context) error {
					return nil
				},				
			)
			m.On("Destination", mock.Anything, mock.Anything).Maybe().Return(
				func(ctx context.Context, name string) *model.Destination {
					return test.testData.destinations[0]
				},
				func(ctx context.Context, name string) error {
					return nil
				},
				
			)
			m.On("DestinationType", mock.Anything, mock.Anything).Maybe().Return(
				func(ctx context.Context, name string) *model.DestinationType {
					return model.NewDestinationType("macOS", nil)
				},
				func(ctx context.Context, name string) error {
					return nil
				},
			)
			got, err := overviewGraph(context.Background(), m)
			if test.testData.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.testData.want, got)
		})
	}
	t.Run("error when store.Configurations fails", func(t *testing.T) {
		mockStore := mocks.NewStore(t)
		mockStore.On("Configurations", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))
		got, err := overviewGraph(context.Background(), mockStore)
		require.Nil(t, got)
		require.Error(t, err)
	})
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
		Spec: model.ParameterizedSpec{
			Type: "macOS",
		},
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

	var activeFlags otel.PipelineTypeFlags
	return &graph.Node{
		ID:    id,
		Label: name,
		Type:  nodeType,
		Attributes: map[string]any{
			"kind":            kind,
			"resourceId":      name,
			"agentCount":      agentCount,
			"activeTypeFlags": activeFlags,
		},
	}
}

func testAttributes() map[string]any {
	var flags otel.PipelineTypeFlags
	return map[string]any{
		"activeTypeFlags": flags,
	}
}
