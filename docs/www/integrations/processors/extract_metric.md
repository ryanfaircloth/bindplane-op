---
title: Extract Metric
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cddea2a005d14423a
slug: extract-metric
hidden: false
---

## Extract Metric Processor

The Extract Metric Processor can look at all logs matching a filter, extract a numerical value from a field, and then create a metric with that value.
Both the name and units of the created metric can be configured. Additionally, fields from matching logs can be preserved as metric attributes.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| | ✓ | |

## Supported Agent Versions

`v1.14.0`+

## Configuration
| Field        | Type     | Default | Description |
| ---          | ---      | ---     | ---         |
| match        | `string`   | `true`  | A boolean [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) used to match which logs to count. By default, all logs are counted. |
| path        | `string`   |   | A boolean [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) used to specify the field to extract from a matching log. |
| metric_name  | `string`   | `log.count` | The name of the metric created. |
| metric_units | `string`   | `{logs}`    | The unit of the metric created. See [Unified Code for Units of Measure](https://unitsofmeasure.org/ucum#section-Alphabetic-Index-By-Name) for available units. |
| attributes   | `map`      | `{}`        | The mapped attributes of the metric created. Each key is an attribute name. Each value is an [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) that extracts data from the log. |


## Expression Language
In order to match or extract values from logs, the following `keys` are reserved and can be used to traverse the logs data model.

| Key               | Description |
| ---               | ---   |
| `body`            | Used to access the body of the log. |
| `attributes`      | Used to access the attributes of the log. |
| `resource`        | Used to access the resource of the log. |
| `severity_enum`   | Used to access the severity enum of the log. |
| `severity_number` | Used to access the severity number of the log. |

In order to access embedded values, use JSON dot notation. For example, `body.example.field` can be used to access a field two levels deep on the log body.

However, if a key already possesses a literal dot, users will need to use bracket notation to access that field. For example, when the field `service.name` exists on the log's resource, users will need to use `resource["service.name"]` to access this value.

For more information about syntax and available operators, see the [Expression Language Definition](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md).

---

## Example Configurations

### Default Configuration

By default, all logs collected by the source will be counted, with the value used to create a new metric called `log.count` with the unit of `{logs}`.

### Break Down HTTP Request Durations by Status

In this configuration, we want to parse our HTTP server logs to create metrics representing how long each request took, broken down by status code. Our logs are JSON with the following structure:

```JSON
{
  "level": "warn",
  "host": "10.0.10.0",
  "datetime":"2022-12-02T10:21",
  "duration": 122,
  "method": "POST",
  "request": "/api/v1/apply",
  "protocol": "HTTP/1.1",
  "status": 200
}
```

The match expression will exclude all logs without a status code in its body:

```expr
body.duration != nil
```

Our path expression will be the path to the `duration` field of the body, which we know is the request duration in milliseconds.

```expr
body.duration
```

We'll name this metric `http.request.duration`, then we'll use the status code for the `status_code` metric attribute on the created metric:

```yaml
attributes:
  status_code: body.status
```

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/extract-metric/extract_metric.png" width="1000px" alt="The processor configured with the values specified above.">
