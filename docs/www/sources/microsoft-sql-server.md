---
title: "Microsoft SQL Server"
slug: "microsoft-sql-server"
hidden: false
createdAt: "2022-08-02T13:47:52.718Z"
updatedAt: "2022-08-10T15:34:29.244Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Windows  | ✓       | ✓    |        |

## Configuration Table

| Parameter           | Type   | Default | Description                                   |
| :------------------ | :----- | :------ | :-------------------------------------------- |
| enable_metrics      | `bool` | true    | Enable to collect metrics.                    |
| collection_interval | `int`  | 60      | How often (seconds) to scrape for metrics.    |
| enable_logs         | `bool` | true    | Enable to collect logs.                       |
| start_at            | `enum` | end     | Start reading file from 'beginning' or 'end'. |