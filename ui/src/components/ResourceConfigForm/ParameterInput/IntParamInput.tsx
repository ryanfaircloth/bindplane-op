import { FormHelperText, TextField } from "@mui/material";
import { isFunction } from "lodash";
import { ChangeEvent, memo } from "react";
import { ParamInputProps } from "./ParameterInput";
import { useValidationContext } from "../ValidationContext";
import { validateIntField } from "../validation-functions";

const IntParamInputComponent: React.FC<ParamInputProps<number>> = ({
  definition,
  value,
  onValueChange,
}) => {
  const { errors, setError, touch, touched } = useValidationContext();

  function handleValueChange(e: ChangeEvent<HTMLInputElement>) {
    const newValue = Number(e.target.value);

    if (isNaN(newValue)) {
      return;
    }

    setError(definition.name, validateIntField(definition, newValue));

    isFunction(onValueChange) && onValueChange(newValue);
  }

  return (
    <TextField
      value={value ?? ""}
      onChange={handleValueChange}
      name={definition.name}
      fullWidth
      size="small"
      label={definition.label}
      onBlur={() => touch(definition.name)}
      FormHelperTextProps={{
        sx: { paddingLeft: "-2px" },
      }}
      helperText={
        <>
          {errors[definition.name] && touched[definition.name] && (
            <FormHelperText sx={{ marginLeft: 0 }} component="span" error>
              {errors[definition.name]}
            </FormHelperText>
          )}
          <FormHelperText sx={{ marginLeft: 0 }} component="span">
            {definition.description}
          </FormHelperText>
        </>
      }
      required={definition.required}
      autoComplete="off"
      autoCorrect="off"
      autoCapitalize="off"
      spellCheck="false"
    />
  );
};

export const IntParamInput = memo(IntParamInputComponent);
