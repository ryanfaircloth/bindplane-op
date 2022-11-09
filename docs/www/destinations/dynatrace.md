---
title: "Dynatrace"
category: 633dd7654359a20031653089
slug: "dynatrace"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       |      |        |

## Configuration Table

| Parameter              | Type     | Default | Description                                                                                                                                                                                                                                                          |
| :--------------------- | :------- | :------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| metric_ingest_endpoint | `string` | `""`    | Dynatrace Metrics Ingest v2 endpoint. Required if OneAgent is not running on the agent host. More information on the endpoint and structure can be found [here](https://www.dynatrace.com/support/help/dynatrace-api/environment-api/metric-v2/post-ingest-metrics). |
| api_token              | `string` |         | API Token that is restricted to `Ingest metrics` scope. Required if `metric_ingest_endpoint` is specified. More information [here](https://www.dynatrace.com/support/help/dynatrace-api/basics/dynatrace-api-authentication).                                        |
| prefix                 | `string` |         | Metric Prefix that will be prepended to each metric name in `prefix.name` format.                                                                                                                                                                                    |
| enable_tls             | `bool`   | false   | Whether or not to use TLS.                                                                                                                                                                                                                                           |
| insecure_skip_verify   | `bool`   | false   | Enable to skip TLS certificate verification.                                                                                                                                                                                                                         |
| ca_file                | `bool`   |         | Certificate authority used to validate the database server's TLS certificate.                                                                                                                                                                                        |
| cert_file              | `bool`   |         | A TLS certificate used for client authentication, if mutual TLS is enabled.                                                                                                                                                                                          |
| key_file               | `bool`   |         | A TLS private key used for client authentication, if mutual TLS is enabled.                                                                                                                                                                                          |
