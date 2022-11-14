---
title: "Microsoft SQL Server"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "microsoft-sql-server"
hidden: false
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
