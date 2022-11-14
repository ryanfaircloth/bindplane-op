---
title: "Apache Common"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "apache-common"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter | Type      | Default                         | Description                                   |
| :-------- | :-------- | :------------------------------ | :-------------------------------------------- |
| file_path | `strings` | ["/var/log/apache2/access.log"] | Path to apache common formatted log file      |
| start_at  | `enum`    | end                             | Start reading logs from 'beginning' or 'end'. |
