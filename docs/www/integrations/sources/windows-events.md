---
title: "Windows Events"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "windows-events"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Windows  |         | ✓    |        |

## Configuration Table

| Parameter            | Type      | Default | Description                           |
| :------------------- | :-------- | :------ | :------------------------------------ |
| system_event_input   | `bool`    | true    | Enable the System event channel.      |
| app_event_input      | `bool`    | true    | Enable the Application event channel. |
| security_event_input | `bool`    | true    | Enable the Security event channel.    |
| custom_channels      | `strings` |         | Custom channels to read events from.  |
