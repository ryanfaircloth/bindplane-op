import { ApolloError, gql, NetworkStatus } from "@apollo/client";
import {
  Alert,
  CircularProgress,
  IconButton,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from "@mui/material";
import { memo, useEffect, useMemo, useState } from "react";
import {
  PipelineType,
  SnapshotQuery,
  useSnapshotQuery,
} from "../../graphql/generated";
import { SnapshotRow } from "./SnapShotRow";
import { RefreshIcon, XIcon } from "../Icons";
import { isEmpty } from "lodash";

import styles from "./snap-shot-console.module.scss";

interface Props {
  agentID: string;
  initialType: PipelineType;
  onClose: () => void;
}

type Metric = SnapshotQuery["snapshot"]["metrics"][0];
type Log = SnapshotQuery["snapshot"]["logs"][0];
type Trace = SnapshotQuery["snapshot"]["traces"][0];

// while the query includes all three pipeline types, only the pipelineType specified will have results
gql`
  query snapshot($agentID: String!, $pipelineType: PipelineType!) {
    snapshot(agentID: $agentID, pipelineType: $pipelineType) {
      metrics {
        name
        timestamp
        value
        unit
        type
        attributes
        resource
      }
      logs {
        timestamp
        body
        severity
        attributes
        resource
      }
      traces {
        name
        traceID
        spanID
        parentSpanID
        start
        end
        attributes
        resource
      }
    }
  }
`;

export const SnapshotConsole: React.FC<Props> = ({
  agentID,
  initialType,
  onClose,
}) => {
  const [logs, setLogs] = useState<Log[] | null>(null);
  const [metrics, setMetrics] = useState<Metric[] | null>(null);
  const [traces, setTraces] = useState<Trace[] | null>(null);
  const [type, setType] = useState<PipelineType>(initialType);
  const [error, setError] = useState<ApolloError>();

  const { loading, refetch, networkStatus } = useSnapshotQuery({
    variables: { agentID, pipelineType: initialType },
    onCompleted: (data) => {
      if (!isEmpty(data.snapshot.logs)) {
        setLogs(data.snapshot.logs.slice().reverse());
      }

      if (!isEmpty(data.snapshot.metrics)) {
        setMetrics(data.snapshot.metrics.slice().reverse());
      }

      if (!isEmpty(data.snapshot.traces)) {
        setTraces(data.snapshot.traces.slice().reverse());
      }
    },
    onError: (error) => {
      setError(error);
    },
    fetchPolicy: "network-only",
    notifyOnNetworkStatusChange: true,
  });

  useEffect(() => {
    switch (type) {
      case PipelineType.Logs:
        if (logs == null) {
          refetch({ agentID, pipelineType: PipelineType.Logs });
        }
        break;
      case PipelineType.Metrics:
        if (metrics == null) {
          refetch({ agentID, pipelineType: PipelineType.Metrics });
        }
        break;
      case PipelineType.Traces:
        if (traces == null) {
          refetch({ agentID, pipelineType: PipelineType.Traces });
        }
        break;
    }
  }, [type, refetch, logs, metrics, traces, agentID]);

  const busy = useMemo(
    () => loading || networkStatus === NetworkStatus.refetch,
    [loading, networkStatus]
  );

  return (
    <>
      <Stack
        direction="row"
        sx={{ width: "100%" }}
        alignItems={"center"}
        justifyContent="space-between"
        marginBottom={2}
      >
        <div /> {/** Used as a spacer */}
        <ToggleButtonGroup
          size={"small"}
          color="primary"
          value={type}
          exclusive
          onChange={(_, value) => {
            if (value != null) {
              setType(value);
            }
          }}
          aria-label="Telemetry Type"
        >
          <ToggleButton sx={{ width: 200 }} value={PipelineType.Logs}>
            Logs
          </ToggleButton>
          <ToggleButton sx={{ width: 200 }} value={PipelineType.Metrics}>
            Metrics
          </ToggleButton>
          <ToggleButton sx={{ width: 200 }} value={PipelineType.Traces}>
            Traces
          </ToggleButton>
        </ToggleButtonGroup>
        <div>
          <IconButton
            color={"primary"}
            disabled={busy}
            onClick={() => refetch({ agentID, pipelineType: type })}
          >
            <RefreshIcon />
          </IconButton>
          <IconButton onClick={onClose}>
            <XIcon />
          </IconButton>
        </div>
      </Stack>

      <MessagesContainer
        type={PipelineType.Logs}
        display={type === PipelineType.Logs}
        loading={busy}
        messages={logs}
      />

      <MessagesContainer
        type={PipelineType.Metrics}
        display={type === PipelineType.Metrics}
        loading={busy}
        messages={metrics}
      />

      <MessagesContainer
        type={PipelineType.Traces}
        display={type === PipelineType.Traces}
        loading={busy}
        messages={traces}
      />

      <Typography marginY={1} color="secondary">
        Showing recent {type}
      </Typography>

      {error && (
        <Alert sx={{ marginTop: 2 }} color="error">
          {error.message}
        </Alert>
      )}
    </>
  );
};

const MessagesContainerComponent: React.FC<{
  messages: (Log | Metric | Trace)[] | null;
  type: PipelineType;
  display: boolean;
  loading?: boolean;
}> = ({ messages, type, display, loading }) => {
  return (
    <div style={{ display: display ? "inline" : "none" }}>
      <div className={styles.container}>
        <div className={styles.console}>
          <div className={styles.header}>
            <div className={styles.ch} />
            <div className={styles.dt}>Time</div>
            <div className={styles.lg}>Message</div>
          </div>

          <div className={styles.stack}>
            {loading ? (
              <Stack
                height="90%"
                width={"100%"}
                justifyContent="center"
                alignItems="center"
              >
                <CircularProgress disableShrink />
              </Stack>
            ) : (
              <>
                {messages?.map((m, ix) => (
                  <SnapshotRow key={`${type}-${ix}`} message={m} type={type} />
                ))}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

const MessagesContainer = memo(MessagesContainerComponent);
