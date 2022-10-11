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
  activateControl: React.Dispatch<React.SetStateAction<boolean>>;
  buttonText: string;
}

export const AddResourceCard: React.FC<AddResourceCardProps> = ({
  activateControl,
  buttonText,
}) => {
  return (
    <Card
      className={styles["ui-control-card"]}
      onClick={() => activateControl(true)}
    >
      <CardActionArea>
        <CardContent>
          <Stack
            direction="row"
            justifyContent="center"
            alignItems="center"
            spacing={1}
          >
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
