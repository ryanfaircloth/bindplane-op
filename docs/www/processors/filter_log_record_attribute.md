---
title: Filter Log Record Attribute
category: 633dd7654359a20031653089
slug: filter-log-record-attribute
hidden: false
---

## Log Record Attribute Filter Processor

The Log Record Attribute Filter processor can be used to include or exclude logs based on matched attributes.

## Supported Types

| Metrics | Logs | Traces |
| :--- | :--- | :--- |
|  | ✓ |  |

## Configuration Table

| Parameter  | Type    | Default   | Description |
| :---       | :---    | :---      | :--- |
| action     | `enum`  | `exclude` | Whether to `include` (retain) or `exclude` (drop) matches. |
| match_type | `enum`  | `strict`  | Method for matching values. Strict matching requires that 'value' be an exact match. Regexp matching uses [re2](https://github.com/google/re2/wiki/Syntax) to match a value. |
| attributes | `map`   | required  | One or more key (attribute name) value (attribute value) pairs to filter on. Logs are filtered if all pairs are matched. |

## Example Configuration

### Web Interface

![filter_log_record_attribute](https://storage.googleapis.com/bindplane-op-doc-images/resources/processor-types/filter_log_record_attribute.png)

### Strict Exclude

Exclude logs that have the following attributes:
- environment: dev
- location: us-east1-b

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: exclude-attributes
  name: exclude-attributes
spec:
  type: filter_log_record_attribute
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

Include logs that have the following attributes. All other logs will be filtered out:
- environment: dev
- location: us-east1-b

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: include-attributes
  name: include-attributes
spec:
  type: filter_log_record_attribute
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

Exclude logs that have the following attribute with any value:
- env: *

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Processor
metadata:
  id: exclude-attributes
  name: exclude-attributes
spec:
  type: filter_log_record_attribute
  parameters:
    - name: action
      value: exclude
    - name: match_type
      value: regexp
    - name: attributes
      value:
        env: "*"
```
