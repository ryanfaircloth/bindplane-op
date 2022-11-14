---
title: "Ubiquiti"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "ubiquiti"
hidden: false
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

<span style="color:red">\*_required field_</span>
