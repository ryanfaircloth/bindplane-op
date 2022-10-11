import { memo } from "react";
import { Handle, Position } from "react-flow-renderer";
import { CardMeasurementContent } from "../../../components/CardMeasurementContent/CardMeasurementContent";
import { InlineDestinationCard } from "../../../components/Cards/InlineDestinationCard";
import { ResourceDestinationCard } from "../../../components/Cards/ResourceDestinationCard";
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
  const { setHoveredNodeAndEdgeSet, hoveredSet } = useOverviewPage();
  return (
    <div
      onMouseEnter={() => setHoveredNodeAndEdgeSet(connectedNodesAndEdges)}
      onMouseLeave={() => setHoveredNodeAndEdgeSet([])}
    >
      {attributes.isInline ? (
        <InlineDestinationCard id={id.replace("destination/", "")} key={id} />
      ) : (
        <ResourceDestinationCard
          name={attributes.resourceId}
          disabled={
            hoveredSet.length > 0 && !hoveredSet.find((elem) => elem === id)
          }
          key={id}
        />
      )}
      <CardMeasurementContent>{metric}</CardMeasurementContent>
      <Handle type="target" position={Position.Left} />
    </div>
  );
}

export default memo(OverviewDestinationNode);
