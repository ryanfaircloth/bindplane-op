---
title: "HAProxy"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "haproxy"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter | Type       | Default                      | Description                                   |
| :-------- | :--------- | :--------------------------- | :-------------------------------------------- |
| file_path | `strings`  | /var/log/haproxy/haproxy.log | Log File paths to tail for logs.              |
| start_at  | `enum`     | true                         | Start reading logs from 'beginning' or 'end'. |
| timezone  | `timezone` | "UTC"                        | The timezone to use when parsing timestamps.  |
