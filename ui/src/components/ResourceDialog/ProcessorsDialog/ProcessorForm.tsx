import { Grid, Typography } from "@mui/material";
import { Parameter, ParameterDefinition } from "../../../graphql/generated";
import { ParameterInput } from "../../ResourceConfigForm/ParameterInput";
import { useResourceFormValues } from "../../ResourceConfigForm/ResourceFormContext";
import { satisfiesRelevantIf } from "../../ResourceConfigForm/satisfiesRelevantIf";

import mixins from "../../../styles/mixins.module.scss";

interface Props {
  parameterDefinitions: ParameterDefinition[];
  parameters?: Parameter[];
}

export const ProcessorForm: React.FC<Props> = ({ parameterDefinitions }) => {
  const { formValues } = useResourceFormValues();
  return (
    <form>
      <Grid container spacing={3} className={mixins["mb-5"]}>
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
    </form>
  );
};
