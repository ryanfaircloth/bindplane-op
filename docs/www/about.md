---
title: About
category: 636c08cbfcc97f00717832c4
slug: about
hidden: false
---

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/bindplaneop.png" width="1000px" alt="bindplaneop.png">

## What is BindPlane OP?

BindPlane OP is an open source observability pipeline that gives you the ability to collect, refine, and ship metrics, logs, and traces to any destination. BindPlane OP provides the controls you need to reduce observability costs and simplify the deployment and management of telemetry agents at scale.  

To get started with BindPlane, see our [Getting Started](doc:getting-started) guide or check out the [Installation](doc:installation) page. Also, checkout the BindPlane OP GitHub repo [here](https://github.com/observIQ/bindplane-op).

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/BindPlane_Architecture_Diagram.jpg" width="1000px" alt="BindPlane_Architecture_Diagram.jpg">

## Features

- Manage the lifecycle of telemetry agents, starting with the [observIQ Distro for OpenTelemetry Collector](https://github.com/observIQ/observiq-otel-collector)
- Build, deploy, and manage telemetry configurations for different Sources and deploy them to your agents
- Ship metric, log, and trace data to one or many Destinations
- Utilize flow controls to adjust the flow of your data in realtime

## Architecture

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/BindPlane_OP_Architecture_Diagram.png" width="1000px" alt="BindPlane OP Architecture Diagram.png">

BindPlane OP is a lightweight web server (no dependencies) that you can deploy anywhere in your environment. It's composed of the following components:

- GraphQL Server: provides configuration and agent details via GraphQL
- REST Server: BindPlane CLI and UI make requests to the server via REST
- WebSocket Server: Agents connect to receive configuration updates via [OpAMP](https://github.com/open-telemetry/opamp-spec)
- Store: pluggable storage manages configuration and Agent state 
- Manager: dispatches configuration changes to Agents
