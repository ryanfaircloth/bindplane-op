---
title: "PgBouncer"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "pgbouncer"
hidden: false
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
