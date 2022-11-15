import { memo } from "react";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../CardMeasurementContent/CardMeasurementContent";
import { InlineSourceCard } from "../../Cards/InlineSourceCard";
import { MinimumRequiredConfig } from "../PipelineGraph";
import { usePipelineGraph } from "../PipelineGraphContext";
import { isNodeDisabled } from "./nodeUtils";

function SourceNode({
  data,
}: {
  data: {
    id: string;
    metric: string;
    attributes: Record<string, any>;
    connectedNodesAndEdges: string[];
    configuration: MinimumRequiredConfig;
    refetchConfiguration: () => void;
  };
}) {
  const { id, metric, attributes, configuration, refetchConfiguration } = data;

  const { hoveredSet, setHoveredNodeAndEdgeSet, selectedTelemetryType } =
    usePipelineGraph();

  const isDisabled = isNodeDisabled(selectedTelemetryType, attributes);
  const isNotInHoverSet =
    hoveredSet.length > 0 && !hoveredSet.find((elem) => elem === data.id);

  return (
    <div
      onMouseEnter={() => {
        setHoveredNodeAndEdgeSet(data.connectedNodesAndEdges);
      }}
      onMouseLeave={() => setHoveredNodeAndEdgeSet([])}
    >
      <InlineSourceCard
        id={id.replace("source/", "")}
        disabled={isDisabled || isNotInHoverSet}
        configuration={configuration}
        refetchConfiguration={refetchConfiguration}
      />
      <CardMeasurementContent>{metric}</CardMeasurementContent>

      <Handle type="source" position={Position.Right} />
    </div>
  );
}

export default memo(SourceNode);
