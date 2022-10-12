---
title: "Nginx"
slug: "nginx"
hidden: false
createdAt: "2022-06-07T18:44:13.769Z"
updatedAt: "2022-08-10T15:34:48.091Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports nginx versions 1.18 and 1.20.

## Configuration Table

| Parameter            | Type      | Default                           | Description                                                                   |
| :------------------- | :-------- | :-------------------------------- | :---------------------------------------------------------------------------- |
| enable_metrics       | `bool`    | true                              | Enable to collect metrics.                                                    |
| endpoint*           | `string`  | "http://localhost:80/status"    | The endpoint of the NGINX server.                                             |
| collection_interval  | `int`     | 60                                | How often (seconds) to scrape for metrics.                                    |
| enable_tls           | `bool`    | false                             | Whether or not to use TLS.                                                    |
| insecure_skip_verify | `bool`    | false                             | Enable to skip TLS certificate verification.                                  |
| ca_file              | `string`  |                                   | Certificate authority used to validate the database server's TLS certificate. |
| cert_file            | `string`  |                                   | A TLS certificate used for client authentication, if mutual TLS is enabled.   |
| key_file             | `string`  |                                   | A TLS private key used for client authentication, if mutual TLS is enabled.   |
| enable_logs          | `bool`    | true                              | Enable to collect logs.                                                       |
| data_flow            | `enum`    | high                              | Enable high flow or reduced low flow.                                         |
| log_format           | `enum`    | default                           |                                                                               |
| access_log_paths     | `strings` | ` - "/var/log/nginx/access.log*"` | Path to NGINX access log file(s).                                             |
| error_log_paths      | `strings` | ` - "/var/log/nginx/error.log*"`  | Path to NGINX error log file(s).                                              |
| start_at             | `enum`    | end                               | Start reading file from 'beginning' or 'end'.                                 |

<span style="color:red">\*_required field_</span>