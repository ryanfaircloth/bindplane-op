import { TableCell, Typography } from "@mui/material";

export const CellLabel: React.FC = ({ children }) => {
  return (
    <TableCell width={200}>
      <Typography component="span" fontSize={12} fontFamily="monospace">
        {children}
      </Typography>
    </TableCell>
  );
};

export const CellValue: React.FC = ({ children }) => {
  return (
    <TableCell>
      <Typography component="span" fontSize={12} fontFamily="monospace">
        {children}
      </Typography>
    </TableCell>
  );
};
