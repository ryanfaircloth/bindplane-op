import { Chip, Stack } from "@mui/material";
import {
  DataGrid,
  DataGridProps,
  GridCellParams,
  GridColumns,
  GridDensityTypes,
  GridSelectionModel,
  GridValueFormatterParams,
  GridValueGetterParams,
} from "@mui/x-data-grid";
import { isFunction } from "lodash";
import React, { memo } from "react";
import {
  ConfigurationTableMetricsSubscription,
  GetConfigurationTableQuery,
} from "../../../graphql/generated";
import { formatMetric } from "../../../utils/graph/utils";
import { SearchLink } from "../../../utils/state";
import { NoMaxWidthTooltip } from "../../Custom/NoMaxWidthTooltip";
import {
  DEFAULT_CONFIGURATION_TABLE_PERIOD,
  TELEMETRY_SIZE_METRICS,
  TELEMETRY_TYPES,
} from "../../MeasurementControlBar/MeasurementControlBar";

export enum ConfigurationsTableField {
  NAME = "name",
  LABELS = "labels",
  DESCRIPTION = "description",
  AGENT_COUNT = "agentCount",
  LOGS = "logs",
  METRICS = "metrics",
  TRACES = "traces",
}

type Configurations =
  GetConfigurationTableQuery["configurations"]["configurations"];
interface ConfigurationsDataGridProps
  extends Omit<DataGridProps, "columns" | "rows"> {
  setSelectionModel?: (configurationIds: GridSelectionModel) => void;
  density?: GridDensityTypes;
  loading: boolean;
  configurations: Configurations;
  configurationMetrics?: ConfigurationTableMetricsSubscription;
  columnFields?: ConfigurationsTableField[];
  minHeight?: string;
  selectionModel?: GridSelectionModel;
}

const ConfigurationsDataGridComponent: React.FC<ConfigurationsDataGridProps> =
  ({
    setSelectionModel,
    loading,
    configurations,
    configurationMetrics,
    columnFields,
    density = GridDensityTypes.Standard,
    minHeight,
    selectionModel,
    ...dataGridProps
  }) => {
    const columns: GridColumns = (columnFields || []).map((field) => {
      switch (field) {
        case ConfigurationsTableField.AGENT_COUNT:
          return {
            field: ConfigurationsTableField.AGENT_COUNT,
            width: 100,
            headerName: "Agents",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.agentCount,
            renderCell: renderAgentCountCell,
          };

        case ConfigurationsTableField.DESCRIPTION:
          return {
            field: ConfigurationsTableField.DESCRIPTION,
            flex: 1,
            headerName: "Description",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.metadata.description,
          };
        case ConfigurationsTableField.LABELS:
          return {
            field: ConfigurationsTableField.LABELS,
            width: 300,
            headerName: "Labels",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.metadata.labels,
            renderCell: renderLabels,
          };
        case ConfigurationsTableField.LOGS:
          return createMetricRateColumn(field, "logs", configurationMetrics);
        case ConfigurationsTableField.METRICS:
          return createMetricRateColumn(field, "metrics", configurationMetrics);
        case ConfigurationsTableField.TRACES:
          return createMetricRateColumn(field, "traces", configurationMetrics);
        default:
          return {
            field: ConfigurationsTableField.NAME,
            headerName: "Name",
            width: 300,
            valueGetter: (params: GridValueGetterParams) =>
              params.row.metadata.name,
            renderCell: renderNameDataCell,
          };
      }
    });

    return (
      <DataGrid
        {...dataGridProps}
        checkboxSelection={isFunction(setSelectionModel)}
        onSelectionModelChange={setSelectionModel}
        components={{
          NoRowsOverlay: () => (
            <Stack height="100%" alignItems="center" justifyContent="center">
              No Configurations
            </Stack>
          ),
        }}
        style={{ minHeight }}
        disableSelectionOnClick
        autoHeight
        getRowId={(row) => row.metadata.name}
        columns={columns}
        rows={configurations}
        selectionModel={selectionModel}
      />
    );
  };

function renderLabels(
  cellParams: GridCellParams<Record<string, string>>
): JSX.Element {
  return (
    <Stack direction="row" spacing={1}>
      {Object.entries(cellParams.value || {}).map(([k, v]) => {
        const formattedLabel = `${k}: ${v}`;
        return <Chip key={k} size="small" label={formattedLabel} />;
      })}
    </Stack>
  );
}

function abreviateName(limit: number, name?: string): string {
  if (!name) return "";
  return name.length > limit ? name.substring(0, limit) + "..." : name;
}

function renderNameDataCell(cellParams: GridCellParams<string>): JSX.Element {
  return (
    <NoMaxWidthTooltip
      title={`${cellParams.value} (Click to view configuration)`}
      enterDelay={1000}
      placement="top-start"
    >
      <div>
        <SearchLink
          path={`/configurations/${cellParams.value || ""}`}
          displayName={abreviateName(40, cellParams.value)}
        />
      </div>
    </NoMaxWidthTooltip>
  );
}

function renderAgentCountCell(
  cellParams: GridCellParams<any, Configurations[0]>
) {
  return <span style={{ margin: "auto" }}>{cellParams.value}</span>;
}

function createMetricRateColumn(
  field: string,
  telemetryType: string,
  configurationMetrics?: ConfigurationTableMetricsSubscription
): GridColumns[0] {
  return {
    field,
    width: 100,
    headerName: TELEMETRY_TYPES[telemetryType],
    valueGetter: (params: GridValueGetterParams) => {
      if (configurationMetrics == null) {
        return "";
      }
      // should probably have a lookup table here rather than interpolate in two places
      const metricName = TELEMETRY_SIZE_METRICS[telemetryType];
      const configurationName = params.row.metadata.name;
      const metric = configurationMetrics.overviewMetrics.metrics.find(
        (m) =>
          m.name === metricName &&
          m.nodeID === `configuration/${configurationName}`
      );
      if (metric == null) {
        return 0;
      }
      // to make this sortable, we use the raw value and provide a valueFormatter implementation to show units
      return metric.value;
    },
    valueFormatter: (params: GridValueFormatterParams<number>): string => {
      if (params.value === 0) {
        return "";
      }
      return formatMetric(
        { value: params.value, unit: "B/s" },
        DEFAULT_CONFIGURATION_TABLE_PERIOD
      );
    },
  };
}

ConfigurationsDataGridComponent.defaultProps = {
  minHeight: "calc(100vh - 300px)",
  density: undefined,
  columnFields: [
    ConfigurationsTableField.NAME,
    ConfigurationsTableField.LABELS,
    ConfigurationsTableField.AGENT_COUNT,
    ConfigurationsTableField.LOGS,
    ConfigurationsTableField.METRICS,
    ConfigurationsTableField.TRACES,
    ConfigurationsTableField.DESCRIPTION,
  ],
};

export const ConfigurationsDataGrid = memo(ConfigurationsDataGridComponent);
