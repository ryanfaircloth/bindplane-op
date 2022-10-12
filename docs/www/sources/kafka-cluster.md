---
title: "Kafka Cluster"
slug: "kafka-cluster"
hidden: false
createdAt: "2022-08-02T13:46:43.406Z"
updatedAt: "2022-08-10T15:33:46.660Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       |      |        |
| Windows  | ✓       |      |        |
| macOS    | ✓       |      |        |

## Configuration Table

| Parameter             | Type      | Default                    | Description                                                                                                                         |
| :-------------------- | :-------- | :------------------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| cluster_name\*        | `string`  |                            | Friendly name used for the resource kafka.cluster.name.                                                                             |
| protocol_version      | `enum`    | "2.0.0"                    | The Kafka protocol version to use when communicating with brokers. Valid values are: `"2.2.1"`, `"2.2.0"`, `"2.0.0"`, or `"1.0.0"`. |
| brokers               | `strings` | localhost:9092             | List of brokers to scrape for metrics.                                                                                              |
| client_id             | `string`  | otel-metrics-receiver      | The consumer client ID that the receiver will use.                                                                                  |
| collection_interval   | `int`     | 60                         | How often (seconds) to scrape for metrics.                                                                                          |
| enable_auth           | `bool`    | false                      |                                                                                                                                     |
| auth_type             | `enum`    | basic                      | `basic`, `sasl`, or `kerberos`                                                                                                      |
| basic_username        | `string`  |                            |                                                                                                                                     |
| basic_password        | `string`  |                            |                                                                                                                                     |
| sasl_username         | `string`  |                            |                                                                                                                                     |
| sasl_password         | `string`  |                            |                                                                                                                                     |
| sasl_mechanism        | `enum`    | SCRAM-SHA-256              | `SCRAM-SHA-256`, `SCRAM-SHA-512`, or `PLAIN`                                                                                        |
| kerberos_service_name | `string`  |                            |                                                                                                                                     |
| kerberos_realm        | `string`  |                            |                                                                                                                                     |
| kerberos_config_file  | `string`  | /etc/krb5.conf             |                                                                                                                                     |
| kerberos_auth_type    | `enum`    | keytab                     | `keytab` or `basic`                                                                                                                 |
| kerberos_keytab_file  | `string`  | /etc/security/kafka.keytab |                                                                                                                                     |
| kerberos_username     | `string`  |                            |                                                                                                                                     |
| kerberos_password     | `string`  |                            |                                                                                                                                     |

<span style="color:red">\*_required field_</span>