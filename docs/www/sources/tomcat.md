---
title: "Tomcat"
slug: "tomcat"
hidden: false
createdAt: "2022-06-09T16:19:42.108Z"
updatedAt: "2022-08-10T15:36:23.987Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports Apache Tomcat versions 9.0.x and 10.x.

## Configuration Table

| Parameter           | Type     | Default                                           | Description                                       |
| :------------------ | :------- | :------------------------------------------------ | :------------------------------------------------ |
| collection_interval | `int`    | 60                                                | How often (seconds) to scrape for metrics.        |
| address             | `string` | localhost                                         | IP address or hostname to scrape for JMX metrics. |
| port                | `int`    | 9012                                              | Port to scrape for JMX metrics.                   |
| jar_path            | `string` | "/opt/opentelemetry-java-contrib-jmx-metrics.jar" | Full path to the JMX metrics jar.                 |