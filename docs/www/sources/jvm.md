---
title: "JVM"
slug: "jvm"
hidden: false
createdAt: "2022-06-08T13:33:38.974Z"
updatedAt: "2022-08-10T15:33:37.137Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       |      |        |
| Windows  | ✓       |      |        |
| macOS    | ✓       |      |        |

## Prerequisites

This source supports Java versions 11 and 16.

## Configuration Table

| Parameter           | Type     | Default                                           | Description                                       |
| :------------------ | :------- | :------------------------------------------------ | :------------------------------------------------ |
| collection_interval | `int`    | 60                                                | How often (seconds) to scrape for metrics.        |
| address             | `string` | localhost                                         | IP address or hostname to scrape for JMX metrics. |
| port                | `int`    | 9999                                              | Port to scrape for JMX metrics.                   |
| jar_path            | `string` | "/opt/opentelemetry-java-contrib-jmx-metrics.jar" | Full path to the JMX metrics jar.                 |