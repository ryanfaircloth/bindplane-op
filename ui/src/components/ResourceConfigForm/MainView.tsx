import { Grid, Button, Typography } from "@mui/material";
import { isFunction } from "lodash";
import { ParameterDefinition } from "../../graphql/generated";
import {
  ParameterInput,
  ResourceNameInput,
  satisfiesRelevantIf,
  useValidationContext,
  isValid,
} from ".";
import { InlineProcessorContainer } from "./InlineProcessorContainer";
import { useResourceFormValues } from "./ResourceFormContext";
import { useResourceDialog } from "../ResourceDialog/ResourceDialogContext";
import { memo, useMemo } from "react";
import { TitleSection } from "../ResourceDialog/TitleSection";
import { ContentSection } from "../ResourceDialog/ContentSection";
import { ActionsSection } from "../ResourceDialog/ActionSection";
import { initFormErrors } from "./init-form-values";

import mixins from "../../styles/mixins.module.scss";

interface MainProps {
  kind: "source" | "destination" | "processor";
  displayName: string;
  description: string;
  formValues: { [key: string]: any };
  includeNameField?: boolean;
  existingResourceNames?: string[];
  parameterDefinitions: ParameterDefinition[];
  enableProcessors?: boolean;
  onBack?: () => void;
  onSave?: (formValues: { [key: string]: any }) => void;
  saveButtonLabel?: string;
  onDelete?: () => void;
  onAddProcessor: () => void;
  onEditProcessor: (editingIndex: number) => void;
  onRemoveProcessor: (removeIndex: number) => void;
  disableSave?: boolean;
}

export const MainViewComponent: React.FC<MainProps> = ({
  kind,
  displayName,
  description,
  formValues,
  includeNameField,
  existingResourceNames,
  parameterDefinitions,
  enableProcessors,
  onBack,
  onSave,
  saveButtonLabel,
  onDelete,
  onAddProcessor,
  onEditProcessor,
  disableSave,
}) => {
  const { touchAll, setErrors } = useValidationContext();
  const { setFormValues } = useResourceFormValues();
  const { purpose, onClose } = useResourceDialog();

  function handleSubmit() {
    const errors = initFormErrors(
      parameterDefinitions,
      formValues,
      kind,
      includeNameField
    );

    if (!isValid(errors)) {
      touchAll();
      setErrors(errors);
      return;
    }

    isFunction(onSave) && onSave(formValues);
  }

  const primaryButton: JSX.Element = (
    <Button
      disabled={disableSave}
      type="submit"
      variant="contained"
      data-testid="resource-form-save"
      onClick={handleSubmit}
    >
      {saveButtonLabel ?? "Save"}
    </Button>
  );

  const backButton: JSX.Element | null = isFunction(onBack) ? (
    <Button variant="contained" color="secondary" onClick={onBack}>
      Back
    </Button>
  ) : null;

  const deleteButton: JSX.Element | null = isFunction(onDelete) ? (
    <Button variant="outlined" color="error" onClick={onDelete}>
      Delete
    </Button>
  ) : null;

  const title = useMemo(() => {
    const capitalizedResource = kind[0].toUpperCase() + kind.slice(1);
    const action = purpose === "create" ? "Add" : "Edit";
    return `${action} ${capitalizedResource}: ${displayName}`;
  }, [displayName, kind, purpose]);

  return (
    <>
      <TitleSection title={title} description={description} onClose={onClose} />

      <ContentSection>
        <form data-testid="resource-form">
          <Grid container spacing={3} className={mixins["mb-5"]}>
            {includeNameField && (
              <Grid item xs={6}>
                <ResourceNameInput
                  kind={kind}
                  value={formValues.name}
                  onValueChange={(v: string) =>
                    setFormValues((prev) => ({ ...prev, name: v }))
                  }
                  existingNames={existingResourceNames}
                />
              </Grid>
            )}

            <Grid item xs={12}>
              <Typography fontWeight={600} fontSize={24}>
                Configure
              </Typography>
            </Grid>

            {parameterDefinitions.length === 0 ? (
              <Grid item>
                <Typography>No additional configuration needed.</Typography>
              </Grid>
            ) : (
              parameterDefinitions.map((p) => {
                if (satisfiesRelevantIf(formValues, p)) {
                  return <ParameterInput key={p.name} definition={p} />;
                }

                return null;
              })
            )}
          </Grid>

          {enableProcessors && (
            <InlineProcessorContainer
              processors={formValues.processors ?? []}
              onAddProcessor={onAddProcessor}
              onEditProcessor={onEditProcessor}
            />
          )}
        </form>
      </ContentSection>

      <ActionsSection>
        {deleteButton}
        {backButton}
        {primaryButton}
      </ActionsSection>
    </>
  );
};

export const MainView = memo(MainViewComponent);
