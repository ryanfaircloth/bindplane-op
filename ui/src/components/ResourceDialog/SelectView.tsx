import { useMemo, useState } from "react";
import { DialogResource, ResourceType } from ".";
import { metadataSatisfiesSubstring } from "../../utils/metadata-satisfies-substring";
import {
  ResourceTypeButton,
  ResourceTypeButtonContainer,
} from "../ResourceTypeButton";
import { ContentSection } from "./ContentSection";
import { useResourceDialog } from "./ResourceDialogContext";
import { TitleSection } from "./TitleSection";

interface SelectViewProps {
  resourceTypes: ResourceType[];
  resources: DialogResource[];
  setSelected: (t: ResourceType) => void;
  setCreateNew: (b: boolean) => void;
  kind: "source" | "destination";
}

export const SelectView: React.FC<SelectViewProps> = ({
  resourceTypes,
  resources,
  setSelected,
  setCreateNew,
  kind,
}) => {
  const [resourceSearchValue, setResourceSearch] = useState("");
  const { onClose } = useResourceDialog();

  const sortedResourceTypes = useMemo(() => {
    const copy = resourceTypes.slice();
    return copy.sort((a, b) =>
      a.metadata
        .displayName!.toLowerCase()
        .localeCompare(b.metadata.displayName!.toLowerCase())
    );
  }, [resourceTypes]);

  return (
    <>
      <TitleSection
        title={kind === "destination" ? "Add Destination" : "Add Source"}
        onClose={onClose}
      />

      <ContentSection>
        <ResourceTypeButtonContainer
          onSearchChange={(v: string) => setResourceSearch(v)}
        >
          {sortedResourceTypes
            // Filter resource types by the resourceSearchValue
            .filter((rt) => metadataSatisfiesSubstring(rt, resourceSearchValue))
            // map the results to resource buttons
            .map((resourceType) => {
              const matchingResourcesExist = resources?.some(
                (resource) => resource.spec.type === resourceType.metadata.name
              );

              // Either we send the directly to the form if there are no existing resources
              // of that type, or we send them to the Choose View by just setting the selected.
              function onSelect() {
                setSelected(resourceType);
                if (!matchingResourcesExist) {
                  setCreateNew(true);
                }
              }
              return (
                <ResourceTypeButton
                  key={resourceType.metadata.name}
                  icon={resourceType.metadata.icon!}
                  displayName={resourceType.metadata.displayName!}
                  onSelect={onSelect}
                  telemetryTypes={resourceType.spec.telemetryTypes}
                />
              );
            })}
        </ResourceTypeButtonContainer>
      </ContentSection>
    </>
  );
};
