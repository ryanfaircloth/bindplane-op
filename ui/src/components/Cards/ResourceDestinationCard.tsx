import { gql } from "@apollo/client";
import {
  Card,
  CardActionArea,
  CardContent,
  Stack,
  Typography,
} from "@mui/material";
import { useSnackbar } from "notistack";
import { memo, useState } from "react";
import { ConfirmDeleteResourceDialog } from "../ConfirmDeleteResourceDialog";
import { EditResourceDialog } from "../ResourceDialog/EditResourceDialog";
import { useGetDestinationWithTypeQuery } from "../../graphql/generated";
import { useConfigurationPage } from "../../pages/configurations/configuration/ConfigurationPageContext";
import { UpdateStatus } from "../../types/resources";
import { BPConfiguration, BPDestination } from "../../utils/classes";
import { FormValues } from "../ResourceConfigForm";
import { classes } from "../../utils/styles";

import styles from "./cards.module.scss";

gql`
  query getDestinationWithType($name: String!) {
    destinationWithType(name: $name) {
      destination {
        metadata {
          name
          id
          labels
        }
        spec {
          type
          parameters {
            name
            value
          }
        }
      }
      destinationType {
        metadata {
          name
          icon
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
            relevantIf {
              name
              operator
              value
            }
            documentation {
              text
              url
            }
            advancedConfig
            validValues
            options {
              multiline
              creatable
              trackUnchecked
              sectionHeader
              gridColumns
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
  }
`;

const ResourceDestinationCardComponent: React.FC<{
  name: string;
  disabled?: boolean;
}> = ({ name, disabled }) => {
  const { data, refetch: refetchDestination } = useGetDestinationWithTypeQuery({
    variables: { name },
    fetchPolicy: "cache-and-network",
  });

  const { configuration, refetchConfiguration } = useConfigurationPage();
  const destinationIndex = configuration.spec.destinations?.findIndex(
    (d) => d.name === name
  );
  const processors =
    destinationIndex !== -1 && destinationIndex != null
      ? configuration.spec.destinations?.[destinationIndex].processors
      : [];

  const { enqueueSnackbar } = useSnackbar();

  const [editing, setEditing] = useState(false);
  const [confirmDeleteOpen, setDeleteOpen] = useState(false);

  function closeEditDialog() {
    setEditing(false);
  }

  async function onSave(formValues: FormValues) {
    const updatedDestination = new BPDestination(
      data!.destinationWithType!.destination!
    );

    updatedDestination.setParamsFromMap(formValues);

    const updatedProcessors = formValues.processors;
    if (updatedProcessors != null) {
      // assign processors to the configuration
      if (destinationIndex == null) {
        enqueueSnackbar("Failed to update destination.", { variant: "error" });
        console.error(
          `Could not find index for destination named ${name} in configuration spec.`,
          configuration
        );
        return;
      }
      const updatedConfig = new BPConfiguration(configuration);
      updatedConfig.replaceDestination(
        {
          name: updatedDestination.name(),
          processors: updatedProcessors,
        },
        destinationIndex
      );

      try {
        const update = await updatedConfig.apply();
        if (update.status === UpdateStatus.INVALID) {
          throw new Error(
            `failed to apply configuration, got status ${update.status}`
          );
        }
      } catch (err) {
        console.error(err);
        enqueueSnackbar("Failed to update configuration.", {
          variant: "error",
        });
      }
    }

    try {
      const update = await updatedDestination.apply();
      if (update.status === UpdateStatus.INVALID) {
        console.error("Update: ", update);
        throw new Error(
          `failed to apply destination, got status ${update.status}`
        );
      }

      enqueueSnackbar("Successfully saved destination.", {
        variant: "success",
      });
      setEditing(false);
      refetchConfiguration();
      refetchDestination();
    } catch (err) {
      console.error(err);
      enqueueSnackbar("Failed to update destination.", { variant: "error" });
    }
  }

  async function onDelete() {
    if (destinationIndex == null) {
      enqueueSnackbar("Failed to delete destination.", { variant: "error" });
      console.error(
        `Could not find index for destination named ${name} in configuration spec.`,
        configuration
      );
      return;
    }

    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.removeDestination(destinationIndex);

    try {
      const update = await updatedConfig.apply();
      if (update.status === UpdateStatus.INVALID) {
        console.error("Update: ", update);
        throw new Error(
          `failed to remove destination from configuration, configuration invalid`
        );
      }

      closeEditDialog();
      closeDeleteDialog();
      refetchConfiguration();
      refetchDestination();
    } catch (err) {
      enqueueSnackbar("Failed to remove destination.", { variant: "error" });
    }
  }

  function closeDeleteDialog() {
    setDeleteOpen(false);
  }

  // Loading
  if (data === undefined) {
    return null;
  }

  if (data.destinationWithType.destination == null) {
    enqueueSnackbar(`Could not retrieve destination ${name}.`, {
      variant: "error",
    });
    return null;
  }

  if (data.destinationWithType.destinationType == null) {
    enqueueSnackbar(
      `Could not retrieve destination type for destination ${name}.`,
      { variant: "error" }
    );
    return null;
  }

  return (
    <div className={disabled ? styles.disabled : undefined}>
      <Card
        className={classes([
          styles["resource-card"],
          disabled ? styles.disabled : undefined,
        ])}
        onClick={() => setEditing(true)}
      >
        <CardActionArea>
          <CardContent>
            <Stack alignItems="center">
              <span
                className={styles.icon}
                style={{
                  backgroundImage: `url(${data?.destinationWithType?.destinationType?.metadata.icon})`,
                }}
              />
              <Typography component="div" fontWeight={600}>
                {name}
              </Typography>
            </Stack>
          </CardContent>
        </CardActionArea>
      </Card>

      <EditResourceDialog
        kind="destination"
        displayName={name}
        description={
          data.destinationWithType.destinationType.metadata.description ?? ""
        }
        enableProcessors
        processors={processors}
        fullWidth
        maxWidth="sm"
        parameters={data.destinationWithType.destination.spec.parameters ?? []}
        parameterDefinitions={
          data.destinationWithType.destinationType.spec.parameters
        }
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
    </div>
  );
};

export const ResourceDestinationCard = memo(ResourceDestinationCardComponent);
