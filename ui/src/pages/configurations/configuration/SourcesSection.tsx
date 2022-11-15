import { Button, Stack, Typography } from "@mui/material";
import { memo, useMemo } from "react";
import { CardContainer } from "../../../components/CardContainer";
import { PlusCircleIcon } from "../../../components/Icons";
import {
  DestinationType,
  SourceType,
  useSourceTypesQuery,
} from "../../../graphql/generated";
import { NewResourceDialog } from "../../../components/ResourceDialog";
import { useSnackbar } from "notistack";
import { ShowPageConfig } from ".";
import {
  BPConfiguration,
  BPResourceConfiguration,
} from "../../../utils/classes";
import { UpdateStatus } from "../../../types/resources";
import { InlineSourceCard } from "../../../components/Cards/InlineSourceCard";

import styles from "./configuration-page.module.scss";
import mixins from "../../../styles/mixins.module.scss";

type ResourceType = SourceType | DestinationType;

const SourcesSectionComponent: React.FC<{
  configuration: NonNullable<ShowPageConfig>;
  refetch: () => {};
  setAddDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
  addDialogOpen: boolean;
}> = ({ configuration, refetch, setAddDialogOpen, addDialogOpen }) => {
  const sources = useMemo(
    () => configuration.spec?.sources || [],
    [configuration.spec?.sources]
  );
  const { enqueueSnackbar } = useSnackbar();
  const { data } = useSourceTypesQuery();

  async function onNewSourceSave(
    values: { [key: string]: any },
    sourceType: ResourceType
  ) {
    const newSourceConfig = new BPResourceConfiguration();
    newSourceConfig.type = sourceType.metadata.name;
    newSourceConfig.setParamsFromMap(values);

    const updatedConfig = new BPConfiguration(configuration);
    updatedConfig.addSource(newSourceConfig);
    try {
      const update = await updatedConfig.apply();
      if (update.status === UpdateStatus.INVALID) {
        console.error(update);
        throw new Error("failed to add source to configuration.");
      }

      setAddDialogOpen(false);
      refetch();
    } catch (err) {
      enqueueSnackbar("Failed to save source.", {
        variant: "error",
      });
      console.error(err);
    }
  }

  return (
    <>
      <CardContainer>
        <div className={styles["title-button-row"]}>
          <Typography variant="h5">Sources</Typography>
          <Button
            onClick={() => setAddDialogOpen(true)}
            variant="contained"
            classes={{ root: mixins["float-right"] }}
            startIcon={<PlusCircleIcon />}
          >
            Add Source
          </Button>
        </div>

        <Stack direction="row" spacing={2}>
          {sources.map((source, ix) => {
            return (
              <InlineSourceCard
                key={`source${ix}`}
                id={`source${ix}`}
                configuration={configuration}
                refetchConfiguration={refetch}
              />
            );
          })}
        </Stack>
      </CardContainer>

      <NewResourceDialog
        kind="source"
        resourceTypes={data?.sourceTypes ?? []}
        open={addDialogOpen}
        onSaveNew={onNewSourceSave}
        onClose={() => setAddDialogOpen(false)}
      />
    </>
  );
};

export const SourcesSection = memo(SourcesSectionComponent);
