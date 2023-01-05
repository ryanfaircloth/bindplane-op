import { GridSelectionModel } from "@mui/x-data-grid";

import { createContext, useContext, useState } from "react";
import {
  CONFIGS_PARAM_NAME,
  DESTINATIONS_PARAM_NAME,
  PERIOD_PARAM_NAME,
  TELEMETRY_TYPE_PARAM_NAME,
  useQueryParamWrapper,
} from "../../utils/state";

import { useQueryParam, StringParam } from "use-query-params";

export interface OverviewPageContextValue {
  hoveredSet: string[];
  setHoveredNodeAndEdgeSet: React.Dispatch<React.SetStateAction<string[]>>;
  selectedTelemetry: string | null | undefined;
  setSelectedTelemetry: (t: string) => void;
  selectedConfigs: GridSelectionModel;
  setSelectedConfigs: (d: GridSelectionModel) => void;
  selectedDestinations: GridSelectionModel;
  setSelectedDestinations: (d: GridSelectionModel) => void;
  selectedPeriod: string | null | undefined;
  setPeriod: (p: string) => void;
  editingDestination: string | null;
  setEditingDestination: (dest: string | null) => void;
  loadTop: boolean;
  setLoadTop: (loadTop: boolean) => void;
}

const defaultContext: OverviewPageContextValue = {
  hoveredSet: [],
  setHoveredNodeAndEdgeSet: () => {},
  selectedTelemetry: "",
  setSelectedTelemetry: () => {},
  selectedConfigs: [],
  setSelectedConfigs: () => {},
  selectedDestinations: [],
  setSelectedDestinations: () => {},
  selectedPeriod: "",
  setPeriod: () => {},
  editingDestination: null,
  setEditingDestination: () => {},
  loadTop: true,
  setLoadTop: () => {},
};

const OverviewPageContext = createContext(defaultContext);

export const OverviewPageProvider: React.FC = ({ children }) => {
  const [selectedPeriod, setPeriodURL] = useQueryParam(
    PERIOD_PARAM_NAME,
    StringParam
  );

  const setPeriod = (p: string) => {
    setPeriodURL(p, "replaceIn");
  };

  // state for knowing which node is being hovered over
  const [hoveredSet, setHoveredNodeAndEdgeSet] = useState<string[]>([]);

  const [selectedTelemetry, setSelectedTelemetryURL] = useQueryParam(
    TELEMETRY_TYPE_PARAM_NAME,
    StringParam
  );
  const setSelectedTelemetry = (t: string) => {
    setSelectedTelemetryURL(t, "replaceIn");
  };

  const [selectedConfigs, setSelectedConfigs] =
    useQueryParamWrapper<GridSelectionModel>(CONFIGS_PARAM_NAME, []);
  const [selectedDestinations, setSelectedDestinations] =
    useQueryParamWrapper<GridSelectionModel>(DESTINATIONS_PARAM_NAME, []);
  const [loadTop, setLoadTop] = useQueryParamWrapper<boolean>("loadTop", true);

  const [editingDestination, setEditingDestination] =
    useState<string | null>(null);
  return (
    <OverviewPageContext.Provider
      value={{
        setHoveredNodeAndEdgeSet,
        hoveredSet,
        selectedTelemetry,
        setSelectedTelemetry,
        selectedConfigs,
        setSelectedConfigs,
        selectedDestinations,
        setSelectedDestinations,
        selectedPeriod,
        setPeriod,
        editingDestination,
        setEditingDestination,
        loadTop,
        setLoadTop,
      }}
    >
      {children}
    </OverviewPageContext.Provider>
  );
};

export function useOverviewPage(): OverviewPageContextValue {
  return useContext(OverviewPageContext);
}
