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

import styles from "./pipeline-graph.module.scss";
import { firstActiveTelemetry } from "./Nodes/nodeUtils";

export type MinimumRequiredConfig = NonNullable<ShowPageConfig>;

interface PipelineGraphProps {
  configuration: MinimumRequiredConfig;
}

export const PipelineGraph: React.FC<PipelineGraphProps> = ({
  configuration,
}) => {
  const [selectedTelemetry, setTelemetry] = useState(
    firstActiveTelemetry(configuration.graph?.attributes) ??
      DEFAULT_TELEMETRY_TYPE
  );
  const [selectedPeriod, setPeriod] = useState(DEFAULT_PERIOD);

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
            />
          </ReactFlowProvider>
        </Card>
      </GraphContainer>
    </PipelineGraphProvider>
  );
};

const GraphContainer: React.FC = ({ children }) => {
  return (
    <Paper classes={{ root: styles.container }} elevation={1}>
      {children}
    </Paper>
  );
};
