import { isEmpty } from "lodash";

type ResourceWithMetadata = {
  metadata: {
    name: string;
    displayName?: string | null;
  };
};

export function metadataSatisfiesSubstring(
  rt: ResourceWithMetadata,
  substring: string
) {
  substring = substring.toLowerCase()
  return isEmpty(substring)
    ? true
    : rt.metadata.name.toLowerCase().includes(substring) ||
        rt.metadata.displayName?.toLowerCase().includes(substring);
}
