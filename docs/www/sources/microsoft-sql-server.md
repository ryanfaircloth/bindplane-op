---
title: "Microsoft SQL Server"
category: 633dd7654359a20031653089
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