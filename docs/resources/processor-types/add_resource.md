## Add Resource Attribute Processor

The Add Resource Attribute processor can be used to enrich telemetry by adding resources to all metrics, traces, and logs in the pipeline.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| ✓ | ✓ | ✓ |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| action     | `enum`  | `upsert` | insert: Add resource if it does not exist. update: Update existing value. upsert: Insert or update. |
| resources  | `map`   | required | One or more key (resource name) value (resource value) pairs to add as resources. |
| telemetry_types | `enums` | The telemetry types to add resources to (Metrics, Traces, Logs). |

## Example Configuration

Add the following key value pairs as resources to metrics, traces and logs:
- office: gr
- user: guest

**Web Interface**

![add_resource](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/add_resource.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: add-resources
  name: add-resources
spec:
  type: add_resource
  parameters:
    - name: action
      value: upsert
    - name: resources
      value:
        office: gr
        user: guest
    - name: telemetry_types
      value:
        - Metrics
        - Traces
        - Logs
```

**Configuration with Embeded Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Configuration
metadata:
  id: add-resource
  name: add-resource
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
        - type: add_resource
          parameters:
            - name: action
              value: upsert
            - name: resources
              value:
                office: gr
                user: guest
            - name: telemetry_types
              value:
                - Metrics
                - Traces
                - Logs
  selector:
    matchLabels:
      configuration: add-resource
```
