import React, { createContext, useContext } from "react";
import { MinimumRequiredConfig } from "../../../components/PipelineGraph/PipelineGraph";

interface ConfigurationPageContextValue {
  configuration: MinimumRequiredConfig;
  refetchConfiguration: () => void;
  setAddSourceDialogOpen: (b: boolean) => void;
  setAddDestDialogOpen: (b: boolean) => void;
}

const defaultContextValue: ConfigurationPageContextValue = {
  configuration: {
    __typename: undefined,
    metadata: {
      __typename: undefined,
      id: "",
      name: "",
      description: undefined,
      labels: undefined,
    },
    spec: {
      __typename: undefined,
      raw: undefined,
      sources: undefined,
      destinations: undefined,
      selector: undefined,
    },
    graph: undefined,
  },
  refetchConfiguration: () => {},
  setAddSourceDialogOpen: () => {},
  setAddDestDialogOpen: () => {},
};

const ConfigurationPageContext = createContext(defaultContextValue);

export const ConfigurationPageContextProvider: React.FC<ConfigurationPageContextValue> =
  ({
    configuration,
    setAddDestDialogOpen,
    setAddSourceDialogOpen,
    refetchConfiguration,
    children,
  }) => {
    return (
      <ConfigurationPageContext.Provider
        value={{
          configuration,
          setAddDestDialogOpen,
          setAddSourceDialogOpen,
          refetchConfiguration,
        }}
      >
        {children}
      </ConfigurationPageContext.Provider>
    );
  };

export function useConfigurationPage(): ConfigurationPageContextValue {
  return useContext(ConfigurationPageContext);
}
