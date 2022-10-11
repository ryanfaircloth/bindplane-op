import { IconButton, Tooltip } from "@mui/material";
import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import { useGetSourceTypeDisplayInfoQuery } from "../../graphql/generated";
import { PlusCircleIcon } from "../Icons";
import { NewProcessorDialog } from "../ResourceDialog/NewProcessorDialog";
import { MinimumRequiredConfig } from "../PipelineGraph/PipelineGraph";

interface AddProcessorCardProps {
  configuration: MinimumRequiredConfig;
  sourceIndex: number;
  refetchConfiguration: () => void;
}

export const AddProcessorCard: React.FC<AddProcessorCardProps> = ({
  configuration,
  sourceIndex,
  refetchConfiguration,
}) => {
  const [open, setOpen] = useState(false);
  const sourceType = configuration.spec.sources![sourceIndex].type;

  const { enqueueSnackbar } = useSnackbar();

  const { data, error } = useGetSourceTypeDisplayInfoQuery({
    variables: { name: sourceType ?? "" },
  });

  useEffect(() => {
    if (error != null) {
      enqueueSnackbar("Failed to retrieve source type.", {
        variant: "error",
        key: "Failed to retrieve source type.",
      });
    }
  }, [error, enqueueSnackbar]);

  if (data == null) {
    return null;
  }

  function onSaveSuccess() {
    refetchConfiguration();
    setOpen(false);
  }

  return (
    <>
      <Tooltip title={"Add a processor"}>
        <IconButton color="primary">
          <PlusCircleIcon onClick={() => setOpen(true)} />
        </IconButton>
      </Tooltip>

      <NewProcessorDialog
        onClose={() => setOpen(false)}
        parentDisplayName={data.sourceType?.metadata.displayName ?? ""}
        open={open}
        onSaveSuccess={onSaveSuccess}
        configuration={configuration}
        sourceIndex={sourceIndex}
      />
    </>
  );
};
