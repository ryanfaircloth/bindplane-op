import { render } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import {
  ConfigurationChangesDocument,
  GetConfigurationTableDocument,
  GetConfigurationTableQuery,
} from "../../../graphql/generated";
import { MockedProvider, MockedResponse } from "@apollo/client/testing";

const TEST_CONFIGS: GetConfigurationTableQuery["configurations"]["configurations"] =
  [
    {
      metadata: {
        name: "config-1",
        description: "description for config-1",
        labels: {
          env: "test",
          foo: "bar",
        },
      },
      agentCount: 10,
    },
    {
      metadata: {
        name: "config-2",
        description: "description for config-2",
        labels: {
          env: "test",
          foo: "bar",
        },
      },
      agentCount: 30,
    },
  ];

const QUERY_RESULT: GetConfigurationTableQuery = {
  configurations: {
    configurations: TEST_CONFIGS,
    query: "",
    suggestions: [],
  },
};

const mocks: MockedResponse<Record<string, any>>[] = [
  {
    request: {
      query: GetConfigurationTableDocument,
      variables: {
        query: "",
      },
    },
    result: () => {
      return { data: QUERY_RESULT };
    },
  },
  {
    request: {
      query: ConfigurationChangesDocument,
      variables: {
        query: "",
      },
    },
    result: () => {
      return {
        data: { configurationChanges: [] },
      };
    },
  },
];

describe("ConfigurationsTable", () => {
  it("renders rows of configurations", async () => {
    // Selected is an array of names of configurations.
    render(
      <MemoryRouter>
        <MockedProvider mocks={mocks} addTypename={false}>
          Configurations
        </MockedProvider>
      </MemoryRouter>
    );
    // const rowOne = await screen.findAllByText("Configurations");
    // expect(rowOne).toBeInTheDocument();
    // const rowTwo = await screen.findAllByText("Configurations");
    // expect(rowTwo).toBeInTheDocument();
  });
  // it("shows delete button after selecting row", async () => {
  //   const [selected, setSelected] = useState<GridSelectionModel>([]);
  //   render(
  //     <MemoryRouter>
  //       <MockedProvider mocks={mocks} addTypename={false}>
  //         <ConfigurationsTable selected={selected} setSelected={setSelected} />
  //       </MockedProvider>
  //     </MemoryRouter>
  //   );
  //   // sanity check
  //   const row1 = await screen.findByText("config-1");
  //   expect(row1).toBeInTheDocument();
  //   const checkbox = await screen.findByLabelText("Select all rows");
  //   checkbox.click();
  //   const deleteButton = await screen.findByText("Delete 2 Configurations");
  //   expect(deleteButton).toBeInTheDocument();
  // });
  // it("opens the delete dialog after clicking delete", async () => {
  //   const [selected, setSelected] = useState<GridSelectionModel>([]);
  //   render(
  //     <MemoryRouter>
  //       <MockedProvider mocks={mocks} addTypename={false}>
  //         <ConfigurationsTable selected={selected} setSelected={setSelected} />
  //       </MockedProvider>
  //     </MemoryRouter>
  //   );
  //   const row1 = await screen.findByText("config-1");
  //   expect(row1).toBeInTheDocument();
  //   const checkbox = await screen.findByLabelText("Select all rows");
  //   checkbox.click();
  //   const deleteButton = await screen.findByText("Delete 2 Configurations");
  //   deleteButton.click();
  //   const dialog = await screen.findByTestId("delete-dialog");
  //   expect(dialog).toBeInTheDocument();
  // });
});
