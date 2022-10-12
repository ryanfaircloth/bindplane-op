---
title: "ZooKeeper"
slug: "zookeeper"
hidden: false
createdAt: "2022-06-17T17:41:12.690Z"
updatedAt: "2022-08-10T15:37:21.839Z"
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