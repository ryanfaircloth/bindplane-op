import { gql } from "@apollo/client";
import { Dialog, DialogProps } from "@mui/material";
import { isEqual } from "lodash";
import { useSnackbar } from "notistack";
import { useEffect, useMemo, useState } from "react";
import {
  useProcessorDialogSourceTypeLazyQuery,
  useProcessorDialogDestinationTypeLazyQuery,
  GetProcessorTypesQuery,
  ResourceTypeKind,
  useUpdateProcessorsMutation,
  ResourceConfiguration,
} from "../../../graphql/generated";
import { BPResourceConfiguration } from "../../../utils/classes";
import { usePipelineGraph } from "../../PipelineGraph/PipelineGraphContext";
import {
  CreateProcessorConfigureView,
  CreateProcessorSelectView,
  EditProcessorView,
  FormValues,
} from "../../ResourceConfigForm";
import { ResourceDialogContextProvider } from "../ResourceDialogContext";
import { AllProcessorsView } from "./AllProcessorsView";

interface ProcessorDialogProps extends DialogProps {
  processors: ResourceConfiguration[];
}

gql`
  query processorDialogSourceType($name: String!) {
    sourceType(name: $name) {
      metadata {
        name
        displayName
        description
      }
      spec {
        telemetryTypes
      }
    }
  }

  query processorDialogDestinationType($name: String!) {
    destinationWithType(name: $name) {
      destinationType {
        metadata {
          name
          displayName
          description
        }
        spec {
          telemetryTypes
        }
      }
    }
  }

  mutation updateProcessors($input: UpdateProcessorsInput!) {
    updateProcessors(input: $input)
  }
`;

enum Page {
  MAIN,
  CREATE_PROCESSOR_SELECT,
  CREATE_PROCESSOR_CONFIGURE,
  EDIT_PROCESSOR,
}

export type ProcessorType = GetProcessorTypesQuery["processorTypes"][0];

export const ProcessorDialog: React.FC = (props) => {
  const {
    editProcessorsInfo,
    configuration,
    editProcessorsOpen,
    closeProcessorDialog,
  } = usePipelineGraph();

  if (editProcessorsInfo === null) {
    return null;
  }

  let processors: ResourceConfiguration[] = [];
  switch (editProcessorsInfo?.resourceType) {
    case "source":
      processors =
        configuration?.spec?.sources?.[editProcessorsInfo.index].processors ??
        [];
      break;
    case "destination":
      processors =
        configuration?.spec?.destinations?.[editProcessorsInfo.index]
          .processors ?? [];
      break;
    default:
      processors = [];
  }

  return (
    <ProcessorDialogComponent
      open={editProcessorsOpen}
      onClose={closeProcessorDialog}
      processors={processors}
    />
  );
};

export const ProcessorDialogComponent: React.FC<ProcessorDialogProps> = ({
  processors: processorsProp,
  ...dialogProps
}) => {
  const {
    editProcessorsInfo,
    configuration,
    closeProcessorDialog,
    refetchConfiguration,
  } = usePipelineGraph();

  const [processors, setProcessors] = useState(processorsProp);
  const [view, setView] = useState<Page>(Page.MAIN);
  const [newProcessorType, setNewProcessorType] =
    useState<ProcessorType | null>(null);
  const [editingProcessorIndex, setEditingProcessorIndex] =
    useState<number>(-1);

  const { enqueueSnackbar } = useSnackbar();

  // Get the type of the source or destination to which we're adding processors
  const resourceTypeName = useMemo(() => {
    try {
      switch (editProcessorsInfo?.resourceType) {
        case "source":
          return configuration?.spec?.sources?.[editProcessorsInfo.index].type;
        case "destination":
          return configuration?.spec?.destinations?.[editProcessorsInfo.index]
            .type;
        default:
          return null;
      }
    } catch (err) {
      return null;
    }
  }, [
    configuration?.spec?.destinations,
    configuration?.spec?.sources,
    editProcessorsInfo?.index,
    editProcessorsInfo?.resourceType,
  ]);

  /* ------------------------ GQL Mutations and Queries ----------------------- */
  const [updateProcessors] = useUpdateProcessorsMutation({
    onCompleted: () => {
      refetchConfiguration();
      enqueueSnackbar("Processors saved", {
        variant: "success",
        key: "save-processors-success",
      });
      closeProcessorDialog();
    },
    onError: (error) => {
      console.error(error);
      enqueueSnackbar("Failed to save processors", {
        variant: "error",
        key: "save-processors-error",
      });
    },
  });

  const [fetchSourceType, { data: sourceTypeData }] =
    useProcessorDialogSourceTypeLazyQuery({
      variables: { name: resourceTypeName ?? "" },
    });
  const [fetchDestinationType, { data: destinationTypeData }] =
    useProcessorDialogDestinationTypeLazyQuery();

  /* ----------------------------- Event Handlers ----------------------------- */
  useEffect(() => {
    // resetState returns the processor to the first page after close.
    function resetState() {
      setView(Page.MAIN);
      setNewProcessorType(null);
      setEditingProcessorIndex(-1);
    }

    let timeout: ReturnType<typeof setTimeout>;
    // Resets the state on close
    if (dialogProps.open === false) {
      timeout = setTimeout(resetState, 500);
    }

    return () => {
      clearTimeout(timeout);
    };
  }, [dialogProps.open]);

  useEffect(() => {
    function fetchData() {
      if (editProcessorsInfo!.resourceType === "source") {
        fetchSourceType({ variables: { name: resourceTypeName ?? "" } });
      } else {
        const destName =
          configuration?.spec?.destinations?.[editProcessorsInfo!.index].name;
        fetchDestinationType({ variables: { name: destName ?? "" } });
      }
    }

    if (editProcessorsInfo == null) {
      return;
    }

    fetchData();
  }, [
    configuration?.spec?.destinations,
    editProcessorsInfo,
    fetchDestinationType,
    fetchSourceType,
    resourceTypeName,
  ]);

  /* -------------------------------- Functions ------------------------------- */

  // handleSelectNewProcessorType is called when a user selects a processor type
  // in the CreateProcessorSelect view.
  function handleSelectNewProcessorType(type: ProcessorType) {
    setNewProcessorType(type);
    setView(Page.CREATE_PROCESSOR_CONFIGURE);
  }

  function handleReturnToAll() {
    setView(Page.MAIN);
    setNewProcessorType(null);
  }

  // handleClose is called when a user clicks off the dialog or the "X" button
  function handleClose() {
    if (!isEqual(processors, processorsProp)) {
      const ok = window.confirm("Discard changes?");
      if (!ok) {
        return;
      }
      // reset form values if chooses to discard.
      setProcessors(processorsProp);
    }

    closeProcessorDialog();
  }

  // handleAddProcessor adds a new processor to the list of processors
  async function handleAddProcessor(formValues: FormValues) {
    const processorConfig = new BPResourceConfiguration();
    processorConfig.setParamsFromMap(formValues);
    processorConfig.type = newProcessorType!.metadata.name;

    const newProcessors = [...processors, processorConfig];
    setProcessors(newProcessors);
    handleReturnToAll();
  }

  // handleSaveExisting saves changes to an existing processor in the list
  async function handleSaveExisting(formValues: FormValues) {
    const processorConfig = new BPResourceConfiguration();
    processorConfig.setParamsFromMap(formValues);
    processorConfig.type = processors[editingProcessorIndex].type;

    const newProcessors = [...processors];
    newProcessors[editingProcessorIndex] = processorConfig;
    setProcessors(newProcessors);

    handleReturnToAll();
  }

  // handleRemoveProcessor removes a processor from the list of processors
  async function handleRemoveProcessor(index: number) {
    const newProcessors = [...processors];
    newProcessors.splice(index, 1);
    setProcessors(newProcessors);

    handleReturnToAll();
  }

  // handleEditProcessorClick sets the editing index and switches to the edit page
  function handleEditProcessorClick(index: number) {
    setEditingProcessorIndex(index);
    setView(Page.EDIT_PROCESSOR);
  }

  // handleSave saves the processors to the backend and closes the dialog.
  async function handleSave() {
    if (isEqual(processorsProp, processors)) {
      closeProcessorDialog();
    }

    await updateProcessors({
      variables: {
        input: {
          configuration: configuration?.metadata?.name!,
          resourceType:
            editProcessorsInfo?.resourceType === "source"
              ? ResourceTypeKind.Source
              : ResourceTypeKind.Destination,
          resourceIndex: editProcessorsInfo?.index!,
          processors: processors,
        },
      },
    });
  }

  // displayName for sources is the sourceType name, for destinations it's the destination name
  const displayName = useMemo(() => {
    if (editProcessorsInfo == null) {
      return undefined;
    }

    if (editProcessorsInfo.resourceType === "source") {
      return sourceTypeData?.sourceType?.metadata.displayName;
    } else {
      return configuration?.spec?.destinations?.[editProcessorsInfo.index].name;
    }
  }, [
    configuration?.spec?.destinations,
    editProcessorsInfo,
    sourceTypeData?.sourceType?.metadata.displayName,
  ]);

  let current: JSX.Element;
  switch (view) {
    case Page.MAIN:
      current = (
        <AllProcessorsView
          processors={processors}
          parentDisplayName={displayName ?? "unknown"}
          onAddProcessor={() => setView(Page.CREATE_PROCESSOR_SELECT)}
          onClose={handleClose}
          onEditProcessor={handleEditProcessorClick}
          onSave={handleSave}
          onProcessorsChange={setProcessors}
        />
      );
      break;
    case Page.CREATE_PROCESSOR_SELECT:
      current = (
        <CreateProcessorSelectView
          onClose={handleClose}
          displayName={displayName ?? "unknown"}
          telemetryTypes={
            editProcessorsInfo?.resourceType === "source"
              ? sourceTypeData?.sourceType?.spec?.telemetryTypes
              : destinationTypeData?.destinationWithType.destinationType?.spec
                  ?.telemetryTypes
          }
          onBack={() => setView(Page.MAIN)}
          onSelect={handleSelectNewProcessorType}
        />
      );
      break;
    case Page.CREATE_PROCESSOR_CONFIGURE:
      current = (
        <CreateProcessorConfigureView
          processorType={newProcessorType!}
          onBack={handleReturnToAll}
          onSave={handleAddProcessor}
          onClose={closeProcessorDialog}
        />
      );
      break;
    case Page.EDIT_PROCESSOR:
      current = (
        <EditProcessorView
          processors={processors}
          editingIndex={editingProcessorIndex}
          onEditProcessorSave={handleSaveExisting}
          onBack={handleReturnToAll}
          onRemove={handleRemoveProcessor}
        />
      );
  }

  return (
    <ResourceDialogContextProvider onClose={handleClose} purpose={"edit"}>
      <Dialog {...dialogProps} fullWidth maxWidth={"md"} onClose={handleClose}>
        {current}
      </Dialog>
    </ResourceDialogContextProvider>
  );
};
