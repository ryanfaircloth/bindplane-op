import { Card, Stack, Typography, TableRow } from "@mui/material";
import Chip from "@mui/material/Chip";
import { getTimestamp } from "../utils";
import { useState, useMemo } from "react";
import { Log, PipelineType } from "../../../graphql/generated";
import { CellLabel, CellValue } from "./Cells";
import { SummaryTable } from "./SummaryTable";
import { RowSummary } from "./RowSummary";
import { DetailsContainer } from "./DetailsContainer";
import { MapValueSummary } from "./MapValueSummary";

import styles from "../snap-shot-console.module.scss";

interface LogRecordRowProps {
  message: Log;
}

export const LogRecordRow: React.FC<LogRecordRowProps> = ({ message }) => {
  const [open, setOpen] = useState(false);
  const timestamp = useMemo(
    () => getTimestamp(message, PipelineType.Logs),
    [message]
  );

  const severity = useMemo(() => {
    switch (message.severity) {
      case "trace":
        return message.severity;
      case "debug":
        return "debug";
      case "info":
        return "info";
      case "warning":
        return "warning";
      case "error":
        return "error";
      case "fatal":
        return "fatal";
      default:
        return "default";
    }
  }, [message]);

  return (
    <Card classes={{ root: styles.card }}>
      <RowSummary
        onClose={() => setOpen((prev) => !prev)}
        open={open}
        timestamp={timestamp}
      >
        <Stack direction="row" spacing={2} alignItems="center" width={"100%"}>
          <Chip label={severity} size={"small"} color={severity} />

          <Typography
            fontFamily="monospace"
            fontSize={12}
            overflow={"hidden"}
            textOverflow="ellipsis"
          >
            {message.body}
          </Typography>
        </Stack>
      </RowSummary>

      <DetailsContainer open={open}>
        <Typography fontWeight={600}>Log</Typography>

        <SummaryTable>
          <TableRow>
            <CellLabel>timestamp</CellLabel>
            <CellValue>{timestamp}</CellValue>
          </TableRow>

          <TableRow>
            <CellLabel>body</CellLabel>
            <CellValue whiteSpace="pre-wrap">{message.body}</CellValue>
          </TableRow>

          <TableRow>
            <CellLabel>severity</CellLabel>
            <CellValue>{severity}</CellValue>
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
