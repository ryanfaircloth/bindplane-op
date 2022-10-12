---
title: "UDP"
slug: "udp"
hidden: false
createdAt: "2022-08-05T15:13:51.354Z"
updatedAt: "2022-08-10T15:36:37.790Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter     | Type     | Default   | Description                                                                        |
| :------------ | :------- | :-------- | :--------------------------------------------------------------------------------- |
| listen_port\* | `int`    |           | Port to listen on.                                                                 |
| listen_ip     | `string` | "0.0.0.0" | IP Address to listen on.                                                           |
| log_type      | `string` | ucp       | Arbitrary for attribute 'log_type'. Useful for filtering between many udp sources. |

<span style="color:red">\*_required field_</span>