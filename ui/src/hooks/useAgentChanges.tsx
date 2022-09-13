import { gql } from "@apollo/client";
import { useContext } from "react";
import { AgentChange, AgentChangesContext } from "../contexts/AgentChanges";

gql`
  subscription AgentChanges($selector: String, $query: String) {
    agentChanges(selector: $selector, query: $query) {
      agent {
        id
        name
        architecture
        operatingSystem
        labels
        hostName
        platform
        version
        macAddress
        home
        type
        status
        connectedAt
        disconnectedAt
        configuration {
          Collector
        }
        configurationResource {
          metadata {
            name
          }
        }
      }
      changeType
    }
  }
`;

export function useAgentChangesContext(): AgentChange[] {
  const { agentChanges } = useContext(AgentChangesContext);
  return agentChanges;
}
