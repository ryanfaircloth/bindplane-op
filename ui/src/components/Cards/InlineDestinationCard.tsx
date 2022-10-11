import { gql } from "@apollo/client";
import {
  Card,
  CardActionArea,
  CardContent,
  Stack,
  Typography,
} from "@mui/material";
import { useSnackbar } from "notistack";
import { useState } from "react";
import { ConfirmDeleteResourceDialog } from "../ConfirmDeleteResourceDialog";
import { EditResourceDialog } from "../ResourceDialog/EditResourceDialog";
import {
  useDestinationTypeQuery,
} from "../../graphql/generated";
import { UpdateStatus } from "../../types/resources";
import { BPConfiguration } from "../../utils/classes/configuration";
import { BPResourceConfiguration } from "../../utils/classes/resource-configuration";

import styles from "./cards.module.scss";
import { useConfigurationPage } from "../../pages/configurations/configuration/ConfigurationPageContext";

gql`
  query DestinationType($name: String!) {
    destinationType(name: $name) {
      metadata {
        displayName
        name
        icon
        displayName
        description
      }
      spec {
        parameters {
          label
          name
          description
          required
          type
          default
          documentation {
            text
            url
          }
          relevantIf {
            name
            operator
            value
          }
          advancedConfig
          validValues
          options {
            creatable
            trackUnchecked
            sectionHeader
            gridColumns
            multiline
            metricCategories {
              label
              column
              metrics {
                name
                description
                kpi
              }
            }
          }
        }
      }
    }
  }
`;

export const InlineDestinationCard: React.FC<{
  id: string;
}> = ({ id }) => {
  const { configuration, refetchConfiguration } = useConfigurationPage();
  const destinationIndex = getDestinatonIndex(id);
  const destination = configuration.spec?.destinations![destinationIndex];

  // We can count on the type existing for an inline Resource
  const { data } = useDestinationTypeQuery({
    variables: { name: destination.type! },
  });
  const { enqueueSnackbar } = useSnackbar();

  const [editing, setEditing] = useState(false);
  const [confirmDeleteOpen, setDeleteOpen] = useState(false);

  if (data?.destinationType == null) {
    return null;
  }

  const icon = data.destinationType.metadata.icon;
  const displayName =
    data.destinationType.metadata.displayName ??
    data.destinationType.metadata.name;
  const description = data.destinationType.metadata.description ?? "";

  async function onDelete() {
    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.removeDestination(destinationIndex);

    try {
      const { status, reason } = await updatedConfig.apply();
      if (status === UpdateStatus.INVALID) {
        throw new Error(
          `failed to update configuration, configuration invalid, ${reason}`
        );
      }

      closeDeleteDialog();
      closeEditDialog();
      refetchConfiguration();
    } catch (err) {
      enqueueSnackbar("Failed to update configuration.", { variant: "error" });
      console.error(err);
    }
  }

  async function onSave(formValues: Record<string, any>) {
    const resourceConfiguration = new BPResourceConfiguration();
    resourceConfiguration.setParamsFromMap(formValues);
    resourceConfiguration.type = data!.destinationType!.metadata.name;

    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.replaceDestination(resourceConfiguration, destinationIndex);

    try {
      const { status, reason } = await updatedConfig.apply();
      if (status === UpdateStatus.INVALID) {
        throw new Error(
          `failed to update configuration, configuration invalid, ${reason}`
        );
      }

      enqueueSnackbar("Saved updated configuration.", { variant: "success" });
      closeEditDialog();
      refetchConfiguration();
    } catch (err) {
      enqueueSnackbar("Failed to save destination.", {
        variant: "error",
      });
      console.error(err);
    }
  }

  function closeEditDialog() {
    setEditing(false);
  }

  function closeDeleteDialog() {
    setDeleteOpen(false);
  }

  return (
    <>
      <Card
        className={styles["resource-card"]}
        onClick={() => setEditing(true)}
      >
        <CardActionArea sx={{ root: { padding: 0 } }}>
          <CardContent>
            <Stack alignItems="center">
              <span
                className={styles.icon}
                style={{ backgroundImage: `url(${icon})` }}
              />
              <Typography component="div" fontWeight={600}>
                {displayName}
              </Typography>
            </Stack>
          </CardContent>
        </CardActionArea>
      </Card>

      <EditResourceDialog
        displayName={displayName}
        description={description}
        kind="destination"
        fullWidth
        maxWidth="sm"
        parameters={destination.parameters ?? []}
        parameterDefinitions={data.destinationType.spec.parameters}
        open={editing}
        onClose={closeEditDialog}
        onCancel={closeEditDialog}
        onDelete={() => setDeleteOpen(true)}
        onSave={onSave}
      />

      <ConfirmDeleteResourceDialog
        open={confirmDeleteOpen}
        onClose={closeDeleteDialog}
        onCancel={closeDeleteDialog}
        onDelete={onDelete}
        action={"remove"}
      >
        <Typography>
          Are you sure you want to remove this destination?
        </Typography>
      </ConfirmDeleteResourceDialog>
    </>
  );
};

function getDestinatonIndex(id: string): number {
  const numberStr = id.split("destination")[1];
  return Number(numberStr);
}
