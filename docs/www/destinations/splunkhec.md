---
title: "Splunk HTTP Event Collector (HEC)"
category: 633dd7654359a20031653089
slug: "splunk-hec"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
|         | ✓    |        |

## Prerequisites

### Setup Splunk HTTP Event Collector

See the "Configure HTTP Event Collector on Splunk Enterprise" section in [Splunk's documentation](https://docs.splunk.com/Documentation/Splunk/9.0.1/Data/UsetheHTTPEventCollector). Configure a "Event Collector token", which will be used during the destination type's configuration.

## Configuration Table

| Parameter        | Type     | Default     | Description |
| :--------------- | :------- | :---------- |:----------- |
| token            | `string` | <required>  | Token used to authenticate with the Event Collector. |
| index            | `string` |             | The index to send logs to (optional). |
| hostname         | `string` | `localhost` | The address or hostname of the Event Collector. |
| port             | `int`    | `8088`      | The TCP port the Event Collector is listening on. |
| path             | `string` | `/services/collector/event` | The HTTP API path the Event Collector is accepting events on. |
| enable_tls       | `bool`   | `false`     | Whether or not to connect to the Event Collector using TLS. |
| insecure_skip_verify | `bool` | `false`     | Whether or not to skip TLS certificate verification. |
| ca_file          | `string` |             | The certificate authority file to use when verifying the Event Collector's TLS certificate (optional). |

## Example Configurations

**Required Values**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Destination
metadata:
    id: hec-basic
    name: hec-basic
    labels:
        platform: linux
spec:
    type: splunkhec
    parameters:
        - name: token
          value: 00000-00000-00000
```


**TLS**

```yaml
apiVersion: bindplane.observiq.com/v1
kind: Destination
metadata:
    id: hec-tls
    name: hec-tls
    labels:
        platform: linux
spec:
    type: splunkhec
    parameters:
        - name: token
          value: 00000-00000-00000
        - name: index
          value: otel
        - name: hostname
          value: hec.corp.net
        - name: port
          value: 8088
        - name: path
          value: /services/collector/event
        - name: enable_tls
          value: true
        - name: insecure_skip_verify
          value: false
        - name: ca_file
          value: "/opt/tls/hec-ca.crt"
```
