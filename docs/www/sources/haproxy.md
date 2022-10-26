---
title: "HAProxy"
category: 633dd7654359a20031653089
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