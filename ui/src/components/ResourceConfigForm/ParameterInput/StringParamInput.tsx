import { FormHelperText, TextField } from "@mui/material";
import { isFunction } from "lodash";
import { ChangeEvent, memo } from "react";
import { validateStringField } from "../validation-functions";
import { useValidationContext } from "../ValidationContext";
import { ParamInputProps } from "./ParameterInput";

const StringParamInputComponent: React.FC<ParamInputProps<string>> = ({
  definition,
  value,
  onValueChange,
}) => {
  const { errors, setError, touched, touch } = useValidationContext();

  function handleValueChange(e: ChangeEvent<HTMLInputElement>) {
    const newValue = e.target.value;
    isFunction(onValueChange) && onValueChange(e.target.value);

    if (!touched[definition.name]) {
      touch(definition.name);
    }

    setError(
      definition.name,
      validateStringField(newValue, definition.required)
    );
  }

  return (
    <TextField
      multiline={definition.options.multiline ?? undefined}
      value={value}
      onChange={handleValueChange}
      onBlur={() => touch(definition.name)}
      name={definition.name}
      fullWidth
      size="small"
      label={definition.label}
      FormHelperTextProps={{
        sx: { paddingLeft: "-2px" },
      }}
      helperText={
        <>
          {errors[definition.name] && touched[definition.name] && (
            <>
              <FormHelperText sx={{ marginLeft: 0 }} component="span" error>
                {errors[definition.name]}
              </FormHelperText>
              <br />
            </>
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

export const StringParamInput = memo(StringParamInputComponent);
