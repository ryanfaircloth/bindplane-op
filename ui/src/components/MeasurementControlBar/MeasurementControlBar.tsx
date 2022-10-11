import { AppBar, Toolbar, Button, Typography, Box } from "@mui/material";

import styles from "./styles.module.scss";

export const PERIODS: { [period: string]: string } = {
  "1m": "1m",
  "5m": "5m",
  "1h": "1h",
  "24h": "24h",
};
export const DEFAULT_PERIOD = "1m";
export const DEFAULT_AGENTS_TABLE_PERIOD = "1m";
export const DEFAULT_CONFIGURATION_TABLE_PERIOD = "1m";
export const TELEMETRY_TYPES: { [telemetryType: string]: string } = {
  logs: "Logs",
  metrics: "Metrics",
  traces: "Traces",
};
export const DEFAULT_TELEMETRY_TYPE = "logs";

export const TELEMETRY_SIZE_METRICS: { [telemetryType: string]: string } = {
  "logs": "log_data_size",
  "metrics": "metric_data_size",
  "traces": "trace_data_size",
}

interface MeasurementControlBarProps {
  telemetry: string;
  onTelemetryTypeChange: (telemetry: string) => void;
  period: string;
  onPeriodChange: (period: string) => void;
}

export const MeasurementControlBar: React.FC<MeasurementControlBarProps> = ({
  onTelemetryTypeChange,
  onPeriodChange,
  telemetry,
  period,
}) => {
  return (
    <AppBar classes={{ root: styles.appbar }} color="primary" position="static">
      <Toolbar>
        {Object.entries(TELEMETRY_TYPES).map(([t, label]) => (
          <Button
            key={t}
            color="inherit"
            onClick={() => onTelemetryTypeChange(t)}
          >
            <Typography style={{ fontWeight: telemetry === t ? 700 : 300 }}>
              {label}
            </Typography>
          </Button>
        ))}
        <Box flex={1} />
        {Object.entries(PERIODS).map(([r, label]) => (
          <Button key={r} color="inherit" onClick={() => onPeriodChange(r)}>
            <Typography style={{ fontWeight: period === r ? 700 : 300 }}>
              {label}
            </Typography>
          </Button>
        ))}
      </Toolbar>
    </AppBar>
  );
};
