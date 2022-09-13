export enum AgentStatus {
  DISCONNECTED = 0,
  CONNECTED = 1,
  ERROR = 2,
  COMPONENT_FAILED = 4,
  DELETED = 5,
  CONFIGURING = 6,
  UPGRADING = 7,
}

export enum AgentFeatures {
  AGENT_SUPPORTS_UPGRADE = 0x0001,
  AGENT_SUPPORTS_SNAPSHOTS = 0x0002,
}

export function hasAgentFeature(
  agent: { features: number },
  feature: AgentFeatures
): boolean {
  return (agent.features & feature) !== 0;
}
