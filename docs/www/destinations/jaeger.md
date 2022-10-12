---
title: "Jaeger"
slug: "jaeger"
hidden: false
createdAt: "2022-08-05T14:38:28.713Z"
updatedAt: "2022-08-10T15:00:51.681Z"
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
|         |      | ✓      |

## Configuration Table

| Parameter            | Type     | Default | Description                                              |
| :------------------- | :------- | :------ | :------------------------------------------------------- |
| hostname\*           | `string` |         | Hostname or IP address of the Jaeger server.             |
| port                 | `int`    | 14250   | Port (gRPC) of the Jaeger server.                        |
| enable_tls           | `bool`   | false   | Whether or not to use TLS.                               |
| insecure_skip_verify | `bool`   | false   | Enable to skip TLS certificate verification.             |
| ca_file              | `string` |         | Certificate authority used to validate TLS certificates. |
| mutual_tls           | `bool`   | false   | Whether or not to use mutual TLS authentication.         |
| cert_file            | `string` |         | A TLS certificate used for client authentication.        |
| key_file             | `bool`   |         | A TLS private key used for client authentication.        |

<span style="color:red">\*_required field_</span>