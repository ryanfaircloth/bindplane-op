---
title: "Windows DHCP"
category: 633dd7654359a20031653089
slug: "windows-dhcp"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Windows  |         | ✓    |        |

## Configuration Table

| Parameter | Type      | Default                                      | Description                                   |
| :-------- | :-------- | :------------------------------------------- | :-------------------------------------------- |
| file_path | `strings` | "C:/Windows/System32/dhcp/DhcpSrvLog-\*.log" | File or directory paths to tail for logs.     |
| start_at  | `enum`    | end                                          | Start reading file from 'beginning' or 'end'. |