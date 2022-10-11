import { memo } from "react";
import { AddResourceCard } from "../../Cards/AddResourceCard";

function UIControlNode({
  data,
}: {
  data: {
    activateControl: React.Dispatch<React.SetStateAction<boolean>>;
    buttonText: string;
  };
}) {
  return (
    <AddResourceCard
      activateControl={data.activateControl}
      buttonText={data.buttonText}
    />
  );
}

export default memo(UIControlNode);
