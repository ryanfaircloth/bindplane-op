import {
  Card,
  CardActionArea,
  CardContent,
  Stack,
  Typography,
} from "@mui/material";
import { useSnackbar } from "notistack";
import { memo } from "react";
import { classes } from "../../utils/styles";
import { useNavigate } from "react-router-dom";
import { useGetDestinationWithTypeQuery } from "../../graphql/generated";

import styles from "./cards.module.scss";

interface ResourceDestinationCardProps {
  name: string;
  // disabled indicates that the card is not active and should be greyed out
  disabled?: boolean;
}

const OverviewDestinationCardComponent: React.FC<ResourceDestinationCardProps> =
  ({ name, disabled }) => {
    const { enqueueSnackbar } = useSnackbar();

    const { data } = useGetDestinationWithTypeQuery({
      variables: { name },
      fetchPolicy: "cache-and-network",
    });

    const navigate = useNavigate();

    // Loading
    if (data === undefined) {
      return null;
    }

    if (data.destinationWithType.destination == null) {
      enqueueSnackbar(`Could not retrieve destination ${name}.`, {
        variant: "error",
      });
      return null;
    }

    if (data.destinationWithType.destinationType == null) {
      enqueueSnackbar(
        `Could not retrieve destination type for destination ${name}.`,
        { variant: "error" }
      );
      return null;
    }

    return (
      <div
        className={classes([
          disabled ? styles.disabled : undefined,
          data.destinationWithType.destination?.spec.disabled
            ? styles.paused
            : undefined,
        ])}
      >
        <Card
          className={classes([
            styles["resource-card"],
            disabled ? styles.disabled : undefined,
            data.destinationWithType.destination?.spec.disabled
              ? styles.paused
              : undefined,
          ])}
          onClick={() => navigate("/destinations/")}
        >
          <CardActionArea className={styles.action}>
            <CardContent>
              <Stack alignItems="center">
                <span
                  className={styles.icon}
                  style={{
                    backgroundImage: `url(${data?.destinationWithType?.destinationType?.metadata.icon})`,
                  }}
                />
                <Typography component="div" fontWeight={600} gutterBottom>
                  {name}
                </Typography>
                {data.destinationWithType.destination?.spec.disabled && (
                  <Typography
                    component="div"
                    fontWeight={400}
                    fontSize={14}
                    variant="overline"
                  >
                    Paused
                  </Typography>
                )}
              </Stack>
            </CardContent>
          </CardActionArea>
        </Card>
      </div>
    );
  };

export const OverviewDestinationCard = memo(OverviewDestinationCardComponent);
