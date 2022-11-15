import { gql } from "@apollo/client";
import {
  Card,
  CardActionArea,
  CardContent,
  Stack,
  Typography,
} from "@mui/material";
import { isNumber } from "lodash";
import { useSnackbar } from "notistack";
import { useState } from "react";
import { ConfirmDeleteResourceDialog } from "../ConfirmDeleteResourceDialog";
import { EditResourceDialog } from "../ResourceDialog/EditResourceDialog";
import { useSourceTypeQuery } from "../../graphql/generated";
import { UpdateStatus } from "../../types/resources";
import { BPConfiguration, BPResourceConfiguration } from "../../utils/classes";
import { classes } from "../../utils/styles";
import { MinimumRequiredConfig } from "../PipelineGraph/PipelineGraph";

import styles from "./cards.module.scss";

gql`
  query SourceType($name: String!) {
    sourceType(name: $name) {
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

export const InlineSourceCard: React.FC<{
  // For inline sources this is expected to be in the form source0, source1, etc
  id: string;
  disabled?: boolean;
  configuration: MinimumRequiredConfig;
  refetchConfiguration: () => void;
}> = ({ id, disabled, configuration, refetchConfiguration }) => {
  const sourceIndex = getSourceIndex(id);

  const source = configuration?.spec?.sources![sourceIndex];
  const name = source?.type || "";

  const { data } = useSourceTypeQuery({
    variables: { name },
  });

  const [editing, setEditing] = useState(false);
  const [confirmDeleteOpen, setDeleteOpen] = useState(false);

  const { enqueueSnackbar } = useSnackbar();

  const icon = data?.sourceType?.metadata.icon;
  const displayName = data?.sourceType?.metadata.displayName ?? "";
  const fontSize = displayName.length > 16 ? 14 : undefined;

  function closeEditDialog() {
    setEditing(false);
  }

  function closeDeleteDialog() {
    setDeleteOpen(false);
  }

  if (data?.sourceType == null) {
    return null;
  }

  async function onSave(values: { [key: string]: any }) {
    const sourceConfig = new BPResourceConfiguration(source);
    sourceConfig.setParamsFromMap(values);

    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.replaceSource(sourceConfig, sourceIndex);

    try {
      const update = await updatedConfig.apply();
      if (update.status === UpdateStatus.INVALID) {
        console.error(update);
        throw new Error("failed to save source on configuration");
      }

      enqueueSnackbar("Successfully saved source!", {
        variant: "success",
        autoHideDuration: 3000,
      });
      closeEditDialog();
      refetchConfiguration();
    } catch (err) {
      enqueueSnackbar("Failed to save source.", {
        variant: "error",
        autoHideDuration: 5000,
      });
      console.error(err);
    }
  }

  /**
   * Toggle `disabled` on the source, replace it in the configuration, and save
   */
  async function onTogglePause() {
    const updatedConfig = new BPConfiguration(configuration);
    const updatedSource = new BPResourceConfiguration(source);
    updatedSource.disabled = !source?.disabled;
    updatedConfig.replaceSource(updatedSource, sourceIndex);

    const action = updatedSource.disabled ? "pause" : "resume";
    try {
      const { status, reason } = await updatedConfig.apply();
      if (status === UpdateStatus.INVALID) {
        throw new Error(
          `failed to ${action} source, configuration invalid, ${reason}`
        );
      }

      enqueueSnackbar(`Successfully ${action}d source!`, {
        variant: "success",
        autoHideDuration: 3000,
      });
      closeEditDialog();
      refetchConfiguration();
    } catch (err) {
      enqueueSnackbar(`Failed to ${action} configuration.`, {
        variant: "error",
      });
      console.error(err);
    }
  }

  async function onDelete() {
    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.removeSource(sourceIndex);

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

  return (
    <>
      <Card
        className={classes([
          styles["resource-card"],
          disabled ? styles.disabled : undefined,
          source?.disabled ? styles.paused : undefined,
        ])}
        onClick={() => setEditing(true)}
      >
        <CardActionArea className={styles.action}>
          <CardContent data-testid={`source-card-${id}`}>
            <Stack alignItems="center" textAlign={"center"} height="100%">
              <span
                className={styles.icon}
                style={{ backgroundImage: `url(${icon})` }}
              />
              <Typography
                component="div"
                fontWeight={600}
                fontSize={fontSize}
                gutterBottom
              >
                {displayName}
              </Typography>
              {/* TODO: horizontal line separating displayName from Paused */}
              {source?.disabled && (
                <Typography
                  component="div"
                  fontWeight={400}
                  fontSize={14}
                  variant="overline"
                >
                  Paused
                </Typography>
              )}
            </Stack>
          </CardContent>
        </CardActionArea>
      </Card>

      <EditResourceDialog
        displayName={displayName}
        description={data?.sourceType?.metadata.description ?? ""}
        kind="source"
        enableProcessors
        processors={source?.processors}
        parameters={source?.parameters ?? []}
        parameterDefinitions={data.sourceType.spec.parameters}
        open={editing}
        onClose={closeEditDialog}
        onCancel={closeEditDialog}
        onDelete={() => setDeleteOpen(true)}
        onSave={onSave}
        paused={source?.disabled}
        onTogglePause={onTogglePause}
      />

      <ConfirmDeleteResourceDialog
        open={confirmDeleteOpen}
        onClose={closeDeleteDialog}
        onCancel={closeDeleteDialog}
        onDelete={onDelete}
        action={"remove"}
      >
        <Typography>Are you sure you want to remove this source?</Typography>
      </ConfirmDeleteResourceDialog>
    </>
  );
};

const REGEX = /^source(?<sourceNum>[0-9]+)$/;
export function getSourceIndex(id: string): number {
  const match = id.match(REGEX);
  if (match?.groups != null) {
    const index = Number(match.groups["sourceNum"]);
    if (isNumber(index)) {
      return index;
    }
  }
  return -1;
}
