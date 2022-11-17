---
title: "StatsD"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "statsd"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |  ✓      |      |        |
| Windows  |  ✓      |      |        |
| macOS    |  ✓      |      |        |

## Configuration Table

| Parameter | Type | Default | Description |
| :---- | :---- | :---- | :---- |
| listen_ip | `string` | "0.0.0.0" | IP Address to listen on. |
| listen_port | `int` | 8125 | Port to listen on and receive metrics from statsd clients. |
| aggregation_interval | `int` | 60 | The aggregation time in seconds that the receiver aggregates the metrics. |
| enable_metric_type | `bool` | false | Enable the statsd receiver to be able to emit the metric type as a label. |
| is_monotonic_counter | `bool` | false | Set all counter-type metrics the statsd receiver received as monotonic. |