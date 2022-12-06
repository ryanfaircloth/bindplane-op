import { intersection, isArray, isEqual } from "lodash";
import { ParameterDefinition, RelevantIfOperatorType } from "../../graphql/generated";

/**
 * Check if form values satisfy the relevantIf conditions of a ParameterDefinition,
 * if any.
 */
export function satisfiesRelevantIf(
  formValues: { [name: string]: any },
  definition: ParameterDefinition
): boolean {
  const relevantIf = definition.relevantIf;
  if (relevantIf == null) {
    return true;
  }

  for (const condition of relevantIf) {
    switch (condition.operator) {
      case RelevantIfOperatorType.Equals:
        if (!isEqual(formValues[condition.name], condition.value)) {
          return false;
        }
        break;

      case RelevantIfOperatorType.NotEquals:
        if (isEqual(formValues[condition.name], condition.value)) {
          return false;
        }
        break;

        case RelevantIfOperatorType.ContainsAny:
          const value = formValues[condition.name];
          if (isArray(value)) {
            if (intersection(value, condition.value).length === 0) {
              return false;
            }
          }
          break;
    }
  }

  return true;
}
