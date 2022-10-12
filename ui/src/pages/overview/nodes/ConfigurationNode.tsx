import { Handle, Position } from "react-flow-renderer";
import { ConfigurationCard } from "../../../components/Cards/ConfigurationCard";
import { isNodeDisabled } from "../../../components/PipelineGraph/Nodes/nodeUtils";
import { useOverviewPage } from "../OverviewPageContext";

export function ConfigurationNode(params: {
  data: {
    id: string;
    label: string;
    attributes: Record<string, any>;
    metric: string;
    connectedNodesAndEdges: string[];
  };
}) {
  const { hoveredSet, setHoveredNodeAndEdgeSet, selectedTelemetry } =
    useOverviewPage();

  const isDisabled = isNodeDisabled(selectedTelemetry, params.data.attributes);
  const isNotInHoverSet =
    hoveredSet.length > 0 &&
    !hoveredSet.find((elem) => elem === params.data.id);

  return (
    <div
      onMouseEnter={() =>
        setHoveredNodeAndEdgeSet(params.data.connectedNodesAndEdges)
      }
      onMouseLeave={() => setHoveredNodeAndEdgeSet([])}
    >
      <ConfigurationCard
        {...params.data}
        disabled={isDisabled || isNotInHoverSet}
      />
      <Handle type={"source"} position={Position.Right} />
    </div>
  );
}

export function formatAgentCount(agentCount: number): string {
  switch (agentCount) {
    case 0:
      return "";
    case 1:
      return "1 agent";
    default:
      return `${agentCount} agents`;
  }
}
