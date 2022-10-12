## Metric Name Filter Processor

The Metric Name Filter processor can be used to include or exclude metrics based on their name.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| ✓ |  |  |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| action     | `enum`  | `exclude` | `eclude` or `include` metrics that match. |
| match_type | `enum`  | `strict`  | Method for matching values. Strict matching requires that 'value' be an exact match. Regexp matching uses [re2](https://github.com/google/re2/wiki/Syntax) to match a value. |
| metric_names  | `strings`   | required | One or more metric names to match on. |

## Example Configuration

### Web Interface

![filter_metric_name](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/filter_metric_name.png)

### Exclude Regexp

Filter out (exclude) metrics that match the expression `k8s.node.*`.

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: filter-name-regexp
  name: filter-name-regexp
spec:
  type: filter_metric_name
  parameters:
    - name: action
      value: exclude
    - name: match_type
      value: regexp
    - name: metric_names
      value:
        - k8s.node.*
```

### Include Strict

Include metrics that match, drop all other metrics.

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: include-name-strict
  name: include-name-strict
spec:
  type: filter_metric_name
  parameters:
    - name: action
      value: include
    - name: match_type
      value: strict
    - name: metric_names
      value:
        - k8s.container.cpu
        - k8s.pod.memory
```
