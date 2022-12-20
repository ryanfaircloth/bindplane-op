import { DialogActions, Button } from "@mui/material";
import { useMemo } from "react";
import { ResourceConfiguration } from "../../../graphql/generated";
import { InlineProcessorContainer } from "./InlineProcessorContainer";
import { ContentSection } from "../ContentSection";
import { TitleSection } from "../TitleSection";
import { usePipelineGraph } from "../../PipelineGraph/PipelineGraphContext";

interface AllProcessorsProps {
  // parentDisplayName is the name of the source or destination
  // for which we're editing processors
  parentDisplayName: string;

  processors: ResourceConfiguration[];
  onAddProcessor: () => void;
  onEditProcessor: (index: number) => void;
  onClose: () => void;
  onSave: () => void;
  onProcessorsChange: (pt: ResourceConfiguration[]) => void;
}

/**
 * AllProcessorsView shows the initial view of the processors dialog, which is a list of all processors,
 * with the ability to add a new processor, reorder processors, and select a processor to edit or delete.
 */
export const AllProcessorsView: React.FC<AllProcessorsProps> = ({
  parentDisplayName,
  processors,
  onAddProcessor,
  onEditProcessor,
  onClose,
  onSave,
  onProcessorsChange,
}) => {
  const { editProcessorsInfo } = usePipelineGraph();

  const title = useMemo(() => {
    const kind =
      editProcessorsInfo?.resourceType === "source" ? "Source" : "Destination";
    return `${kind} ${parentDisplayName}: Processors`;
  }, [editProcessorsInfo?.resourceType, parentDisplayName]);

  const description =
    "Processors are run on data after it's received and prior to being sent to a destination. They will be executed in the order they appear below.";

  return (
    <>
      <TitleSection title={title} description={description} onClose={onClose} />
      <ContentSection>
        <InlineProcessorContainer
          processors={processors}
          onAddProcessor={onAddProcessor}
          onEditProcessor={onEditProcessor}
          onProcessorsChange={onProcessorsChange}
        />
      </ContentSection>

      <DialogActions>
        <Button variant="contained" onClick={onSave}>
          Save
        </Button>
      </DialogActions>
    </>
  );
};
