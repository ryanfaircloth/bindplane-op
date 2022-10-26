---
title: "Host Metrics"
category: 633dd7654359a20031653089
slug: "host-metrics"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       |      |        |
| Windows  | ✓       |      |        |
| macOS    | ✓       |      |        |

## Configuration Table

| Parameter           | Type   | Default | Description                                                                                                                                   |
| :------------------ | :----- | :------ | :-------------------------------------------------------------------------------------------------------------------------------------------- |
| collection_interval | `int`  | 60      | How often (seconds) to scrape for metrics.                                                                                                    |
| enable_load         | `bool` | true    | Enable to collect load metrics. Compatible with all platforms.                                                                                |
| enable_filesystem   | `bool` | true    | Enable to collect filesystem metrics. Compatible with all platforms.                                                                          |
| enable_memory       | `bool` | true    | Enable to collect memory metrics. Compatible with all platforms.                                                                              |
| enable_network      | `bool` | true    | Enable to collect network metrics. Compatible with all platforms.                                                                             |
| enable_paging       | `bool` | true    | Enable to collect paging metrics. Compatible with all platforms.                                                                              |
| enable_cpu          | `bool` | false   | Enable to collect CPU metrics. Compatible with Linux and Windows.                                                                             |
| enable_disk         | `bool` | false   | Enable to collect disk metrics. Compatible with Linux and Windows.                                                                            |
| enable_process      | `bool` | false   | Enable to collect disk metrics. Compatible with Linux and Windows. The collector must be running as root (Linux) and Administrator (Windows). |
| enable_processes    | `bool` | false   | Enable to collect process count metrics. Compatible with Linux only.                                                                          |

## Metrics

| Metric     | Supported OS     | Description                                            |
| :--------- | :--------------- | :----------------------------------------------------- |
| cpu        | All except macOS | CPU utilization metrics                                |
| disk       | All except macOS | Disk I/O metrics                                       |
| load       | All              | CPU load metrics                                       |
| filesystem | All              | File System utilization metrics                        |
| memory     | All              | Memory utilization metrics                             |
| network    | All              | Network interface I/O metrics & TCP connection metrics |
| paging     | All              | Paging/Swap space utilization and I/O metrics          |
| processes  | Linux            | Process count metrics                                  |
| process    | Linux & Windows  | Per process CPU, Memory, and Disk I/O metrics          |