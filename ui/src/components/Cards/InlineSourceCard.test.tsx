import { MockedProvider, MockedResponse } from "@apollo/client/testing";
import { cleanup, render, screen } from "@testing-library/react";
import { SnackbarProvider } from "notistack";
import { SourceTypeDocument } from "../../graphql/generated";
import { ConfigurationPageContextProvider } from "../../pages/configurations/configuration/ConfigurationPageContext";
import { MinimumRequiredConfig } from "../PipelineGraph/PipelineGraph";
import { InlineSourceCard } from "./InlineSourceCard";

describe("InlineSourceCard", () => {
  afterEach(() => {
    cleanup();
  });

  it("indicates whether the source is paused", async () => {
    const configMock: MockedResponse<Record<string, any>>[] = [
      {
        request: {
          operationName: "SourceType",
          query: SourceTypeDocument,
          variables: { name: "redis" },
        },
        result: () => {
          return {
            data: {
              sourceType: {
                apiVersion: "bindplane.observiq.com/v1",
                kind: "SourceType",
                metadata: {
                  id: "252423b1-e3de-4e35-b6b6-f1ecffa66106",
                  name: "redis",
                  displayName: "Redis",
                  description: "Collect metrics and logs from Redis.",
                  icon: "/icons/sources/redis.svg",
                  __typename: "Metadata",
                },
                spec: {
                  parameters: [
                    {
                      name: "enable_metrics",
                      label: "Enable Metrics",
                      description: "Enable to collect metrics.",
                      relevantIf: null,
                      documentation: null,
                      advancedConfig: false,
                      required: false,
                      type: "bool",
                      validValues: null,
                      default: true,
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "endpoint",
                      label: "Endpoint",
                      description: "The endpoint of the Redis server.",
                      relevantIf: [
                        {
                          name: "enable_metrics",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: false,
                      required: false,
                      type: "string",
                      validValues: null,
                      default: "localhost:6379",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "transport",
                      label: "Transport",
                      description:
                        "The transport protocol being used to connect to Redis.",
                      relevantIf: [
                        {
                          name: "enable_metrics",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: false,
                      required: false,
                      type: "enum",
                      validValues: ["tcp", "unix"],
                      default: "tcp",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "password",
                      label: "Password",
                      description:
                        "The password used to access the Redis instance; must match the password specified in the requirepass server configuration option.",
                      relevantIf: [
                        {
                          name: "enable_metrics",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "string",
                      validValues: null,
                      default: "",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "collection_interval",
                      label: "Collection Interval",
                      description: "How often (seconds) to scrape for metrics.",
                      relevantIf: [
                        {
                          name: "enable_metrics",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "int",
                      validValues: null,
                      default: 60,
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "enable_tls",
                      label: "Enable TLS",
                      description: "Whether or not to use TLS.",
                      relevantIf: [
                        {
                          name: "enable_metrics",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "bool",
                      validValues: null,
                      default: false,
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: true,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "insecure_skip_verify",
                      label: "Skip TLS Certificate Verification",
                      description:
                        "Enable to skip TLS certificate verification.",
                      relevantIf: [
                        {
                          name: "enable_tls",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "bool",
                      validValues: null,
                      default: false,
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "ca_file",
                      label: "TLS Certificate Authority File",
                      description:
                        "Certificate authority used to validate TLS certificates.",
                      relevantIf: [
                        {
                          name: "enable_tls",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "string",
                      validValues: null,
                      default: "",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "cert_file",
                      label: "Mutual TLS Client Certificate File",
                      description:
                        "A TLS certificate used for client authentication, if mutual TLS is enabled.",
                      relevantIf: [
                        {
                          name: "enable_tls",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "string",
                      validValues: null,
                      default: "",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "key_file",
                      label: "Mutual TLS Client Private Key File",
                      description:
                        "A TLS private key used for client authentication, if mutual TLS is enabled.",
                      relevantIf: [
                        {
                          name: "enable_tls",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "string",
                      validValues: null,
                      default: "",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "enable_logs",
                      label: "Enable Logs",
                      description: "Enable to collect logs.",
                      relevantIf: null,
                      documentation: null,
                      advancedConfig: false,
                      required: false,
                      type: "bool",
                      validValues: null,
                      default: true,
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "file_path",
                      label: "Log Paths",
                      description: "Path to Redis log file(s).",
                      relevantIf: [
                        {
                          name: "enable_logs",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "strings",
                      validValues: null,
                      default: [
                        "/var/log/redis/redis-server.log",
                        "/var/log/redis_6379.log",
                        "/var/log/redis/redis.log",
                        "/var/log/redis/default.log",
                        "/var/log/redis/redis_6379.log",
                      ],
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: 12,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                    {
                      name: "start_at",
                      label: "Start At",
                      description:
                        "Start reading logs from 'beginning' or 'end'.",
                      relevantIf: [
                        {
                          name: "enable_logs",
                          operator: "equals",
                          value: true,
                          __typename: "RelevantIfCondition",
                        },
                      ],
                      documentation: null,
                      advancedConfig: true,
                      required: false,
                      type: "enum",
                      validValues: ["beginning", "end"],
                      default: "end",
                      options: {
                        creatable: false,
                        multiline: false,
                        trackUnchecked: false,
                        sectionHeader: null,
                        gridColumns: null,
                        metricCategories: null,
                        __typename: "ParameterOptions",
                      },
                      __typename: "ParameterDefinition",
                    },
                  ],
                  supportedPlatforms: [],
                  version: "0.0.1",
                  telemetryTypes: ["logs", "metrics"],
                  __typename: "ResourceTypeSpec",
                },
                __typename: "SourceType",
              },
            },
          };
        },
      },
    ];

    const config: MinimumRequiredConfig = {
      metadata: {
        id: "test",
        name: "test",
        description: "",
        labels: {
          platform: "macos",
        },
        __typename: "Metadata",
      },
      spec: {
        raw: "",
        sources: [
          {
            type: "file",
            name: "",
            parameters: [
              {
                name: "file_path",
                value: ["/tmp/test.log"],
                __typename: "Parameter",
              },
              {
                name: "exclude_file_path",
                value: [],
                __typename: "Parameter",
              },
              {
                name: "log_type",
                value: "file",
                __typename: "Parameter",
              },
              {
                name: "parse_format",
                value: "none",
                __typename: "Parameter",
              },
              {
                name: "regex_pattern",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "multiline_line_start_pattern",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "encoding",
                value: "utf-8",
                __typename: "Parameter",
              },
              {
                name: "start_at",
                value: "end",
                __typename: "Parameter",
              },
            ],
            processors: null,
            disabled: false,
            __typename: "ResourceConfiguration",
          },
          {
            type: "redis",
            name: "",
            parameters: [
              {
                name: "enable_metrics",
                value: true,
                __typename: "Parameter",
              },
              {
                name: "endpoint",
                value: "localhost:6379",
                __typename: "Parameter",
              },
              {
                name: "transport",
                value: "tcp",
                __typename: "Parameter",
              },
              {
                name: "disable_metrics",
                value: [],
                __typename: "Parameter",
              },
              {
                name: "password",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "collection_interval",
                value: 10,
                __typename: "Parameter",
              },
              {
                name: "enable_tls",
                value: false,
                __typename: "Parameter",
              },
              {
                name: "insecure_skip_verify",
                value: false,
                __typename: "Parameter",
              },
              {
                name: "ca_file",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "cert_file",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "key_file",
                value: "",
                __typename: "Parameter",
              },
              {
                name: "enable_logs",
                value: true,
                __typename: "Parameter",
              },
              {
                name: "file_path",
                value: [
                  "/var/log/redis/redis-server.log",
                  "/var/log/redis_6379.log",
                  "/var/log/redis/redis.log",
                  "/var/log/redis/default.log",
                  "/var/log/redis/redis_6379.log",
                ],
                __typename: "Parameter",
              },
              {
                name: "start_at",
                value: "end",
                __typename: "Parameter",
              },
            ],
            processors: null,
            disabled: true,
            __typename: "ResourceConfiguration",
          },
        ],
        destinations: [],
        selector: {
          matchLabels: {
            configuration: "test",
          },
          __typename: "AgentSelector",
        },
        __typename: "ConfigurationSpec",
      },
      graph: {
        attributes: {
          activeTypeFlags: 1,
        },
        sources: [
          {
            id: "source/source0",
            type: "sourceNode",
            label: "file",
            attributes: {
              activeTypeFlags: 1,
              kind: "Source",
              resourceId: "source0",
              supportedTypeFlags: 1,
            },
            __typename: "Node",
          },
          {
            id: "source/source1",
            type: "sourceNode",
            label: "redis",
            attributes: {
              activeTypeFlags: 0,
              kind: "Source",
              resourceId: "source1",
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
        ],
        intermediates: [
          {
            id: "source/source0/processors",
            type: "processorNode",
            label: "Processors",
            attributes: {
              activeTypeFlags: 1,
              supportedTypeFlags: 1,
            },
            __typename: "Node",
          },
          {
            id: "source/source1/processors",
            type: "processorNode",
            label: "Processors",
            attributes: {
              activeTypeFlags: 0,
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
          {
            id: "destination/otlphttp/processors",
            type: "processorNode",
            label: "Processors",
            attributes: {
              activeTypeFlags: 0,
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
          {
            id: "destination/otlphttp-2/processors",
            type: "processorNode",
            label: "Processors",
            attributes: {
              activeTypeFlags: 1,
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
        ],
        targets: [
          {
            id: "destination/otlphttp",
            type: "destinationNode",
            label: "otlphttp",
            attributes: {
              activeTypeFlags: 0,
              isInline: false,
              kind: "Destination",
              resourceId: "otlphttp",
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
          {
            id: "destination/otlphttp-2",
            type: "destinationNode",
            label: "otlphttp-2",
            attributes: {
              activeTypeFlags: 1,
              isInline: false,
              kind: "Destination",
              resourceId: "otlphttp-2",
              supportedTypeFlags: 0,
            },
            __typename: "Node",
          },
        ],
        edges: [
          {
            id: "source/source0|source/source0/processors",
            source: "source/source0",
            target: "source/source0/processors",
            __typename: "Edge",
          },
          {
            id: "source/source1|source/source1/processors",
            source: "source/source1",
            target: "source/source1/processors",
            __typename: "Edge",
          },
          {
            id: "source/source0/processors|destination/otlphttp/processors",
            source: "source/source0/processors",
            target: "destination/otlphttp/processors",
            __typename: "Edge",
          },
          {
            id: "source/source1/processors|destination/otlphttp/processors",
            source: "source/source1/processors",
            target: "destination/otlphttp/processors",
            __typename: "Edge",
          },
          {
            id: "destination/otlphttp/processors|destination/otlphttp",
            source: "destination/otlphttp/processors",
            target: "destination/otlphttp",
            __typename: "Edge",
          },
          {
            id: "source/source0/processors|destination/otlphttp-2/processors",
            source: "source/source0/processors",
            target: "destination/otlphttp-2/processors",
            __typename: "Edge",
          },
          {
            id: "source/source1/processors|destination/otlphttp-2/processors",
            source: "source/source1/processors",
            target: "destination/otlphttp-2/processors",
            __typename: "Edge",
          },
          {
            id: "destination/otlphttp-2/processors|destination/otlphttp-2",
            source: "destination/otlphttp-2/processors",
            target: "destination/otlphttp-2",
            __typename: "Edge",
          },
        ],
        __typename: "Graph",
      },
      __typename: "Configuration",
    };

    render(
      <ConfigurationPageContextProvider
        configuration={config}
        refetchConfiguration={jest.fn()}
        setAddDestDialogOpen={jest.fn()}
        setAddSourceDialogOpen={jest.fn()}
      >
        <SnackbarProvider>
          <MockedProvider mocks={configMock} addTypename={false}>
            <InlineSourceCard
              id="source1"
              configuration={config}
              refetchConfiguration={() => {}}
            />
          </MockedProvider>
        </SnackbarProvider>
      </ConfigurationPageContextProvider>
    );

    const sourceBtn = await screen.findByTestId("source-card-source1");
    expect(sourceBtn).toBeEnabled();
    expect(sourceBtn).toHaveTextContent("Paused");
  });
});
