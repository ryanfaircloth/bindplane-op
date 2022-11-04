import { createContext, useContext, useState } from "react";
import { DEFAULT_TELEMETRY_TYPE } from "../MeasurementControlBar/MeasurementControlBar";

export interface PipelineGraphContextValue {
  selectedTelemetryType: string;
  hoveredSet: string[];
  setHoveredNodeAndEdgeSet: React.Dispatch<React.SetStateAction<string[]>>;  
}
export interface PipelineGraphProviderProps {
  selectedTelemetryType: string;
}

const defaultValue: PipelineGraphContextValue = {
  selectedTelemetryType: DEFAULT_TELEMETRY_TYPE,
  hoveredSet: [],
  setHoveredNodeAndEdgeSet: () => {},
};

const PipelineContext = createContext(defaultValue);

export const PipelineGraphProvider: React.FC<PipelineGraphProviderProps> = ({
  children,
  selectedTelemetryType,
}) => {
  const [hoveredSet, setHoveredNodeAndEdgeSet] = useState<string[]>([]);
  return (
    <PipelineContext.Provider value={{ 
        setHoveredNodeAndEdgeSet,
        hoveredSet,
        selectedTelemetryType }}>
      {children}
    </PipelineContext.Provider>
  );
};

export function usePipelineGraph(): PipelineGraphContextValue {
  return useContext(PipelineContext);
}
