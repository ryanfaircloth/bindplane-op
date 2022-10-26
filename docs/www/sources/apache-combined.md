---
title: "Apache Combined"
category: 633dd7654359a20031653089
slug: "apache-combined"
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
| file_path | `strings` | ["/var/log/apache_combined.log"] | Paths to Apache combined formatted log files  |
| start_at  | `enum`    | end                              | Start reading logs from 'beginning' or 'end'. |