---
title: "Windows DHCP"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
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
