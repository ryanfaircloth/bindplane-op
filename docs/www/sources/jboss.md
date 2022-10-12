---
title: "JBoss"
slug: "jboss"
hidden: false
createdAt: "2022-08-10T15:21:56.836Z"
updatedAt: "2022-08-10T15:33:18.744Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter | Type       | Default                                 | Description                                   |
| :-------- | :--------- | :-------------------------------------- | :-------------------------------------------- |
| file_path | `strings`  | /usr/local/JBoss/EAP-_/_/log/server.log | File paths to tail for logs.                  |
| timezone  | `timezone` | "UTC"                                   | The timezone to use when parsing timestamps.  |
| start_at  | `enum`     | end                                     | Start reading file from 'beginning' or 'end'. |