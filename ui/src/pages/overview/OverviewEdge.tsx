import React, { memo } from "react";
import { EdgeProps } from "react-flow-renderer";
import { useOverviewPage } from "./OverviewPageContext";

import CustomEdge, {
  CustomEdgeData,
} from "../../components/PipelineGraph/Nodes/CustomEdge";
import { DEFAULT_TELEMETRY_TYPE } from "../../components/MeasurementControlBar/MeasurementControlBar";

const OverviewEdge: React.FC<EdgeProps<CustomEdgeData>> = (props) => {
  const { hoveredSet, selectedTelemetry } = useOverviewPage();

  return CustomEdge({
    ...props,
    hoveredSet: hoveredSet,
    className: "overview-metric",
    telemetry: selectedTelemetry || DEFAULT_TELEMETRY_TYPE,
  });
};

export default memo(OverviewEdge);
