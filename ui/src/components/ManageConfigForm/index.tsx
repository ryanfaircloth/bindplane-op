import { Button, Typography } from "@mui/material";
import { GetAgentAndConfigurationsQuery } from "../../graphql/generated";
import { classes } from "../../utils/styles";
import { patchConfigLabel } from "../../utils/patch-config-label";
import { Link } from "react-router-dom";
import { useSnackbar } from "notistack";
import { Config } from "./types";

import mixins from "../../styles/mixins.module.scss";
import styles from "./apply-config-form.module.scss";

interface ManageConfigFormProps {
  agent: NonNullable<GetAgentAndConfigurationsQuery["agent"]>;
  configurations: Config[];
  onImport: () => void;
  editing: boolean;
  setEditing: React.Dispatch<React.SetStateAction<boolean>>;
  selectedConfig: Config | undefined;
  setSelectedConfig: React.Dispatch<React.SetStateAction<Config | undefined>>;
}

export const ManageConfigForm: React.FC<ManageConfigFormProps> = ({
  agent,
  configurations,
  onImport,
  editing,
  setEditing,
  selectedConfig,
  setSelectedConfig,
}) => {
  const snackbar = useSnackbar();

  const configResourceName = agent?.configurationResource?.metadata.name;

  async function onApplyConfiguration() {
    try {
      await patchConfigLabel(agent.id, selectedConfig!.metadata.name);

      setEditing(false);
    } catch (err) {
      snackbar.enqueueSnackbar("Failed to patch label.", {
        color: "error",
        autoHideDuration: 5000,
      });
    }
  }

  function onCancelEdit() {
    setEditing(false);
    setSelectedConfig(
      configurations.find((c) => c.metadata.name === configResourceName)
    );
  }

  const ShowConfiguration: React.FC = () => {
    return (
      <>
        {configResourceName ? (
          <>
            <Link to={`/configurations/${configResourceName}`}>
              {configResourceName}
            </Link>
          </>
        ) : (
          <>
            <Typography variant={"body2"} classes={{ root: mixins["mb-2"] }}>
              This agent configuration is not currently managed by BindPlane.
              Click import to pull this agent&apos;s configuration in as a new
              managed configuration.
            </Typography>
          </>
        )}
      </>
    );
  };

  return (
    <>
      <div
        className={classes([
          mixins.flex,
          mixins["align-center"],
          mixins["mb-3"],
        ])}
      >
        <Typography variant="h6">
          Configuration - {editing ? <></> : <ShowConfiguration />}
        </Typography>

        <div className={styles["title-button-group"]}>
          {editing ? (
            <>
              <Button variant="outlined" onClick={onCancelEdit}>
                Cancel
              </Button>
              <Button
                variant="contained"
                onClick={onApplyConfiguration}
                classes={{ root: mixins["ml-2"] }}
              >
                Apply
              </Button>
            </>
          ) : (
            <>
              {configResourceName == null && (
                <>
                  <Button variant="contained" onClick={onImport}>
                    Import
                  </Button>
                </>
              )}
              {configurations.length > 0 && (
                <Button
                  className={classes([mixins["ml-2"], styles["choose-button"]])}
                  variant="text"
                  onClick={() => setEditing(true)}
                >
                  Choose Another Configuration
                </Button>
              )}
            </>
          )}
        </div>
      </div>
    </>
  );
};
