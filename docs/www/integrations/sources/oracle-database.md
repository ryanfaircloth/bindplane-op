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
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Metrics Requirements

To collect metrics from OracleDB, a user with `SELECT` access to the relevant views is required. To create a new user with those permissions, run the following SQL script as a user with sufficient permissions connected to the Oracle DB instance as SYSDBA or SYSOPER.

```sql
  -- Create the monitoring user "bindplane"
  CREATE USER bindplane IDENTIFIED BY <authentication password>;

  -- Grant the "bindplane" user the required permissions
  GRANT CONNECT TO bindplane;
  GRANT SELECT ON SYS.GV_$DATABASE to bindplane;
  GRANT SELECT ON SYS.GV_$INSTANCE to bindplane;
  GRANT SELECT ON SYS.GV_$PROCESS to bindplane;
  GRANT SELECT ON SYS.GV_$RESOURCE_LIMIT to bindplane;
  GRANT SELECT ON SYS.GV_$SYSMETRIC to bindplane;
  GRANT SELECT ON SYS.GV_$SYSSTAT to bindplane;
  GRANT SELECT ON SYS.GV_$SYSTEM_EVENT to bindplane;
  GRANT SELECT ON SYS.V_$RMAN_BACKUP_JOB_DETAILS to bindplane;
  GRANT SELECT ON SYS.V_$SORT_SEGMENT to bindplane;
  GRANT SELECT ON SYS.V_$TABLESPACE to bindplane;
  GRANT SELECT ON SYS.V_$TEMPFILE to bindplane;
  GRANT SELECT ON SYS.DBA_DATA_FILES to bindplane;
  GRANT SELECT ON SYS.DBA_FREE_SPACE to bindplane;
  GRANT SELECT ON SYS.DBA_TABLESPACE_USAGE_METRICS to bindplane;
  GRANT SELECT ON SYS.DBA_TABLESPACES to bindplane;
  GRANT SELECT ON SYS.GLOBAL_NAME to bindplane;
```

## Configuration Table

| Parameter           | Type      | Default                                                                    | Description                                                                    |
| :------------------ | :-------- | :------------------------------------------------------------------------- | :----------------------------------------------------------------------------- |
| enable_metrics      | `bool`    | true                                                                       | Enable to collect metrics.                                                     |
| host                | `string`  | localhost                                                                  | Host to scrape metrics from.                                                   |
| port                | `int`     | 1521                                                                       | Port of host to scrape metrics from.                                           |
| username*           | `string`  |                                                                            | Database user to run metric queries with.                                      |
| password            | `string`  |                                                                            | Password for user.                                                             |
| sid                 | `string`  |                                                                            | Site Identifier. One or both of sid or service_name must be specified.         |
| service_name        | `string`  |                                                                            | OracleDB Service Name. One or both of sid or service_name must be specified.   |
| wallet              | `string`  |                                                                            | OracleDB Wallet file location (must be URL encoded).                           |
| collection_interval | `int`     | 60                                                                         | How often (seconds) to scrape for metrics.                                     |
| enable_logs         | `bool`    | true                                                                       | Enable to collect logs.                                                        |
| enable_audit_log    | `bool`    | true                                                                       | Enable to collect audit logs.                                                  |
| audit_log_path      | `strings` | "/u01/app/oracle/product/_/dbhome_1/admin/_/adump/\*.aud"                  | File paths to audit logs.                                                      |
| enable_alert_log    | `bool`    | true                                                                       |                                                                                |
| alert_log_path      | `strings` | "/u01/app/oracle/product/_/dbhome_1/diag/rdbms/_/\_/trace/alert\_\_.log"   | File paths to alert logs.                                                      |
| enable_listener_log | `bool`    | true                                                                       |                                                                                |
| listener_log_path   | `strings` | "/u01/app/oracle/product/_/dbhome_1/diag/tnslsnr/_/listener/alert/log.xml" | File paths to listener logs.                                                   |
| start_at            | `enum`    | end                                                                        | Start reading file from 'beginning' or 'end'.                                  |

<span style="color:red">\*_required field_</span>
