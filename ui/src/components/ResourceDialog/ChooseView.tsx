import { Button, Stack } from "@mui/material";
import { DialogResource, ResourceType } from ".";
import { ResourceTypeButton } from "../ResourceTypeButton";
import { ActionsSection } from "./ActionSection";
import { ContentSection } from "./ContentSection";
import { useResourceDialog } from "./ResourceDialogContext";
import { TitleSection } from "./TitleSection";

interface ChooseViewProps {
  resources: DialogResource[];
  selected: ResourceType;
  kind: "source" | "destination";
  clearResource: () => void;
  handleSaveExisting: (r: DialogResource) => void;
  setCreateNew: (b: boolean) => void;
}

export const ChooseView: React.FC<ChooseViewProps> = ({
  resources,
  selected,
  handleSaveExisting,
  setCreateNew,
  clearResource,
}) => {
  const { onClose } = useResourceDialog();

  const matchingResources = resources?.filter(
    (r) => r.spec.type === selected!.metadata.name
  );

  return (
    <>
      <TitleSection title={"Choose Existing or Create New"} onClose={onClose} />

      <ContentSection>
        <Stack spacing={1}>
          {matchingResources?.map((resource) => {
            return (
              <ResourceTypeButton
                key={resource.metadata.name}
                icon={selected?.metadata.icon!}
                displayName={resource.metadata.name}
                onSelect={() => handleSaveExisting(resource)}
              />
            );
          })}
          <Button
            variant="contained"
            color="primary"
            onClick={() => setCreateNew(true)}
          >
            Create New
          </Button>
        </Stack>
      </ContentSection>

      <ActionsSection>
        <Button
          variant="contained"
          color="secondary"
          onClick={() => clearResource()}
        >
          Back
        </Button>
      </ActionsSection>
    </>
  );
};
