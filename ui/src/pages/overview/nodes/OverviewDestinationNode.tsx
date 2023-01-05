import { memo } from "react";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../../components/CardMeasurementContent/CardMeasurementContent";
import { InlineDestinationCard } from "../../../components/Cards/InlineDestinationCard";
import { OverviewDestinationCard } from "../../../components/Cards/OverviewDestinationCard";
import { DEFAULT_TELEMETRY_TYPE } from "../../../components/MeasurementControlBar/MeasurementControlBar";
import { isNodeDisabled } from "../../../components/PipelineGraph/Nodes/nodeUtils";
import { useOverviewPage } from "../OverviewPageContext";

export function OverviewDestinationNode(params: {
  data: {
    pipelineType: string;
    id: string;
    label: string;
    attributes: Record<string, any>;
    metric: string;
    connectedNodesAndEdges: string[];
  };
}): JSX.Element {
  const { id, attributes, metric, connectedNodesAndEdges } = params.data;
  const { setHoveredNodeAndEdgeSet, hoveredSet, selectedTelemetry } =
    useOverviewPage();

  const isDisabled = isNodeDisabled(
    selectedTelemetry || DEFAULT_TELEMETRY_TYPE,
    params.data.attributes
  );
  const isNotInHoverSet =
    hoveredSet.length > 0 &&
    !hoveredSet.find((elem) => elem === params.data.id);
  return (
    <div
      onMouseEnter={() => setHoveredNodeAndEdgeSet(connectedNodesAndEdges)}
      onMouseLeave={() => setHoveredNodeAndEdgeSet([])}
    >
      {attributes.isInline ? (
        <InlineDestinationCard id={id.replace("destination/", "")} key={id} />
      ) : (
        <OverviewDestinationCard
          name={attributes.resourceId}
          disabled={isDisabled || isNotInHoverSet}
          key={id}
        />
      )}
      <CardMeasurementContent>{metric}</CardMeasurementContent>
      <Handle type="target" position={Position.Left} />
    </div>
  );
}

export default memo(OverviewDestinationNode);
