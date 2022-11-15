import { gql } from "@apollo/client";
import { IconButton, Stack, Typography } from "@mui/material";
import React, { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { CardContainer } from "../../../components/CardContainer";
import {
  GetConfigurationQuery,
  useGetConfigurationQuery,
} from "../../../graphql/generated";
import { GridDensityTypes } from "@mui/x-data-grid";
import { AgentsTable } from "../../../components/Tables/AgentsTable";
import { AgentsTableField } from "../../../components/Tables/AgentsTable/AgentsDataGrid";
import { PlusCircleIcon } from "../../../components/Icons";
import { selectorString } from "../../../types/configuration";
import { ApplyConfigDialog } from "./ApplyConfigDialog";
import { DetailsSection } from "./DetailsSection";
import { ConfigurationSection } from "./ConfigurationSection";
import { SourcesSection } from "./SourcesSection";
import { DestinationsSection } from "./DestinationsSection";
import { useSnackbar } from "notistack";
import { withRequireLogin } from "../../../contexts/RequireLogin";
import { withNavBar } from "../../../components/NavBar";
import { PipelineGraph } from "../../../components/PipelineGraph/PipelineGraph";
import { ConfigurationPageContextProvider } from "./ConfigurationPageContext";

import styles from "./configuration-page.module.scss";
import { RawOrTopologyControl } from "../../../components/PipelineGraph/RawOrTopologyControl";

gql`
  query GetConfiguration($name: String!) {
    configuration(name: $name) {
      metadata {
        id
        name
        description
        labels
      }
      spec {
        raw
        sources {
          type
          name
          parameters {
            name
            value
          }
          processors {
            type
            parameters {
              name
              value
            }
            disabled
          }
          disabled
        }
        destinations {
          type
          name
          parameters {
            name
            value
          }
          processors {
            type
            parameters {
              name
              value
            }
            disabled
          }
          disabled
        }
        selector {
          matchLabels
        }
      }
      graph {
        attributes
        sources {
          id
          type
          label
          attributes
        }
        intermediates {
          id
          type
          label
          attributes
        }
        targets {
          id
          type
          label
          attributes
        }
        edges {
          id
          source
          target
        }
      }
      rendered
    }
  }
`;

export type ShowPageConfig = GetConfigurationQuery["configuration"];

const ConfigPageContent: React.FC = () => {
  const { name } = useParams();

  // Get Configuration Data
  const { data, refetch } = useGetConfigurationQuery({
    variables: { name: name ?? "" },
    fetchPolicy: "cache-and-network",
  });

  function toast(msg: string, variant: "error" | "success") {
    enqueueSnackbar(msg, { variant: variant, autoHideDuration: 3000 });
  }

  const [showApplyDialog, setShowApply] = useState(false);
  const [addSourceDialogOpen, setAddSourceDialogOpen] = useState(false);
  const [addDestDialogOpen, setAddDestDialogOpen] = useState(false);
  const [rawOrTopology, setTopologyOrRaw] =
    useState<"topology" | "raw">("topology");

  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const isRaw = (data?.configuration?.spec?.raw?.length || 0) > 0;
  const hasPipeline = useMemo(() => {
    if (
      data?.configuration?.spec == null ||
      data?.configuration.spec.sources == null ||
      data?.configuration.spec.destinations == null
    ) {
      return false;
    }

    return (
      data.configuration.spec.sources.length > 0 &&
      data.configuration.spec.destinations.length > 0
    );
  }, [data?.configuration?.spec]);

  function openApplyDialog() {
    setShowApply(true);
  }

  function closeApplyDialog() {
    setShowApply(false);
  }

  function onApplySuccess() {
    toast("Saved configuration!", "success");
    closeApplyDialog();
  }

  if (data?.configuration === undefined) {
    return null;
  }

  if (data.configuration === null) {
    enqueueSnackbar(`No configuration with name ${name} found.`, {
      variant: "error",
    });
    navigate("/configurations");
    return null;
  }

  return (
    <ConfigurationPageContextProvider
      configuration={data.configuration!}
      setAddDestDialogOpen={setAddDestDialogOpen}
      setAddSourceDialogOpen={setAddSourceDialogOpen}
      refetchConfiguration={refetch}
    >
      <section>
        <DetailsSection
          configuration={data.configuration}
          refetch={refetch}
          onSaveDescriptionError={() =>
            toast("Failed to save description.", "error")
          }
          onSaveDescriptionSuccess={() =>
            toast("Saved description.", "success")
          }
        />
      </section>

      {isRaw && (
        <section>
          <ConfigurationSection
            configuration={data.configuration}
            refetch={refetch}
            onSaveSuccess={() => toast("Saved configuration!", "success")}
            onSaveError={() => toast("Failed to save configuration.", "error")}
          />
        </section>
      )}

      {!isRaw && hasPipeline && (
        <CardContainer>
          <Stack spacing={2}>
            <RawOrTopologyControl
              rawOrTopology={rawOrTopology}
              setTopologyOrRaw={setTopologyOrRaw}
            />
            <PipelineGraph
              configuration={data.configuration}
              refetchConfiguration={refetch}
              agent={""}
              rawOrTopology={rawOrTopology}
              yamlValue={data.configuration.rendered ?? ""}
            />
          </Stack>
        </CardContainer>
      )}

      {!isRaw && (
        <section>
          <SourcesSection
            configuration={data.configuration}
            refetch={refetch}
            setAddDialogOpen={setAddSourceDialogOpen}
            addDialogOpen={addSourceDialogOpen}
          />
        </section>
      )}

      {!isRaw && (
        <section>
          <DestinationsSection
            configuration={data.configuration}
            destinations={data.configuration.spec.destinations ?? []}
            refetch={refetch}
            setAddDialogOpen={setAddDestDialogOpen}
            addDialogOpen={addDestDialogOpen}
          />
        </section>
      )}

      <section>
        <CardContainer>
          <div className={styles["title-button-row"]}>
            <Typography variant="h5">Agents</Typography>
            <IconButton onClick={openApplyDialog} color="primary">
              <PlusCircleIcon />
            </IconButton>
          </div>

          <AgentsTable
            selector={selectorString(data.configuration.spec.selector)}
            columnFields={[
              AgentsTableField.NAME,
              AgentsTableField.STATUS,
              AgentsTableField.OPERATING_SYSTEM,
            ]}
            density={GridDensityTypes.Compact}
            minHeight="300px"
          />
        </CardContainer>
      </section>

      {showApplyDialog && (
        <ApplyConfigDialog
          configuration={data.configuration}
          maxWidth="lg"
          fullWidth
          open={showApplyDialog}
          onError={() => toast("Failed to apply configuration.", "error")}
          onSuccess={onApplySuccess}
          onClose={closeApplyDialog}
          onCancel={closeApplyDialog}
        />
      )}
    </ConfigurationPageContextProvider>
  );
};

export const ViewConfiguration = withRequireLogin(
  withNavBar(ConfigPageContent)
);
