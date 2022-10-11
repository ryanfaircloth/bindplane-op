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

package graph

import (
	"fmt"
)

// Graph represents a horizontal directed flow graph that groups nodes by their position. Sources will appear on the left and
// should only be the source of edges, intermediates will appear it the middle and can be the target and source of edges, and
// targets will appear on the right and should only be the target of edges.
type Graph struct {
	Sources       []*Node
	Intermediates []*Node
	Targets       []*Node
	Edges         []*Edge
}

// Node represents a single node in a graph and will be connected to other nodes by Edges.
type Node struct {
	// ID is used to connect this Node to other Nodes via Edges
	ID string

	// Type is used to identify similar nodes and may impact rendering, e.g. "Processor"
	Type string

	// Label will be displayed on the Node in a rendering of the graph
	Label string

	// Attributes allow arbitrary information to be included with the Node
	Attributes map[string]any
}

// Edge represents an single edge in the graph and connects a source node to a target node.
type Edge struct {
	ID     string
	Source string
	Target string
}

// NewEdge constructs a new Edge with the specified source and target and sets the id to "source|target"
func NewEdge(source, target string) *Edge {
	return &Edge{
		ID:     fmt.Sprintf("%s|%s", source, target),
		Source: source,
		Target: target,
	}
}

// NewGraph constructs a new graph with empty nodes and edges
func NewGraph() *Graph {
	return &Graph{
		Sources:       []*Node{},
		Intermediates: []*Node{},
		Targets:       []*Node{},
		Edges:         []*Edge{},
	}
}

// AddSource adds another source to the graph
func (f *Graph) AddSource(source *Node) {
	f.Sources = append(f.Sources, source)
}

// AddIntermediate adds an intermediate node (not a source or target) to the graph
func (f *Graph) AddIntermediate(intermediate *Node) {
	f.Intermediates = append(f.Intermediates, intermediate)
}

// AddTarget adds another target to the graph
func (f *Graph) AddTarget(target *Node) {
	f.Targets = append(f.Targets, target)
}

// AddEdge adds an edge to the graph
func (f *Graph) AddEdge(edge *Edge) {
	if !f.HasEdge(edge) {
		f.Edges = append(f.Edges, edge)
	}
}

// HasEdge adds an edge to the graph
func (f *Graph) HasEdge(edge *Edge) bool {
	for _, e := range f.Edges {
		if e.ID == edge.ID {
			return true
		}
	}
	return false
}

// Connect will construct an Edge that joins two Nodes
func (f *Graph) Connect(source, target *Node) {
	f.AddEdge(NewEdge(source.ID, target.ID))
}

// NodeAttributes allow additional information, like the kind of the node and active telemetry types, to be associated with nodes.
type NodeAttributes map[string]any

// MakeAttributes adds the expected attributes to a map for a Node
func MakeAttributes(resourceKind, resourceID string) NodeAttributes {
	return map[string]any{
		"kind":       resourceKind,
		"resourceId": resourceID,
	}
}
