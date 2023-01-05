import { AppBar, Toolbar, Button, Typography, Box, Stack } from "@mui/material";
import { classes } from "../../utils/styles";

import styles from "./measurement-control-bar.module.scss";

export const PERIODS: { [period: string]: string } = {
  "10s": "10s",
  "1m": "1m",
  "5m": "5m",
  "1h": "1h",
  "24h": "24h",
};
export const DEFAULT_PERIOD = "10s";
export const DEFAULT_AGENTS_TABLE_PERIOD = "10s";
export const DEFAULT_CONFIGURATION_TABLE_PERIOD = "10s";
export const DEFAULT_OVERVIEW_GRAPH_PERIOD = "1h";
export const TELEMETRY_TYPES: { [telemetryType: string]: string } = {
  logs: "Logs",
  metrics: "Metrics",
  traces: "Traces",
};
export const DEFAULT_TELEMETRY_TYPE = "logs";

export const TELEMETRY_SIZE_METRICS: { [telemetryType: string]: string } = {
  logs: "log_data_size",
  metrics: "metric_data_size",
  traces: "trace_data_size",
};

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
      <Toolbar classes={{ root: styles.toolbar }}>
        {Object.entries(TELEMETRY_TYPES).map(([t, label]) => (
          <Stack key={t}>
            <Button
              variant="outlined"
              className={classes([
                styles["menu-button"],
                telemetry === t ? styles["menu-button-selected"] : undefined,
              ])}
              key={t}
              color="inherit"
              onClick={() => onTelemetryTypeChange(t)}
            >
              <Typography
                className={
                  telemetry === t
                    ? styles["selected-text"]
                    : styles["regular-text"]
                }
              >
                {label}
              </Typography>
            </Button>

            <Box
              className={classes([
                styles["regular-box"],
                telemetry === t ? styles["selected-box"] : undefined,
              ])}
            />
          </Stack>
        ))}
        <Box className={styles["flex-box"]} />

        {Object.entries(PERIODS).map(([p, label]) => (
          <Stack key={p}>
            <Button
              variant="outlined"
              className={classes([
                styles["menu-button"],
                period === p ? styles["menu-button-selected"] : undefined,
              ])}
              key={p}
              color="inherit"
              onClick={() => onPeriodChange(p)}
            >
              <Typography
                className={
                  period === p
                    ? styles["selected-text"]
                    : styles["regular-text"]
                }
              >
                {label}
              </Typography>
            </Button>
            <Box
              className={classes([
                styles["regular-box"],
                period === p ? styles["selected-box"] : undefined,
              ])}
            />
          </Stack>
        ))}
      </Toolbar>
    </AppBar>
  );
};
