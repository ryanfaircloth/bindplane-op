import { AgentSelector, PipelineType } from "../graphql/generated";

export function selectorString(sel: AgentSelector | undefined | null): string {
  if (sel == null) {
    return "";
  }
  return Object.entries(sel.matchLabels)
    .map(([k, v]) => `${k}=${v}`)
    .join(",");
}

export enum PipelineTypeFlags {
  Logs = 0x1,
  Metrics = 0x2,
  Traces = 0x4,
  ALL = 0xf,
}

export function isPipelineTypeFlagSet(
  flags: PipelineTypeFlags,
  flag: PipelineTypeFlags
): boolean {
  return (flags & flag) !== 0;
}

export function getPipelineTypeFlags(telemetryType: string): PipelineTypeFlags {
  switch (telemetryType) {
    case PipelineType.Logs:
      return PipelineTypeFlags.Logs;
    case PipelineType.Metrics:
      return PipelineTypeFlags.Metrics;
    case PipelineType.Traces:
      return PipelineTypeFlags.Traces;
  }
  return 0;
}
