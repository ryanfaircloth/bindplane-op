import { Paper, PaperProps } from "@mui/material";
import React from "react";
import { classes } from "../../utils/styles";

import styles from "./card-container.module.scss";

export const CardContainer: React.FC<PaperProps> = ({ children, ...rest }) => {
  return (
    <Paper
      className={classes([styles.root, styles.bordered])}
      elevation={1}
      {...rest}
    >
      {children}
    </Paper>
  );
};

export const BorderlessCardContainer: React.FC = ({ children }) => {
  return (
    <Paper
      classes={{ root: styles.root }}
      elevation={1}
      sx={{ borderWidth: 0, backgroundColor: "#fcfcfc" }}
    >
      {children}
    </Paper>
  );
};
