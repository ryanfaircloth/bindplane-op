import { memo } from 'react';
import { EdgeProps } from "react-flow-renderer";
import { usePipelineGraph } from '../PipelineGraphContext';
import CustomEdge, { CustomEdgeData } from './CustomEdge';

const ConfigurationEdge: React.FC<EdgeProps<CustomEdgeData>> = (props) => {
  const { selectedTelemetryType, hoveredSet } = usePipelineGraph()

  return CustomEdge({
    ...props,
    hoveredSet: hoveredSet,
    telemetry: selectedTelemetryType,
  });
};

export default memo(ConfigurationEdge);
