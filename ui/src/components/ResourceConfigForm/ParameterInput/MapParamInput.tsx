import {
  FormHelperText,
  Grid,
  Typography,
  Stack,
  OutlinedInput,
  IconButton,
  Box,
  Button,
} from "@mui/material";
import { useState, useMemo, memo } from "react";
import { TrashIcon, PlusCircleIcon } from "../../Icons";
import { validateMapField } from "../validation-functions";
import { useValidationContext } from "../ValidationContext";
import { ParamInputProps } from "./ParameterInput";

const MapParamInputComponent: React.FC<
  ParamInputProps<Record<string, string>>
> = ({ definition, value, onValueChange }) => {
  const initValue = valueToTupleArray(value);
  const [controlValue, setControlValue] = useState<Tuple[]>(initValue);

  const { errors, setError, touched, touch } = useValidationContext();

  const onChangeInput = useMemo(() => {
    return function (
      e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>,
      row: number,
      index: number
    ) {
      setControlValue((prev) => {
        const newVal = [...prev];
        newVal[row][index] = e.target.value;
        return newVal;
      });
    };
  }, []);

  function handleBlur() {
    if (!touched[definition.name]) {
      touch(definition.name);
    }

    const mapValue = tupleArrayToMap(controlValue);
    onValueChange && onValueChange(mapValue);
    setError(definition.name, validateMapField(mapValue, definition.required));
  }

  function handleDeleteRow(rowIndex: number) {
    const newControlValue = removeRow(controlValue, rowIndex);
    setControlValue(newControlValue);

    const mapValue = tupleArrayToMap(newControlValue);
    onValueChange && onValueChange(mapValue);

    setError(definition.name, validateMapField(mapValue, definition.required));
  }

  // Special handling for enter key in Key fields
  function handleKeyFieldEnter(
    e: React.KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>,
    rowIndex: number
  ) {
    if (e.key !== "Enter") {
      return;
    }

    e.preventDefault();

    // go to the next input
    const nextInput = document.querySelector(
      `#${definition.name}-input-${rowIndex * 2 + 1}`
    );

    if (nextInput != null) {
      (nextInput as HTMLElement).focus();
    }
  }

  // Special handling for enter key on Value fields
  function handleValueFieldEnter(
    e: React.KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>,
    rowIndex: number
  ) {
    if (e.key !== "Enter") {
      return;
    }

    e.preventDefault();

    // try to find the next input
    const nextInput = document.querySelector(
      `#${definition.name}-input-${rowIndex * 2 + 2}`
    );

    if (nextInput != null) {
      (nextInput as HTMLElement).focus();
    } else {
      setControlValue((prev) => addRow(prev));
    }
  }

  return (
    <>
      <label aria-required={definition.required} htmlFor={definition.name}>
        {definition.label}
        {definition.required && " *"}
      </label>

      {touched[definition.name] && errors[definition.name] && (
        <FormHelperText key={"error-text"} error>
          {errors[definition.name]}
        </FormHelperText>
      )}

      <FormHelperText key={"description-text"}>
        {definition.description}
      </FormHelperText>

      <Grid container spacing={1} marginY={1}>
        <Grid item xs={6}>
          <Typography marginLeft={4} fontWeight={600}>
            Key
          </Typography>
        </Grid>
        <Grid item xs={6}>
          <Typography marginLeft={2} fontWeight={600}>
            Value
          </Typography>
        </Grid>
      </Grid>

      <Stack spacing={1}>
        {controlValue.map(([k, v], rowIndex) => {
          return (
            <Stack
              key={`${definition.name}-row-${rowIndex}`}
              direction="row"
              spacing={1}
            >
              <OutlinedInput
                autoFocus={rowIndex === controlValue.length - 1}
                id={`${definition.name}-input-${rowIndex * 2}`}
                key={`${definition.name}-${rowIndex}-0-input`}
                data-testid={`${definition.name}-${rowIndex}-0-input`}
                size="small"
                type="text"
                value={k}
                onChange={(e) => onChangeInput(e, rowIndex, 0)}
                onBlur={handleBlur}
                onKeyDown={(e) => handleKeyFieldEnter(e, rowIndex)}
              />

              <OutlinedInput
                id={`${definition.name}-input-${rowIndex * 2 + 1}`}
                key={`${definition.name}-${rowIndex}-1-input`}
                data-testid={`${definition.name}-${rowIndex}-1-input`}
                size="small"
                type="text"
                value={v}
                onChange={(e) => onChangeInput(e, rowIndex, 1)}
                onBlur={handleBlur}
                onKeyDown={(e) => handleValueFieldEnter(e, rowIndex)}
              />

              <IconButton
                key={`${definition.name}-${rowIndex}-remove-button`}
                size={"small"}
                onClick={() => handleDeleteRow(rowIndex)}
                data-testid={`${definition.name}-${rowIndex}-remove-button`}
              >
                <TrashIcon
                  key={`${definition.name}-${rowIndex}-remove-icon`}
                  width={18}
                />
              </IconButton>
            </Stack>
          );
        })}
      </Stack>

      <Box marginLeft={1} marginTop={1}>
        <Button
          startIcon={<PlusCircleIcon />}
          onClick={() => setControlValue((prev) => addRow(prev))}
        >
          New Row
        </Button>
      </Box>
    </>
  );
};

export const MapParamInput = memo(MapParamInputComponent);

// Utility functions
export type Tuple = [string, string];

export function valueToTupleArray(value: any): Tuple[] {
  try {
    const tuples = Object.entries(value);

    tuples.push(["", ""]);
    return tuples as Tuple[];
  } catch (err) {
    return [["", ""]];
  }
}

export function tupleArrayToMap(tuples: Tuple[]): Record<string, string> {
  const mapValue: Record<string, string> = {};
  for (const [k, v] of tuples) {
    if (k === "") {
      continue;
    }

    mapValue[k] = v;
  }

  return mapValue;
}

function addRow(tuples: Tuple[]): Tuple[] {
  const newTuples = [...tuples];
  newTuples.push(["", ""]);
  return newTuples;
}

function removeRow(tuples: Tuple[], removeIndex: number): Tuple[] {
  const newTuples = [...tuples];
  newTuples.splice(removeIndex, 1);
  return newTuples;
}
