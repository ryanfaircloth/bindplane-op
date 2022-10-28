import {
  Maybe,
  Parameter,
  ResourceConfiguration,
} from "../../graphql/generated";

export class BPResourceConfiguration implements ResourceConfiguration {
  name?: Maybe<string> | undefined;
  type?: Maybe<string> | undefined;
  parameters?: Maybe<Parameter[]> | undefined;
  processors?: Maybe<ResourceConfiguration[]> | undefined;
  disabled: boolean;
  constructor(rc?: ResourceConfiguration) {
    this.name = rc?.name;
    this.type = rc?.type;
    this.parameters = rc?.parameters;
    this.processors = rc?.processors;
    this.disabled = rc?.disabled ?? false;
  }

  isInline(): boolean {
    return this.name == null;
  }

  hasConfigurationParameters(): boolean {
    return this.parameters != null && this.parameters.length > 0;
  }

  // setParamsFromMap will set the parameters from Record<string, any>.
  // If the "name" key is specified it will set the name field of the ResourceConfiguration.
  // If the "processors" key is specified it will set the processors value.
  // It will not set undefined or null values to parameters.
  setParamsFromMap(map: Record<string, any>) {
    // Set name field if present
    if (map.name != null && map.name !== "") {
      this.name = map.name;
      delete map.name;
    }

    // Set processors field if present
    if (map.processors != null) {
      this.processors = map.processors;
      delete map.processors;
    }

    // Set the parameters only if their values are not nullish.
    const parameters = Object.entries(map).reduce<Parameter[]>(
      (params, [name, value]) => {
        if (value != null) {
          params.push({ name, value });
        }
        return params;
      },
      []
    );

    this.parameters = parameters;
  }
}
