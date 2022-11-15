import { Card, Paper } from "@mui/material";
import { useState } from "react";
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

import styles from "./pipeline-graph.module.scss";

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
  const [selectedTelemetry, setTelemetry] = useState(
    firstActiveTelemetry(configuration?.graph?.attributes) ??
      DEFAULT_TELEMETRY_TYPE
  );
  const [selectedPeriod, setPeriod] = useState(DEFAULT_PERIOD);

  if (rawOrTopology === "topology") {
    return (
      <PipelineGraphProvider selectedTelemetryType={selectedTelemetry}>
        <GraphContainer>
          <Card style={{ border: 0 }}>
            <MeasurementControlBar
              telemetry={selectedTelemetry}
              onTelemetryTypeChange={(t: string) => setTelemetry(t)}
              period={selectedPeriod}
              onPeriodChange={(r: string) => setPeriod(r)}
            />
            <ReactFlowProvider>
              <ConfigurationFlow
                period={selectedPeriod}
                selectedTelemetry={selectedTelemetry}
                configuration={configuration}
                refetchConfiguration={refetchConfiguration}
                agent={agent}
              />
            </ReactFlowProvider>
          </Card>
        </GraphContainer>
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
