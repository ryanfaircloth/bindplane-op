import { createContext, useContext, useState } from "react";
import { DEFAULT_TELEMETRY_TYPE } from "../../components/MeasurementControlBar/MeasurementControlBar";

export interface OverviewPageContextValue {
  hoveredSet: string[];
  setHoveredNodeAndEdgeSet: React.Dispatch<React.SetStateAction<string[]>>;
  selectedTelemetry: string;
  onTelemetryTypeChange: (t: string) => void;
}

const defaultContext: OverviewPageContextValue = {
  hoveredSet: [],
  setHoveredNodeAndEdgeSet: () => {},
  selectedTelemetry: DEFAULT_TELEMETRY_TYPE,
  onTelemetryTypeChange: () => {},
};

const OverviewPageContext = createContext(defaultContext);

export const OverviewPageProvider: React.FC = ({ children }) => {
  // state for knowing which node is being hovered over
  const [hoveredSet, setHoveredNodeAndEdgeSet] = useState<string[]>([]);
  const [selectedTelemetry, setSelectedTelemetry] = useState<string>(
    DEFAULT_TELEMETRY_TYPE
  );

  function onTelemetryTypeChange(t: string) {
    setSelectedTelemetry(t);
  }

  return (
    <OverviewPageContext.Provider
      value={{
        setHoveredNodeAndEdgeSet,
        hoveredSet,
        selectedTelemetry,
        onTelemetryTypeChange,
      }}
    >
      {children}
    </OverviewPageContext.Provider>
  );
};

export function useOverviewPage(): OverviewPageContextValue {
  return useContext(OverviewPageContext);
}
