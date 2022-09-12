import { Dialog, DialogProps } from "@mui/material";
import React, { memo, useEffect, useMemo, useState } from "react";
import { isFunction } from "lodash";
import { ResourceDialogContextProvider } from "./ResourceDialogContext";
import { SelectView } from "./SelectView";
import { ChooseView } from "./ChooseView";
import { ConfigureView } from "./ConfigureView";
import { DialogResource, ResourceType } from ".";

interface ResourceDialogProps extends DialogProps {
  kind: "destination" | "source";

  // Either SourceType[] or DestinationType[]
  resourceTypes: ResourceType[];

  // If present, it will allow users to choose existing Resources
  resources?: DialogResource[];

  // Callback for the save button when creating a new Resource.
  onSaveNew?: (
    parameters: { [key: string]: any },
    resourceType: ResourceType
  ) => void;

  // Callback for saving an existing resource.
  onSaveExisting?: (resource: DialogResource) => void;

  onClose: () => void;

  // The supported telemetry types of the resource type that is
  // being configured.  a subset of ['logs', 'metrics', 'traces']
  telemetryTypes?: string[];
}

// There are three possible views to this dialog -
// 1. The "Select" view is the first view, which displays all available ResourceTypes.
//    This is displayed anytime selected == null
//    To go do this step its sufficient to set selected to null.
//
// 2. The "Choose" view is an optional second view, shown after the Select view
//    and there are existing resources from the resources prop
//    such that resource.spec.type === selected.metadata.name
//
// 3. The "Configure" view.  The final step for creation, contains the ResourceTypeForm
//    This is displayed when selected != null and createNew == true.
//
export const ResourceDialogComponent: React.FC<ResourceDialogProps> = ({
  kind,
  resourceTypes,
  resources,
  onSaveNew,
  onSaveExisting,
  onClose,
  ...dialogProps
}) => {
  const [selected, setSelected] = useState<ResourceType | null>(null);
  const [createNew, setCreateNew] = useState(false);

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
    setCreateNew(false);
  }

  const View: JSX.Element = useMemo(() => {
    function handleSaveNew(
      parameters: { [key: string]: any },
      resourceType: ResourceType
    ) {
      isFunction(onSaveNew) && onSaveNew(parameters, resourceType);
      clearResource();
    }

    function handleSaveExisting(resource: DialogResource) {
      isFunction(onSaveExisting) && onSaveExisting(resource);
      clearResource();
    }
    if (selected == null) {
      // Step one is to select a ResourceType
      return (
        <SelectView
          resourceTypes={resourceTypes}
          resources={resources ?? []}
          setSelected={setSelected}
          setCreateNew={setCreateNew}
          kind={kind}
        />
      );
    } else if (
      resources?.some((r) => r.spec.type === selected.metadata.name) &&
      !createNew
    ) {
      // There are existing resources that match the selected type.
      // Show the "Choose View" to allow users to pick existing resources.
      return (
        <ChooseView
          resources={resources}
          selected={selected}
          kind={kind}
          clearResource={clearResource}
          handleSaveExisting={handleSaveExisting}
          setCreateNew={setCreateNew}
        />
      );
    } else {
      // Show the form to configure a new source
      return (
        <ConfigureView
          selected={selected}
          kind={kind}
          createNew={createNew}
          clearResource={clearResource}
          handleSaveNew={handleSaveNew}
          resources={resources ?? []}
        />
      );
    }
  }, [
    createNew,
    kind,
    onSaveExisting,
    onSaveNew,
    resourceTypes,
    resources,
    selected,
  ]);

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

const MemoizedDialog = memo(ResourceDialogComponent);

export const NewResourceDialog: React.FC<ResourceDialogProps> = (props) => {
  return (
    <ResourceDialogContextProvider purpose="create" onClose={props.onClose}>
      <MemoizedDialog {...props} />
    </ResourceDialogContextProvider>
  );
};
