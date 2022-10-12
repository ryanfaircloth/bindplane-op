## Severity Filter Processor

The Severity Filter processor can be used to filter out logs that do not meet a given severity threshold.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
|  | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| severity   | `enum`  | `TRACE` | Minimum severity to match. Log entries with lower severities will be filtered. |

Valud severity levels:
- TRACE
- INFO
- WARN
- ERROR
- FATAL

## Example Configuration

Filter out INFO and TRACE logs.

**Web Interface**

![filter_severity](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/filter_severity.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: severity-filter
  name: severity-filter
spec:
  type: filter_severity
  parameters:
    - name: severity
      value: WARN
```

**Configuration with Embedded Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Configuration
metadata:
  id: severity-filter
  name: severity-filter
  labels:
    platform: linux
spec:
  sources:
    - type: journald
      parameters:
        - name: units
          value: []
        - name: directory
          value: ""
        - name: priority
          value: info
        - name: start_at
          value: end
      processors:
        - type: filter_severity
          parameters:
            - name: severity
              value: WARN
  selector:
    matchLabels:
      configuration: severity-filter
```
