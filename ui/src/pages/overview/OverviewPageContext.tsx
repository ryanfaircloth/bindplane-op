import { createContext, useContext, useState } from "react";

export interface OverviewPageContextValue {
  hoveredSet: string[];
  setHoveredNodeAndEdgeSet: React.Dispatch<React.SetStateAction<string[]>>;
}

const defaultContext: OverviewPageContextValue = {
  hoveredSet: [],
  setHoveredNodeAndEdgeSet: () => {},
};

const OverviewPageContext = createContext(defaultContext);

export const OverviewPageProvider: React.FC = ({ children }) => {
  // state for knowing which node is being hovered over
  const [hoveredSet, setHoveredNodeAndEdgeSet] = useState<string[]>([]);
  return (
    <OverviewPageContext.Provider
      value={{ setHoveredNodeAndEdgeSet, hoveredSet }}
    >
      {children}
    </OverviewPageContext.Provider>
  );
};

export function useOverviewPage(): OverviewPageContextValue {
  return useContext(OverviewPageContext);
}
