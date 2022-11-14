---
title: "Aerospike"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "aerospike"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter               | Type     | Default   | Description                                                      |
| :---------------------- | :------- | :-------- | :--------------------------------------------------------------- |
| enable_metrics          | `bool`   | true      | Enable to send metrics.                                          |
| hostname\*              | `string` | localhost | The hostname or IP address of the Aerospike system.              |
| port                    | `int`    | 3000      | The TCP port of the Aerospike system.                            |
| collection_interval     | `int`    | 60        | How often (seconds) to scrape for metrics.                       |
| collect_cluster_metrics | `bool`   | false     | Whether discovered peer nodes should be collected.               |
| aerospike_enterprise    | `bool`   | false     | Enable Aerospike enterprise authentication.                      |
| username\*              | `string` |           | The username to use when connecting to Aerospike.                |
| password\*              | `string` |           | The password to use when connecting to Aerospike.                |
| enable_logs             | `bool`   | true      | Enable to collect Aerospike logs from Journald.                  |
| start_at                | `enum`   | end       | Start reading Aerospike Journald logs from 'beginning' or 'end'. |

<span style="color:red">\*_required field_</span>
