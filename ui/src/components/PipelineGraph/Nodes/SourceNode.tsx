import { memo } from "react";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../CardMeasurementContent/CardMeasurementContent";
import { InlineSourceCard } from "../../Cards/InlineSourceCard";
import { usePipelineGraph } from "../PipelineGraphContext";
import { isNodeDisabled } from "./nodeUtils";

function SourceNode({
  data,
}: {
  data: {
    id: string;
    metric: string;
    attributes: Record<string, any>;
  };
}) {
  const { id, metric, attributes } = data;

  const { selectedTelemetryType } = usePipelineGraph();

  return (
    <>
      <InlineSourceCard
        id={id.replace("source/", "")}
        disabled={isNodeDisabled(selectedTelemetryType, attributes)}
      />
      <CardMeasurementContent>{metric}</CardMeasurementContent>

      <Handle type="source" position={Position.Right} />
    </>
  );
}

export default memo(SourceNode);
