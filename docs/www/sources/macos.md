---
title: "macOS"
slug: "macos"
hidden: false
createdAt: "2022-06-06T18:56:16.318Z"
updatedAt: "2022-08-10T15:34:14.269Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter                | Type     | Default                | Description                                   |
| :----------------------- | :------- | :--------------------- | :-------------------------------------------- |
| enable_metrics           | `bool`   | true                   | Enable to collect metrics.                    |
| host_collection_interval | `int`    | 60                     | How often (seconds) to scrape for metrics.    |
| enable_logs              | `bool`   | true                   | Enable to collect logs.                       |
| enable_system_log        | `bool`   | true                   | Enable to collect macOS system logs.          |
| system_log_path          | `string` | "/var/log/system.log"  | The absolute path to the System log.          |
| enable_install_log       | `bool`   | true                   | Enable to collect macOS install logs.         |
| install_log_path         | `string` | "/var/log/install.log" | The absolute path to the Install log.         |
| start_at                 | `enum`   | end                    | Start reading file from 'beginning' or 'end'. |