import { Box, Button, Card, Grid, Tooltip } from "@mui/material";
import { ReactFlowProvider } from "react-flow-renderer";
import { withNavBar } from "../../components/NavBar";

import { ConfigurationsTable } from "../../components/Tables/ConfigurationTable";
import { withRequireLogin } from "../../contexts/RequireLogin";
import {
  useConfigurationTableMetricsSubscription,
  useDestinationsInConfigsQuery,
} from "../../graphql/generated";
import { OverviewGraph } from "./OverviewGraph";
import { OverviewPageProvider, useOverviewPage } from "./OverviewPageContext";
import mixins from "../../styles/mixins.module.scss";
import {
  DEFAULT_OVERVIEW_GRAPH_PERIOD,
  DEFAULT_TELEMETRY_TYPE,
  MeasurementControlBar,
  TELEMETRY_SIZE_METRICS,
} from "../../components/MeasurementControlBar/MeasurementControlBar";
import global from "../../styles/global.module.scss";
import { gql } from "@apollo/client";

import { DestinationsTableField } from "../../components/Tables/DestinationsTable/DestinationsDataGrid";
import { ConfigurationsTableField } from "../../components/Tables/ConfigurationTable/ConfigurationsDataGrid";
import { DestinationsPageContent } from "../destinations/DestinationsPage";
import { useCallback, useEffect } from "react";

gql`
  query DestinationsInConfigs {
    destinationsInConfigs {
      kind
      metadata {
        name
      }
      spec {
        type
      }
    }
  }
`;

const OverviewPageContent: React.FC = () => {
  const {
    selectedTelemetry,
    selectedConfigs,
    setSelectedConfigs,
    selectedDestinations,
    setSelectedDestinations,
    setSelectedTelemetry,
    selectedPeriod,
    setPeriod,
    editingDestination,
    setEditingDestination,
    loadTop,
    setLoadTop,
  } = useOverviewPage();

  // we need these metrics to select the top three configs on load
  const { data: configurationMetrics } =
    useConfigurationTableMetricsSubscription({
      variables: { period: "1h" }, // TODO: selectedPeriod?
    });

  const selectTopResources = useCallback(
    (count: number, resourceType: "configuration" | "destination") => {
      const filteredMetrics =
        configurationMetrics?.overviewMetrics.metrics
          .filter((metric) => metric.nodeID.startsWith(resourceType))
          .filter(
            (metric) =>
              metric.name ===
              TELEMETRY_SIZE_METRICS[
                selectedTelemetry || DEFAULT_TELEMETRY_TYPE
              ]
          )
          .sort((a, b) => {
            return b.value - a.value;
          }) || [];

      const topResources = filteredMetrics.slice(0, count).map((metric) => {
        return metric.nodeID.split("/")[1];
      });
      return topResources;
    },
    [configurationMetrics, selectedTelemetry]
  );

  const selectTopConfigs = useCallback(
    (count: number) => {
      const topConfigs = selectTopResources(count, "configuration");
      if (topConfigs) {
        setSelectedConfigs(topConfigs);
      }
    },
    [setSelectedConfigs, selectTopResources]
  );

  const selectTopDestinations = useCallback(
    (count: number) => {
      const topDests = selectTopResources(count, "destination").map((name) => {
        // map from name to "Destination|name"
        return `Destination|${name}`;
      });
      if (topDests) {
        setSelectedDestinations(topDests);
      }
    },
    [setSelectedDestinations, selectTopResources]
  );

  useEffect(() => {
    // select top three configs on load
    if (loadTop && configurationMetrics && selectedTelemetry) {
      selectTopConfigs(3);
      selectTopDestinations(3);
      setLoadTop(false);
    }
  }, [
    configurationMetrics,
    loadTop,
    setLoadTop,
    selectTopConfigs,
    selectTopDestinations,
    selectedTelemetry,
  ]);

  return (
    <Grid container spacing={1} alignItems="center" wrap={"nowrap"}>
      <Grid item md={"auto"} lg={"auto"}>
        <Box
          sx={{
            width: "300px",
          }}
        >
          <Tooltip
            enterDelay={1000}
            title="Limit the displayed configurations to the three receiving the most data of the selected telemetry type over the selected period."
          >
            <Button
              variant="contained"
              classes={{ root: mixins["float-right"] }}
              onClick={() => selectTopConfigs(3)}
            >
              Top Three
            </Button>
          </Tooltip>

          <ConfigurationsTable
            selected={selectedConfigs}
            setSelected={setSelectedConfigs}
            enableDelete={false}
            minHeight="calc(100vh - 231px)"
            columns={[ConfigurationsTableField.NAME]}
            onlyDeployedConfigurations
          />
        </Box>
      </Grid>
      <Grid item md={true} lg={true}>
        <Card
          className={global["graph"]}
          style={{
            height: "calc(100vh - 120px)",
            width: "100%",
          }}
        >
          <MeasurementControlBar
            telemetry={selectedTelemetry || DEFAULT_TELEMETRY_TYPE}
            onTelemetryTypeChange={setSelectedTelemetry}
            period={selectedPeriod || DEFAULT_OVERVIEW_GRAPH_PERIOD}
            onPeriodChange={setPeriod}
          />
          <ReactFlowProvider>
            <OverviewGraph />
          </ReactFlowProvider>
        </Card>
      </Grid>

      <Grid item md={"auto"} lg={"auto"}>
        <Box
          sx={{
            width: "300px",
          }}
        >
          <Tooltip
            enterDelay={1000}
            title="Limit the displayed destinations to the three receiving the most data of the selected telemetry type over the selected period."
          >
            <Button
              variant="contained"
              classes={{ root: mixins["float-right"] }}
              onClick={() => selectTopDestinations(3)}
            >
              Top Three
            </Button>
          </Tooltip>

          <DestinationsPageContent
            selected={selectedDestinations}
            setSelected={setSelectedDestinations}
            enableDelete={false}
            destinationsQuery={useDestinationsInConfigsQuery}
            columnFields={[DestinationsTableField.NAME]}
            minHeight="calc(100vh - 181px)"
            editingDestination={editingDestination}
            setEditingDestination={setEditingDestination}
          />
        </Box>
      </Grid>
    </Grid>
  );
};

export const OverviewPage: React.FC = withRequireLogin(
  withNavBar((props) => {
    return (
      <OverviewPageProvider>
        <OverviewPageContent />
      </OverviewPageProvider>
    );
  })
);
