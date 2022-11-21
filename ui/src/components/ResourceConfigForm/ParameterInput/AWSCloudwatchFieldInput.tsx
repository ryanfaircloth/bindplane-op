import { ParamInputProps } from "./ParameterInput";
import { StringsParamInput } from "./StringsParamInput";
import { StringParamInput } from "./StringParamInput";
import { useValidationContext } from "../ValidationContext";
import {
  Button,
  Stack,
  IconButton,
  FormLabel,
  FormHelperText,
} from "@mui/material";
import { ParameterType } from "../../../graphql/generated";
import { PlusCircleIcon, TrashIcon } from "../../Icons";
import { cloneDeep } from "lodash";
import { useMemo, useState } from "react";
import { validateAWSNamedField } from "../validation-functions";

export type InputValue = ItemValue[];

interface ItemValue {
  id: string;
  prefixes: string[];
  names: string[];
}

export const AWSCloudwatchInput: React.FC<ParamInputProps<InputValue>> = ({
  definition,
  value: paramValue,
  onValueChange,
}) => {
  const initValue =
    paramValue != null && paramValue.length > 0
      ? paramValue
      : [
          {
            id: "",
            prefixes: [],
            names: [],
          },
        ];
  const [controlValue, setControlValue] = useState<InputValue>(initValue);
  const { errors, touched, touch, setError } = useValidationContext();

  const onIDValueChange = useMemo(() => {
    return function (v: string, index: number) {
      if (!touched[definition.name]) {
        touch(definition.name);
      }

      const newValue = cloneDeep(controlValue);
      newValue[index].id = v;

      setControlValue(newValue);
      onValueChange && onValueChange(newValue);
      setError(definition.name, validateAWSNamedField(newValue));
    };
  }, [controlValue, definition.name, onValueChange, setError, touch, touched]);

  const onStringsValueChange = useMemo(() => {
    return function (v: string[], index: number, field: "names" | "prefixes") {
      const newValue = cloneDeep(controlValue);
      newValue[index][field] = v;

      setControlValue(newValue);
      onValueChange && onValueChange(newValue);
    };
  }, [controlValue, onValueChange]);

  function addNewField() {
    const defaultItem: ItemValue = {
      id: "",
      names: [],
      prefixes: [],
    };

    const curValue = cloneDeep(controlValue) ?? [];
    curValue.push(defaultItem);
    setControlValue(curValue);
  }

  function removeField(index: number) {
    if (controlValue.length === 1) {
      return;
    }

    const curValue = cloneDeep(controlValue) ?? [];
    curValue.splice(index, 1);
    setControlValue(curValue);

    onValueChange && onValueChange(curValue);
  }

  return (
    <>
      <FormLabel filled={true}>{definition.label}</FormLabel>
      <FormHelperText filled={true}>{definition.description}</FormHelperText>
      {errors[definition.name] && touched[definition.name] && (
        <Stack>
          <FormHelperText sx={{ marginLeft: 0 }} component="span" error>
            {errors[definition.name]}
          </FormHelperText>
        </Stack>
      )}
      {controlValue.map((itemValue, index) => {
        return (
          <Stack
            key={`aws-parameter-item-${index}`}
            direction={"row"}
            spacing={2}
            padding={2}
            alignItems={"center"}
            justifyContent="center"
          >
            <Stack
              border="2px solid rgba(0, 0, 0, 0.2)"
              borderRadius={5}
              spacing={2}
              padding={2}
            >
              <StringParamInput
                definition={{
                  name: "named_groups_" + index,
                  label: "ID",
                  type: ParameterType.String,
                  required: true,
                  description: "Specified log group to collect",
                  options: {},
                }}
                value={itemValue.id}
                onValueChange={(v) => onIDValueChange(v, index)}
              />

              <StringsParamInput
                definition={{
                  name: "aws.names",
                  label: "Names",
                  type: ParameterType.Strings,
                  required: false,
                  description:
                    "A list of full log stream names to filter the discovered log groups to collect from.",
                  options: {},
                }}
                value={itemValue.names}
                onValueChange={(v) => onStringsValueChange(v, index, "names")}
              />

              <StringsParamInput
                definition={{
                  name: "aws.prefixes",
                  label: "Prefixes",
                  type: ParameterType.Strings,
                  required: false,
                  description:
                    "A list of prefixes to filter the discovered log groups to collect from.",
                  options: {},
                }}
                value={itemValue.prefixes}
                onValueChange={(v) =>
                  onStringsValueChange(v, index, "prefixes")
                }
              />
            </Stack>
            {!(index === 0 && controlValue.length === 1) && (
              <IconButton size={"small"} onClick={() => removeField(index)}>
                <TrashIcon width={18} />
              </IconButton>
            )}
          </Stack>
        );
      })}
      <Stack width={"100%"} alignItems="center">
        <Button startIcon={<PlusCircleIcon />} onClick={addNewField}>
          New field
        </Button>
      </Stack>
    </>
  );
};
