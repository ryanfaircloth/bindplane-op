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

import styles from "./pipeline-graph.module.scss";
import { PipelineGraphProvider } from "./PipelineGraphContext";

export type MinimumRequiredConfig = NonNullable<ShowPageConfig>;

export const PipelineGraph: React.FC = () => {
  const [selectedTelemetry, setTelemetry] = useState(DEFAULT_TELEMETRY_TYPE);
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
