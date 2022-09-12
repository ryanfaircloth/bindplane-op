import { DialogActions } from "@mui/material";
import { memo } from "react";

import styles from "./resource-dialog.module.scss";

interface Props {}

const ActionSectionComponent: React.FC<React.PropsWithChildren<Props>> = ({
  children,
}) => {
  return (
    <DialogActions
      classes={{
        root: styles.actions,
      }}
    >
      {children}
    </DialogActions>
  );
};

export const ActionsSection = memo(ActionSectionComponent);
