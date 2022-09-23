import { FormControl, InputLabel, FormHelperText } from "@mui/material";
import { isEmpty, isFunction } from "lodash";
import { useState, ChangeEvent, memo } from "react";
import { YamlEditor } from "../../YamlEditor";
import { validateYamlField } from "../validation-functions";
import { useValidationContext } from "../ValidationContext";
import { ParamInputProps } from "./ParameterInput";

const YamlParamInputComponent: React.FC<ParamInputProps<string>> = ({
  definition,
  value,
  onValueChange,
}) => {
  const [isFocused, setFocused] = useState(false);

  const { touch, touched, errors, setError } = useValidationContext();

  const shrinkLabel = isFocused || !isEmpty(value);

  function handleValueChange(e: ChangeEvent<HTMLTextAreaElement>) {
    const value = e.target.value;
    isFunction(onValueChange) && onValueChange(value);

    setError(definition.name, validateYamlField(value, definition.required));
    touch(definition.name);
  }

  return (
    <FormControl fullWidth required={definition.required}>
      <InputLabel
        shrink={shrinkLabel}
        htmlFor={definition.name}
        style={{
          backgroundColor: "#fff",
          color: shrinkLabel ? "#4abaeb" : undefined,
          padding: shrinkLabel ? "0 10px 0 5px" : undefined,
        }}
      >
        {definition.label}
      </InputLabel>
      <YamlEditor
        required={definition.required}
        name={definition.name}
        value={value ?? ""}
        onValueChange={handleValueChange}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        minHeight={200}
      />
      {touched[definition.name] && errors[definition.name] && (
        <FormHelperText error>{errors[definition.name]}</FormHelperText>
      )}
      <FormHelperText>{definition.description}</FormHelperText>
      {(definition.documentation ?? []).map((d) => {
        return (
          <FormHelperText key={d.text}>
            <a href={d.url} rel="noreferrer" target="_blank">
              {d.text}
            </a>
          </FormHelperText>
        );
      })}
    </FormControl>
  );
};

export const YamlParamInput = memo(YamlParamInputComponent);
