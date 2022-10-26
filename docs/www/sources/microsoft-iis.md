---
title: "Microsoft IIS"
category: 633dd7654359a20031653089
slug: "microsoft-iis"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Windows  | ✓       | ✓    |        |

## Prerequisites

This source supports IIS versions 8.5 and 10.0.

## Configuration Table

| Parameter             | Type      | Default                                        | Description                                                |
| :-------------------- | :-------- | :--------------------------------------------- | :--------------------------------------------------------- |
| enable_metrics        | `bool`    | true                                           | Enable to send metrics.                                    |
| collection_interval   | `int`     | 60                                             | How often (seconds) to scrape for metrics.                 |
| enable_logs           | `bool`    | true                                           | Enable to send logs.                                       |
| file_path             | `strings` | `["C:/inetpub/logs/LogFiles/W3SVC_/**/_.log"]` | File or directory paths to tail for logs.                  |
| exclude_file_log_path | `strings` |                                                | File or directory paths to exclude.                        |
| timezone              | `enum`    | UTC                                            | RFC3164 only. The timezone to use when parsing timestamps. |
| start_at              | `enum`    | end                                            | Start reading file from 'beginning' or 'end'.              |

## Metrics

| Metric                       | Unit          | Description                                          |
| :--------------------------- | :------------ | :--------------------------------------------------- |
| iis.connection.active        | {connections} | Number of active connections.                        |
| iis.connection.anonymous     | {connections} | Number of connections established anonymously.       |
| iis.connection.attempt.count | {attempts}    | Total number of attempts to connect to the server.   |
| iis.network.blocked          | By            | Number of bytes blocked due to bandwidth throttling. |
| iis.network.file.count       | {files}       | Number of transmitted files.                         |
| iis.network.io               | By            | Total amount of bytes sent and received.             |
| iis.request.count            | {requests}    | Total number of requests of a given type.            |
| iis.request.queue.age.max    | ms            | Age of oldest request in the queue.                  |
| iis.request.queue.count      | {requests}    | Current number of requests in the queue.             |
| iis.request.rejected         | {requests}    | Total number of requests rejected.                   |
| iis.thread.active            | {threads}     | Current number of active threads.                    |
| iis.uptime                   | s             | The amount of time the server has been up.           |