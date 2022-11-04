import { Edge, Position, MarkerType, Node } from "react-flow-renderer";
import { TELEMETRY_SIZE_METRICS } from "../../components/MeasurementControlBar/MeasurementControlBar";
import { isSourceID } from "../../components/PipelineGraph/Nodes/ProcessorNode";
import { Graph, GraphMetric } from "../../graphql/generated";

export const GRAPH_NODE_OFFSET = 160;
export const GRAPH_PADDING = 120;

export const enum Page {
  Overview,
  Configuration,
}

export const enum MetricPosition {
  SourceBeforeProcessors, // s0
  SourceAfterProcessors, // s1
  DestinationBeforeProcessors, // d0
  DestinationAfterProcessors, // d1
  Configuration,
}

export function getNodesAndEdges(
  page: Page,
  graph: Graph,
  targetOffsetMultiplier: number
): {
  nodes: Node[];
  edges: Edge[];
} {
  const nodes: Node[] = [];
  const edges: Edge[] = [];

  const offsets = pipelineOffsets(graph.edges || []);

  let y = 0;

  // This number gives the vertical spacing between cards.
  // TODO: find the Card's bounding box, use that number to govern
  //       the spacing.

  const offset = GRAPH_NODE_OFFSET;
  const sourceOffsetMultiplier = 250;

  // This number gives the horizontal spacing between sources and targets
  // TODO: make function of Cards bounding box
  const processorYoffset = 35;
  const targetProcOffsetMultiplier = 270;

  for (let i = 0; i < (graph.sources ?? []).length; i++) {
    const n = graph.sources[i];
    const x = offsets[n.id] * sourceOffsetMultiplier;

    nodes.push({
      id: n.id,
      data: {
        id: n.id,
        label: n.label,
        attributes: n.attributes,
        connectedNodesAndEdges: [n.id],
      },
      position: { x, y },
      sourcePosition: Position.Right,
      type: n.type,
    });

    y += offset;
  }

  y = 0;

  // layout source processors
  for (let i = 0; i < (graph.intermediates ?? []).length; i++) {
    const n = graph.intermediates[i];

    const x = offsets[n.id] * sourceOffsetMultiplier;

    if (isSourceID(n.id)) {
      nodes.push({
        id: `${n.id}`,
        data: {
          id: n.id,
          attributes: n.attributes,
        },
        position: { x, y: y + processorYoffset },
        type: n.type,
      });
    }
    y += offset;
  }

  y =
    (((graph.sources?.length || 1) - (graph.targets?.length || 1)) * offset) /
    2;

  // layout destination processors
  for (let i = 0; i < (graph.intermediates ?? []).length; i++) {
    const n = graph.intermediates[i];

    const x = offsets[n.id] * targetProcOffsetMultiplier;

    if (!isSourceID(n.id)) {
      nodes.push({
        id: `${n.id}`,
        data: {
          id: n.id,
          attributes: n.attributes,
        },
        position: { x, y: y + processorYoffset },
        type: n.type,
      });
      y += offset;
    }
  }

  // save the y position to align the add source & add destination buttons
  const bottomY = y;

  y =
    (((graph.sources?.length || 1) - (graph.targets?.length || 1)) * offset) /
    2;
  for (let i = 0; i < (graph.targets ?? []).length; i++) {
    const n = graph.targets[i];
    const x = offsets[n.id] * targetOffsetMultiplier;

    nodes.push({
      id: n.id,
      data: {
        id: n.id,
        label: n.label,
        attributes: n.attributes,
        connectedNodesAndEdges: [n.id],
      },
      position: { x, y },
      targetPosition: Position.Left,
      type: n.type,
    });
    y += offset;
  }

  // find max pipeline position
  let max = 0;

  // This seems like it should be Object.keys(offsets), but alas...
  for (const key in offsets) {
    if (offsets[key] > max) {
      max = offsets[key];
    }
  }

  y = Math.max(bottomY, y);
  y += offset / 4;

  for (const e of graph.edges || []) {
    const edge: Edge<any> & { key: string } = {
      key: e.id,
      id: e.id,
      source: e.source,
      target: e.target,
      markerEnd: {
        type: MarkerType.ArrowClosed,
      },
      data: {
        connectedNodesAndEdges: [e.id],
      },
      type: page === Page.Configuration ? "configurationEdge" : "overviewEdge",
    };

    edges.push(edge);
    nodes.forEach((node) => {
      if (node.id === e.source) {
        edge.data.attributes = node.data.attributes;
      }
    });
  }
  
  // First create a map from node id to all child nodes and edges
  const nodeMap: NodeMap = {};
  for (const edge of edges) {
    if (!nodeMap[edge.source]) {
      nodeMap[edge.source] = [];
    }
    nodeMap[edge.source].push(edge.target, edge.id);
  }
  nodes.forEach((node) => {
    node.data.connectedNodesAndEdges = getConnectedNodesAndEdges(node.id, nodeMap);
  });

  // Next, create a reverse map
  const reverseNodeMap: NodeMap = {};
  for (const edge of edges) {
    if (!reverseNodeMap[edge.target]) {
      reverseNodeMap[edge.target] = [];
    }
    reverseNodeMap[edge.target].push(edge.source, edge.id);
  }
  // Then recursively find all connected nodes and edges
  nodes.forEach((node) => {
    node.data.connectedNodesAndEdges.push(...getConnectedNodesAndEdges(node.id, reverseNodeMap));
  });

  return { nodes, edges };
}

type NodeMap = { [id: string]: string[] };
function getConnectedNodesAndEdges(
  nodeID: string,
  nodeMap: NodeMap
): string[] {
  const connectedIDS : string[] = [];
  if (!nodeMap[nodeID]) {
    return connectedIDS;
  }
  connectedIDS.push(nodeID)
  nodeMap[nodeID].forEach((id) => {
    connectedIDS.push(id);
    connectedIDS.push(...getConnectedNodesAndEdges(id, nodeMap));
  });
  return connectedIDS;
}

/**
 *
 * @param edges array of GraphEdge
 * @returns [id: string]: number
 * @remarks
 * map of source names : offsets
 *
 * For example:
 *  ____________       ____________      ____________
 *  |          |       |          |      |          |
 *  | source 1 | --->  | inter  1 | ---> | dest   1 |
 *  |          |       |          |      |          |
 *  ____________       ____________      ____________
 *
 * offset  0                1                  2
 *
 * {"source1": 0, "inter 1": 1, "dest 1": 2}
 */
export function pipelineOffsets(edges: { source: string; target: string }[]): {
  [id: string]: number;
} {
  const lens: { [source: string]: number } = {};

  function pipelineLength(source: string): number {
    if (lens[source] != null) {
      return lens[source];
    }
    const lengths: number[] = [0];
    for (const edge of edges) {
      if (source === edge.source) {
        lengths.push(pipelineLength(edge.target) + 1);
      }
    }

    const l = Math.max(...lengths);
    lens[source] = l;
    return l;
  }

  const max = Math.max(...edges.map((e) => pipelineLength(e.source)));

  const result: { [id: string]: number } = {};
  for (const [id, len] of Object.entries(lens)) {
    if (len === 0) {
      result[id] = max;
    } else {
      result[id] = max - len;
    }
  }

  return result;
}

function getMetricPosition(nodeID: string): MetricPosition | undefined {
  if (nodeID.startsWith("source/")) {
    if (nodeID.endsWith("/processors")) {
      return MetricPosition.SourceAfterProcessors;
    } else {
      return MetricPosition.SourceBeforeProcessors;
    }
  } else if (nodeID.startsWith("destination/")) {
    if (nodeID.endsWith("/processors")) {
      return MetricPosition.DestinationBeforeProcessors;
    } else {
      return MetricPosition.DestinationAfterProcessors;
    }
  } else if (nodeID.startsWith("configuration/")) {
    return MetricPosition.Configuration;
  }
}

export function updateMetricData(
  page: Page,
  nodes: Node<any>[],
  edges: Edge<any>[],
  metrics: GraphMetric[],
  rate: string,
  telemetryType: string
) {
  for (const node of nodes) {
    const metric = metrics.find(
      (m) =>
        m.nodeID === node.id && m.name === TELEMETRY_SIZE_METRICS[telemetryType]
    );
    if (metric != null) {
      const formattedMetric = formatMetric(metric, rate);

      var startOffset = "50%";
      var textAnchor = "middle";
      var edge: Edge<any> | undefined;

      // put this metric on the associated edge:
      //
      // Configuration page:
      //
      // s0 => s1 ====> d0 => d1
      //    A     B  C     D
      //
      // sX metrics go on the edges that match the source
      //
      // dX metrics go on the edges that match the target
      //
      const position = getMetricPosition(node.id);
      switch (position) {
        case MetricPosition.SourceBeforeProcessors:
          // A
          edge = edges.find((e) => e.source === node.id);
          break;

        case MetricPosition.SourceAfterProcessors:
          // B
          edge = edges.find((e) => e.source === node.id);
          textAnchor = "start";
          startOffset = "4%";
          break;

        case MetricPosition.DestinationBeforeProcessors:
          // C
          edge = edges.find((e) => e.target === node.id);
          textAnchor = "end";
          startOffset = "97%";
          break;

        case MetricPosition.DestinationAfterProcessors:
          // D
          edge = edges.find((e) => e.target === node.id);
          if (page === Page.Overview) {
            textAnchor = "end";
            startOffset = "97%";
          }
          break;

        case MetricPosition.Configuration:
          edge = edges.find((e) => e.source === node.id);
          textAnchor = "start";
          startOffset = "4%";
          break;
      }

      if (edge != null) {
        edge.data.metrics ||= [];
        edge.data.metrics.push({
          startOffset,
          textAnchor,
          value: formattedMetric,
        });
      }
    } else {
      node.data.metric = "";
    }
  }
}

export interface MetricValue {
  value: number;
  unit: string;
}

export function formatMetric(metric: MetricValue, rate: string): string {
  let { value, unit } = metric;
  let units: string[];

  if (rate.endsWith("m")) {
    unit = unit.replace("/s", "/m");
    value *= 60;
    units = getUnits("/m");
  } else if (rate.endsWith("h")) {
    unit = unit.replace("/s", "/h");
    value *= 60 * 60;
    units = getUnits("/h");
  } else {
    units = getUnits("/s");
  }

  const converted = convertUnits({ value, unit }, units);
  return `${converted.value} ${converted.unit}`;
}

function getUnits(rate: string): string[] {
  return ["B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"].map(
    (size) => `${size}${rate}`
  );
}
export function convertUnits(
  metric: MetricValue,
  units: string[]
): MetricValue {
  const unitIndex = units.findIndex((u) => u === metric.unit) ?? 0;
  if (metric.value > 1024) {
    return convertUnits(
      {
        value: metric.value / 1024,
        unit: units[unitIndex + 1],
      },
      units
    );
  }
  const round = 10;
  return { value: Math.round(metric.value * round) / round, unit: metric.unit };
}
