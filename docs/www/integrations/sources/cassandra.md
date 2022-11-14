---
title: "Cassandra"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "cassandra"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports Apache Cassandra versions 3.11 and 4.0.

## Configuration Table

| Parameter           | Type     | Default                                           | Description                                       |
| :------------------ | :------- | :------------------------------------------------ | :------------------------------------------------ |
| collection_interval | `int`    | 60                                                | How often (seconds) to scrape for metrics.        |
| address             | `string` | localhost                                         | IP address or hostname to scrape for JMX metrics. |
| port                | `int`    | 7199                                              | Port to scrape for JMX metrics.                   |
| jar_path            | `string` | "/opt/opentelemetry-java-contrib-jmx-metrics.jar" | Full path to the JMX metrics jar.                 |
