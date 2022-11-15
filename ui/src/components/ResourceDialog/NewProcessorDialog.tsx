import { Dialog, DialogProps } from "@mui/material";
import React, { memo, useEffect, useMemo, useState, useCallback } from "react";
import { ResourceDialogContextProvider } from "./ResourceDialogContext";
import { GetProcessorTypesQuery } from "../../graphql/generated";
import {
  CreateProcessorConfigureView,
  CreateProcessorSelectView,
  FormValues,
} from "../ResourceConfigForm";
import { MinimumRequiredConfig } from "../PipelineGraph/PipelineGraph";
import { UpdateStatus } from "../../types/resources";
import { BPConfiguration, BPResourceConfiguration } from "../../utils/classes";
import { useSnackbar } from "notistack";

export type ProcessorType = GetProcessorTypesQuery["processorTypes"][0];

interface NewProcessorDialogProps extends DialogProps {
  onClose: () => void;
  // The supported telemetry types of the parent resource type that is
  // being configured.  a subset of ['logs', 'metrics', 'traces']
  telemetryTypes?: string[];
  // The display name of the Source or Destination on which we're adding the processor.
  parentDisplayName: string;
  configuration: MinimumRequiredConfig;
  sourceIndex: number;
  onSaveSuccess: () => void;
}

export const NewProcessorDialogComponent: React.FC<NewProcessorDialogProps> = ({
  onClose,
  telemetryTypes,
  parentDisplayName,
  configuration,
  sourceIndex,
  onSaveSuccess,
  ...dialogProps
}) => {
  const [selected, setSelected] = useState<ProcessorType | null>(null);
  const { enqueueSnackbar } = useSnackbar();

  // This resets the form after close.
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout>;
    if (dialogProps.open === false) {
      timer = setTimeout(() => clearResource(), 500);
    }

    return () => clearTimeout(timer);
  }, [dialogProps.open]);

  function clearResource() {
    setSelected(null);
  }

  const onSave = useCallback(
    async (formValues: FormValues) => {
      // Create a configuration
      const bpConfig = new BPConfiguration(configuration);

      // Copy the source we're editing
      const editedSource = new BPResourceConfiguration(
        configuration?.spec?.sources![sourceIndex]
      );

      // initialize a new processor
      const newProcessor = new BPResourceConfiguration();

      // Set the parameters and type
      newProcessor.setParamsFromMap(formValues);
      newProcessor.type = selected?.metadata.name;

      // Add processor to the source
      editedSource.processors = [
        ...(editedSource.processors ?? []),
        newProcessor,
      ];

      // Replace the source in the configuration
      bpConfig.replaceSource(editedSource, sourceIndex);
      const update = await bpConfig.apply();

      if (update.status !== UpdateStatus.CONFIGURED) {
        enqueueSnackbar("Failed to add processor.", { variant: "error" });
        console.error("Got unexpected resource status:", update);
        return;
      }

      enqueueSnackbar("Processor Added!", { variant: "success" });
      onSaveSuccess();
    },
    [
      configuration,
      enqueueSnackbar,
      onSaveSuccess,
      selected?.metadata.name,
      sourceIndex,
    ]
  );

  const View: JSX.Element = useMemo(() => {
    function handleSaveNew(formValues: FormValues) {
      onSave(formValues);
      clearResource();
    }

    if (selected == null) {
      // Step one is to select a ProcessorType
      return (
        <CreateProcessorSelectView
          displayName={parentDisplayName}
          telemetryTypes={telemetryTypes}
          onSelect={(pt: ProcessorType) => setSelected(pt)}
        />
      );
    } else {
      // Show the form to configure a new processor
      return (
        <CreateProcessorConfigureView
          processorType={selected}
          onBack={clearResource}
          onSave={handleSaveNew}
        />
      );
    }
  }, [onSave, parentDisplayName, selected, telemetryTypes]);

  return (
    <Dialog
      {...dialogProps}
      fullWidth
      maxWidth="md"
      data-testid="resource-dialog"
      onClose={onClose}
    >
      {View}
    </Dialog>
  );
};

const MemoizedDialog = memo(NewProcessorDialogComponent);

export const NewProcessorDialog: React.FC<NewProcessorDialogProps> = (
  props
) => {
  return (
    <ResourceDialogContextProvider purpose="create" onClose={props.onClose}>
      <MemoizedDialog {...props} />
    </ResourceDialogContextProvider>
  );
};
