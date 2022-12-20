import { Card, Chip, Stack } from "@mui/material";
import { ProcessorIcon } from "../Icons";
import { usePipelineGraph } from "../PipelineGraph/PipelineGraphContext";

import styles from "./cards.module.scss";

interface ProcessorCardProps {
  processorCount: number;
  resourceType: "source" | "destination";
  resourceIndex: number;
}

export const ProcessorCard: React.FC<ProcessorCardProps> = ({
  processorCount,
  resourceType,
  resourceIndex,
}) => {
  const { editProcessors } = usePipelineGraph();

  return (
    <>
      <Card
        className={styles["processor-card"]}
        onClick={() => editProcessors(resourceType, resourceIndex)}
      >
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
      {processorCount > 0 && (
        <Chip
          classes={{
            root: styles["count-chip"],
            label: styles["count-chip-label"],
          }}
          size="small"
          label={processorCount}
        />
      )}
    </>
  );
};
