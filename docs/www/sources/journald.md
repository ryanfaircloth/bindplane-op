---
title: "Journald"
slug: "journald"
hidden: false
createdAt: "2022-06-06T18:41:24.733Z"
updatedAt: "2022-08-10T15:33:25.131Z"
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