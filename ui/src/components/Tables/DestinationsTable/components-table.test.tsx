import { GridSelectionModel } from "@mui/x-data-grid";
import { render, screen, waitFor } from "@testing-library/react";
import { resourcesFromSelected } from "../../../pages/destinations/DestinationsPage";

import {
  Destination1,
  Destination2,
} from "../../ResourceConfigForm/__test__/dummyResources";
import { DestinationsDataGrid } from "./DestinationsDataGrid";

describe("resourcesFromSelected", () => {
  it("Destination|gcp", () => {
    const selected = ["Destination|gcp"];

    const want = [
      {
        kind: "Destination",
        metadata: {
          name: "gcp",
        },
      },
    ];

    const got = resourcesFromSelected(selected);

    expect(got).toEqual(want);
  });
});

describe("DestinationsDataGrid", () => {
  const destinationData = [Destination1, Destination2];

  it("renders without error", () => {
    render(
      <DestinationsDataGrid
        loading={false}
        queryData={{ destinations: destinationData }}
        setSelectionModel={() => {}}
        disableSelectionOnClick
        checkboxSelection
        onEditDestination={() => {}}
      />
    );
  });

  it("displays destinations", () => {
    render(
      <DestinationsDataGrid
        loading={false}
        queryData={{ destinations: destinationData }}
        setSelectionModel={() => {}}
        disableSelectionOnClick
        checkboxSelection
        onEditDestination={() => {}}
      />
    );

    screen.getByText(Destination1.metadata.name);
    screen.getByText(Destination2.metadata.name);
  });

  it("uses the expected GridSelectionModel", () => {
    function onDestinationsSelected(m: GridSelectionModel) {
      expect(m).toEqual([
        `Destination|${Destination1.metadata.name}`,
        `Destination|${Destination2.metadata.name}`,
      ]);
    }
    render(
      <DestinationsDataGrid
        loading={false}
        queryData={{ destinations: destinationData }}
        setSelectionModel={onDestinationsSelected}
        disableSelectionOnClick
        checkboxSelection
        onEditDestination={() => {}}
      />
    );

    screen.getByLabelText("Select all rows").click();
  });

  it("calls onEditDestination when destinations are selected", async () => {
    let editCalled: boolean = false;
    function onEditDestination() {
      editCalled = true;
    }
    render(
      <DestinationsDataGrid
        loading={false}
        queryData={{ destinations: destinationData }}
        setSelectionModel={() => {}}
        disableSelectionOnClick
        checkboxSelection
        onEditDestination={onEditDestination}
      />
    );

    screen.getByText(Destination1.metadata.name).click();

    await waitFor(() => expect(editCalled).toEqual(true));
  });
});
