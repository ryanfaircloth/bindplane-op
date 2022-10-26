---
title: "VMware vCenter"
category: 633dd7654359a20031653089
slug: "vmware-vcenter"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This receiver has been built to support ESXi and vCenter versions:

- 7.5
- 7.0
- 6.7

A “Read Only” user assigned to a vSphere with permissions to the vCenter server, cluster and all subsequent resources being monitored must be specified in order for the receiver to retrieve information about them.

## Configuration Table

| Parameter           | Type     | Default                                           | Description                                       |
| :------------------ | :------- | :------------------------------------------------ | :------------------------------------------------ |
| endpoint* | `string`    |                       | Endpoint to the vCenter Server or ESXi host that has the sdk path enabled. Required. The expected format is `<protocol>://<hostname>`  \n  \ni.e: `https://vcsa.hostname.localnet` |
| username* | `string` |                                       | Username used to authenticate. |
| password* | `string`    |          | Password used to authenticate. |
| tls | `string` |  | Not Required. Will use defaults for [configtls.TLSClientSetting](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md). By default insecure settings are rejected and certificate verification is on. |
| collection_interval | `int` | 60 | How often (seconds) to scrape for metrics. |

<span style="color:red">\*_required field_</span>