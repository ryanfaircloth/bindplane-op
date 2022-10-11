import { Edge, Position, MarkerType, Node } from "react-flow-renderer";
import { TELEMETRY_SIZE_METRICS } from '../../components/MeasurementControlBar/MeasurementControlBar';
import { isSourceID } from "../../components/PipelineGraph/Nodes/ProcessorNode";
import { Graph, GraphMetric } from "../../graphql/generated";

export const GRAPH_NODE_OFFSET = 150;
export const GRAPH_PADDING = 120;

export function getNodesAndEdges(
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
  const sourceOffsetMultiplier = 200;

  // This number gives the horizontal spacing between sources and targets
  // TODO: make function of Cards bounding box
  const processorYoffset = 25;
  const targetProcOffsetMultiplier = 300;

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
    let edge = {
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
      type: "overviewEdge",
    };

    edges.push(edge);
    nodes.forEach((node) => {
      if (node.id === e.source) {
        node.data.connectedNodesAndEdges?.push(e.target, e.id);
      }
      if (node.id === e.target) {
        node.data.connectedNodesAndEdges?.push(e.source, e.id);
      }
    });
  }

  return { nodes, edges };
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

export function updateMetricData(
  nodes: Node<any>[],
  metrics: GraphMetric[],
  rate: string,
  telemetryType: string
) {
  for (const node of nodes) {
    const metric = metrics.find(
      (m) => m.nodeID === node.id && m.name === TELEMETRY_SIZE_METRICS[telemetryType]
    );
    if (metric != null) {
      node.data.metric = formatMetric(metric, rate);
    } else if (
      node.id.startsWith("destination/") &&
      node.id.endsWith("/processors")
    ) {
      // The destination processor node will look like destination/destinationName/processors. If that couldn't be found,
      // then there are no destination processors so just use the metric from the destination.
      const destinationId = node.id.replace("/processors", "");
      const metric = metrics.find(
        (m) =>
          m.nodeID === destinationId && m.name === TELEMETRY_SIZE_METRICS[telemetryType]
      );
      if (metric != null) {
        node.data.metric = formatMetric(metric, rate);
      } else {
        node.data.metric = "";
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
