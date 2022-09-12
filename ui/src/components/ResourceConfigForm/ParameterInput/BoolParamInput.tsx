import { Checkbox, FormControlLabel, Switch, Typography } from "@mui/material";
import { isFunction } from "lodash";
import { ParamInputProps } from "./ParameterInput";
import { memo } from "react";
import { ParameterDefinition } from "../../../graphql/generated";

import styles from "./parameter-input.module.scss";

const BoolParamInputComponent: React.FC<ParamInputProps<boolean>> = ({
  definition,
  value,
  onValueChange,
}) => {
  return isTelemetryHeader(definition) ? (
    <TelemetrySectionHeader
      definition={definition}
      value={value}
      onValueChange={onValueChange}
    />
  ) : (
    <FormControlLabel
      control={
        <Checkbox
          onChange={(e) => {
            isFunction(onValueChange) && onValueChange(e.target.checked);
          }}
          name={definition.name}
          checked={value}
        />
      }
      label={
        definition.options.sectionHeader ? (
          <Typography fontSize={18} fontWeight={600}>
            {definition.label}
          </Typography>
        ) : (
          definition.label
        )
      }
    />
  );
};

const TelemetrySectionHeader: React.FC<ParamInputProps<boolean>> = ({
  definition,
  value,
  onValueChange,
}) => {
  return (
    <FormControlLabel
      control={
        <Switch
          onChange={(e) => {
            isFunction(onValueChange) && onValueChange(e.target.checked);
          }}
          name={definition.name}
          checked={value}
          classes={{
            root: styles["telemetry-header-switch"],
          }}
        />
      }
      label={
        <Typography fontSize={18} fontWeight={600}>
          {definition.label}
        </Typography>
      }
      labelPlacement={"start"}
      classes={{
        root: styles["telemetry-header-root"],
        label: "label",
      }}
    />
  );
};

function isTelemetryHeader(definition: ParameterDefinition) {
  return ["enable_metrics", "enable_logs", "enable_traces"].includes(
    definition.name
  );
}

export const BoolParamInput = memo(BoolParamInputComponent);
