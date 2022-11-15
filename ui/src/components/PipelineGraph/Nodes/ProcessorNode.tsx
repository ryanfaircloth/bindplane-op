import { isNumber } from "lodash";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../CardMeasurementContent/CardMeasurementContent";
import { ProcessorCard } from "../../Cards/ProcessorCard";
import { MinimumRequiredConfig } from "../PipelineGraph";

export function ProcessorNode({
  data,
}: {
  data: {
    // ID will have the form source/source0/processors, source/source1/processors, dest-name/destination/processors, etc
    id: string;
    metric: string;
    attributes: Record<string, any>;
    configuration: MinimumRequiredConfig;
  };
}) {
  const { id, metric, configuration } = data;

  let processorCount = 0;

  if (isSourceID(id)) {
    const source = configuration?.spec?.sources![getSourceIndex(id)];
    processorCount = source?.processors?.length ?? 0;
  } else {
    const destination = configuration?.spec?.destinations?.find(
      (d) => d.name === getDestinationName(id)
    );
    processorCount = destination?.processors?.length ?? 0;
  }
  return (
    <>
      <Handle type="target" position={Position.Left} />
      <ProcessorCard processors={processorCount} />
      <CardMeasurementContent>{metric}</CardMeasurementContent>
      <Handle type="source" position={Position.Right} />
    </>
  );
}

export function isSourceID(id: string): boolean {
  return id.startsWith("source/");
}

export function getSourceIndex(id: string): number {
  const REGEX = /^source\/source(?<sourceNum>[0-9]+)\/processors$/;
  const match = id.match(REGEX);
  if (match?.groups != null) {
    const index = Number(match.groups["sourceNum"]);
    if (isNumber(index)) {
      return index;
    }
  }
  return -1;
}

export function getDestinationName(id: string): string {
  // /destination/name/processors
  const REGEX = /^destination\/(?<name>.*)\/processors$/;
  const match = id.match(REGEX);
  return match?.groups?.["name"] ?? "";
}
