---
title: "OpenTelemetry (OTLP)"
category: 633dd7654359a20031653089
slug: "opentelemetry"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    | ✓      |
| Windows  | ✓       | ✓    | ✓      |
| macOS    | ✓       | ✓    | ✓      |

## Configuration Table

| Parameter            | Type     | Default   | Description                                                                                                                           |
| :------------------- | :------- | :-------- | :------------------------------------------------------------------------------------------------------------------------------------ |
| listen_address       | `string` | "0.0.0.0" | The IP address to listen on.                                                                                                          |
| grpc_port            | `int`    | 4317      | TCP port to receive OTLP telemetry using the gRPC protocol. The port used must not be the same as the HTTP port. Set to 0 to disable. |
| http_port            | `int`    | 4318      | TCP port to receive OTLP telemetry using the HTTP protocol. The port used must not be the same as the gRPC port. Set to 0 to disable. |
| enable_tls           | `bool`   | false     | Whether or not to use TLS.                                                                                                            |
| insecure_skip_verify | `bool`   | false     | Enable to skip TLS certificate verification.                                                                                          |
| ca_file              | `string` |           | Certificate authority used to validate the database server's TLS certificate.                                                         |
| cert_file            | `string` |           | A TLS certificate used for client authentication, if mutual TLS is enabled.                                                           |
| key_file             | `string` |           | A TLS private key used for client authentication, if mutual TLS is enabled.                                                           |