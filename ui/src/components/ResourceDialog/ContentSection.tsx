import { DialogContent } from "@mui/material";
import { memo } from "react";

import styles from "./resource-dialog.module.scss";

const ContentSectionComponent: React.FC<React.PropsWithChildren<{}>> = ({
  children,
}) => {
  return (
    <DialogContent classes={{ root: styles.content }}>{children}</DialogContent>
  );
};

export const ContentSection = memo(ContentSectionComponent);
