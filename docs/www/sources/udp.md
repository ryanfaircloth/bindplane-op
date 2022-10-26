---
title: "UDP"
category: 633dd7654359a20031653089
slug: "udp"
hidden: false
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