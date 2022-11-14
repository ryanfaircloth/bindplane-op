---
title: "Oracle Database"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "oracle-database"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter           | Type      | Default                                                                    | Description                                   |
| :------------------ | :-------- | :------------------------------------------------------------------------- | :-------------------------------------------- |
| enable_audit_log    | `bool`    | true                                                                       |                                               |
| audit_log_path      | `strings` | "/u01/app/oracle/product/_/dbhome_1/admin/_/adump/\*.aud"                  | File paths to audit logs.                     |
| enable_alert_log    | `bool`    | true                                                                       |                                               |
| alert_log_path      | `strings` | "/u01/app/oracle/product/_/dbhome_1/diag/rdbms/_/\_/trace/alert\_\_.log"   | File paths to alert logs.                     |
| enable_listener_log | `bool`    | true                                                                       |                                               |
| listener_log_path   | `strings` | "/u01/app/oracle/product/_/dbhome_1/diag/tnslsnr/_/listener/alert/log.xml" | File paths to listener logs.                  |
| start_at            | `enum`    | end                                                                        | Start reading file from 'beginning' or 'end'. |

<span style="color:red">\*_required field_</span>
