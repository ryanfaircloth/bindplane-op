import { Card, Chip, Stack } from "@mui/material";
import { ProcessorIcon } from "../Icons";

import styles from "./cards.module.scss";

interface ProcessorCardProps {
  processors: number;
}

export const ProcessorCard: React.FC<ProcessorCardProps> = ({ processors }) => {
  return (
    <>
      <Card className={styles["processor-card"]}>
        <Stack
          width="100%"
          height="100%"
          justifyContent="center"
          alignItems="center"
        >
          <ProcessorIcon
            fill={"#d9d9d9"}
            stroke={"#222222"}
            strokeWidth="2px"
          />
        </Stack>
      </Card>
      {processors > 0 && <Chip
        classes={{
          root: styles["count-chip"],
          label: styles["count-chip-label"],
        }}
        size="small"
        label={processors}
      />}
    </>
  );
};
