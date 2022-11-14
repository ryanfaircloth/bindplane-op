---
title: "Prometheus"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0cd9f114009de9ba78
slug: "prometheus-1"
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       |      |        |

## Configuration Table

| Parameter      | Type     | Default     | Description                                                                                   |
| :------------- | :------- | :---------- | :-------------------------------------------------------------------------------------------- |
| listen_port    | `int`    | 9000        | The TCP port the Prometheus exporter should listen on, to be scraped by a Prometheus server   |
| listen_address | `string` | "127.0.0.1" | The IP address the Prometheus exporter should listen on, to be scraped by a Prometheus server |
| namespace      | `string` |             | When set, exports metrics under the provided value                                            |
