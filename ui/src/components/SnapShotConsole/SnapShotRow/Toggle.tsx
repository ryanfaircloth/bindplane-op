import { IconButton } from "@mui/material";
import { ChevronDown } from "../../Icons";

import styles from "../snap-shot-console.module.scss";

interface ToggleProps {
  open: boolean;
  onClick: () => void;
}

export const Toggle: React.FC<ToggleProps> = ({ open, onClick }) => {
  return (
    <IconButton size="small" className={styles.ch} onClick={onClick}>
      <div className={open ? styles.ch : styles["ch-open"]}>
        <ChevronDown className={styles.chevron} />
      </div>
    </IconButton>
  );
};
