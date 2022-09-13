import { memo } from "react";
import { Log, Metric, Trace, PipelineType } from "../../../graphql/generated";
import { LogRecordRow } from "./LogRecordRow";
import { MetricsRecordRow } from "./MetricsRecordRow";
import { TraceRecordRow } from "./TraceRecordRow";

const SnapShotRowComponent: React.FC<{
  message: Log | Metric | Trace;
  type: PipelineType;
}> = ({ type, message }) => {
  switch (type) {
    case PipelineType.Logs:
      return <LogRecordRow message={message! as Log} />;
    case PipelineType.Metrics:
      return <MetricsRecordRow message={message! as Metric} />;
    case PipelineType.Traces:
      return <TraceRecordRow message={message! as Trace} />;
  }
};

export const SnapshotRow = memo(SnapShotRowComponent);
