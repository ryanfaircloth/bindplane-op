import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Grid,
  Typography,
} from "@mui/material";
import { ChevronDown } from "../Icons";
import { ParameterGroup } from "./MainView";
import { ParameterInput } from "./ParameterInput";
import { useResourceFormValues } from "./ResourceFormContext";
import { satisfiesRelevantIf } from "./satisfiesRelevantIf";

interface ParameterSectionProps {
  group: ParameterGroup;
}

export const ParameterSection: React.FC<ParameterSectionProps> = ({
  group,
}) => {
  const { formValues } = useResourceFormValues();

  const inputs: JSX.Element[] = [];
  for (const p of group.parameters) {
    if (satisfiesRelevantIf(formValues, p)) {
      inputs.push(<ParameterInput key={p.name} definition={p} />);
    }
  }

  if (group.advanced && inputs.length > 0) {
    return (
      <Grid item xs={12}>
        <Accordion>
          <AccordionSummary expandIcon={<ChevronDown />}>
            <Typography>Advanced</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Grid container spacing={3}>
              {inputs}
            </Grid>
          </AccordionDetails>
        </Accordion>
      </Grid>
    );
  } else {
    return <>{inputs}</>;
  }
};
