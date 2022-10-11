import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { ResourceConfigForm } from ".";
import {
  GetProcessorTypeDocument,
  GetProcessorTypesDocument,
  ParameterDefinition,
  ParameterType,
  PipelineType,
  RelevantIfOperatorType,
} from "../../graphql/generated";
import {
  MapParamInput,
  ParameterInput,
  Tuple,
  tupleArrayToMap,
  valueToTupleArray,
} from "./ParameterInput";
import { satisfiesRelevantIf } from "./satisfiesRelevantIf";
import {
  ProcessorTypeMetric,
  ProcessorTypeSeverity,
  ResourceType1,
  ResourceType2,
} from "./__test__/dummyResources";
import renderer from "react-test-renderer";
import { MockedProvider, MockedResponse } from "@apollo/client/testing";
import { SnackbarProvider } from "notistack";

describe("satisfiesRelevantIf", () => {
  const formValues: { [key: string]: any } = {
    one: true,
    two: "bar",
    three: 25,
  };

  const param1: ParameterDefinition = {
    name: "string_name",
    label: "String Input",
    description: "Here is the description.",
    required: false,
    options: {},

    type: ParameterType.String,

    relevantIf: [
      {
        name: "one",
        operator: RelevantIfOperatorType.Equals,
        value: true,
      },
    ],
  };

  const param2: ParameterDefinition = {
    name: "string_name",
    label: "String Input",
    description: "Here is the description.",
    required: false,
    options: {},

    type: ParameterType.String,
    relevantIf: [
      {
        name: "one",
        operator: RelevantIfOperatorType.Equals,
        value: false,
      },
    ],

    default: "default-value",
  };

  it("param1 matches", () => {
    const got = satisfiesRelevantIf(formValues, param1);
    expect(got).toEqual(true);
  });
  it("param2 does not match", () => {
    const got = satisfiesRelevantIf(formValues, param2);
    expect(got).toEqual(false);
  });
});

describe("ResourceForm component", () => {
  it("does not display field if relevantIf isn't satisfied", () => {
    render(
      <ResourceConfigForm
        kind="destination"
        displayName={ResourceType2.metadata.displayName!}
        description={ResourceType2.metadata.description!}
        parameterDefinitions={ResourceType2.spec.parameters}
      />
    );
    const stringInput = screen.queryByText("String Input");
    expect(stringInput).toBeNull();
  });

  it("will render input when relevantIf is satisfied", () => {
    render(
      <ResourceConfigForm
        kind="destination"
        displayName={ResourceType2.metadata.displayName!}
        description={ResourceType2.metadata.description!}
        parameterDefinitions={ResourceType2.spec.parameters}
      />
    );
    let stringInput = screen.queryByLabelText("String Input");
    expect(stringInput).toBeNull();

    screen.getByRole("checkbox").click();
    stringInput = screen.getByLabelText("String Input");
    expect(stringInput).toBeInTheDocument();
  });

  it("maintains stateful formValues as correctType", async () => {
    const expectedValues = {
      name: "some-name",
      string_name: "default-value",
      string_required_name: "default-required-value",
      enum_name: "option1",
      strings_name: ["option1", "option2"],
      int_name: 25,
      bool_name: true,
    };

    let saveDone = false;

    let values: { [key: string]: any } = {};
    function onSave(formValues: { [key: string]: any }) {
      values = Object.assign({}, formValues);
      saveDone = true;
    }
    render(
      <ResourceConfigForm
        onSave={onSave}
        kind="source"
        displayName={ResourceType1.metadata.displayName!}
        description={ResourceType1.metadata.description!}
        parameterDefinitions={ResourceType1.spec.parameters}
        includeNameField
      />
    );

    fireEvent.change(screen.getByLabelText("Name *"), {
      target: { value: "some-name" },
    });

    screen.getByText("Save").click();

    await waitFor(() => saveDone === true);
    expect(values).toEqual(expectedValues);
  });

  it("maintains stateful formValues as correctType after change", async () => {
    const expectedValues = {
      name: "some-name",
      string_name: "default-value",
      string_required_name: "default-required-value",
      enum_name: "option1",
      strings_name: ["option1", "option2"],
      int_name: 50,
      bool_name: true,
    };

    let saveDone = false;

    let values: { [key: string]: any } = {};
    function onSave(formValues: { [key: string]: any }) {
      values = Object.assign({}, formValues);
      saveDone = true;
    }
    render(
      <ResourceConfigForm
        onSave={onSave}
        kind="source"
        displayName={ResourceType1.metadata.displayName!}
        description={ResourceType1.metadata.description!}
        parameterDefinitions={ResourceType1.spec.parameters}
        includeNameField
      />
    );

    fireEvent.change(screen.getByLabelText("Int Input"), {
      target: { value: 50 },
    });

    fireEvent.change(screen.getByLabelText("Name *"), {
      target: { value: "some-name" },
    });

    screen.getByText("Save").click();

    await waitFor(() => saveDone === true);
    expect(values).toEqual(expectedValues);
  });

  it("will not submit when the form has an error", async () => {
    var saveCalled: boolean = false;
    render(
      <ResourceConfigForm
        onSave={() => {
          saveCalled = true;
        }}
        kind="destination"
        displayName={ResourceType1.metadata.displayName!}
        description={ResourceType1.metadata.description!}
        parameterDefinitions={ResourceType1.spec.parameters}
        includeNameField
      />
    );

    screen.getByText("Save").click();
    screen.getByText("Required.");

    fireEvent.change(screen.getByLabelText("Name *"), {
      target: { value: "some-name" },
    });

    screen.getByText("Save").click();
    await waitFor(() => expect(saveCalled).toBeTruthy());
  });

  it("strings type parameter validation", async () => {
    var saveCalled: boolean = false;
    const label = "p1 label";
    const definitions: ParameterDefinition[] = [
      {
        name: "p1",
        label: label,
        type: ParameterType.Strings,
        required: true,
        description: "",
        options: {},
      },
    ];

    render(
      <ResourceConfigForm
        onSave={() => {
          saveCalled = true;
        }}
        kind="destination"
        displayName={"Title"}
        description={ResourceType1.metadata.description!}
        parameterDefinitions={definitions}
      />
    );

    screen.getByText("Save").click();
    screen.getByText("Required.");

    const input = screen.getByLabelText(`${label} *`);
    fireEvent.change(input, { target: { value: "/tmp.file.log" } });
    fireEvent.keyDown(input, { key: "Enter", code: "Enter", charCode: 13 });

    screen.getByText("Save").click();
    await waitFor(() => expect(saveCalled).toBeTruthy());
  });

  describe("map type parameter validation", () => {
    const mapParameter: ParameterDefinition = {
      required: true,
      label: "Label",
      description: "description",
      type: ParameterType.Map,
      name: "map_type_param",
      options: {},
    };

    it("blocks submission initially if required", () => {
      var saveCalled: boolean = false;

      render(
        <ResourceConfigForm
          onSave={() => {
            saveCalled = true;
          }}
          kind="destination"
          displayName={"Title"}
          description={ResourceType1.metadata.description!}
          parameterDefinitions={[mapParameter]}
        />
      );

      screen.getByText("Save").click();
      expect(saveCalled).toBeFalsy();
    });

    it("enables form submission when one non empty key is specified", async () => {
      var saveCalled: boolean = false;
      render(
        <ResourceConfigForm
          onSave={() => {
            saveCalled = true;
          }}
          kind="destination"
          displayName={"Title"}
          description={ResourceType1.metadata.description!}
          parameterDefinitions={[mapParameter]}
        />
      );

      const firstKey = screen.getAllByRole("textbox")[0];
      fireEvent.change(firstKey, { target: { value: "blah" } });
      fireEvent.blur(firstKey);

      screen.getByText("Save").click();
      await waitFor(() => expect(saveCalled).toBeTruthy());
    });
  });
});

describe("MapParamInput", () => {
  const mapParameter: ParameterDefinition = {
    required: true,
    label: "Label",
    description: "description",
    type: ParameterType.Map,
    name: "map_type_param",
    options: {},
  };

  it("valueToTupleArray", () => {
    const tests = [
      {
        value: {
          foo: "bar",
          blah: "baz",
        },
        expect: [
          ["foo", "bar"],
          ["blah", "baz"],
          ["", ""],
        ],
      },
      {
        value: null,
        expect: [["", ""]],
      },
      {
        value: {},
        expect: [["", ""]],
      },
    ];

    for (const test of tests) {
      const got = valueToTupleArray(test.value);
      expect(got).toEqual(test.expect);
    }
  });

  it("tupleArrayToMap", () => {
    const tests: { tuples: Tuple[]; expect: any }[] = [
      {
        tuples: [
          ["one", "two"],
          ["three", "four"],
        ],
        expect: {
          one: "two",
          three: "four",
        },
      },
      {
        tuples: [
          ["", "blah"],
          ["three", "four"],
          ["some", "thing"],
          ["", ""],
        ],
        expect: {
          three: "four",
          some: "thing",
        },
      },
      {
        tuples: [["", ""]],
        expect: {},
      },
    ];

    for (const test of tests) {
      const got = tupleArrayToMap(test.tuples);
      expect(got).toEqual(test.expect);
    }
  });

  it("renders correctly", () => {
    const tree = renderer.create(<ParameterInput definition={mapParameter} />);
    expect(tree).toMatchSnapshot();
  });

  it("renders map values", () => {
    const value: Record<string, string> = {
      one: "two",
      three: "four",
      five: "six",
    };
    render(<MapParamInput definition={mapParameter} value={value} />);
    screen.getByDisplayValue("one");
    screen.getByDisplayValue("two");
    screen.getByDisplayValue("three");
    screen.getByDisplayValue("four");
    screen.getByDisplayValue("five");
    screen.getByDisplayValue("six");
  });

  it("can add key value pairs", () => {
    render(<ParameterInput definition={mapParameter} />);

    screen.getByText("New Row").click();
    screen.getByText("New Row").click();

    // We should have three rows
    screen.getByTestId(`${mapParameter.name}-0-0-input`);
    screen.getByTestId(`${mapParameter.name}-1-0-input`);
    screen.getByTestId(`${mapParameter.name}-2-0-input`);
  });

  it("can delete key value pairs", () => {
    render(<ParameterInput definition={mapParameter} />);

    screen.getByText("New Row").click();
    screen.getByText("New Row").click();

    // We should have three rows
    screen.getByTestId(`${mapParameter.name}-0-0-input`);
    screen.getByTestId(`${mapParameter.name}-1-0-input`);
    screen.getByTestId(`${mapParameter.name}-2-0-input`);

    // Delete one
    screen.getByTestId(`${mapParameter.name}-1-remove-button`).click();

    // We should have two rows
    screen.getByTestId(`${mapParameter.name}-0-0-input`);
    screen.getByTestId(`${mapParameter.name}-1-0-input`);
  });
});

describe("EnumsParameter", () => {
  it("renders correctly", () => {
    const enumsParameter: ParameterDefinition = {
      required: true,
      label: "Label",
      description: "description",
      type: ParameterType.Enums,
      default: {},
      validValues: ["one", "two", "three", "four"],
      name: "enums_type_param",
      options: {},
    };

    const tree = renderer.create(
      <ParameterInput definition={enumsParameter} />
    );
    expect(tree).toMatchSnapshot();
  });
});

describe("YamlParameter", () => {
  it("renders correctly", () => {
    const yamlParameter: ParameterDefinition = {
      required: true,
      label: "Label",
      description: "description",
      type: ParameterType.Yaml,
      default: "",
      name: "yaml_type_param",
      options: {},
    };

    const tree = renderer.create(<ParameterInput definition={yamlParameter} />);
    expect(tree).toMatchSnapshot();
  });
});

describe("EnumParameter", () => {
  it("renders correctly", () => {
    const enumParam: ParameterDefinition = {
      required: true,
      label: "Label",
      description: "description",
      type: ParameterType.Enum,
      default: "",
      name: "yaml_type_param",
      validValues: ["one", "two", "three"],
      options: {},
    };

    const tree = renderer.create(<ParameterInput definition={enumParam} />);
    expect(tree).toMatchSnapshot();
  });

  it("renders creatable enum correctly", () => {
    const creatableEnum: ParameterDefinition = {
      required: true,
      label: "Label",
      description: "description",
      type: ParameterType.Enum,
      name: "yaml_type_param",
      validValues: ["one", "two", "three"],
      options: {},
    };

    const tree = renderer.create(<ParameterInput definition={creatableEnum} />);
    expect(tree).toMatchSnapshot();
  });
});

describe("Processors section", () => {
  it("displays a processors section with enableProcessors", () => {
    render(
      <ResourceConfigForm
        displayName={ResourceType2.metadata.displayName!}
        kind="source"
        description={ResourceType2.metadata.description!}
        parameterDefinitions={ResourceType2.spec.parameters}
        enableProcessors
      />
    );

    screen.getByText("Processors");
    screen.getByText(
      "Processors are run on data after it's received and prior to being sent to a destination. They will be executed in the order they appear below."
    );
  });

  const processorMocks: MockedResponse<Record<string, any>>[] = [
    {
      request: {
        query: GetProcessorTypesDocument,
      },
      result: () => {
        return {
          data: {
            processorTypes: [ProcessorTypeSeverity, ProcessorTypeMetric],
          },
        };
      },
    },
    {
      request: {
        query: GetProcessorTypeDocument,
        variables: {
          type: "metric_processor",
        },
      },
      result: () => {
        return {
          data: {
            processorType: ProcessorTypeMetric,
          },
        };
      },
    },
  ];

  it("can add a processor to a source", async () => {
    render(
      <SnackbarProvider>
        <MockedProvider mocks={processorMocks} addTypename={false}>
          <ResourceConfigForm
            kind="source"
            displayName={ResourceType2.metadata.displayName!}
            description={ResourceType2.metadata.description!}
            parameterDefinitions={ResourceType2.spec.parameters}
            enableProcessors
          />
        </MockedProvider>
      </SnackbarProvider>
    );

    // Click Add processor button
    screen.getByText("Add processor").click();

    // Displays processor types in a list
    await screen.findByText(ProcessorTypeMetric.metadata.displayName!);
    await screen.findByText(ProcessorTypeSeverity.metadata.displayName!);

    // Can select a processor type
    screen.getByText(ProcessorTypeMetric.metadata.displayName!).click();

    // On the configuration page
    await screen.findByText(
      `Add Processor: ${ProcessorTypeMetric.metadata.displayName}`
    );
    screen.getByText("Save").click();

    // Expect to see it in the main view
    await screen.findByText("Processors");
    await screen.findByText(ProcessorTypeMetric.metadata.displayName!);
  });

  it("filters processors in the select view by telemetry type", async () => {
    render(
      <SnackbarProvider>
        <MockedProvider mocks={processorMocks} addTypename={false}>
          <ResourceConfigForm
            kind="source"
            displayName={ResourceType2.metadata.displayName!}
            description={ResourceType2.metadata.description!}
            parameterDefinitions={ResourceType2.spec.parameters}
            enableProcessors
            telemetryTypes={[PipelineType.Metrics]}
          />
        </MockedProvider>
      </SnackbarProvider>
    );

    // Click Add processor button
    screen.getByText("Add processor").click();

    await screen.findByText(ProcessorTypeMetric.metadata.displayName!);
    expect(
      screen.queryByText(ProcessorTypeSeverity.metadata.displayName!)
    ).not.toBeInTheDocument();
  });
});
