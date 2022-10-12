---
title: "Zipkin"
slug: "zipkin"
hidden: false
createdAt: "2022-08-05T14:47:01.044Z"
updatedAt: "2022-08-10T15:03:12.350Z"
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
|         |      | ✓      |

## Configuration Table

| Parameter            | Type     | Default         | Description                                                                                                                                    |
| :------------------- | :------- | :-------------- | :--------------------------------------------------------------------------------------------------------------------------------------------- |
| hostname\*           | `string` |                 | Hostname or IP address of the Zipkin server.                                                                                                   |
| port                 | `int`    | 14250           | Port (gRPC) of the Zipkin server.                                                                                                              |
| path\*               | `string` | "/api/v2/spans" | API path to send traces to.                                                                                                                    |
| enable_tls           | `bool`   | false           | Whether or not to use TLS.                                                                                                                     |
| insecure_skip_verify | `bool`   | false           | Enable to skip TLS certificate verification.                                                                                                   |
| ca_file              | `string` |                 | Certificate authority used to validate TLS certificates. Required only if the underlying operating system does not trust Zipkin's certificate. |
| mutual_tls           | `bool`   | false           | Whether or not to use mutual TLS authentication.                                                                                               |
| cert_file            | `string` |                 | A TLS certificate used for client authentication.                                                                                              |
| key_file             | `bool`   |                 | A TLS private key used for client authentication.                                                                                              |

<span style="color:red">\*_required field_</span>