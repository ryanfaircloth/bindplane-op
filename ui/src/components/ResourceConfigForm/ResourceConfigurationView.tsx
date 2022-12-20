import { Maybe } from "graphql/jsutils/Maybe";
import { isEqual } from "lodash";
import { memo } from "react";
import {
  initFormValues,
  ConfigureResourceView,
  ValidationContextProvider,
} from ".";
import {
  ParameterDefinition,
  Parameter,
  ResourceConfiguration,
  GetProcessorTypesQuery,
} from "../../graphql/generated";
import { initFormErrors } from "./init-form-values";
import {
  FormValueContextProvider,
  useResourceFormValues,
} from "./ResourceFormContext";

export type ProcessorType = GetProcessorTypesQuery["processorTypes"][0];

export interface FormValues {
  // The name of the Source or Destination
  name?: string;
  // The values for the Parameters
  [key: string]: any;
  // The inline processors configured for the Source or Destination
  processors?: ResourceConfiguration[];
}

interface ResourceConfigurationViewProps {
  // Display name for the resource
  displayName: string;

  description: string;

  // Used to determine some form values.
  kind: "destination" | "source" | "processor";

  // The supported telemetry types of the resource type that is
  // being configured.  a subset of ['logs', 'metrics', 'traces']
  telemetryTypes?: string[];

  parameterDefinitions: ParameterDefinition[];

  // If present the form will use these values as defaults
  parameters?: Maybe<Parameter[]>;

  // If present the form will have a name field at the top and will be sent
  // as the formValues["name"] key.
  includeNameField?: boolean;

  // Used to validate the name field if includeNameField is present.
  existingResourceNames?: string[];

  // If present the form will display a "delete" button which calls
  // the onDelete callback.
  onDelete?: () => void;

  // The callback when the resource is saved.
  onSave?: (formValues: FormValues) => void;
  // The copy on the primary button, defaults to "Save"
  saveButtonLabel?: string;

  // The callback when cancel is clicked.
  onBack?: () => void;

  // If present, whether the resource is paused
  paused?: boolean;

  // Callback for when the Pause/Resume button is clicked
  onTogglePause?: () => void;
}

interface ComponentProps extends ResourceConfigurationViewProps {
  initValues: Record<string, any>;
}

const ResourceConfigurationViewComponent: React.FC<ComponentProps> = ({
  displayName,
  description,
  parameters,
  parameterDefinitions,
  includeNameField,
  existingResourceNames,
  kind,
  onDelete,
  onSave,
  saveButtonLabel,
  paused,
  onTogglePause,
  onBack,
  initValues,
}) => {
  const { formValues } = useResourceFormValues();

  // This is passed down to determine whether to enable the primary save button.
  // If no parameters are passed down, then the form is new and is "dirty".
  const isDirty = parameters == null || !isEqual(initValues, formValues);

  return (
    <ConfigureResourceView
      displayName={displayName}
      description={description}
      kind={kind}
      formValues={formValues}
      includeNameField={includeNameField}
      existingResourceNames={existingResourceNames}
      parameterDefinitions={parameterDefinitions}
      onBack={onBack}
      onSave={onSave}
      saveButtonLabel={saveButtonLabel}
      onDelete={onDelete}
      disableSave={!isDirty}
      paused={paused}
      onTogglePause={onTogglePause}
    />
  );
};

const MemoizedComponent = memo(ResourceConfigurationViewComponent);

export const ResourceConfigurationView: React.FC<ResourceConfigurationViewProps> =
  (props) => {
    const { parameterDefinitions, parameters, includeNameField } = props;

    const initValues = initFormValues(
      parameterDefinitions,
      parameters,
      includeNameField
    );

    const initErrors = initFormErrors(
      parameterDefinitions,
      initValues,
      props.kind,
      props.includeNameField,
      props.existingResourceNames
    );

    return (
      <FormValueContextProvider initValues={initValues}>
        <ValidationContextProvider
          initErrors={initErrors}
          definitions={props.parameterDefinitions}
          includeNameField={includeNameField}
        >
          <MemoizedComponent initValues={initValues} {...props} />
        </ValidationContextProvider>
      </FormValueContextProvider>
    );
  };
