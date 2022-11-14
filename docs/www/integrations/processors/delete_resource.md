---
title: Delete Resource
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cddea2a005d14423a
slug: delete-resource
hidden: false
---

## Delete Resource Processor

The Delete Resource processor can be used to remove resources from metrics, traces, and logs.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| ✓ | ✓ | ✓ |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| resources | `strings`   | required | One or more resource names to remove. |
| telemetry_types | `enums` | The telemetry types to remove resource from (Metrics, Traces, Logs). |


## Example Configuration

Remove the following resources:
- environment
- location

**Web Interface**

![delete_resource](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/delete_resource.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: delete-resources
  name: delete-resources
spec:
  type: delete_resource
  parameters:
    - name: resources
      value:
        - environment
        - location
    - name: telemetry_types
      value:
        - Metrics
        - Traces
        - Logs
```

**Configuration with Embedded Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Configuration
metadata:
  id: delete-resource
  name: delete-resource
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
        - type: delete_resource
          parameters:
            - name: resources
              value:
                - environment
                - location
            - name: telemetry_types
              value:
                - Metrics
                - Traces
                - Logs
  selector:
    matchLabels:
      configuration: delete-resource
```
