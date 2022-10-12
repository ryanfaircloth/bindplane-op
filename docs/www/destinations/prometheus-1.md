---
title: "Prometheus"
slug: "prometheus-1"
hidden: false
createdAt: "2022-06-07T18:47:19.095Z"
updatedAt: "2022-08-10T15:02:44.293Z"
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