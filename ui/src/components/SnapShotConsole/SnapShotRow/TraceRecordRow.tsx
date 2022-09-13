import { Card, Chip, TableRow, Typography } from "@mui/material";
import { parseISO } from "date-fns";
import { useMemo, useState } from "react";
import { PipelineType, Trace } from "../../../graphql/generated";
import { getTimestamp } from "../utils";
import { RowSummary } from "./RowSummary";
import { differenceInMilliseconds } from "date-fns";
import { SummaryTable } from "./SummaryTable";
import { CellLabel, CellValue } from "./Cells";
import { DetailsContainer } from "./DetailsContainer";
import { MapValueSummary } from "./MapValueSummary";

import styles from "../snap-shot-console.module.scss";

interface TraceRecordRowProps {
  message: Trace;
}

export const TraceRecordRow: React.FC<TraceRecordRowProps> = ({ message }) => {
  const [open, setOpen] = useState(false);
  const timestamp = useMemo(
    () => getTimestamp(message, PipelineType.Traces),
    [message]
  );

  const diff = useMemo(
    function calcSpan() {
      const [start, end] = [parseISO(message.start), parseISO(message.end)];

      return `${differenceInMilliseconds(end, start)} ms`;
    },
    [message]
  );

  return (
    <Card classes={{ root: styles.card }}>
      <RowSummary
        open={open}
        onClose={() => setOpen((prev) => !prev)}
        timestamp={timestamp}
      >
        <Chip size="small" label={message.name} />
        <Typography fontFamily="monospace" fontSize={12}>
          {diff}
        </Typography>
      </RowSummary>

      <DetailsContainer open={open}>
        <Typography fontWeight={600}>Span</Typography>
        <SummaryTable>
          <TableRow>
            <CellLabel>start</CellLabel>
            <CellValue>{message.start}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>end</CellLabel>
            <CellValue>{message.end}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>name</CellLabel>
            <CellValue>{message.name}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>id</CellLabel>
            <CellValue>{message.spanID}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>parent</CellLabel>
            <CellValue>{message.parentSpanID}</CellValue>
          </TableRow>
          <TableRow>
            <CellLabel>trace</CellLabel>
            <CellValue>{message.traceID}</CellValue>
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
