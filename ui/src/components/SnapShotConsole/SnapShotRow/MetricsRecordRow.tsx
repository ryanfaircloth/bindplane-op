import { Card, Stack, Typography, Chip, TableRow } from "@mui/material";
import { useMemo, useState } from "react";
import { Metric, PipelineType } from "../../../graphql/generated";
import { CellLabel, CellValue } from "./Cells";
import { SummaryTable } from "./SummaryTable";
import { RowSummary } from "./RowSummary";
import { getTimestamp } from "../utils";
import { DetailsContainer } from "./DetailsContainer";
import { MapValueSummary } from "./MapValueSummary";

import styles from "../snap-shot-console.module.scss";

interface MetricsRecordRowProps {
  message: Metric;
}

export const MetricsRecordRow: React.FC<MetricsRecordRowProps> = ({
  message,
}) => {
  const [open, setOpen] = useState(false);
  const timestamp = useMemo(
    () => getTimestamp(message, PipelineType.Logs),
    [message]
  );

  return (
    <Card classes={{ root: styles.card }}>
      <RowSummary
        open={open}
        onClose={() => setOpen((prev) => !prev)}
        timestamp={timestamp}
      >
        <Stack direction="row" spacing={2} alignItems="center">
          <Chip label={message.name} size={"small"} />

          <Typography
            fontFamily="monospace"
            fontSize={12}
            overflow={"hidden"}
            textOverflow="ellipsis"
          >
            {message.value} {message.unit}
          </Typography>
        </Stack>
      </RowSummary>

      <DetailsContainer open={open}>
        <Typography fontWeight={600}>Metric</Typography>
        <SummaryTable>
          <TableRow>
            <CellLabel>timestamp</CellLabel>
            <CellValue>{timestamp}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>name</CellLabel>
            <CellValue>{message.name}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>value</CellLabel>
            <CellValue>{message.value}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>type</CellLabel>
            <CellValue>{message.type}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>unit</CellLabel>
            <CellValue>{message.unit}</CellValue>
          </TableRow>
        </SummaryTable>

        <Typography fontWeight={600} marginTop={2}>
          Attributes
        </Typography>
        <MapValueSummary
          value={message.attributes}
          emptyMessage="No attribute values"
        />

        <Typography fontWeight={600} marginTop={2}>
          Resource
        </Typography>
        <MapValueSummary
          value={message.resource}
          emptyMessage="No resource values"
        />
      </DetailsContainer>
    </Card>
  );
};
