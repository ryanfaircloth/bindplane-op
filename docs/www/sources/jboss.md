---
title: "JBoss"
category: 633dd7654359a20031653089
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