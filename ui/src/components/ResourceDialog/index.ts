import {
  DestinationType,
  Maybe,
  Parameter,
  SourceType,
} from "../../graphql/generated";

export { NewResourceDialog } from "./NewResourceDialog";
export type ResourceType = DestinationType | SourceType;
export type DialogResource = {
  metadata: {
    name: string;
  };
  spec: { type: string; parameters?: Maybe<Parameter[]>, disabled?: boolean };
};
