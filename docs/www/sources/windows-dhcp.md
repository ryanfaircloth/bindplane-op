---
title: "Windows DHCP"
slug: "windows-dhcp"
hidden: false
createdAt: "2022-08-10T14:45:03.212Z"
updatedAt: "2022-08-10T15:37:08.010Z"
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