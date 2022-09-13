import { Stack, Typography } from "@mui/material";
import { formatLogDate } from "../utils";
import { Toggle } from "./Toggle";

import styles from "../snap-shot-console.module.scss";

interface RowSummaryProps {
  open: boolean;
  onClose: () => void;
  timestamp: string;
}

export const RowSummary: React.FC<RowSummaryProps> = ({
  open,
  onClose,
  timestamp,
  children,
}) => {
  return (
    <Stack direction="row" alignItems="center" width={"100%"}>
      <div className={styles.ch}>
        <Toggle open={open} onClick={onClose} />
      </div>

      <div className={styles.dt}>
        <Typography fontFamily="monospace" fontSize={12} marginLeft={0.2}>
          {formatLogDate(new Date(timestamp))}
        </Typography>
      </div>

      <div className={styles.summary}>
        <Stack direction="row" spacing={2} alignItems="center">
          {children}
        </Stack>
      </div>
    </Stack>
  );
};
