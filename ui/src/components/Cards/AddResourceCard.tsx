import {
  Card,
  CardActionArea,
  CardContent,
  Stack,
  Typography,
} from "@mui/material";
import React from "react";
import { PlusCircleIcon } from "../Icons";

import styles from "./cards.module.scss";

interface AddResourceCardProps {
  onClick: React.Dispatch<React.SetStateAction<boolean>>;
  buttonText: string;
}

export const AddResourceCard: React.FC<AddResourceCardProps> = ({
  onClick,
  buttonText,
}) => {
  return (
    <Card className={styles["ui-control-card"]} onClick={() => onClick(true)}>
      <CardActionArea>
        <CardContent>
          <Stack justifyContent="center" alignItems="center" gap={1}>
            <PlusCircleIcon className={styles["ui-control-icon"]} />
            <Typography className={styles["ui-control-text"]}>
              {buttonText}
            </Typography>
          </Stack>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};
