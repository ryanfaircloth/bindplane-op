import { memo } from 'react';
import { EdgeProps } from "react-flow-renderer";
import CustomEdge, { CustomEdgeData } from './CustomEdge';

const ConfigurationEdge: React.FC<EdgeProps<CustomEdgeData>> = (props) => {
  return CustomEdge({
    ...props,
  });
};

export default memo(ConfigurationEdge);