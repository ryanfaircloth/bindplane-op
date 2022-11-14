---
title: "SAP HANA"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "sap-hana"
hidden: false
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
