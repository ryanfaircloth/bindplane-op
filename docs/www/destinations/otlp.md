---
title: "OpenTelemetry (OTLP)"
slug: "otlp"
hidden: false
createdAt: "2022-06-07T18:47:04.111Z"
updatedAt: "2022-08-10T15:02:21.179Z"
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Configuration Table

| Parameter            | Type     | Default | Description                                                                              |
| :------------------- | :------- | :------ | :--------------------------------------------------------------------------------------- |
| protocol             | `enum`   | grpc    | The OTLP protocol to use when sending OTLP telemetry. Valid values are `grpc` or `http`. |
| hostname\*           | `string` |         | Hostname or IP address to which the exporter is going to send OTLP data.                 |
| grpc_port            | `int`    | 4317    | TCP port to which the exporter is going to send OTLP data.                               |
| http_port            | `int`    | 4318    | TCP port to which the exporter is going to send OTLP data.                               |
| enable_tls           | `bool`   | false   | Whether or not to use TLS.                                                               |
| insecure_skip_verify | `bool`   | false   | Enable to skip TLS certificate verification.                                             |
| ca_file              | `bool`   |         | Certificate authority used to validate the database server's TLS certificate.            |
| cert_file            | `bool`   |         | A TLS certificate used for client authentication, if mutual TLS is enabled.              |
| key_file             | `bool`   |         | A TLS private key used for client authentication, if mutual TLS is enabled.              |