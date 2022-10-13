import { gql } from "@apollo/client";
import {
  Dialog,
  DialogContent,
  Grid,
  Stack,
  Typography,
  Alert,
  AlertTitle,
  Button,
  Tooltip,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { CardContainer } from "../../components/CardContainer";
import { ManageConfigForm } from "../../components/ManageConfigForm";
import { AgentTable } from "../../components/Tables/AgentTable";
import { useGetAgentAndConfigurationsQuery } from "../../graphql/generated";
import { useAgentChangesContext } from "../../hooks/useAgentChanges";
import { RawConfigWizard } from "../configurations/wizards/RawConfigWizard";
import { useSnackbar } from "notistack";
import { labelAgents } from "../../utils/rest/label-agents";
import { RawConfigFormValues } from "../../types/forms";
import {
  hasAgentFeature,
  AgentFeatures,
  AgentStatus,
} from "../../types/agents";
import { withRequireLogin } from "../../contexts/RequireLogin";
import { withNavBar } from "../../components/NavBar";
import { AgentChangesProvider } from "../../contexts/AgentChanges";
import { RecentTelemetryDialog } from "../../components/RecentTelemetryDialog/RecentTelemetryDialog";

import mixins from "../../styles/mixins.module.scss";

gql`
  query GetAgentAndConfigurations($agentId: ID!) {
    agent(id: $agentId) {
      id
      name
      architecture
      operatingSystem
      labels
      hostName
      platform
      version
      macAddress
      remoteAddress
      home
      status
      connectedAt
      disconnectedAt
      errorMessage
      configuration {
        Collector
      }
      configurationResource {
        metadata {
          name
        }
      }
      upgrade {
        status
        version
        error
      }
      upgradeAvailable
      features
    }
    configurations {
      configurations {
        metadata {
          name
          labels
        }
        spec {
          raw
        }
      }
    }
  }
`;

const AgentPageContent: React.FC = () => {
  const { id } = useParams();
  const snackbar = useSnackbar();
  const [importOpen, setImportOpen] = useState(false);
  const [recentTelemetryOpen, setRecentTelemetryOpen] = useState(false);

  // AgentChanges subscription to trigger a refetch.
  const agentChanges = useAgentChangesContext();

  const { data, refetch } = useGetAgentAndConfigurationsQuery({
    variables: { agentId: id ?? "" },
    fetchPolicy: "network-only",
  });

  const navigate = useNavigate();

  async function onImportSuccess(values: RawConfigFormValues) {
    if (data?.agent != null) {
      try {
        await labelAgents(
          [data.agent.id],
          { configuration: values.name },
          true
        );
      } catch (err) {
        snackbar.enqueueSnackbar("Failed to apply label to agent.", {
          variant: "error",
        });
      }
    }

    setImportOpen(false);
  }

  useEffect(() => {
    if (agentChanges.length > 0) {
      const thisAgent = agentChanges
        .map((c) => c.agent)
        .find((a) => a.id === id);
      if (thisAgent != null) {
        refetch();
      }
    }
  }, [agentChanges, id, refetch]);

  const currentConfig = useMemo(() => {
    if (data?.agent == null || data?.configurations == null) {
      return null;
    }

    const configName = data.agent.configurationResource?.metadata.name;
    if (configName == null) {
      return null;
    }

    return data.configurations.configurations.find(
      (c) => c.metadata.name === configName
    );
  }, [data?.agent, data?.configurations]);

  const viewTelemetryButton = useMemo(() => {
    if (currentConfig?.spec.raw !== "") {
      return null;
    }

    let disableReason: string | null = null;

    if (
      data?.agent == null ||
      data.agent?.status === AgentStatus.DISCONNECTED
    ) {
      disableReason = "Cannot view recent telemetry, agent is disconnected.";
    }

    if (
      disableReason == null &&
      !hasAgentFeature(data!.agent!, AgentFeatures.AGENT_SUPPORTS_SNAPSHOTS)
    ) {
      disableReason =
        "Upgrade Agent to v1.8.0 or later to view recent telemetry.";
    }

    if (disableReason == null && data?.agent?.configurationResource == null) {
      disableReason =
        "Cannot view recent telemetry for an agent with an unmanaged configuration.";
    }

    if (disableReason != null) {
      return (
        <Tooltip title={disableReason} disableInteractive>
          <div style={{ display: "inline-block" }}>
            <Button
              variant="contained"
              size="large"
              onClick={() => setRecentTelemetryOpen(true)}
              sx={{ marginTop: 3 }}
              disabled
            >
              View Recent Telemetry
            </Button>
          </div>
        </Tooltip>
      );
    } else {
      return (
        <Button
          variant="contained"
          size="large"
          onClick={() => setRecentTelemetryOpen(true)}
          sx={{ marginTop: 3 }}
        >
          View Recent Telemetry
        </Button>
      );
    }
  }, [currentConfig?.spec.raw, data]);

  // Here we use the distinction between graphql returning null vs undefined.
  // If the agent is null then this agent doesn't exist, redirect to agents.
  if (data?.agent === null) {
    navigate("/agents");
    return null;
  }

  // Data is loading, return null for now.
  if (data === undefined || data.agent == null) {
    return null;
  }

  return (
    <>
      <CardContainer>
        <Typography variant="h5" marginRight={3}>
          Agent - {data.agent.name}
        </Typography>

        <Grid container spacing={5}>
          <Grid item xs={12} lg={6}>
            <Typography variant="h6" classes={{ root: mixins["mb-2"] }}>
              Details
            </Typography>

            <AgentTable agent={data.agent} />

            {data.agent.errorMessage && (
              <Alert severity="error" classes={{ root: mixins["mt-3"] }}>
                <AlertTitle>Error</AlertTitle>
                {data.agent.errorMessage}
              </Alert>
            )}

            {viewTelemetryButton}
          </Grid>
          <Grid item xs={12} lg={6}>
            <ManageConfigForm
              agent={data.agent}
              configurations={data.configurations.configurations ?? []}
              onImport={() => setImportOpen(true)}
            />
          </Grid>
        </Grid>
      </CardContainer>

      {/** Raw Config wizard for importing an agents config */}
      <Dialog
        open={importOpen}
        onClose={() => setImportOpen(false)}
        PaperComponent={EmptyComponent}
        scroll={"body"}
      >
        <DialogContent>
          <Stack justifyContent="center" alignItems="center" height="100%">
            <RawConfigWizard
              onClose={() => setImportOpen(false)}
              initialValues={{
                name: data.agent.name,
                description: `Imported config from agent ${data.agent.name}.`,
                fileName: "",
                rawConfig: data.agent.configuration?.Collector ?? "",
                platform: configPlatformFromAgentPlatform(data.agent.platform),
              }}
              onSuccess={onImportSuccess}
              fromImport
            />
          </Stack>
        </DialogContent>
      </Dialog>

      {currentConfig?.spec.raw === "" && (
        <RecentTelemetryDialog
          open={recentTelemetryOpen}
          onClose={() => setRecentTelemetryOpen(false)}
          agentID={id!}
        />
      )}
    </>
  );
};

const EmptyComponent: React.FC = ({ children }) => {
  return <>{children}</>;
};

function configPlatformFromAgentPlatform(platform: string | null | undefined) {
  if (platform == null) return "linux";
  if (platform === "darwin") return "macos";
  return platform;
}

export const AgentPage = withRequireLogin(
  withNavBar(() => (
    <AgentChangesProvider>
      <AgentPageContent />
    </AgentChangesProvider>
  ))
);
