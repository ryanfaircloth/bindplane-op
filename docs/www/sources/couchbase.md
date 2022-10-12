---
title: "Couchbase"
slug: "couchbase"
hidden: false
createdAt: "2022-08-02T13:44:21.414Z"
updatedAt: "2022-08-10T15:32:18.042Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter                     | Type      | Default                                                          | Description                                      |
| :---------------------------- | :-------- | :--------------------------------------------------------------- | :----------------------------------------------- |
| enable_metrics                | `bool`    | true                                                             | Enable to collect metrics.                       |
| hostname\*                    | `string`  | "localhost"                                                      | The hostname or IP address of the Couchbase API. |
| port                          | `int`     | 8091                                                             | The TCP port of the Couchbase API.               |
| username\*                    | `string`  |                                                                  | Username used to authenticate.                   |
| password\*                    | `string`  |                                                                  | Password used to authenticate.                   |
| collection_interval           | `int`     | 60                                                               | How often (seconds) to scrape for metrics.       |
| enable_logs                   | `bool`    | true                                                             | Enable to collect logs.                          |
| enable_error_log              | `bool`    | true                                                             | Enable to read error logs.                       |
| error_log_path                | `strings` | "/opt/couchbase/var/lib/couchbase/logs/error.log"                | Log File paths to tail for error logs.           |
| enable_info_log               | `bool`    | false                                                            | Enable to read info logs.                        |
| info_log_path                 | `strings` | "/opt/couchbase/var/lib/couchbase/logs/info.log"                 | Log File paths to tail for info logs.            |
| enable_debug_log              | `bool`    | false                                                            | Enable to read debug logs.                       |
| debug_log_path                | `strings` | "/opt/couchbase/var/lib/couchbase/logs/debug.log"                | Log File paths to tail for debug logs.           |
| enable_access_log             | `bool`    | false                                                            | Enable to read http access logs.                 |
| http_access_log_path          | `strings` | "/opt/couchbase/var/lib/couchbase/logs/http_access.log"          | Log File paths to tail for http access logs.     |
| enable_internal_access_log    | `bool`    | false                                                            | Enable to read internal access logs.             |
| http_internal_access_log_path | `strings` | "/opt/couchbase/var/lib/couchbase/logs/http_access_internal.log" | Log File paths to tail for internal access logs. |
| enable_babysitter_log         | `bool`    | false                                                            | Enable to read baby sitter logs.                 |
| babysitter_log_path           | `strings` | "/opt/couchbase/var/lib/couchbase/logs/babysitter.log"           | Log File paths to tail for baby sitter logs.     |
| enable_xdcr_log               | `bool`    | false                                                            | Enable to read xdcr logs.                        |
| xdcr_log_path                 | `strings` | "/opt/couchbase/var/lib/couchbase/logs/goxdcr.log"               | Log File paths to tail for xdcr logs.            |
| start_at                      | `enum`    | end                                                              | Start reading logs from 'beginning' or 'end'.    |

<span style="color:red">\*_required field_</span>