---
title: "Apache Combined"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
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
