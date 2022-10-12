---
title: "Prometheus"
slug: "prometheus"
hidden: false
createdAt: "2022-06-07T18:44:46.852Z"
updatedAt: "2022-08-10T15:35:25.879Z"
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