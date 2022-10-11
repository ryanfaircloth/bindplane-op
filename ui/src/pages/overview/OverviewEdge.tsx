import React from "react";
import { EdgeProps, getBezierPath } from "react-flow-renderer";
import { useOverviewPage } from "./OverviewPageContext";

interface EdgeData {
  connectedNodesAndEdges: string[];
}

export const OverviewEdge: React.FC<EdgeProps<EdgeData>> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  markerEnd,
}) => {
  const { hoveredSet } = useOverviewPage();

  const edgePath = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const dimmed =
    hoveredSet.length > 0 && !hoveredSet.find((elem) => elem === id);
  const highlight =
    hoveredSet.length > 0 && hoveredSet.find((elem) => elem === id);
  return (
    <>
      <path
        id={id}
        style={{
          strokeWidth: highlight ? 1.5 : undefined,
          strokeOpacity: dimmed ? 0.3 : undefined,
        }}
        className="react-flow__edge-path"
        d={edgePath}
        markerEnd={markerEnd}
      />
    </>
  );
};
