---
title: "Common Event Format"
slug: "common-event-format"
hidden: false
createdAt: "2022-08-10T14:45:27.009Z"
updatedAt: "2022-08-10T15:32:11.304Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter             | Type       | Default | Description                                                                                                                                                                          |
| :-------------------- | :--------- | :------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| file_log_path\*       | `strings`  |         | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (\*) to read multiple files within a directory.                                       |
| exclude_file_log_path | `strings`  | ""      | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (\*) to exclude multiple files from being read within a directory. |
| log_type              | `string`   | "cef"   | Adds the specified 'Type' as a log record attribute to each log message.                                                                                                             |
| location              | `timezone` | "UTC"   | The geographic location (timezone) to use when parsing logs that contain a timestamp                                                                                                 |
| timezone              | `timezone` | "UTC"   | The timezone to use when parsing timestamps.                                                                                                                                         |
| start_at              | `enum`     | end     | Start reading file from 'beginning' or 'end'.                                                                                                                                        |

<span style="color:red">\*_required field_</span>