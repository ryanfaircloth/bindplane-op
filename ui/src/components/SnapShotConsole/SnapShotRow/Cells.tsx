import { TableCell, Typography } from "@mui/material";
import { CSSProperties } from "react";

export const CellLabel: React.FC = ({ children }) => {
  return (
    <TableCell width={200}>
      <Typography component="span" fontSize={12} fontFamily="monospace">
        {children}
      </Typography>
    </TableCell>
  );
};

interface CellValueProps {
  whiteSpace?: CSSProperties["whiteSpace"];
}

export const CellValue: React.FC<CellValueProps> = ({ children, whiteSpace }) => {
  return (
    <TableCell>
      <Typography component="span" fontSize={12} fontFamily="monospace" whiteSpace={whiteSpace}>
        {children}
      </Typography>
    </TableCell>
  );
};
