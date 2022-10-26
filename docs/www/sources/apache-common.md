---
title: "Apache Common"
category: 633dd7654359a20031653089
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