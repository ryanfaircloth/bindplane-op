---
title: "SAP HANA"
category: 633dd7654359a20031653089
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