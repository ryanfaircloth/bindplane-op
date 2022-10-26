---
title: "Elasticsearch"
category: 633dd7654359a20031653089
slug: "elasticsearch"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This receiver supports Elasticsearch versions 7.9+

If Elasticsearch security features are enabled, you must have either the `monitor` or `manage` cluster privilege. See the [Elasticsearch docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/authorization.html) for more information on authorization and [Security privileges](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-privileges.html).

## Configuration Table

| Parameter            | Type      | Default                           | Description                                                                   |
| :------------------- | :-------- | :-------------------------------- | :---------------------------------------------------------------------------- |
| enable_metrics       | `bool`    | true                              | Enable to collect metrics.                                                    |
| hostname*           | `string`  | "localhost"    | The hostname or IP address of the Elasticsearch API.                                             |
| port  | `int`     | 9200                                | The TCP port of the Elasticsearch API.                                    |
| username           | `string`    |                              | Username used to authenticate.                                                   |
| password | `string`    |                            | Password used to authenticate.                                  |
| collection_interval             | `int`  | 60                                  | How often (seconds) to scrape for metrics. |
| nodes           | `strings`  | _node                                  | Filters that define which nodes are scraped for node-level metrics. Should be set to '_node' if collector is installed on all nodes.  '_all' if single collector is scraping the entire collector. <https://www.elastic.co/guide/en/elasticsearch/reference/7.9/cluster.html#cluster-nodes>. |
| skip_cluster_metrics             | `bool`  | false                                  | Enable to disable the collection of cluster level metrics. |
| enable_logs          | `bool`    | true                              | Enable to collect logs.                                                       |
| json_log_paths           | `strings`    | `- \"/var/log/elasticsearch/__server.json\"`  \n`- \"/var/log/elasticsearch/__deprecation.json\"`  \n`- \"/var/log/elasticsearch/__index_search_slowlog.json\"`  \n`- \"/var/log/elasticsearch/__index_indexing_slowlog.json\"`  \n`- \"/var/log/elasticsearch/*_audit.json\"` | File paths for the JSON formatted logs. |
| gc_log_paths | `strings` | `- \"/var/log/elasticsearch/gc.log*\"` | File paths for the garbage collection logs. |
| start_at             | `enum`    | end                               | Start reading file from 'beginning' or 'end'.                                 |

<span style="color:red">\*_required field_</span>