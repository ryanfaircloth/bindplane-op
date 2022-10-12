---
title: "Cisco ASA"
slug: "cisco-asa"
hidden: false
createdAt: "2022-08-10T14:32:12.736Z"
updatedAt: "2022-08-10T15:31:52.863Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter   | Type     | Default   | Description                                                                     |
| :---------- | :------- | :-------- | :------------------------------------------------------------------------------ |
| listen_port | `int`    | 5140      | A TCP port which the agent will listen for syslog messages.                     |
| listen_ip   | `string` | "0.0.0.0" | An IP address for the agent to bind. Typically 0.0.0.0 for most configurations. |