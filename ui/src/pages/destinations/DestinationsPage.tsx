import { gql } from "@apollo/client";
import { Typography, FormControl, Button } from "@mui/material";
import { GridSelectionModel } from "@mui/x-data-grid";
import { useSnackbar } from "notistack";
import { useState, useEffect } from "react";
import { CardContainer } from "../../components/CardContainer";
import { ConfirmDeleteResourceDialog } from "../../components/ConfirmDeleteResourceDialog";
import { withNavBar } from "../../components/NavBar";
import { DestinationsDataGrid } from "../../components/Tables/DestinationsTable/DestinationsDataGrid";
import { EditDestinationDialog } from "../../components/Tables/DestinationsTable/EditDestinationDialog";
import { FailedDeleteDialog } from "../../components/Tables/DestinationsTable/FailedDeleteDialog";
import { withRequireLogin } from "../../contexts/RequireLogin";
import { useDestinationsQuery } from "../../graphql/generated";
import { ResourceStatus, ResourceKind } from "../../types/resources";
import {
  deleteResources,
  MinimumDeleteResource,
} from "../../utils/rest/delete-resources";

import mixins from "../../styles/mixins.module.scss";

gql`
  query Destinations {
    destinations {
      kind
      metadata {
        name
      }
      spec {
        type
      }
    }
  }
`;

export const DestinationsPageContent: React.FC = () => {
  // Selected is an array of names of destinations in the form
  // <Kind>|<Name>
  const [selected, setSelected] = useState<GridSelectionModel>([]);

  // Used to control the delete confirmation modal.
  const [open, setOpen] = useState<boolean>(false);

  const [editingDestination, setEditingDestination] =
    useState<string | null>(null);

  const [failedDeletes, setFailedDeletes] = useState<ResourceStatus[]>([]);
  const [failedDeletesOpen, setFailedDeletesOpen] = useState(false);

  const { enqueueSnackbar } = useSnackbar();

  const { data, loading, refetch, error } = useDestinationsQuery({
    fetchPolicy: "network-only",
  });

  useEffect(() => {
    if (error != null) {
      enqueueSnackbar("There was an error retrieving data.", {
        variant: "error",
      });
    }
  }, [enqueueSnackbar, error]);

  useEffect(() => {
    if (failedDeletes.length > 0) {
      setFailedDeletesOpen(true);
    }
  }, [failedDeletes, setFailedDeletesOpen]);

  function onDestinationsSelected(s: GridSelectionModel) {
    setSelected(s);
  }

  function onAcknowledge() {
    setFailedDeletesOpen(false);
  }

  function handleEditSaveSuccess() {
    refetch();
    setEditingDestination(null);
  }

  async function deleteDestinations() {
    try {
      const items = resourcesFromSelected(selected);
      const { updates } = await deleteResources(items);
      setOpen(false);

      const failures = updates.filter((u) => u.status !== "deleted");
      setFailedDeletes(failures);

      refetch();
    } catch (err) {
      console.error(err);
      enqueueSnackbar("Failed to delete components.", { variant: "error" });
    }
  }

  return (
    <CardContainer>
      <div className={mixins.flex}>
        <Typography variant="h5" className={mixins["mb-5"]}>
          Destinations
        </Typography>
        {selected.length > 0 && (
          <FormControl classes={{ root: mixins["ml-5"] }}>
            <Button
              variant="contained"
              color="error"
              onClick={() => setOpen(true)}
            >
              Delete {selected.length} Component
              {selected.length > 1 && "s"}
            </Button>
          </FormControl>
        )}
      </div>
      <DestinationsDataGrid
        loading={loading}
        queryData={data ?? { destinations: [] }}
        onDestinationsSelected={onDestinationsSelected}
        disableSelectionOnClick
        checkboxSelection
        onEditDestination={(name: string) => setEditingDestination(name)}
      />
      <ConfirmDeleteResourceDialog
        open={open}
        onClose={() => setOpen(false)}
        onDelete={deleteDestinations}
        onCancel={() => setOpen(false)}
        action={"delete"}
      >
        <Typography>
          Are you sure you want to delete {selected.length} component
          {selected.length > 1 && "s"}?
        </Typography>
      </ConfirmDeleteResourceDialog>

      <FailedDeleteDialog
        open={failedDeletesOpen}
        failures={failedDeletes}
        onAcknowledge={onAcknowledge}
        onClose={() => {}}
      />

      {editingDestination && (
        <EditDestinationDialog
          name={editingDestination}
          onCancel={() => setEditingDestination(null)}
          onSaveSuccess={handleEditSaveSuccess}
        />
      )}
    </CardContainer>
  );
};

export function resourcesFromSelected(
  selected: GridSelectionModel
): MinimumDeleteResource[] {
  return selected.reduce<MinimumDeleteResource[]>((prev, cur) => {
    if (typeof cur !== "string") {
      console.error(`Unexpected type for GridRowId: ${typeof cur}"`);
      return prev;
    }
    const [kind, name] = cur.split("|");

    if (kind == null || name == null) {
      console.error(`Malformed grid row ID: ${cur}`);
      return prev;
    }

    prev.push({ kind: ResourceKind.DESTINATION, metadata: { name } });
    return prev;
  }, []);
}

export const DestinationsPage = withRequireLogin(
  withNavBar(DestinationsPageContent)
);
