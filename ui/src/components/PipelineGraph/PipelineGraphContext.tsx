import { createContext, useContext } from "react";
import { DEFAULT_TELEMETRY_TYPE } from "../MeasurementControlBar/MeasurementControlBar";

export interface PipelineGraphContextValue {
  selectedTelemetryType: string;
}

const defaultValue: PipelineGraphContextValue = {
  selectedTelemetryType: DEFAULT_TELEMETRY_TYPE,
};

const PipelineContext = createContext(defaultValue);

export const PipelineGraphProvider: React.FC<PipelineGraphContextValue> = ({
  children,
  selectedTelemetryType,
}) => {
  return (
    <PipelineContext.Provider value={{ selectedTelemetryType }}>
      {children}
    </PipelineContext.Provider>
  );
};

export function usePipelineGraph(): PipelineGraphContextValue {
  return useContext(PipelineContext);
}
