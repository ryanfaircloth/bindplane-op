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
import { DestinationsQuery } from "../../../graphql/generated";
import { DestinationTypeCell } from "./cells";

import styles from "./cells.module.scss";

export enum DestinationsTableField {
  NAME = "name",
  KIND = "kind",
  TYPE = "type",
}

interface DestinationsDataGridProps
  extends Omit<DataGridProps, "columns" | "rows"> {
  onDestinationsSelected?: (names: GridSelectionModel) => void;
  onEditDestination: (name: string) => void;
  queryData: DestinationsQuery;
  loading: boolean;
}

export const DestinationsDataGrid: React.FC<DestinationsDataGridProps> = ({
  onDestinationsSelected,
  queryData,
  onEditDestination,
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

  const columns: GridColumns = [
    {
      field: DestinationsTableField.NAME,
      flex: 1,
      headerName: "Name",
      valueGetter: (params: GridValueGetterParams) => params.row.metadata.name,
      renderCell: renderNameCell,
    },
    // TODO removing this column breaks the tests in an inscrutable way.
    {
      field: DestinationsTableField.KIND,
      flex: 1,
      headerName: "Kind",
      valueGetter: (params: GridValueGetterParams) => params.row.kind,
      renderCell: renderStringCell,
    },
    {
      field: DestinationsTableField.TYPE,
      flex: 1,
      headerName: "Type",
      valueGetter: (params: GridValueGetterParams) => params.row.spec.type,
      renderCell: renderTypeCell,
    },
  ];

  function handleSelect(s: GridSelectionModel) {
    isFunction(onDestinationsSelected) && onDestinationsSelected(s);
  }

  const rows = [...queryData.destinations];

  return (
    <DataGrid
      disableSelectionOnClick
      {...dataGridProps}
      onSelectionModelChange={handleSelect}
      components={{
        NoRowsOverlay: () => (
          <Stack height="100%" alignItems="center" justifyContent="center">
            No Components
          </Stack>
        ),
      }}
      autoHeight
      getRowId={(row) => `${row.kind}|${row.metadata.name}`}
      columns={columns}
      rows={rows}
    />
  );
};

function renderTypeCell(cellParams: GridCellParams<string>): JSX.Element {
  return <DestinationTypeCell type={cellParams.value ?? ""} />;
}

function renderStringCell(cellParams: GridCellParams<string>): JSX.Element {
  return <>{cellParams.value}</>;
}
