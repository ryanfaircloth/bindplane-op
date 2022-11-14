---
title: "SignalFx"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "signalfx"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Configuration Table

| Parameter      | Type     | Default | Description                                                                                                                          |
| :------------- | :------- | :------ | :----------------------------------------------------------------------------------------------------------------------------------- |
| token\*        | `string` |         | Token used to authenticate with the Splunk (SignalFx) metric, trace, and log APIs                                                    |
| realm          | `enum`   | us0     | The Splunk API realm (region) to use when sending metrics, traces, and logs. Valid values are: `us0`, `us1`, `us2`, `eu0`, or `jp0`. |
| enable_metrics | `bool`   | true    |                                                                                                                                      |
| enable_logs    | `bool`   | true    |                                                                                                                                      |
| enable_traces  | `bool`   | true    |                                                                                                                                      |

<span style="color:red">\*_required field_</span>
