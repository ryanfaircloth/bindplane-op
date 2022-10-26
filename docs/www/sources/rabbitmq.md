---
title: "RabbitMQ"
category: 633dd7654359a20031653089
slug: "rabbitmq"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

Supports RabbitMQ versions 3.8 and 3.9.

The RabbitMQ Management Plugin must be enabled by following the [official instructions](https://www.rabbitmq.com/management.html#getting-started).

Also, a user with at least [monitoring](https://www.rabbitmq.com/management.html#permissions) level permissions must be used for monitoring.

## Configuration Table

| Parameter            | Type      | Default                               | Description                                                                   |
| :------------------- | :-------- | :------------------------------------ | :---------------------------------------------------------------------------- |
| enable_metrics       | `bool`    | true                                  | Enable to collect metrics.                                                    |
| username\*           | `string`  |                                       | Username used to authenticate.                                                |
| password\*           | `string`  |                                       | Password used to authenticate.                                                |
| endpoint             | `string`  | <http://localhost:15672>              | The endpoint of the Rabbitmq server.                                          |
| collection_interval  | `int`     | 60                                    | How often (seconds) to scrape for metrics.                                    |
| enable_tls           | `bool`    | false                                 | Whether or not to use TLS.                                                    |
| insecure_skip_verify | `bool`    | false                                 | Enable to skip TLS certificate verification.                                  |
| ca_file              | `string`  |                                       | Certificate authority used to validate the database server's TLS certificate. |
| cert_file            | `string`  |                                       | A TLS certificate used for client authentication, if mutual TLS is enabled.   |
| key_file             | `string`  |                                       | A TLS private key used for client authentication, if mutual TLS is enabled.   |
| enable_logs          | `bool`    | true                                  | Enable to collect logs.                                                       |
| daemon_log_paths     | `strings` | ` - "/var/log/rabbitmq/rabbit@*.log"` | Path to Rabbitmq log file(s).                                                 |
| start_at             | `enum`    | end                                   | Start reading file from 'beginning' or 'end'.                                 |

<span style="color:red">\*_required field_</span>