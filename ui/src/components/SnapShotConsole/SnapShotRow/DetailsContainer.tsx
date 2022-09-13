import { Collapse, Stack } from "@mui/material";

interface DetailsContainerProps {
  open: boolean;
}

export const DetailsContainer: React.FC<DetailsContainerProps> = ({
  open,
  children,
}) => {
  return (
    <Collapse in={open}>
      <Stack paddingX={4} paddingY={2}>
        {open && children}
      </Stack>
    </Collapse>
  );
};
