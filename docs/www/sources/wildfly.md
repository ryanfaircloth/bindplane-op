---
title: "WildFly"
category: 633dd7654359a20031653089
slug: "wildfly"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |

## Configuration Table

| Parameter                | Type       | Default                                         | Description                                    |
| :----------------------- | :--------- | :---------------------------------------------- | :--------------------------------------------- |
| standalone_file_path     | `strings`  | /opt/wildfly/standalone/log/server.log          | File paths to tail for standalone server logs. |
| enable_domain_server     | `bool`     | true                                            | Enable to read domain server logs.             |
| domain_server_path       | `strings`  | '/opt/wildfly/domain/servers/\*/log/server.log' | File paths to tail for domain server logs.     |
| enable_domain_controller | `bool`     | true                                            | Enable to read domain controller logs.         |
| domain_controller_path   | `strings`  | '/opt/wildfly/domain/log/\*.log'                | File paths to tail for domain controller logs. |
| start_at                 | `enum`     | end                                             | Start reading file from 'beginning' or 'end'.  |
| timezone                 | `timezone` | "UTC"                                           | The timezone to use when parsing timestamps.   |