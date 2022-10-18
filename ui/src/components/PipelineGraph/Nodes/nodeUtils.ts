import { PipelineType } from "../../../graphql/generated";
import {
  getPipelineTypeFlags,
  isPipelineTypeFlagSet,
  PipelineTypeFlags,
} from "../../../types/configuration";

// AttributeName specifies the names of attributes set on graph nodes
export enum AttributeName {
  ActiveTypeFlags = "activeTypeFlags",
  IsInline = "isInline",
  Kind = "kind",
  ResourceId = "resourceId",
}

export function isNodeDisabled(
  telemetryType: string,
  attributes: Record<string, any>
): boolean {
  const typeFlag = getPipelineTypeFlags(telemetryType);
  const activeFlags =
    attributes[AttributeName.ActiveTypeFlags] ?? PipelineTypeFlags.ALL;
  return !isPipelineTypeFlagSet(activeFlags, typeFlag);
}

export function firstActiveTelemetry(
  attributes: Record<string, any> | null
): PipelineType | null {
  if (attributes == null) return null;

  const activeFlags = attributes[AttributeName.ActiveTypeFlags];
  if (activeFlags & PipelineTypeFlags.Logs) {
    return PipelineType.Logs;
  }
  if (activeFlags & PipelineTypeFlags.Metrics) {
    return PipelineType.Metrics;
  }
  if (activeFlags && PipelineTypeFlags.Traces) {
    return PipelineType.Traces;
  }

  return null;
}
