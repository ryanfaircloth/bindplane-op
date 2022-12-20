import { Button, Stack } from "@mui/material";
import { ResourceConfiguration } from "../../../graphql/generated";
import { PlusCircleIcon } from "../../Icons";
import { InlineProcessorLabel } from "./InlineProcessorLabel";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import { useCallback, useState } from "react";

import mixins from "../../../styles/mixins.module.scss";

interface Props {
  processors: ResourceConfiguration[];
  onAddProcessor: () => void;
  onEditProcessor: (index: number) => void;
  onProcessorsChange: (ps: ResourceConfiguration[]) => void;
}

export const InlineProcessorContainer: React.FC<Props> = ({
  processors: processorsProp,
  onProcessorsChange,
  onAddProcessor,
  onEditProcessor,
}) => {
  // manage state internally
  const [processors, setProcessors] = useState(processorsProp);

  function handleDrop() {
    onProcessorsChange(processors);
  }

  const moveProcessor = useCallback(
    (dragIndex: number, hoverIndex: number) => {
      if (dragIndex === hoverIndex) {
        return;
      }

      const newProcessors = [...processors];

      const dragItem = newProcessors[dragIndex];
      const hoverItem = newProcessors[hoverIndex];

      // Swap places of dragItem and hoverItem in the array
      newProcessors[dragIndex] = hoverItem;
      newProcessors[hoverIndex] = dragItem;

      setProcessors(newProcessors);
    },
    [processors, setProcessors]
  );

  return (
    <>
      <DndProvider backend={HTML5Backend}>
        {processors.map((p, ix) => {
          return (
            <InlineProcessorLabel
              moveProcessor={moveProcessor}
              key={`${p.name}-${ix}`}
              processor={p}
              onEdit={() => onEditProcessor(ix)}
              index={ix}
              onDrop={handleDrop}
            />
          );
        })}
        <Stack alignItems="center">
          <Button
            variant="text"
            startIcon={<PlusCircleIcon />}
            classes={{ root: mixins["mb-2"] }}
            onClick={onAddProcessor}
          >
            Add processor
          </Button>
        </Stack>
      </DndProvider>
    </>
  );
};
