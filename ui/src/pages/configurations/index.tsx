import { Button } from "@mui/material";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import { CardContainer } from "../../components/CardContainer";
import { ConfigurationsTable } from "../../components/Tables/ConfigurationTable";
import { PlusCircleIcon } from "../../components/Icons";
import { withRequireLogin } from "../../contexts/RequireLogin";
import { withNavBar } from "../../components/NavBar";

import mixins from "../../styles/mixins.module.scss";
import { GridSelectionModel } from "@mui/x-data-grid";

const ConfigurationsPageContent: React.FC = () => {
  // Selected is an array of names of configurations.
  const [selected, setSelected] = useState<GridSelectionModel>([]);
  return (
    <CardContainer>
      <Button
        component={Link}
        to="/configurations/new"
        variant="contained"
        classes={{ root: mixins["float-right"] }}
        startIcon={<PlusCircleIcon />}
      >
        New Configuration
      </Button>

      <ConfigurationsTable
        selected={selected}
        setSelected={setSelected}
        enableDelete={true}
      />
    </CardContainer>
  );
};

export const ConfigurationsPage = withRequireLogin(
  withNavBar(ConfigurationsPageContent)
);
