---
title: "Logz.io"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "logzio"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Configuration Table

| Parameter       | Type     | Default                           | Description                                                                                   |
| :-------------- | :------- | :-------------------------------- | :-------------------------------------------------------------------------------------------- |
| enable_logs     | `bool`   | true                              | Enable to send logs to Logz.io                                                                |
| logs_token\*    | `string` |                                   | Your logz.io account token for your logs account                                              |
| enable_metrics  | `bool`   | true                              | Enable to send metrics to Logz.io                                                             |
| metrics_token\* | `string` |                                   | Your logz.io account token for your metrics account                                           |
| listener_url\*  | `string` | "<https://listener.logz.io:8053>" | The URL of the Logz.io listener in your region                                                |
| enable_tracing  | `bool`   | true                              | Enable to send spans to Logz.io                                                               |
| tracing_token   | `string` |                                   | Your logz.io account token for your tracing account                                           |
| region\*        | `enum`   | "us"                              | Your logz.io account region code. Valid options are: `us`, `eu`, `uk`, `nl`, `wa`, `ca`, `au` |
| timeout         | `int`    | 30                                | Time to wait per individual attempt to send data to a backend                                 |

<span style="color:red">\*_required field_</span>
