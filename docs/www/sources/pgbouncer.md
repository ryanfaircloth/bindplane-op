---
title: "PgBouncer"
category: 633dd7654359a20031653089
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