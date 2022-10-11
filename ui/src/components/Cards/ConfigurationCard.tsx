import {
  Card,
  CardActionArea,
  CardContent,
  Chip,
  Stack,
  Typography,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { CardMeasurementContent } from "../CardMeasurementContent/CardMeasurementContent";
import { SlidersIcon } from "../Icons";

import styles from "./cards.module.scss";

interface ConfigurationCardProps {
  id: string;
  label: string;
  attributes: Record<string, any>;
  metric: string;
  disabled?: boolean;
  connectedNodesAndEdges: string[];
}

export const ConfigurationCard: React.FC<ConfigurationCardProps> = ({
  id,
  label,
  attributes,
  metric,
  disabled,
  connectedNodesAndEdges,
}) => {
  const navigate = useNavigate();
  const configurationURL = `/configurations/${id.split("/").pop()}`;
  const agentCount = attributes["agentCount"] ?? 0;

  return (
    <div className={disabled ? styles.disabled : undefined}>
      <Card className={styles["resource-card"]}>
        <CardActionArea
          classes={{ root: styles.action }}
          onClick={() => navigate(configurationURL)}
        >
          <CardContent>
            <Stack justifyContent="center" alignItems="center">
              <SlidersIcon height="40px" width="40px" />
              <Typography fontWeight={600}>{label}</Typography>
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
