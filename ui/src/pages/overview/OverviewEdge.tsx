import React, { memo } from "react";
import { EdgeProps } from "react-flow-renderer";
import { useOverviewPage } from "./OverviewPageContext";

import CustomEdge, { CustomEdgeData } from '../../components/PipelineGraph/Nodes/CustomEdge';

const OverviewEdge: React.FC<EdgeProps<CustomEdgeData>> = (props) => {
  const { hoveredSet } = useOverviewPage();
  return CustomEdge({
    ...props,
    hoveredSet: hoveredSet,
    className: 'overview-metric',
  });
}

export default memo(OverviewEdge)
