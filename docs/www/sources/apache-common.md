---
title: "Apache Common"
slug: "apache-common"
hidden: false
createdAt: "2022-08-10T14:40:28.761Z"
updatedAt: "2022-08-10T15:30:57.235Z"
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