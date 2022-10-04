## Log Sampling Processor

The Log Sampling processor can be used to filter out logs with a configured "drop ratio".

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
|  | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| drop_ratio   | `enum`  | `"0.50"` | The probability an entry is dropped (used for sampling). A value of 1.0 will drop 100% of matching entries, while a value of 0.0 will drop 0%. |

Valud drop ratio's range from "0.0" (0%) to "1.00" (100%) with 5% increments. Note that the drop ratio value is a string.

Valid drop ratio values:
- "1.00"
- "0.95"
- "0.90"
- "0.85"
- "0.80"
- "0.75"
- "0.70"
- "0.65"
- "0.60"
- "0.55"
- "0.50"
- "0.45"
- "0.40"
- "0.35"
- "0.30"
- "0.25"
- "0.20"
- "0.15"
- "0.10"
- "0.05"
- "0.0"

## Example Configuration

Filter out 65%.

**Web Interface**

![sample_log](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/sample_log.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: sampling
  name: sampling
spec:
  type: sampling
  parameters:
    - name: drop_ratio
      value: "0.65"
```
