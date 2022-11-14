---
title: "Journald"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "journald"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |

## Configuration Table

| Parameter | Type      | Default | Description                                                                                                |
| :-------- | :-------- | :------ | :--------------------------------------------------------------------------------------------------------- |
| units     | `strings` |         | Service Units to filter on. If not set, all units will be read.                                            |
| directory | `string`  |         | The directory containing Journald's log files. If not set, /run/log/journal and /run/journal will be used. |
| priority  | `enum`    | "info"  | Set log level priority. Valid values are: `trace`, `info`, `warn`, `error`, and `fatal`.                   |
| start_at  | `enum`    | end     | Start reading journal from 'beginning' or 'end'.                                                           |
