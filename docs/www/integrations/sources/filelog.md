---
title: "Filelog"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "filelog"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓   |        |
| Windows  |         | ✓   |        |
| macOS    |         | ✓   |        |

## Configuration Table

| Parameter                    | Type      | Default                          | Description                                                                                                                                                                                                                                                                            |
| :--------------------------- | :-------- | :------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| file_path\*                  | `strings` |                                  | File or directory paths to tail for logs.                                                                                                                                                                                                                                              |
| exclude_file_path            | `strings` | ""                               | File or directory paths to exclude.                                                                                                                                                                                                                                                    |
| log_type                     | `string`  | "file"                           | A friendly name that will be added to each log entry as an attribute.                                                                                                                                                                                                                  |
| parse_format                 | `enum`    | none                             | Method to use when parsing. Valid values are `none`, `json`, and `regex`. When regex is selected, 'Regex Pattern' must be set.                                                                                                                                                         |
| regex_pattern                | `string`  |                                  | The regex pattern used when parsing log entries.                                                                                                                                                                                                                                       |
| multiline_line_start_pattern | `string`  |                                  | Regex pattern that matches beginning of a log entry, for handling multiline logs.                                                                                                                                                                                                      |
| multiline_line_end_pattern   | `string`  |                                  | Regex pattern that matches end of a log entry, useful for terminating parsing of multiline logs.                                                                                                                                                                                       |
| encoding                     | `enum`    | utf-8                            | The encoding of the file being read. Valid values are `nop`, `utf-8`, `utf-16le`, `utf-16be`, `ascii`, and `big5`.                                                                                                                                                                     |
| include_file_name_attribute  | `bool`    | true                             | Whether to add the file name as the attribute `log.file.name`.                                                                                                                                                                                                                         |
| include_file_path_attribute  | `bool`    | false                            | Whether to add the file path as the attribute `log.file.path`.                                                                                                                                                                                                                         |
| include_file_name_resolved   | `bool`    | false                            | Whether to add the file name after symlinks resolution as the attribute `log.file.name_resolved`.                                                                                                                                                                                      |
| include_file_path_resolved   | `bool`    | false                            | Whether to add the file path after symlinks resolution as the attribute `log.file.path_resolved`.                                                                                                                                                                                      |
| offset_storage_dir           | `string`  | $OIQ_OTEL_COLLECTOR_HOME/storage | The directory that the offset storage file will be created. It is okay if multiple receivers use the same directory. By default the [observIQ Distro for OpenTelemetry Collector](https://github.com/observIQ/observiq-otel-collector) sets `$OIQ_OTEL_COLLECTOR_HOME` in its runtime. |
| poll_interval                | `int`     | 200                              | The duration of time in milliseconds between filesystem polls.                                                                                                                                                                                                                         |
| max_concurrent_files         | `int`     | 1024                             | The maximum number of log files fromw hich logs will be read concurrently. If the number of files matched exceeds this number, then files will be processed in batches.                                                                                                                |
| parse_to                     | `string`  | body                             | The [field](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/stanza/docs/types/field.md)  that the log will be parsed to. Some exporters handle logs favorably when parsed to `attributes` over `body` and vice verca.                                  |
| start_at                     | `enum`    | end                              | Start reading file from 'beginning' or 'end'.                                                                                                                                                                                                                                          |

<span style="color:red">\*_required field_</span>
