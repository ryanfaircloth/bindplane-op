---
title: "Kafka Node"
category: 633dd7654359a20031653089
slug: "kafka-node"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter               | Type      | Default                                           | Description                                       |
| :---------------------- | :-------- | :------------------------------------------------ | :------------------------------------------------ |
| enable_metrics          | `bool`    | true                                              |                                                   |
| collection_interval     | `int`     | 60                                                | How often (seconds) to scrape for metrics.        |
| address                 | `string`  | localhost                                         | IP address or hostname to scrape for JMX metrics. |
| port                    | `int`     | 9999                                              | Port to scrape for JMX metrics.                   |
| jar_path                | `string`  | "/opt/opentelemetry-java-contrib-jmx-metrics.jar" | Full path to the JMX metrics jar.                 |
| enable_logs             | `bool`    | true                                              |                                                   |
| enable_server_log       | `bool`    | true                                              |                                                   |
| server_log_path         | `strings` | /home/kafka/kafka/logs/server.log                 | File paths to tail for server logs.               |
| enable_controller_log   | `bool`    | true                                              |                                                   |
| controller_log_path     | `strings` | /home/kafka/kafka/logs/controller.log             | File paths to tail for controller logs.           |
| enable_state_change_log | `bool`    | true                                              |                                                   |
| state_change_log_path   | `strings` | /home/kafka/kafka/logs/state-change.log           | File paths to tail for stage change logs.         |
| enable_log_cleaner_log  | `bool`    | true                                              |                                                   |
| log_cleaner_log_path    | `strings` | /home/kafka/kafka/logs/state-cleaner.log          | File paths to tail for log cleaner logs.          |
| start_at                | `enum`    | end                                               | Start reading file from 'beginning' or 'end'.     |