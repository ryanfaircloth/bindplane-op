---
title: "Redis"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "redis"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports Redis version 6.2.

## Configuration Table

| Parameter           | Type  | Default | Description                                |
| :------------------ | :---- | :------ | :----------------------------------------- |
| enable_metrics | `bool` | true | Enable to collect metrics. |
| endpoint | `string` | "localhost:6379" | The endpoint of the Redis server. |
| transport | `enum` | tcp | The transport protocol being used to connect to Redis. Valid values are `tcp` or `unix`. |
| password | `string` | | The password used to access the Redis instance; must match the password specified in the requirepass server configuration option. |
| collection_interval | `int` | 60 | How often (seconds) to scrape for metrics. |
| enable_tls | `bool` | false | Whether or not to use TLS. |
| insecure_skip_verify | `bool` | false | Enable to skip TLS certificate verification. |
| ca_file | `string` | | Certificate authority used to validate the database server's TLS certificate. |
| cert_file | `string` | | A TLS certificate used for client authentication, if mutual TLS is enabled. |
| key_file | `string` | | A TLS private key used for client authentication, if mutual TLS is enabled. |
| enable_logs | `bool` | true | Enable to collect logs. |
| file_path | `strings` | One-click installer: `- \"/var/log/redis/redis_6379.log\"`  \nUbuntu / Debian: `- \"/var/log/redis/redis-server.log\"`  \nsrc: `- \"/var/log/redis_6379.log\"`  \nCentOS / RHEL: `- \"/var/log/redis/redis.log\"`  \nSLES: `- \"/var/log/redis/default.log\"` | Path to Redis log file(s). |
| start_at | `enum` | end | Start reading file from 'beginning' or 'end'. |
