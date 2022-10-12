---
title: "SAP HANA"
slug: "sap-hana"
hidden: false
createdAt: "2022-08-10T14:46:20.961Z"
updatedAt: "2022-08-10T15:35:47.484Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter | Type       | Default                         | Description                                          |
| :-------- | :--------- | :------------------------------ | :--------------------------------------------------- |
| file_path | `strings`  | "/usr/sap/_/HDB_/_/trace/_.trc" | File paths to logs.                                  |
| timezone  | `timezone` | "UTC"                           | The hostname or IP address of the Elasticsearch API. |
| start_at  | `enum`     | end                             | Start reading file from 'beginning' or 'end'.        |