---
title: Kafka
category: 633dd7654359a20031653089
slug: kafka
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Configuration Table

| Parameter             | Type      | Default                    | Description                                                                                                                         |
| :-------------------- | :-------- | :------------------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| protocol_version      | `enum`    | "2.0.0"                    | The Kafka protocol version to use when communicating with brokers. Valid values are: `"2.2.1"`, `"2.2.0"`, `"2.0.0"`, or `"1.0.0"`. |
| brokers               | `strings` | localhost:9092             | List of brokers to connect to when sending metrics, traces and logs.                                                                |
| timeout               | `int`     | 5                          | Timeout (seconds) for every attempt to send data to the backend.                                                                    |
| enable_metrics        | `bool`    | true                       |                                                                                                                                     |
| metric_topic          | `string`  | otlp_metrics               | The name of the topic to export metrics to.                                                                                         |
| enable_logs           | `bool`    | true                       |                                                                                                                                     |
| log_topic             | `string`  | otlp_logs                  | The name of the topic to export logs to.                                                                                            |
| enable_traces         | `bool`    | true                       |                                                                                                                                     |
| trace_topic           | `string`  | otlp_spans                 | The name of the topic to export traces to.                                                                                          |
| enable_auth           | `bool`    | false                      |                                                                                                                                     |
| auth_type             | `enum`    | basic                      | `basic`, `sasl`, or `kerberos`                                                                                                      |
| basic_username        | `string`  |                            |                                                                                                                                     |
| basic_password        | `string`  |                            |                                                                                                                                     |
| sasl_username         | `string`  |                            |                                                                                                                                     |
| sasl_password         | `enum`    |                            |                                                                                                                                     |
| sasl_mechanism        | `string`  | SCRAM-SHA-256              | `SCRAM-SHA-256`, `SCRAM-SHA-512`, or `PLAIN`                                                                                        |
| kerberos_service_name | `string`  |                            |                                                                                                                                     |
| kerberos_realm        | `string`  |                            |                                                                                                                                     |
| kerberos_config_file  | `string`  | /etc/krb5.conf             |                                                                                                                                     |
| kerberos_auth_type    | `enum`    | keytab                     | `keytab` or `basic`                                                                                                                 |
| kerberos_keytab_file  | `string`  | /etc/security/kafka.keytab |                                                                                                                                     |
| kerberos_username     | `string`  |                            |                                                                                                                                     |
| kerberos_password     | `string`  |                            |                                                                                                                                     |