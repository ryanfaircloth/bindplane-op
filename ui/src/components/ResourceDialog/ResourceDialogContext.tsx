import { createContext, useContext } from "react";

interface ResourceDialogContextValue {
  onClose: () => void;
  purpose: "create" | "edit" | null;
}

const defaults: ResourceDialogContextValue = {
  onClose: () => {},
  purpose: null,
};

const ResourceDialogContext = createContext(defaults);

export const useResourceDialog = () => {
  return useContext(ResourceDialogContext);
};

interface Props {
  onClose: () => void;
  purpose: "create" | "edit";
}

export const ResourceDialogContextProvider: React.FC<Props> = ({
  children,
  onClose,
  purpose,
}) => {
  return (
    <ResourceDialogContext.Provider
      value={{
        purpose,
        onClose,
      }}
    >
      {children}
    </ResourceDialogContext.Provider>
  );
};
