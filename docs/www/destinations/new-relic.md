---
title: "New Relic"
slug: "new-relic"
hidden: false
createdAt: "2022-06-07T18:46:33.613Z"
updatedAt: "2022-08-10T15:01:57.656Z"
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Configuration Table

| Parameter     | Type     | Default                    | Description                                                                                                                                                                                                                                                         |
| :------------ | :------- | :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| endpoint\*    | `enum`   | <https://otlp.nr-data.net> | Endpoint where the exporter sends data to New Relic. Endpoints are region-specific, so use the one according to where your account is based. Valid values are `<https://otlp.nr-data.net`>, `<https://otlp.eu01.nr-data.net`>, or `<https://gov-otlp.nr-data.net`>. |
| license_key\* | `string` |                            | License key used for data ingest.                                                                                                                                                                                                                                   |

<span style="color:red">\*_required field_</span>