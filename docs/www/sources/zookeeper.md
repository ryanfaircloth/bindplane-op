---
title: "ZooKeeper"
category: 633dd7654359a20031653089
slug: "zookeeper"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter           | Type     | Default   | Description                                     |
| :------------------ | :------- | :-------- | :---------------------------------------------- |
| collection_interval | `int`    | 60        | How often (seconds) to scrape for metrics.      |
| address             | `string` | localhost | IP address or hostname of the ZooKeeper system. |
| port                | `int`    | 2181      | Port of the ZooKeeper system.                   |