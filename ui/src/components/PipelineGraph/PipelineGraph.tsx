import { Card, Paper } from "@mui/material";
import { ReactFlowProvider } from "react-flow-renderer";
import { ShowPageConfig } from "../../pages/configurations/configuration";
import {
  DEFAULT_PERIOD,
  DEFAULT_TELEMETRY_TYPE,
  MeasurementControlBar,
} from "../MeasurementControlBar/MeasurementControlBar";
import { ConfigurationFlow } from "./ConfigurationFlow";
import { PipelineGraphProvider } from "./PipelineGraphContext";
import { firstActiveTelemetry } from "./Nodes/nodeUtils";
import { YamlEditor } from "../YamlEditor";
import { ProcessorDialog } from "../ResourceDialog/ProcessorsDialog";

import { useQueryParam, StringParam } from "use-query-params";

import styles from "./pipeline-graph.module.scss";
import {
  PERIOD_PARAM_NAME,
  TELEMETRY_TYPE_PARAM_NAME,
} from "../../utils/state";
import { useEffect } from "react";

export type MinimumRequiredConfig = Partial<ShowPageConfig>;
interface PipelineGraphProps {
  configuration: MinimumRequiredConfig;
  refetchConfiguration: () => void;
  agent: string;
  yamlValue: string;
  rawOrTopology: string;
}

export const PipelineGraph: React.FC<PipelineGraphProps> = ({
  configuration,
  refetchConfiguration,
  agent,
  yamlValue,
  rawOrTopology,
}) => {
  const [selectedPeriod, setPeriodURL] = useQueryParam(
    PERIOD_PARAM_NAME,
    StringParam
  );
  const setPeriod = (p: string) => {
    setPeriodURL(p, "replaceIn");
  };

  const [selectedTelemetry, setSelectedTelemetryURL] = useQueryParam(
    TELEMETRY_TYPE_PARAM_NAME,
    StringParam
  );
  const setSelectedTelemetry = (t: string) => {
    setSelectedTelemetryURL(t, "replaceIn");
  };

  const activeTelemetry =
    firstActiveTelemetry(configuration?.graph?.attributes) ??
    DEFAULT_TELEMETRY_TYPE;
  useEffect(() => {
    if (!selectedPeriod) {
      setPeriod(DEFAULT_PERIOD);
    }
    if (!selectedTelemetry) {
      setSelectedTelemetry(activeTelemetry);
    }
  });

  if (rawOrTopology === "topology") {
    return (
      <PipelineGraphProvider
        selectedTelemetryType={selectedTelemetry || DEFAULT_PERIOD}
        configuration={configuration}
        refetchConfiguration={refetchConfiguration}
      >
        <GraphContainer>
          <Card style={{ border: 0 }}>
            <MeasurementControlBar
              telemetry={selectedTelemetry || activeTelemetry}
              onTelemetryTypeChange={setSelectedTelemetry}
              period={selectedPeriod || DEFAULT_PERIOD}
              onPeriodChange={setPeriod}
            />
            <ReactFlowProvider>
              <ConfigurationFlow
                period={selectedPeriod || DEFAULT_PERIOD}
                selectedTelemetry={selectedTelemetry || activeTelemetry}
                configuration={configuration}
                refetchConfiguration={refetchConfiguration}
                agent={agent}
              />
            </ReactFlowProvider>
          </Card>
        </GraphContainer>
        <ProcessorDialog />
      </PipelineGraphProvider>
    );
  } else {
    return (
      <YamlEditor
        value={yamlValue === "" ? "Raw configuration unavailable" : yamlValue}
        readOnly
        limitHeight
      />
    );
  }
};

const GraphContainer: React.FC = ({ children }) => {
  return (
    <Paper classes={{ root: styles.container }} elevation={1}>
      {children}
    </Paper>
  );
};
