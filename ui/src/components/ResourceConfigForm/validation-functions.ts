import { isEmpty } from "lodash";
import { ParameterDefinition } from "../../graphql/generated";

const REQUIRED_ERROR_MSG = "Required.";

export function validateStringField(
  value: string | null,
  required?: boolean
): string | null {
  if (required && isEmpty(value)) {
    return REQUIRED_ERROR_MSG;
  }
  return null;
}

export function validateStringsField(
  value: string[],
  required?: boolean
): string | null {
  if (required && isEmpty(value)) {
    return REQUIRED_ERROR_MSG;
  }

  return null;
}

export function validateYamlField(
  value: string | null,
  required?: boolean
): string | null {
  if (required && isEmpty(value)) {
    return REQUIRED_ERROR_MSG;
  }

  return null;
}

export function validateMapField(
  value: Record<string, string> | null,
  required?: boolean
): string | null {
  if (required) {
    if (value == null) {
      return REQUIRED_ERROR_MSG;
    }

    const entries = Object.entries(value);
    if (isEmpty(entries)) {
      return REQUIRED_ERROR_MSG;
    }

    let nonEmptyKeyFound = false;
    for (const entry of entries) {
      if (!isEmpty(entry[0])) {
        nonEmptyKeyFound = true;
        break;
      }
    }

    if (!nonEmptyKeyFound) {
      return REQUIRED_ERROR_MSG;
    }
  }

  return null;
}

export function validateIntField(
  definition: ParameterDefinition,
  value?: number
): string | null {
  if (definition.required && value == null) {
    return REQUIRED_ERROR_MSG;
  }

  return null;
}
