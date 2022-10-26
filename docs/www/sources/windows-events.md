---
title: "Windows Events"
category: 633dd7654359a20031653089
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