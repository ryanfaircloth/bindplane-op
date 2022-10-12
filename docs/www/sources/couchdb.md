---
title: "CouchDB"
slug: "couchdb"
hidden: false
createdAt: "2022-08-05T15:32:12.484Z"
updatedAt: "2022-08-10T15:32:23.615Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter           | Type      | Default                        | Description                                                                                                                                         |
| :------------------ | :-------- | :----------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| enable_metrics      | `bool`    | true                           | Enable to send metrics.                                                                                                                             |
| hostname\*          | `string`  |                                | The hostname or IP address of the CouchDB system.                                                                                                   |
| port                | `int`     | 5984                           | The TCP port of the CouchDB system.                                                                                                                 |
| username\*          | `string`  |                                | The username to use when connecting to CouchDB.                                                                                                     |
| password\*          | `string`  |                                | The password to use when connecting to CouchDB.                                                                                                     |
| collection_interval | `int`     | 60                             | How often (seconds) to scrape for metrics.                                                                                                          |
| enable_tls          | `bool`    | false                          | Whether or not to use TLS when connecting to CouchDB.                                                                                               |
| strict_tls_verify   | `bool`    | false                          | Enable to require TLS certificate verification.                                                                                                     |
| ca_file             | `string`  |                                | Certificate authority used to validate TLS certificates. Not required if the collector's operating system already trusts the certificate authority. |
| mutual_tls          | `bool`    | false                          | Enable to require TLS mutual authentication.                                                                                                        |
| cert_file           | `string`  |                                | A TLS certificate used for client authentication, if mutual TLS is enabled.                                                                         |
| key_file            | `string`  |                                | A TLS private key used for client authentication, if mutual TLS is enabled.                                                                         |
| enable_logs         | `bool`    | true                           | Enable to collect logs.                                                                                                                             |
| log_paths           | `strings` | "/var/log/couchdb/couchdb.log" | Path to CouchDB log file(s).                                                                                                                        |
| start_at            | `enum`    | end                            | Start reading file from 'beginning' or 'end'.                                                                                                       |

<span style="color:red">\*_required field_</span>