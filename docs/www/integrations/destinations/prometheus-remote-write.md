---
title: "Prometheus Remote Write"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "prometheus-remote-write"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       |      |        |

## Configuration Table

| Parameter                               | Type     | Default | Description                                                                                                                                                                            |
| :-------------------------------------- | :------- | :------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| hostname\*                              | `string` |         | The hostname or IP address for the remote write backend.                                                                                                                               |
| port\*                                  | `int`    | 9009    | The port remote write backend.                                                                                                                                                         |
| path\*                                  | `string` |         | The API Path of the remote write URL. Ex: `api/v1/metrics`.                                                                                                                           |
| headers                                 | `map`    |         | Additional headers to attach to each HTTP Request. The following headers cannot be changed: `Content-Encoding`, `Content-Type`, `X-Prometheus-Remote-Write-Version`, and `User-Agent`. |
| external_labels                         | `map`    |         | Label names and values to be attached as metric attributes.                                                                                                                            |
| namespace                               | `string` | ""      | Prefix to attach to each metric name.                                                                                                                                                  |
| enable_resource_to_telemetry_conversion | `bool`   | false   | When enabled will convert all resource attributes to metric attributes.                                                                                                                |
| enable_write_ahead_log                  | `bool`   | false   | Whether or not to enable a Write Ahead Log for the exporter.                                                                                                                           |
| wal_buffer_size                         | `int`    | 300     | Number of objects to store in Write Ahead Log before truncating. Applicable if `enable_write_ahead_log` is `true`.                                                                     |
| wal_truncate_frequency                  | `int`    | 60      | How often, in seconds, the Write Ahead Log should be truncated. Applicable if `enable_write_ahead_log` is `true`.                                                                      |
| enable_tls                              | `bool`   | false   | Whether or not to use TLS.                                                                                                                                                             |
| strict_tls_verify                       | `bool`   | false   | Strict TLS Certificate Verification.                                                                                                                                                   |
| ca_file                                 | `string` |         | Certificate authority used to validate TLS certificates. Not required if the collector's operating system already trusts the certificate authority.                                    |
| cert_file                               | `string` |         | A TLS certificate used for client authentication, if mutual TLS is enabled.                                                                                                            |
| key_file                                | `string` |         | A TLS private key used for client authentication, if mutual TLS is enabled.                                                                                                            |

<span style="color:red">\*_required field_</span>
