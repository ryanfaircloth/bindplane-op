import { createContext, useContext, useState } from "react";
import { ParameterDefinition } from "../../graphql/generated";

interface ValidationContextValue {
  // errors is a map from parameter name to error message
  errors: Record<string, string | null>;
  // setError sets the parameter to the error value in errors
  setError: (parameterName: string, error: string | null) => void;
  // setErrors sets the errors state
  setErrors: (errors: Record<string, null | string>) => void;

  // touched is a map from parameter name to boolean
  touched: Record<string, boolean>;
  // touch sets the parameter with given name to true in the touched map
  touch: (parameterName: string) => void;
  // touchAll sets all parameters to touched to display their error messages
  touchAll: () => void;
}

const defaultValue: ValidationContextValue = {
  errors: {},
  setError: () => {},
  setErrors: () => {},
  touched: {},
  touch: () => {},
  touchAll: () => {},
};

const ValidationContext = createContext(defaultValue);

interface Props {
  initErrors: Record<string, string | null>;
  definitions: ParameterDefinition[];
  includeNameField?: boolean;
}

export const ValidationContextProvider: React.FC<Props> = ({
  children,
  initErrors,
  definitions,
  includeNameField,
}) => {
  const [errors, setErrors] =
    useState<Record<string, string | null>>(initErrors);
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  function setError(parameterName: string, error: string | null) {
    setErrors((prev) => ({
      ...prev,
      [parameterName]: error,
    }));
  }

  function touch(parameterName: string) {
    setTouched((prev) => ({
      ...prev,
      [parameterName]: true,
    }));
  }

  function touchAll() {
    const allFields: Record<string, boolean> = {};

    if (includeNameField) {
      allFields.name = true;
    }

    for (const def of definitions) {
      allFields[def.name] = true;
    }
    setTouched(allFields);
  }

  return (
    <ValidationContext.Provider
      value={{ errors, setError, touched, touch, touchAll, setErrors }}
    >
      {children}
    </ValidationContext.Provider>
  );
};

export function useValidationContext(): ValidationContextValue {
  return useContext(ValidationContext);
}

export function isValid(errors: Record<string, string | null>): boolean {
  for (const error of Object.values(errors)) {
    if (error != null) {
      return false;
    }
  }

  return true;
}
