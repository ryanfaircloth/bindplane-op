---
title: "PgBouncer"
slug: "pgbouncer"
hidden: false
createdAt: "2022-08-02T13:47:14.865Z"
updatedAt: "2022-08-10T15:35:10.900Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter | Type      | Default                          | Description                                   |
| :-------- | :-------- | :------------------------------- | :-------------------------------------------- |
| file_path | `strings` | /var/log/pgbouncer/pgbouncer.log | Path to log file(s).                          |
| start_at  | `enum`    | end                              | Start reading file from 'beginning' or 'end'. |