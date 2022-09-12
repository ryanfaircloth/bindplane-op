import { DialogResource, ResourceType } from ".";
import { FormValues, ResourceConfigForm } from "../ResourceConfigForm";

interface ConfigureViewProps {
  selected: ResourceType;
  kind: "source" | "destination";
  createNew: boolean;
  clearResource: () => void;
  handleSaveNew: (formValues: FormValues, selected: ResourceType) => void;
  resources: DialogResource[];
}

export const ConfigureView: React.FC<ConfigureViewProps> = ({
  selected,
  kind,
  createNew,
  clearResource,
  handleSaveNew,
  resources,
}) => {
  if (selected === null) {
    return <></>;
  }

  return (
    <ResourceConfigForm
      kind={kind}
      includeNameField={kind === "destination" && createNew}
      existingResourceNames={resources?.map((r) => r.metadata.name)}
      onBack={clearResource}
      onSave={(fv) => handleSaveNew(fv, selected)}
      displayName={selected.metadata.displayName ?? ""}
      description={selected.metadata.description ?? ""}
      parameterDefinitions={selected.spec.parameters ?? []}
    />
  );
};
