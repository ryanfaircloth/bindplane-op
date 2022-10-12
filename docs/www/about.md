---
title: "About"
slug: "about"
hidden: false
createdAt: "2022-06-01T15:49:31.618Z"
updatedAt: "2022-08-11T14:49:32.906Z"
---

![](https://files.readme.io/48f216c-bindplaneop.png "bindplaneop.png")



## What is BindPlane OP?

BindPlane OP is an open source observability pipeline that gives you the ability to collect, refine, and ship metrics, logs, and traces to any destination. BindPlane OP provides the controls you need to reduce observability costs and simplify the deployment and management of telemetry agents at scale.  

To get started with BindPlane, see our [Getting Started](doc:getting-started) guide or check out the [Installation](doc:installation) page. Also, checkout the BindPlane OP GitHub repo [here](https://github.com/observIQ/bindplane-op).

![](https://files.readme.io/1659f34-BindPlane_Architecture_Diagram.jpg "BindPlane_Architecture_Diagram.jpg")



## Features

- Manage the lifecycle of telemetry agents, starting with the [observIQ Distro for OpenTelemetry Collector](https://github.com/observIQ/observiq-otel-collector)
- Build, deploy, and manage telemetry configurations for different Sources and deploy them to your agents
- Ship metric, log, and trace data to one or many Destinations
- Utilize flow controls to adjust the flow of your data in realtime

## Architecture

![](https://files.readme.io/c1a71e4-BindPlane_OP_Architecture_Diagram.png "BindPlane OP Architecture Diagram.png")



BindPlane OP is a lightweight web server (no dependencies) that you can deploy anywhere in your environment. It's composed of the following components:

- GraphQL Server: provides configuration and agent details via GraphQL
- REST Server: BindPlane CLI and UI make requests to the server via REST
- WebSocket Server: Agents connect to receive configuration updates via [OpAMP](https://github.com/open-telemetry/opamp-spec)
- Store: pluggable storage manages configuration and Agent state 
- Manager: dispatches configuration changes to Agents