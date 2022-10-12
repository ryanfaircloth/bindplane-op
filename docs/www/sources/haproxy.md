---
title: "HAProxy"
slug: "haproxy"
hidden: false
createdAt: "2022-08-02T13:45:55.728Z"
updatedAt: "2022-08-10T15:32:59.387Z"
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