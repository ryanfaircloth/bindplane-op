---
title: "JBoss"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "jboss"
hidden: false
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
