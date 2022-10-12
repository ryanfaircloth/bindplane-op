---
title: "Solr"
slug: "solr"
hidden: false
createdAt: "2022-06-17T17:38:39.828Z"
updatedAt: "2022-08-10T15:36:00.170Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter           | Type     | Default                                           | Description                                       |
| :------------------ | :------- | :------------------------------------------------ | :------------------------------------------------ |
| collection_interval | `int`    | 60                                                | How often (seconds) to scrape for metrics.        |
| address             | `string` | localhost                                         | IP address or hostname to scrape for JMX metrics. |
| port                | `int`    | 9012                                              | Port to scrape for JMX metrics.                   |
| jar_path            | `string` | "/opt/opentelemetry-java-contrib-jmx-metrics.jar" | Full path to the JMX metrics jar.                 |