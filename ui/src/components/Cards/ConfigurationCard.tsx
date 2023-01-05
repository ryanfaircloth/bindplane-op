import {
  Card,
  CardActionArea,
  CardContent,
  Chip,
  Stack,
  Typography,
} from "@mui/material";
import { useLocation, useNavigate } from "react-router-dom";
import { CardMeasurementContent } from "../CardMeasurementContent/CardMeasurementContent";
import { SlidersIcon } from "../Icons";

import styles from "./cards.module.scss";

interface ConfigurationCardProps {
  id: string;
  label: string;
  attributes: Record<string, any>;
  metric: string;
  disabled?: boolean;
}

export const ConfigurationCard: React.FC<ConfigurationCardProps> = ({
  id,
  label,
  attributes,
  metric,
  disabled,
}) => {
  const navigate = useNavigate();
  const location = useLocation();
  const isEverything = id === "everything/configuration";
  const configurationURL = isEverything
    ? "/configurations"
    : `/configurations/${id.split("/").pop()}`;
  const agentCount = attributes["agentCount"] ?? 0;

  const cardLabel = isEverything ? "Other Configurations" : label;

  return (
    <div className={disabled ? styles.disabled : undefined}>
      <Card className={styles["resource-card"]}>
        <CardActionArea
          classes={{ root: styles.action }}
          onClick={() =>
            navigate({ pathname: configurationURL, search: location.search })
          }
        >
          <CardContent>
            <Stack justifyContent="center" alignItems="center" spacing={2}>
              {isEverything ? (
                <Stack direction="row" spacing={2}>
                  <SlidersIcon color="#b3b3b3" height="20px" width="20px" />
                  <SlidersIcon color="#b3b3b3" height="20px" width="20px" />
                  <SlidersIcon color="#b3b3b3" height="20px" width="20px" />
                </Stack>
              ) : (
                <SlidersIcon height="40px" width="40px" />
              )}
              <Typography align="center" fontWeight={600}>
                {cardLabel}
              </Typography>
            </Stack>
          </CardContent>
        </CardActionArea>
      </Card>
      <Chip
        classes={{
          root: styles["overview-count-chip"],
          label: styles["overview-count-chip-label"],
        }}
        size="small"
        label={formatAgentCount(agentCount)}
      />
      <CardMeasurementContent>{metric}</CardMeasurementContent>
    </div>
  );
};
export function formatAgentCount(agentCount: number): string {
  switch (agentCount) {
    case 0:
      return "";
    case 1:
      return "1 agent";
    default:
      return `${agentCount} agents`;
  }
}
