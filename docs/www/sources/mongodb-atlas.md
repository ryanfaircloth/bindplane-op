---
title: "MongDB Atlas"
category: 633dd7654359a20031653089
slug: "mogondbatlas"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter                    | Type      | Default        | Description                                                                                                                                                                                                       |
| :--------------------------- | :-------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------  |
| public_key\*                 | `string`  | ""             | API Public Key with at least Organization Read Only permissions.                                                                                                                                                  |
| private_key\*                | `string`  | ""             | API Private Key.                                                                                                                                                                                                  |
| enable_metrics               | `bool`    | true           | Enable to collect metrics.                                                                                                                                                                                        |
| collection_interval          | `int`     | 60             | How often (seconds) to scrape for                                                                                                                                                                                 |
| granularity                  | `enum`    | PT1M           | Duration interval between measurement data points. Read more [here](https://www.mongodb.com/docs/atlas/reference/api/process-measurements/#request-query-parameters). Valid values: `PT1M`, `PT5M`, `PT1H`, `P1D` |
| enable_logs                  | `bool`    | true           | Enable to collect MongoDB Atlas logs from the API.                                                                                                                                                                |
| log_project_name\*           | `string`  | ""             | Project to collect logs for.                                                                                                                                                                                      |
| collect_audit_logs           | `bool`    | false          | Enable to collect Audit Logs. Must be enabled on project and API Key must have Organization Owner permissions.                                                                                                    |
| log_filter_mode\*            | `enum`    | All            | Mode of filtering clusters. Either collect from all clusters or specify an inclusive list or exclusive list. Valid values: `All`, `Inclusive`, `Exclusive`                                                        |
| log_include_clusters         | `strings` |                | Clusters in the project to collect logs from. Applicable if `log_filter_mode` is `Inclusive`                                                                                                                      |
| log_exclude_clusters         | `strings` |                | Clusters in the project to exclude from log collection. Applicable if `log_filter_mode` is `Exclusive`                                                                                                            |
| enable_alerts                | `bool`    | false          | Enable to collect alerts.                                                                                                                                                                                         |
| alert_collection_mode\*      | `enum`    | poll           | Method of collecting alerts. In `poll` mode alerts are scrapped from the API. In `listen` mode a sever is setup to listen for incoming alerts. Valid values: `poll`, `listen`.                                    |
| alert_project_name\*         | `string`  | ""             | Project to collect alerts from. Applicable if `alert_collection_mode` is `poll`.                                                                                                                                  |
| alert_filter_mode\*          | `enum`    | All            | Mode of filtering clusters. Either collect from all clusters or specify an inclusive list or exclusive list. Applicable if `alert_collection_mode` is `poll`. Valid values: `All`, `Inclusive`, `Exclusive`.      |
| alert_include_clusters       | `strings` |                | Clusters in the project to collect alerts from. Applicable if `log_filter_mode` is `Inclusive` and `alert_collection_mode` is `poll`. |
| alert_exclude_clusters       | `strings` |                | Clusters in the project to exclude from alert collection. Applicable if `log_filter_mode` is `Exclusive` and `alert_collection_mode` is `poll`. |
| page_size                    | `int`     | 100            | The number of alerts to collect per API request. Applicable if `alert_collection_mode` is `poll`. |
| max_pages                    | `int`     | 10             | The limit of how many pages of alerts will request per project. Applicable if `alert_collection_mode` is `poll`. |
| listen_secret\*              | `string`  | ""             | Secret key configured for push notifications. Applicable if `alert_collection_mode` is `listen`. |
| listen_endpoint\*            | `string`  | "0.0.0.0:4396" | Local "ip:port" to bind to, to listen for incoming webhooks. Applicable if `alert_collection_mode` is `listen`. |
| enable_listen_tls            | `bool`    | false          | Enable TLS for alert webhook server. Applicable if `alert_collection_mode` is `listen`. |
| listen_tls_key_file          | `string`  | ""             | Local path to the TLS key file. Applicable if `enable_listen_tls` is true and `alert_collection_mode` is `listen`. |
| listen_tls_cert_file         | `string`  | ""             | Local path to the TLS cert file. Applicable if `enable_listen_tls` is true and `alert_collection_mode` is `listen`. |

<span style="color:red">\*_required field_</span>
