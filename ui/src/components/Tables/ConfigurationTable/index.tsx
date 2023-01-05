import { gql } from "@apollo/client";
import { Button, FormControl, Stack, Typography } from "@mui/material";
import { GridSelectionModel } from "@mui/x-data-grid";
import { debounce } from "lodash";
import React, { useEffect, useMemo, useState } from "react";
import {
  ConfigurationChangesDocument,
  ConfigurationChangesSubscription,
  EventType,
  GetConfigurationTableQuery,
  useConfigurationTableMetricsSubscription,
  useGetConfigurationTableQuery,
} from "../../../graphql/generated";
import mixins from "../../../styles/mixins.module.scss";
import { SearchBar } from "../../SearchBar";
import {
  ConfigurationsDataGrid,
  ConfigurationsTableField,
} from "./ConfigurationsDataGrid";
import { DeleteDialog } from "./DeleteDialog";

gql`
  query GetConfigurationTable(
    $selector: String
    $query: String
    $onlyDeployedConfigurations: Boolean
  ) {
    configurations(
      selector: $selector
      query: $query
      onlyDeployedConfigurations: $onlyDeployedConfigurations
    ) {
      configurations {
        metadata {
          name
          labels
          description
        }
        agentCount
      }
      query
      suggestions {
        query
        label
      }
    }
  }

  subscription ConfigurationChanges($selector: String, $query: String) {
    configurationChanges(selector: $selector, query: $query) {
      configuration {
        metadata {
          name
          description
          labels
        }
        agentCount
      }
      eventType
    }
  }

  subscription ConfigurationTableMetrics($period: String!) {
    overviewMetrics(period: $period) {
      metrics {
        name
        nodeID
        pipelineType
        value
        unit
      }
    }
  }
`;

type TableConfig =
  GetConfigurationTableQuery["configurations"]["configurations"][0];

function mergeConfigs(
  currentConfigs: TableConfig[],
  configurationUpdates:
    | ConfigurationChangesSubscription["configurationChanges"]
    | undefined
): TableConfig[] {
  const newConfigs: TableConfig[] = [...currentConfigs];

  for (const update of configurationUpdates || []) {
    const config = update.configuration;
    const configIndex = currentConfigs.findIndex(
      (c) => c.metadata.name === config.metadata.name
    );
    if (update.eventType === EventType.Remove) {
      // remove the agent if it exists
      if (configIndex !== -1) {
        newConfigs.splice(configIndex, 0);
      }
    } else if (configIndex === -1) {
      newConfigs.push(config);
    } else {
      newConfigs[configIndex] = config;
    }
  }
  return newConfigs;
}

interface ConfigurationTableProps {
  selector?: string;
  initQuery?: string;
  columns?: ConfigurationsTableField[];
  setSelected: (selected: GridSelectionModel) => void;
  selected: GridSelectionModel;
  enableDelete?: boolean;
  minHeight?: string;
  onlyDeployedConfigurations?: boolean;
}

export const ConfigurationsTable: React.FC<ConfigurationTableProps> = ({
  initQuery = "",
  selector,
  setSelected,
  selected,
  columns,
  enableDelete = true,
  minHeight,
  onlyDeployedConfigurations = false,
  ...dataGridProps
}) => {
  const { data, loading, refetch, subscribeToMore } =
    useGetConfigurationTableQuery({
      variables: { selector, query: initQuery, onlyDeployedConfigurations },
      fetchPolicy: onlyDeployedConfigurations
        ? "cache-and-network"
        : "network-only",
      nextFetchPolicy: "cache-only",
    });

  const { data: configurationMetrics } =
    useConfigurationTableMetricsSubscription({
      variables: { period: "1m" },
    });

  // Used to control the delete confirmation modal.
  const [open, setOpen] = useState<boolean>(false);

  const [subQuery, setSubQuery] = useState<string>(initQuery);
  const debouncedRefetch = useMemo(() => debounce(refetch, 100), [refetch]);

  useEffect(() => {
    subscribeToMore({
      document: ConfigurationChangesDocument,
      variables: { query: subQuery, selector },
      updateQuery: (prev, { subscriptionData, variables }) => {
        if (
          subscriptionData == null ||
          variables?.query !== subQuery ||
          variables.selector !== selector
        ) {
          return prev;
        }
        const { data } = subscriptionData as unknown as {
          data: ConfigurationChangesSubscription;
        };

        return {
          configurations: {
            __typename: "Configurations",
            suggestions: prev.configurations.suggestions,
            query: prev.configurations.query,
            configurations: mergeConfigs(
              prev.configurations.configurations,
              data.configurationChanges
            ),
          },
        };
      },
    });
  }, [selector, subQuery, subscribeToMore]);

  function onQueryChange(query: string) {
    debouncedRefetch({ selector, query });
    setSubQuery(query);
  }

  function openModal() {
    setOpen(true);
  }

  function closeModal() {
    setOpen(false);
  }

  return (
    <>
      <div className={mixins.flex}>
        <Typography variant="h5" className={mixins["mb-5"]}>
          Configurations
        </Typography>
        {selected.length > 0 && enableDelete && (
          <FormControl classes={{ root: mixins["ml-5"] }}>
            <Button variant="contained" color="error" onClick={openModal}>
              Delete {selected.length} Configuration
              {selected.length > 1 && "s"}
            </Button>
          </FormControl>
        )}
      </div>

      <Stack spacing={1}>
        <SearchBar
          suggestions={data?.configurations.suggestions}
          onQueryChange={onQueryChange}
          suggestionQuery={data?.configurations.query}
          initialQuery={initQuery}
        />

        <ConfigurationsDataGrid
          {...dataGridProps}
          setSelectionModel={setSelected}
          loading={loading}
          configurations={data?.configurations.configurations ?? []}
          configurationMetrics={configurationMetrics}
          columnFields={columns}
          selectionModel={selected}
          minHeight={minHeight}
        />
      </Stack>
      <DeleteDialog
        onClose={closeModal}
        selected={selected}
        open={open}
        onDeleteSuccess={refetch}
      />
    </>
  );
};
