import { Typography, CardContent, Stack } from "@mui/material";

import styles from "./styles.module.scss";

export const CardMeasurementContent: React.FC = ({ children }) => {
  return children ? (
    <CardContent classes={{ root: styles.content }}>
      <Stack
        justifyContent="center"
        alignItems="center"
        width="100%"
        height="100%"
      >
        <Typography classes={{ root: styles.metric }}>{children}</Typography>
      </Stack>
    </CardContent>
  ) : (
    <></>
  );
};
