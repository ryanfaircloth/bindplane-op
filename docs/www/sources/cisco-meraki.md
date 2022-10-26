---
title: "Cisco Meraki"
category: 633dd7654359a20031653089
slug: "cisco-meraki"
hidden: false
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
| listen_port | `int`    | 5140      | A UDP port which the agent will listen for syslog messages.                     |
| listen_ip   | `string` | "0.0.0.0" | An IP address for the agent to bind. Typically 0.0.0.0 for most configurations. |