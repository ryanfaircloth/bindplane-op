---
title: Elasticsearch
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: elasticsearch-1
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
|         | ✓    |        |

## Configuration Table

| Parameter            | Type      | Default                | Description                                                                                                               |
| :------------------- | :-------- | :--------------------- | :------------------------------------------------------------------------------------------------------------------------ |
| endpoints            | `strings` |                        | List of Elasticsearch URLs. If endpoints and cloudid is missing, the ELASTICSEARCH_URL environment variable will be used. |
| cloudid              | `string`  |                        | ID of the Elastic Cloud Cluster to publish events to. The cloudid can be used instead of endpoints.                       |
| index                | `string`  | "logs-generic-default" | The index or datastream name to publish events to.                                                                        |
| pipeline             | `string`  |                        | Optional Ingest Node pipeline ID used for processing documents published by the exporter.                                 |
| user                 | `string`  |                        | Username used for HTTP Basic Authentication.                                                                              |
| password             | `string`  |                        | Password used for HTTP Basic Authentication.                                                                              |
| api_key              | `string`  |                        | Authorization API Key.                                                                                                    |
| enable_tls           | `bool`    | false                  | Whether or not to use TLS.                                                                                                |
| insecure_skip_verify | `bool`    | false                  | Enable to skip TLS certificate verification.                                                                              |
| ca_file              | `string`  |                        | Certificate authority used to validate the database server's TLS certificate.                                             |
| cert_file            | `string`  |                        | A TLS certificate used for client authentication, if mutual TLS is enabled.                                               |
| key_file             | `string`  |                        | A TLS private key used for client authentication, if mutual TLS is enabled.                                               |
