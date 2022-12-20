import { Button } from "@mui/material";
import {
  FormValues,
  initFormValues,
  isValid,
  ProcessorType,
  useValidationContext,
  ValidationContextProvider,
} from "../../ResourceConfigForm";
import { ActionsSection } from "../ActionSection";
import { ContentSection } from "../ContentSection";
import { TitleSection } from "../TitleSection";
import { initFormErrors } from "../../ResourceConfigForm/init-form-values";
import { ProcessorForm } from "./ProcessorForm";
import {
  FormValueContextProvider,
  useResourceFormValues,
} from "../../ResourceConfigForm/ResourceFormContext";

interface CreateProcessorConfigureViewProps {
  processorType: ProcessorType;
  onBack: () => void;
  onSave: (formValues: FormValues) => void;
  onClose: () => void;
}

const CreateProcessorConfigureViewComponent: React.FC<CreateProcessorConfigureViewProps> =
  ({ processorType, onSave, onBack, onClose }) => {
    const { formValues } = useResourceFormValues();
    const { touchAll, setErrors } = useValidationContext();

    function handleSave() {
      const errors = initFormErrors(
        processorType.spec.parameters,
        formValues,
        "processor"
      );

      if (!isValid(errors)) {
        setErrors(errors);
        touchAll();
        return;
      }

      onSave(formValues);
    }

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

          <Button variant="contained" color="primary" onClick={handleSave}>
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
          definitions={props.processorType.spec.parameters}
        >
          <CreateProcessorConfigureViewComponent {...props} />
        </ValidationContextProvider>
      </FormValueContextProvider>
    );
  };
