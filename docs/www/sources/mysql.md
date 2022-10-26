---
title: "MySQL"
category: 633dd7654359a20031653089
slug: "mysql"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports MySQL versions 5.7 and 8.0. 

## Configuration Table

| Parameter           | Type  | Default | Description                                |
| :------------------ | :---- | :------ | :----------------------------------------- |
| enable_metrics | `bool` | true | Enable to collect metrics. |
| username* | `string` | | Username used to authenticate. |
| password* | `string` | | Password used to authenticate. |
| endpoint | `string` | localhost:3306 | The endpoint of the mysql server. |
| transport | `enum` | tcp | The transport protocol being used to connect to mysql. |
| database | `string` | | The database name. If not specified, metrics will be collected for all databases. |
| collection_interval | `int` | 60 | How often (seconds) to scrape for metrics. |
| enable_logs | `bool` | true | Enable to collect logs. |
| enable_general_log | `bool` | false | Enable to read and parse the general log file. |
| general_log_paths | `strings` | ` - \"/var/log/mysql/general.log\"` | Path to the general log file(s). |
| enable_slow_log | `bool` | true | Enable to read and parse the slow query log. |
| slow_query_log_paths | `strings` | ` - \"/var/log/mysql/slow*.log\"` | Path to the slow query log file(s). |
| enable_error_log | `bool` | true | Enable to read and parse the error log. |
| error_log_paths | `strings` | For CentOS / RHEL: `- \"/var/log/mysqld.log\"`  \nFor SLES: `- \"/var/log/mysql/mysqld.log\"`  \nFor Debian / Ubuntu: `- \"/var/log/mysql/error.log\"` | Path to the error log file(s). |
| start_at | `enum` | end | Start reading file from 'beginning' or 'end'. |

<span style="color:red">\*_required field_</span>