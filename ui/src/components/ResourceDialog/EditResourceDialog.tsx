import { Dialog, DialogProps } from "@mui/material";
import {
  Maybe,
  Parameter,
  ParameterDefinition,
  ResourceConfiguration,
} from "../../graphql/generated";
import { ResourceConfigForm } from "../ResourceConfigForm";
import { ResourceDialogContextProvider } from "./ResourceDialogContext";
import { isFunction } from "lodash";

interface EditResourceDialogProps extends DialogProps {
  displayName: string;
  description: string;
  onSave: (values: { [key: string]: any }) => void;
  onDelete?: () => void;
  onCancel: () => void;
  parameters: Parameter[];
  parameterDefinitions: ParameterDefinition[];
  processors?: Maybe<ResourceConfiguration[]>;
  enableProcessors?: boolean;
  includeNameField?: boolean;
  kind: "source" | "destination";
  // The supported telemetry types of the resource type that is
  // being configured.  a subset of ['logs', 'metrics', 'traces']
  telemetryTypes?: string[];
}

const EditResourceDialogComponent: React.FC<EditResourceDialogProps> = ({
  onSave,
  onDelete,
  onCancel,
  displayName,
  description,
  parameters,
  processors,
  enableProcessors,
  parameterDefinitions,
  kind,
  telemetryTypes,
  includeNameField = false,
  ...dialogProps
}) => {
  return (
    <Dialog {...dialogProps} onClose={onCancel} fullWidth maxWidth="md">
      <ResourceConfigForm
        includeNameField={includeNameField}
        displayName={displayName}
        description={description}
        kind={kind}
        parameterDefinitions={parameterDefinitions}
        parameters={parameters}
        processors={processors}
        enableProcessors={enableProcessors}
        onSave={onSave}
        onDelete={onDelete}
        telemetryTypes={telemetryTypes}
      />
    </Dialog>
  );
};

export const EditResourceDialog: React.FC<EditResourceDialogProps> = (
  props
) => {
  function handleClose() {
    if (isFunction(props.onClose)) {
      props.onClose({}, "backdropClick");
    }
  }
  return (
    <ResourceDialogContextProvider purpose="edit" onClose={handleClose}>
      <EditResourceDialogComponent {...props} />
    </ResourceDialogContextProvider>
  );
};
