---
title: "Filelog"
category: 633dd7654359a20031653089
slug: "filelog"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter                    | Type      | Default | Description                                                                                                                    |
| :--------------------------- | :-------- | :------ | :----------------------------------------------------------------------------------------------------------------------------- |
| file_path\*                  | `strings` |         | File or directory paths to tail for logs.                                                                                      |
| exclude_file_path            | `strings` | ""      | File or directory paths to exclude.                                                                                            |
| log_type                     | `string`  | "file"  | A friendly name that will be added to each log entry as an attribute.                                                          |
| parse_format                 | `enum`    | none    | Method to use when parsing. Valid values are `none`, `json`, and `regex`. When regex is selected, 'Regex Pattern' must be set. |
| regex_pattern                | `string`  |         | The regex pattern used when parsing log entries.                                                                               |
| start_at                     | `enum`    | end     | Start reading file from 'beginning' or 'end'.                                                                                  |
| multiline_line_start_pattern | `string`  |         | Regex pattern that matches beginning of a log entry, for handling multiline logs.                                              |
| encoding                     | `enum`    | utf-8   | The encoding of the file being read. Valid values are `nop`, `utf-8`, `utf-16le`, `utf-16be`, `ascii`, and `big5`.             |

<span style="color:red">\*_required field_</span>