import { Typography } from "@mui/material";
import { version } from "./utils";

import styles from "./version.module.scss";

// Version displays the server version received from the /version endpoint.
export const BindplaneVersion: React.FC = () => {
  return (
    <Typography
      variant="body2"
      fontWeight={300}
      classes={{ root: styles.root }}
    >
      BindPlane OP {version()}
    </Typography>
  );
};
