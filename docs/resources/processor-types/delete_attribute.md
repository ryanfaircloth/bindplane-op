## Delete Log Record Attribute Processor

The Delete Log Record Attribute processor can be used to remove attributes from all
[log records](https://opentelemetry.io/docs/reference/specification/logs/data-model/#log-and-event-record-definition) in the pipeline.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
|  | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| attributes | `strings`   | required | One or more attribute names to remove from all log records. |

## Example Configuration

Remove the following attributes:
- environment
- location

**Web Interface**

![delete_attribute](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/delete_attribute.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: delete-attributes
  name: delete-attributes
spec:
  type: delete_attribute
  parameters:
    - name: attributes
      value:
        - environment
        - location
```

**Configuration with Embeded Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Configuration
metadata:
  id: delete-attribute
  name: delete-attribute
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
        - type: delete_attribute
          parameters:
            - name: attributes
              value:
                - environment
                - location
  selector:
    matchLabels:
      configuration: delete-attribute
```
