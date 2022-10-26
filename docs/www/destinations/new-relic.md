---
title: "New Relic"
category: 633dd7654359a20031653089
slug: "new-relic"
hidden: false
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