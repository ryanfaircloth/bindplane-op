import { Table, TableBody } from "@mui/material";

/**
 * SummaryTable wraps the children in a Table and TableBody
 */
export const SummaryTable: React.FC = ({ children }) => {
  return (
    <Table size="small">
      <TableBody>{children}</TableBody>
    </Table>
  );
};
