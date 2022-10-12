---
title: "Apache Combined"
slug: "apache-combined"
hidden: false
createdAt: "2022-08-10T14:31:52.541Z"
updatedAt: "2022-08-10T15:30:51.065Z"
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