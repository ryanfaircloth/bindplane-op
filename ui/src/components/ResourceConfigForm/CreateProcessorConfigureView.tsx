import { Button } from "@mui/material";
import {
  FormValues,
  initFormValues,
  isValid,
  ProcessorType,
  useValidationContext,
  ValidationContextProvider,
} from ".";
import { ParameterDefinition } from "../../graphql/generated";
import { ActionsSection } from "../ResourceDialog/ActionSection";
import { ContentSection } from "../ResourceDialog/ContentSection";
import { useResourceDialog } from "../ResourceDialog/ResourceDialogContext";
import { TitleSection } from "../ResourceDialog/TitleSection";
import { initFormErrors } from "./init-form-values";
import { ProcessorForm } from "./ProcessorForm";
import {
  FormValueContextProvider,
  useResourceFormValues,
} from "./ResourceFormContext";

interface CreateProcessorConfigureViewProps {
  processorType: ProcessorType;
  parameterDefinitions: ParameterDefinition[];
  onBack: () => void;
  onSave: (formValues: FormValues) => void;
}

const CreateProcessorConfigureViewComponent: React.FC<CreateProcessorConfigureViewProps> =
  ({ processorType, onSave, onBack }) => {
    const { formValues } = useResourceFormValues();
    const { onClose } = useResourceDialog();
    const { errors } = useValidationContext();

    return (
      <>
        <TitleSection
          title={`Add Processor: ${processorType.metadata.displayName}`}
          description={processorType.metadata.description ?? ""}
          onClose={onClose}
        />

        <ContentSection>
          <ProcessorForm parameterDefinitions={processorType.spec.parameters} />
        </ContentSection>

        <ActionsSection>
          <Button variant="contained" color="secondary" onClick={onBack}>
            Back
          </Button>

          <Button
            variant="contained"
            color="primary"
            onClick={() => onSave(formValues)}
            disabled={!isValid(errors)}
          >
            Save
          </Button>
        </ActionsSection>
      </>
    );
  };

export const CreateProcessorConfigureView: React.FC<CreateProcessorConfigureViewProps> =
  (props) => {
    const initValues = initFormValues(props.processorType.spec.parameters);
    const initErrors = initFormErrors(
      props.processorType.spec.parameters,
      initValues,
      "processor"
    );
    return (
      <FormValueContextProvider initValues={initValues}>
        <ValidationContextProvider
          initErrors={initErrors}
          definitions={props.parameterDefinitions}
        >
          <CreateProcessorConfigureViewComponent {...props} />
        </ValidationContextProvider>
      </FormValueContextProvider>
    );
  };
