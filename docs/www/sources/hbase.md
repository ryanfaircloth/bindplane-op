---
title: "HBase"
slug: "hbase"
hidden: false
createdAt: "2022-08-02T13:46:16.664Z"
updatedAt: "2022-08-10T15:33:06.160Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter            | Type      | Default                                             | Description                                       |
| :------------------- | :-------- | :-------------------------------------------------- | :------------------------------------------------ |
| enable_metrics       | `bool`    | true                                                | Enable to send metrics.                           |
| collection_interval  | `int`     | 60                                                  | How often (seconds) to scrape for metrics.        |
| address              | `string`  | localhost                                           | IP address or hostname to scrape for JMX metrics. |
| jar_path             | `string`  | "/opt/opentelemetry-java-contrib-jmx-metrics.jar"   | Full path to the JMX metrics jar.                 |
| enable_master_jmx    | `bool`    | true                                                | Enable to scrape master server's JMX port.        |
| master_jmx_port      | `int`     | 10101                                               | Master server's JMX Port.                         |
| enable_region_jmx    | `bool`    | true                                                | Enable to scrape region server's JMX port.        |
| region_jmx_port      | `int`     | 10102                                               | Region server's JMX Port.                         |
| enable_logs          | `bool`    | true                                                | Enable to send logs.                              |
| enable_master_log    | `bool`    | true                                                | Enable to read master logs.                       |
| master_log_path      | `strings` | "/usr/local/hbase_/logs/hbase_-master-\*.log"       | File paths to tail for master logs.               |
| enable_region_log    | `bool`    | true                                                | Enable to read region server logs.                |
| region_log_path      | `strings` | "/usr/local/hbase_/logs/hbase_-regionserver-\*.log" | File paths to tail for region server logs.        |
| enable_zookeeper_log | `bool`    | false                                               | Enable to read zookeeper logs.                    |
| zookeeper_log_path   | `strings` | "/usr/local/hbase_/logs/hbase_-zookeeper-\*.log"    | File paths to tail for zookeeper logs.            |
| start_at             | `enum`    | end                                                 | Start reading file from 'beginning' or 'end'.     |