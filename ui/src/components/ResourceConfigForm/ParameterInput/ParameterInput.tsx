import { Grid } from "@mui/material";
import { useMemo } from "react";
import {
  BoolParamInput,
  EnumParamInput,
  EnumsParamInput,
  IntParamInput,
  MapParamInput,
  MetricsParamInput,
  StringParamInput,
  StringsParamInput,
  TimezoneParamInput,
  YamlParamInput,
} from ".";
import { ParameterDefinition, ParameterType } from "../../../graphql/generated";
import { useResourceFormValues } from "../ResourceFormContext";
import { AWSCloudwatchInput } from "./AWSCloudwatchFieldInput";

export interface ParamInputProps<T> {
  definition: ParameterDefinition;
  value?: T;
  onValueChange?: (v: T) => void;
}

export const ParameterInput: React.FC<{ definition: ParameterDefinition }> = ({
  definition,
}) => {
  const { formValues, setFormValues } = useResourceFormValues();
  const onValueChange = useMemo(
    () => (newValue: any) => {
      setFormValues((prev) => ({ ...prev, [definition.name]: newValue }));
    },
    [definition.name, setFormValues]
  );

  const Control: JSX.Element = useMemo(() => {
    switch (definition.type) {
      case ParameterType.String:
        return (
          <StringParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Strings:
        return (
          <StringsParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Enum:
        return (
          <EnumParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Enums:
        return (
          <EnumsParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Bool:
        return (
          <BoolParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Int:
        return (
          <IntParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Map:
        return (
          <MapParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Yaml:
        return (
          <YamlParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Timezone:
        return (
          <TimezoneParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.Metrics:
        return (
          <MetricsParamInput
            definition={definition}
            value={formValues[definition.name]}
            onValueChange={onValueChange}
          />
        );

      case ParameterType.AwsCloudwatchNamedField:
        return (
          <>
            <AWSCloudwatchInput
              definition={definition}
              value={formValues[definition.name]}
              onValueChange={onValueChange}
            />
          </>
        );
    }
  }, [definition, formValues, onValueChange]);

  const gridColumns = useMemo(() => {
    if (isMetricsType(definition)) {
      return 12;
    }

    if (isTelemetryHeader(definition)) {
      return 12;
    }

    if (isSectionHeader(definition)) {
      return 12;
    }

    if (isAWSCloudwatch(definition)) {
      return 12;
    }

    return definition.options.gridColumns ?? 6;
  }, [definition]);

  return (
    <Grid item xs={gridColumns}>
      {Control}
    </Grid>
  );
};

function isTelemetryHeader(definition: ParameterDefinition) {
  return ["enable_metrics", "enable_logs", "enable_traces"].includes(
    definition.name
  );
}

function isSectionHeader(definition: ParameterDefinition) {
  return definition.options.sectionHeader === true;
}

function isMetricsType(definition: ParameterDefinition) {
  return definition.type === "metrics";
}

function isAWSCloudwatch(definition: ParameterDefinition) {
  return definition.type === "awsCloudwatchNamedField";
}
