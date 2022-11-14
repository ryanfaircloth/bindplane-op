---
title: Filter Resource Attribute
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cddea2a005d14423a
slug: filter-resource-attribute
hidden: false
---

## Resource Attribute Filter Processor

The Resource Attribute Filter processor can be used to include or exclude logs based on matched resources.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
| ✓ | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default   | Description |
| :---       | :---    | :---      | :--- |
| action     | `enum`  | `exclude` | Whether to `include` (retain) or `exclude` (drop) matches. |
| match_type | `enum`  | `strict`  | Method for matching values. Strict matching requires that 'value' be an exact match. Regexp matching uses [re2](https://github.com/google/re2/wiki/Syntax) to match a value. |
| attributes | `map`   | required  | One or more key (resource name) value (resource value) pairs to filter on. |

## Example Configuration

### Web Interface

![filter_resource_attribute](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/filter_resource_attribute.png)

### Strict Exclude

Exclude metrics and logs that have the following resources:
- environment: dev
- location: us-east1-b

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: exclude-resources
  name: exclude-resources
spec:
  type: filter_resource_record_attribute
  parameters:
    - name: action
      value: exclude
    - name: match_type
      value: strict
    - name: attributes
      value:
        environment: dev
        location: us-east1-b
```

### Strict Include

Include metrics and logs that have the following resources:
- environment: dev
- location: us-east1-b

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: include-resources
  name: include-resources
spec:
  type: filter_resource_record_attribute
  parameters:
    - name: action
      value: exclude
    - name: match_type
      value: strict
    - name: attributes
      value:
        environment: dev
        location: us-east1-b
```

### Regexp Exclude

Exclude metrics and logs that have the following attribute with any value:
- env: *

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: exclude-resources
  name: exclude-resources
spec:
  type: filter_resource_record_attribute
  parameters:
    - name: action
      value: exclude
    - name: match_type
      value: regexp
    - name: attributes
      value:
        env: "*"
```
