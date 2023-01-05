import { Stack } from "@mui/material";
import {
  DataGrid,
  DataGridProps,
  GridCellParams,
  GridColumns,
  GridSelectionModel,
  GridValueGetterParams,
} from "@mui/x-data-grid";
import { isFunction } from "lodash";
import { memo } from "react";
import {
  DestinationsQuery,
  DestinationsInConfigsQuery,
} from "../../../graphql/generated";
import { DestinationTypeCell } from "./cells";

import styles from "./cells.module.scss";

export enum DestinationsTableField {
  NAME = "name",
  KIND = "kind",
  TYPE = "type",
}

interface DestinationsDataGridProps
  extends Omit<DataGridProps, "columns" | "rows"> {
  setSelectionModel?: (names: GridSelectionModel) => void;
  onEditDestination: (name: string) => void;
  queryData: DestinationsQuery | DestinationsInConfigsQuery;
  loading: boolean;
  columnFields?: DestinationsTableField[];
  minHeight?: string;
  selectionModel?: GridSelectionModel;
}

export const DestinationsDataGrid: React.FC<DestinationsDataGridProps> = memo(
  ({
    setSelectionModel,
    queryData,
    onEditDestination,
    columnFields,
    minHeight,
    selectionModel,
    ...dataGridProps
  }) => {
    function renderNameCell(cellParams: GridCellParams<string>): JSX.Element {
      if (cellParams.row.kind === "Destination") {
        return (
          <button
            onClick={() => onEditDestination(cellParams.value!)}
            className={styles.link}
          >
            {cellParams.value}
          </button>
        );
      }

      return renderStringCell(cellParams);
    }

    const columns: GridColumns = (columnFields || []).map((field) => {
      switch (field) {
        case DestinationsTableField.NAME:
          return {
            field: DestinationsTableField.NAME,
            width: 300,

            headerName: "Name",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.metadata.name,
            renderCell: renderNameCell,
          };
        case DestinationsTableField.KIND:
          return {
            // TODO removing this column breaks the tests in an inscrutable way.

            field: DestinationsTableField.KIND,
            flex: 1,
            headerName: "Kind",
            valueGetter: (params: GridValueGetterParams) => params.row.kind,
            renderCell: renderStringCell,
          };
        case DestinationsTableField.TYPE:
          return {
            field: DestinationsTableField.TYPE,
            flex: 1,
            headerName: "Type",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.spec.type,
            renderCell: renderTypeCell,
          };
        default:
          return { field: DestinationsTableField.TYPE };
      }
    });

    const rows =
      (queryData as DestinationsQuery).destinations !== undefined
        ? [...(queryData as DestinationsQuery).destinations]
        : [...(queryData as DestinationsInConfigsQuery).destinationsInConfigs];

    return (
      <DataGrid
        {...dataGridProps}
        checkboxSelection={isFunction(setSelectionModel)}
        onSelectionModelChange={setSelectionModel}
        components={{
          NoRowsOverlay: () => (
            <Stack height="100%" alignItems="center" justifyContent="center">
              No Components
            </Stack>
          ),
        }}
        style={{ minHeight }}
        disableSelectionOnClick
        autoHeight
        getRowId={(row) => `${row.kind}|${row.metadata.name}`}
        columns={columns}
        rows={rows}
        selectionModel={selectionModel}
      />
    );
  }
);

function renderTypeCell(cellParams: GridCellParams<string>): JSX.Element {
  return <DestinationTypeCell type={cellParams.value ?? ""} />;
}

function renderStringCell(cellParams: GridCellParams<string>): JSX.Element {
  return <>{cellParams.value}</>;
}

DestinationsDataGrid.defaultProps = {
  minHeight: "calc(100vh - 250px)",
  columnFields: [
    DestinationsTableField.NAME,
    DestinationsTableField.KIND,
    DestinationsTableField.TYPE,
  ],
};
