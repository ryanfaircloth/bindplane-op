import { AgentsTableAgent, AgentsTableConfiguration } from "..";

function createAgent(): AgentsTableAgent {
  const now = new Date();
  const connectedAt = new Date(now.getTime() - 5 * 60 * 1000);
  const id = makeId();
  const configurationResource: AgentsTableConfiguration = {
    metadata: {
      name: "configuration-name",
    },
  };
  return {
    id: id,
    name: `test-agent-${id}`,
    status: 1,
    __typename: "Agent",
    architecture: "amd64",
    connectedAt: connectedAt.toString(),
    hostName: "host-name",
    configurationResource,
    labels: {},
    macAddress: "",
    operatingSystem: "Ubuntu",
    home: "/path/to/home",
    disconnectedAt: null,
    type: "otel",
    version: "3.0.0",
    platform: "linux",
  };
}

export function generateAgents(length: number): AgentsTableAgent[] {
  const agents: AgentsTableAgent[] = [];
  for (let i = 0; i < length; i++) {
    agents.push(createAgent());
  }

  return agents;
}

function makeId() {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < 5; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
