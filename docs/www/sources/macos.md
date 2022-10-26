---
title: "macOS"
category: 633dd7654359a20031653089
slug: "macos"
hidden: false
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