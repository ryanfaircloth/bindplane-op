import { ReactFlowProvider } from "react-flow-renderer";
import { withNavBar } from "../../components/NavBar";
import { withRequireLogin } from "../../contexts/RequireLogin";
import { OverviewGraph } from "./OverviewGraph";
import { OverviewPageProvider } from "./OverviewPageContext";

export const OverviewPage: React.FC = withRequireLogin(
  withNavBar((props) => {
    return (
      <OverviewPageProvider>
        <ReactFlowProvider>
          <OverviewGraph />
        </ReactFlowProvider>
      </OverviewPageProvider>
    );
  })
);
