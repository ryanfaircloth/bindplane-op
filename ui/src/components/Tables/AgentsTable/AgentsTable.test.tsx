import { mergeAgents } from "./merge-agents";
import { AgentChangeType } from "../../../graphql/generated";
import { cloneDeep } from "@apollo/client/utilities";
import { AgentsTableAgent } from ".";
import { AgentChange } from "../../../contexts/AgentChanges";

const a1: AgentsTableAgent = {
  id: "1",
  name: "",
  status: 1,
  configurationResource: {
    metadata: {
      name: "",
    },
  },
};
const a2: AgentsTableAgent = {
  id: "2",
  name: "",
  status: 1,
  configurationResource: {
    metadata: {
      name: "",
    },
  },
};
const a3: AgentsTableAgent = {
  id: "3",
  name: "",
  status: 1,
  configurationResource: {
    metadata: {
      name: "",
    },
  },
};
const a4: AgentsTableAgent = {
  id: "4",
  name: "",
  status: 1,
  configurationResource: {
    metadata: {
      name: "",
    },
  },
};
const a5: AgentsTableAgent = {
  id: "5",
  name: "",
  status: 1,
  configurationResource: {
    metadata: {
      name: "",
    },
  },
};

describe("mergeAgents", () => {
  it("removes agents when it gets event type remove", () => {
    const updates: AgentChange[] = [
      {
        agent: a1,
        changeType: AgentChangeType.Remove,
      },
      {
        agent: a2,
        changeType: AgentChangeType.Remove,
      },
      {
        agent: a3,
        changeType: AgentChangeType.Remove,
      },
      {
        agent: a4,
        changeType: AgentChangeType.Remove,
      },
      {
        agent: a5,
        changeType: AgentChangeType.Remove,
      },
    ];

    const newAgents = mergeAgents([a1, a2, a3, a4, a5], updates);
    expect(newAgents).toEqual([]);
  });

  it("adds an agent with event type insert", () => {
    const current = [a1, a2, a3, a4];
    const updates = [{ agent: a5, changeType: AgentChangeType.Insert }];

    const merged = mergeAgents(current, updates);
    expect(merged).toEqual([a1, a2, a3, a4, a5]);
  });

  it("replaces an agent with event type update", () => {
    const a1Updated = cloneDeep(a1);
    a1Updated.status = 0;

    const current = [a1, a2];
    const updates = [
      {
        agent: a1Updated,
        changeType: AgentChangeType.Update,
      },
    ];

    const merged = mergeAgents(current, updates);
    expect(merged).toEqual([a1Updated, a2]);
  });

  it("will not re-insert an existing agent", () => {
    const current = [a1, a2];
    const updates = [
      {
        agent: a2,
        changeType: AgentChangeType.Insert,
      },
    ];

    const merged = mergeAgents(current, updates);
    expect(merged).toEqual([a1, a2]);
  });

  it("will insert an updated agent if not present", () => {
    const current = [a1, a2];
    const updates = [
      {
        agent: a3,
        changeType: AgentChangeType.Update,
      },
    ];

    const merged = mergeAgents(current, updates);
    expect(merged).toEqual([a1, a2, a3]);
  });
});
