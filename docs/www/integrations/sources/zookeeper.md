---
title: "ZooKeeper"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
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
