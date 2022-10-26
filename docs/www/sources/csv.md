---
title: "CSV"
category: 633dd7654359a20031653089
slug: "csv"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter         | Type      | Default | Description                                                                                                         |
| :---------------- | :-------- | :------ | :------------------------------------------------------------------------------------------------------------------ |
| header\*          | `string`  |         | A comma delimited list of keys assigned to each of the columns.                                                     |
| file_path\*       | `strings` |         | File or directory paths to tail for logs.                                                                           |
| exclude_file_path | `strings` |         | File or directory paths to exclude.                                                                                 |
| log_type          | `string`  | csv     | A friendly name that will be added to each log entry as an attribute.                                               |
| start_at          | `enum`    | end     | Start reading file from 'beginning' or 'end'.                                                                       |
| encoding          | `enum`    | utf-8   | The encoding of the file being read. Valid values include: `nop`, `utf-8`, `utf-16le`, `utf-16be`, `ascii`, `big5`. |