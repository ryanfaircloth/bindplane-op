---
title: "Prometheus"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
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
