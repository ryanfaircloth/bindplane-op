---
title: "Prometheus"
category: 633dd7654359a20031653089
slug: "prometheus"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       |      |        |
| Windows  | ✓       |      |        |
| macOS    | ✓       |      |        |

## Configuration Table

| Parameter           | Type      | Default | Description                                                              |
| :------------------ | :-------- | :------ | :----------------------------------------------------------------------- |
| job_name\*          | `string`  |         | The name of the scraper job. Will be set as service.name resource label. |
| static_targets\*    | `strings` |         | List of endpoints to scrape.                                             |
| collection_interval | `int`     | 60      | How often (seconds) to scrape for metrics.                               |

<span style="color:red">\*_required field_</span>