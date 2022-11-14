---
title: "OpenTelemetry (OTLP)"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "otlp"
hidden: false
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
