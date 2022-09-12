import { DialogTitle, IconButton, Stack, Typography } from "@mui/material";
import { memo } from "react";
import { XIcon } from "../Icons";

import styles from "./resource-dialog.module.scss";

interface Props {
  description?: string;
  title?: string;
  onClose: () => void;
}

const TitleSectionComponent: React.FC<Props> = ({
  description,
  title,
  onClose,
}) => {
  return (
    <DialogTitle
      classes={{
        root: styles.title,
      }}
    >
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Stack>
          <Typography fontSize={28} fontWeight={600}>
            {title}
          </Typography>

          <Typography fontSize={18}>{description}</Typography>
        </Stack>

        <IconButton className={styles.close} onClick={onClose}>
          <XIcon strokeWidth={"3"} width={"28"} />
        </IconButton>
      </Stack>
    </DialogTitle>
  );
};

export const TitleSection = memo(TitleSectionComponent);
