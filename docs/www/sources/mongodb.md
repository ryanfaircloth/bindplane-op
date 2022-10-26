---
title: "MongoDB"
category: 633dd7654359a20031653089
slug: "mongodb"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports MongoDB versions 2.6, 3.x, 4.x, and 5.0.

## Configuration Table

| Parameter           | Type  | Default | Description                                |
| :------------------ | :---- | :------ | :----------------------------------------- |
| enable_metrics | `bool` | true      | Enable to collect metrics. |
| hosts | `strings` | "localhost:27017" | List of host:port or unix domain socket endpoints.  \n  \n- For standalone MongoDB deployments this is the hostname and port of the mongod instance.\n- For replica sets specify the hostnames and ports of the mongod instances that are in the replica set configuration. If the replica_set field is specified, nodes will be autodiscovered.\n- For a sharded MongoDB deployment, please specify a list of the mongos hosts. |
| username | `string` | | If authentication is required, specify a username with \"clusterMonitor\" permission. |
| password | `string` | | The password user's password. |
| collection_interval | `int` | 60 | How often (seconds) to scrape for metrics. |
| enable_tls | `bool` | false | Whether or not to use TLS. |
| insecure_skip_verify | `bool` | false | Enable to skip TLS certificate verification. |
| ca_file | `string` | | Certificate authority used to validate the database server's TLS certificate. |
| cert_file | `string` | | A TLS certificate used for client authentication, if mutual TLS is enabled. |
| key_file | `string` | | A TLS private key used for client authentication, if mutual TLS is enabled. |
| enable_logs | `bool` | true | Enable to collect logs. |
| log_paths | `strings` | `- \"/var/log/mongodb/mongodb.log_\"`  \n` - \"/var/log/mongodb/mongod.log_\"` | Path to Mongodb log file(s). |
| start_at | `enum` | end | Start reading file from 'beginning' or 'end'. |