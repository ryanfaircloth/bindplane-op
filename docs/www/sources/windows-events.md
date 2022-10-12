---
title: "Windows Events"
slug: "windows-events"
hidden: false
createdAt: "2022-06-08T13:34:16.236Z"
updatedAt: "2022-08-10T15:37:12.705Z"
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