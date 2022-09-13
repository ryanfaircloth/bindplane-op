import { TableRow, Typography } from "@mui/material";
import { isEmpty, isNil } from "lodash";
import { CellLabel, CellValue } from "./Cells";
import { SummaryTable } from "./SummaryTable";

interface MapValueSummaryProps {
  value: Record<string, any> | undefined;
  // The message to display when value is empty or undefined
  emptyMessage: string;
}

export const MapValueSummary: React.FC<MapValueSummaryProps> = ({
  value,
  emptyMessage,
}) => {
  return isNil(value) || isEmpty(value) ? (
    <Typography
      marginLeft={2}
      marginTop={1}
      fontWeight={300}
      fontSize={14}
      component="span"
      color="gray"
    >
      {emptyMessage}
    </Typography>
  ) : (
    <SummaryTable>
      {Object.entries(value).map(([k, v], ix) => {
        return (
          <TableRow key={k}>
            <CellLabel>{k}</CellLabel>
            <CellValue>{String(v)}</CellValue>
          </TableRow>
        );
      })}
    </SummaryTable>
  );
};
