import { memo } from "react";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../CardMeasurementContent/CardMeasurementContent";
import { InlineDestinationCard } from "../../Cards/InlineDestinationCard";
import { ResourceDestinationCard } from "../../Cards/ResourceDestinationCard";
import { usePipelineGraph } from "../PipelineGraphContext";
import { isNodeDisabled } from "./nodeUtils";

function DestinationNode(params: {
  data: {
    pipelineType: string;
    id: string;
    label: string;
    attributes: Record<string, any>;
    metric: string;
    connectedNodesAndEdges: string[];
  };
}): JSX.Element {
  const { id, attributes, metric } = params.data;
  const { selectedTelemetryType } = usePipelineGraph();
  return (
    <>
      {attributes.isInline ? (
        <InlineDestinationCard id={id.replace("destination/", "")} key={id} />
      ) : (
        <ResourceDestinationCard
          name={attributes.resourceId}
          key={id}
          disabled={isNodeDisabled(selectedTelemetryType, attributes)}
        />
      )}
      <CardMeasurementContent>{metric}</CardMeasurementContent>
      <Handle type="target" position={Position.Left} />
    </>
  );
}

export default memo(DestinationNode);
