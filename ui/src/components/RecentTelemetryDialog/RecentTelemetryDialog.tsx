import { Dialog, DialogContent, DialogProps } from "@mui/material";
import { isFunction } from "lodash";
import { PipelineType } from "../../graphql/generated";
import { SnapshotConsole } from "../SnapShotConsole/SnapShotConsole";

interface RecentTelemetryDialogProps extends DialogProps {
  agentID: string;
}

export const RecentTelemetryDialog: React.FC<RecentTelemetryDialogProps> = ({
  agentID,
  ...dialogProps
}) => {
  function handleClose() {
    isFunction(dialogProps.onClose) && dialogProps.onClose({}, "backdropClick");
  }

  return (
    <Dialog fullWidth maxWidth={"xl"} {...dialogProps}>
      <DialogContent>
        <SnapshotConsole
          onClose={handleClose}
          agentID={agentID}
          initialType={PipelineType.Logs}
        />
      </DialogContent>
    </Dialog>
  );
};
