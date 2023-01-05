import { Card, Stack } from "@mui/material";

import { Handle, Position } from "react-flow-renderer";
import { ProcessorIcon } from "../../Icons";

import styles from "../../Cards/cards.module.scss";

export function DummyProcessorNode() {
  return (
    <>
      <Handle type="target" position={Position.Left} />
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
      <Handle type="source" position={Position.Right} />
    </>
  );
}
