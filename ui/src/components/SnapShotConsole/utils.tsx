import { Log, Metric, Trace, PipelineType } from "../../graphql/generated";

const logDateFormat = new Intl.DateTimeFormat(undefined, {
  month: "short",
  day: "2-digit",
  year: undefined,
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  timeZoneName: "short",
});

export function formatLogDate(date: Date): string {
  return logDateFormat.format(date);
}

export function getTimestamp(
  message: Log | Metric | Trace,
  type: PipelineType
): any {
  switch (type) {
    case PipelineType.Logs:
      return (message as Log).timestamp;
    case PipelineType.Metrics:
      return (message as Metric).timestamp;
    case PipelineType.Traces:
      return (message as Trace).end;
  }
}
