import { memo } from "react";
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

type ResourceType = SourceType | DestinationType;

const AddSourcesComponent: React.FC<{
  configuration: NonNullable<ShowPageConfig>;
  refetch: () => {};
  setAddDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
  addDialogOpen: boolean;
}> = ({ configuration, refetch, setAddDialogOpen, addDialogOpen }) => {
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
    <NewResourceDialog
      kind="source"
      resourceTypes={data?.sourceTypes ?? []}
      open={addDialogOpen}
      onSaveNew={onNewSourceSave}
      onClose={() => setAddDialogOpen(false)}
    />
  );
};

export const AddSourcesSection = memo(AddSourcesComponent);
