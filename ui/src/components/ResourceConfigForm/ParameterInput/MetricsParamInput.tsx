import {
  Checkbox,
  FormControlLabel,
  Grid,
  Stack,
  Typography,
} from "@mui/material";
import { memo, useMemo } from "react";
import { ParamInputProps } from "./ParameterInput";
import { MetricCategory } from "../../../graphql/generated";
import { union, without } from "lodash";

import styles from "./parameter-input.module.scss";

const MetricsParamInputComponent: React.FC<ParamInputProps<string[]>> = ({
  value,
  onValueChange,
  definition,
}) => {
  // Get the columns
  const column1 = useMemo(() => {
    return (
      definition.options.metricCategories?.filter((c) => c.column === 0) ?? []
    );
  }, [definition.options.metricCategories]);

  const column2 = useMemo(() => {
    return (
      definition.options.metricCategories?.filter((c) => c.column === 1) ?? []
    );
  }, [definition.options.metricCategories]);

  function handleToggleValue(toggleValue: string) {
    const newValue = [...(value ?? [])];

    if (!newValue.includes(toggleValue)) {
      // Make sure that toggleValue is in new value array
      newValue.push(toggleValue);
    } else {
      // Remove the toggle value from the array
      const atIndex = newValue.findIndex((v) => v === toggleValue);
      if (atIndex > -1) {
        newValue.splice(atIndex, 1);
      }
    }

    onValueChange && onValueChange(newValue);
  }

  function handleDeselectAll(
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    selectingValues: string[]
  ) {
    e.preventDefault();
    const newValue = union(selectingValues, value);
    onValueChange && onValueChange(newValue);
  }

  function handleSelectAll(
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    deselectingValues: string[]
  ) {
    e.preventDefault();
    const newValue = without(value, ...deselectingValues);
    onValueChange && onValueChange(newValue);
  }

  return (
    <>
      <Grid container spacing={10}>
        {[column1, column2].map((c, ix) => (
          <CategoryStack
            key={`${definition.name}-column-${ix}`}
            handleToggle={handleToggleValue}
            metricCategories={c}
            value={value ?? []}
            handleDeselectAll={handleDeselectAll}
            handleSelectAll={handleSelectAll}
          />
        ))}
      </Grid>
    </>
  );
};

const CategoryStack: React.FC<{
  metricCategories: MetricCategory[];
  value: string[];
  handleToggle: (toggleValue: string) => void;
  handleDeselectAll: (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    deselectingValues: string[]
  ) => void;
  handleSelectAll: (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    deselectingValues: string[]
  ) => void;
}> = ({
  metricCategories,
  value,
  handleToggle,
  handleSelectAll,
  handleDeselectAll,
}) => {
  return (
    <Grid item xs={6}>
      {metricCategories?.map((category) => {
        const allValues = category.metrics.map((m) => m.name);
        return (
          <>
            <div className={styles["metric-category-label"]}>
              <Stack
                direction="row"
                justifyContent="space-between"
                alignItems="center"
              >
                <Typography fontSize={16} fontWeight={600}>
                  {category.label}
                </Typography>

                <div>
                  <button
                    className={styles["metric-category-button"]}
                    onClick={(e) => handleDeselectAll(e, allValues)}
                  >
                    Disable All
                  </button>
                  |
                  <button
                    className={styles["metric-category-button"]}
                    onClick={(e) => handleSelectAll(e, allValues)}
                  >
                    Enable All
                  </button>
                </div>
              </Stack>
            </div>
            <Stack>
              {category.metrics.map((m) => {
                return (
                  <FormControlLabel
                    control={
                      <Checkbox
                        onChange={() => {
                          handleToggle(m.name);
                        }}
                        name={m.name}
                        checked={!value.includes(m.name)}
                      />
                    }
                    classes={{ root: styles["metric-label"] }}
                    label={<Typography fontSize={18}>{m.name}</Typography>}
                  />
                );
              })}
            </Stack>
          </>
        );
      })}
    </Grid>
  );
};

export const MetricsParamInput = memo(MetricsParamInputComponent);
