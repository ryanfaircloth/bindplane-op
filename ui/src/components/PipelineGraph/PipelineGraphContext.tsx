import { createContext, useContext, useState } from "react";
import { DEFAULT_TELEMETRY_TYPE } from "../MeasurementControlBar/MeasurementControlBar";
import { MinimumRequiredConfig } from "./PipelineGraph";

export interface PipelineGraphContextValue {
  configuration: MinimumRequiredConfig;
  refetchConfiguration: () => void;

  selectedTelemetryType: string;
  hoveredSet: string[];
  setHoveredNodeAndEdgeSet: React.Dispatch<React.SetStateAction<string[]>>;

  editProcessorsInfo: EditProcessorsInfo | null;
  // editProcessor opens up the editing dialog for a source or destination
  editProcessors: (
    resourceType: "source" | "destination",
    index: number
  ) => void;
  closeProcessorDialog(): void;
  editProcessorsOpen: boolean;
}
export interface PipelineGraphProviderProps {
  configuration: MinimumRequiredConfig;
  refetchConfiguration: () => void;
  selectedTelemetryType: string;
}

interface EditProcessorsInfo {
  resourceType: "source" | "destination";
  index: number;
}

const defaultValue: PipelineGraphContextValue = {
  configuration: {
    __typename: undefined,
    metadata: {
      __typename: undefined,
      id: "",
      name: "",
      description: undefined,
      labels: undefined,
    },
    spec: {
      __typename: undefined,
      raw: undefined,
      sources: undefined,
      destinations: undefined,
      selector: undefined,
    },
    graph: undefined,
  },
  refetchConfiguration: () => {},
  selectedTelemetryType: DEFAULT_TELEMETRY_TYPE,
  hoveredSet: [],
  setHoveredNodeAndEdgeSet: () => {},
  editProcessors: () => {},
  closeProcessorDialog: () => {},
  editProcessorsInfo: null,
  editProcessorsOpen: false,
};

export const PipelineContext = createContext(defaultValue);

export const PipelineGraphProvider: React.FC<PipelineGraphProviderProps> = ({
  children,
  selectedTelemetryType,
  configuration,
  refetchConfiguration,
}) => {
  const [hoveredSet, setHoveredNodeAndEdgeSet] = useState<string[]>([]);
  const [editProcessorsInfo, setEditingProcessors] =
    useState<{
      resourceType: "source" | "destination";
      index: number;
    } | null>(null);

  const [editProcessorsOpen, setEditProcessorsOpen] = useState(false);

  function editProcessors(
    resourceType: "source" | "destination",
    index: number
  ) {
    setEditingProcessors({ resourceType, index });
    setEditProcessorsOpen(true);
  }
  function closeProcessorDialog() {
    setEditProcessorsOpen(false);
    // Reset the editing processors on a timeout to avoid a flash of empty state.
    setTimeout(() => {
      setEditingProcessors(null);
    }, 300);
  }
  return (
    <PipelineContext.Provider
      value={{
        configuration,
        refetchConfiguration,
        setHoveredNodeAndEdgeSet,
        hoveredSet,
        selectedTelemetryType,
        editProcessors,
        closeProcessorDialog,
        editProcessorsInfo,
        editProcessorsOpen,
      }}
    >
      {children}
    </PipelineContext.Provider>
  );
};

export function usePipelineGraph(): PipelineGraphContextValue {
  return useContext(PipelineContext);
}
