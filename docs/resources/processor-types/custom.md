## Custom Processor

The Custom processor can be used to inject a custom processor configuration into a pipeline. A list of supported processors can be found [here](https://github.com/observIQ/observiq-otel-collector/blob/main/docs/processors.md).

The Custom processor is useful for solving usecases not covered by BindPlane OP's other processor types.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| ✓ | ✓ | ✓ |

The custom processor type can support all telemetry types, however, it is up to the user to enable / disable the correct types
based on the processor being used.

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| configuration     | `yaml`  | required | Enter any supported Processor and the YAML will be inserted into the configuration. |
| telemetry_types   | `enums` | The pipelines the processor should be used with. |

## Example Configuration

Inject the following [resource processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor) configuration:

```yaml
resource:
  attributes:
    - action: upsert
      key: custom
      value: true
```

**Web Interface**

![custom](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/custom.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: custom
  name: custom
spec:
  type: custom
  parameters:
    - name: configuration
      value: |
        resource:
          attributes:
            - action: upsert
              key: custom
              value: true
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
  id: custom
  name: custom
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
        - type: custom
          parameters:
            - name: configuration
              value: |
                resource:
                  attributes:
                    - action: upsert
                      key: custom
                      value: true
            - name: telemetry_types
              value:
                - Metrics
                - Traces
                - Logs
  selector:
    matchLabels:
      configuration: custom
```
