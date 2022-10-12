---
title: "Cisco Catalyst"
slug: "cisco-catalyst"
hidden: false
createdAt: "2022-08-10T14:36:16.599Z"
updatedAt: "2022-08-10T15:32:02.217Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter   | Type       | Default   | Description                                                                     |
| :---------- | :--------- | :-------- | :------------------------------------------------------------------------------ |
| listen_port | `int`      | 5140      | A UDP port which the agent will listen for syslog messages.                     |
| listen_ip   | `string`   | "0.0.0.0" | An IP address for the agent to bind. Typically 0.0.0.0 for most configurations. |
| timezone    | `timezone` | "UTC"     | The timezone to use when parsing timestamps.                                    |