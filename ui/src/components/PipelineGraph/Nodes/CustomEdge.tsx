import { EdgeProps, getBezierPath } from "react-flow-renderer";
import { classes } from '../../../utils/styles';
import { isNodeDisabled } from "./nodeUtils";

import styles from "./graph.styles.module.scss";

export interface CustomEdgeData {
  connectedNodesAndEdges: string[];
  metrics: {
    value: string;
    startOffset: string;
    textAnchor: string;
  }[];
  active: boolean;
  attributes: Record<string, any>;
}

export interface CustomEdgeProps extends EdgeProps<CustomEdgeData> {
  hoveredSet?: string[];
  className?: string;
  telemetry?: string;
};

const CustomEdge: React.FC<CustomEdgeProps> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  markerEnd,
  hoveredSet,
  className,
  telemetry,
}) => {
  hoveredSet ||= [];

  const edgePath = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const inactive = isNodeDisabled(telemetry || "", data?.attributes || {});
  const dimmed = hoveredSet.length > 0 && !hoveredSet.includes(id);
  const highlight = hoveredSet.includes(id);

  const metrics: JSX.Element[] = [];
  const strokeWidth = highlight && !inactive ? 1.5 : 1;
  const strokeOpacity = dimmed ? 0.3 : (inactive ? 0.2 : undefined);

  var index = 0;
  for (const m of data?.metrics || []) {
    const metric = (
      // transform moves the metric a few pixels off the line
      <g transform={`translate(0 -3)`} key={`metric${index++}`}>
        <text>
          <textPath
            className={classes([styles["metric"], className ? styles[className] : undefined])}
            href={`#${id}`}
            startOffset={m.startOffset}
            textAnchor={m.textAnchor}
            spacing="auto"
          >
            {m.value}
          </textPath>
        </text>
      </g>
    );
    metrics.push(metric);
  }

  return (
    <>
      <path
        id={id}
        style={{
          strokeWidth,
          strokeOpacity,
        }}
        className="react-flow__edge-path"
        d={edgePath}
        markerEnd={markerEnd}
      />
      {metrics}
    </>
  );
};

export default CustomEdge;