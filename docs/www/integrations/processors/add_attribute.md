---
title: Add Attribute
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cddea2a005d14423a
slug: add-attribute
hidden: false
---

## Add Log Record Attribute Processor

The Add Log Record Attribute processor can be used to enrich logs by adding attributes to all
[log records](https://opentelemetry.io/docs/reference/specification/logs/data-model/#log-and-event-record-definition) in the pipeline.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
|  | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default  | Description |
| :---       | :---    | :---     | :--- |
| action     | `enum`  | `upsert` | `insert`: Add attribute if it does not exist. `update`: Update existing value. `upsert`: Insert or update. |
| attributes | `map`   | required | One or more key (attribute name) value (attribute value) pairs to add to all log records as attributes. |

## Example Configuration

Add the following key value pairs as attributes:
- environment: dev
- location: us-east1-b

**Web Interface**

![add_attribute](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/add_attribute.png)

**Standalone Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: add-attributes
  name: add-attributes
spec:
  type: add_attribute
  parameters:
    - name: action
      value: upsert
    - name: attributes
      value:
        environment: dev
        location: us-east1-b
```

**Configuration with Embedded Processor**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Configuration
metadata:
  id: add-attribute
  name: add-attribute
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
        - type: add_attribute
          parameters:
            - name: action
              value: upsert
            - name: attributes
              value:
                environment: dev
                location: us-east1-b
  selector:
    matchLabels:
      configuration: add-attribute
```
