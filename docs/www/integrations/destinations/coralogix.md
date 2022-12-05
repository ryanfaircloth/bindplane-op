---
title: "Coralogix"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "coralogix"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

| Parameter                   | Type      | Default                          | Description                                                                                                                                                                                                                                                          |
| :-------------------------- | :-------- | :------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| private_key\*               | `string`  | `""`                             | API Private Key. More information on finding your key can be found [here](https://coralogix.com/docs/private-key/).                                                                                    |
| application_name\*          | `string`  | `""`                             | OTel objects that are sent to Coralogix are tagged wit this Application Name. More on application names [here](https://coralogix.com/docs/application-and-subsystem-names).                            |
| region\*                    | `string`  | `EUROPE1`                        | Region of your Coralogix account associated with the provided `private_key`. See the reference table to see [telemetry ingress endpoints](#coralogix-region-ingress-endpoints) related to each region. |
| enable_metrics              | `bool`    | true                             |                                                                                                                                                                                                        |
| enable_logs                 | `bool`    | true                             |                                                                                                                                                                                                        |
| enable_traces               | `bool`    | true                             |                                                                                                                                                                                                        |
| subsystem_name              | `string`  | `""`                             | OTel objects that are sent to Coralogix are tagged wit this Subsystem Name. More on application names [here](https://coralogix.com/docs/application-and-subsystem-names).                              |
| timeout                     | `int`     | 5                                | Seconds to wait per individual attempt to send data to a backend                                                                                                                                       |
| resource_attributes         | `bool`    | `false`                          |                                                                                                                                                                                                        |
| application_name_attributes | `strings` | `[]`                             | Ordered list of resource attributes that are used for Coralogix AppName.                                                                                                                               |
| subsystem_name_attributes   | `strings` | `[]`                             | Ordered list of resource attributes that are used for Coralogix SubSystem.                                                                                                                             |

<span style="color:red">\*_required field_</span>

## Coralogix Region Ingress Endpoints

| Region  | Traces Endpoint                     | Metrics Endpoint                     | Logs Endpoint                     |
| :------ | :---------------------------------- | :----------------------------------- | :-------------------------------- |
| USA1    | `otel-traces.coralogix.us:443`      | `otel-metrics.coralogix.us:443`      | `otel-logs.coralogix.us:443`      |
| APAC1   | `otel-traces.app.coralogix.in:443`  | `otel-metrics.coralogix.in:443`      | `otel-logs.coralogix.in:443`      | 
| APAC2   | `otel-traces.coralogixsg.com:443`   | `otel-metrics.coralogixsg.com:443`   | `otel-logs.coralogixsg.com:443`   |
| EUROPE1 | `otel-traces.coralogix.com:443`     | `otel-metrics.coralogix.com:443`     | `otel-logs.coralogix.com:443`     |
| EUROPE2 | `otel-traces.eu2.coralogix.com:443` | `otel-metrics.eu2.coralogix.com:443` | `otel-logs.eu2.coralogix.com:443` |
