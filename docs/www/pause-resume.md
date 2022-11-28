---
title: Pausing Sources and Destinations
excerpt: Pause the collection or reporting of telemetry
category: 636c08d51eb043000f8ce20e
slug: pause-resume
hidden: false
---

# Pausing Sources and Destinations

By default, telemetry is collected from all sources in a configuration, and sent to all supported destinations in that configuration. However, there may be times when you don't want to collect from a certain source or to export to a destination. To support those situations, BindPlane allows you to "Pause" sources and destinations. When a source is paused, agents simply won't try to collect from it. When a destination is paused, no collected telemetry will be sent to it.

An important detail to remember is that pausing or resuming a source or destination in a configuration will update all agents using that configuration to pause/resume that source/destination.

## Pausing a Source

Sources can be paused from the page of either an Agent or its configuration. To pause a source, click on the card for the source you want to pause in the topology view. Its current status will be shown in the bottom left corner, either `Running` or `Paused`. If running, clicking the `Pause` button will pause collection of that source.

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/pause-resume/pause_source_prompt.png" width="1000px" alt="The 'Edit Source' menu shows the source is running and has a button to pause it.">

After clicking "Pause", the topology reflects that the Redis source has been paused.

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/pause-resume/paused_source_topology.png" width="1000px" alt="The configuration's topology shows the Redis source is now paused.">


## Pausing a Destination

Destinations can be paused just like sources, by clicking on the appropriate card in the topology view of a configuration. Its current status will be shown in the bottom left corner, either `Running` or `Paused`. If running, clicking the `Pause` button will pause the sending of telemetry to that destination. If the only destination in a configuration is paused, the agent will also pause collecting all sources as the telemetry has nowhere to go.

A major difference between pausing sources and destinations is that, while pausing a source only affects that configuration, pausing a destination will pause it in all configurations including it. For example, imagine a single Google Cloud Monitoring destination, `example-gcp-project`, is being used in several configurations to send telemetry to that GCP project. If you need to stop all telemetry from being sent, pausing the `example-gcp-project` destination in one configuration will pause it in all other configurations.

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/pause-resume/pause_destination_prompt.png" width="1000px" alt="The 'Edit Destination' menu shows the destination is running and has a button to pause it.">

After clicking "Pause", the topology reflects that the `otlphttp` destination has been paused.

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/pause-resume/paused_destination_topology.png" width="1000px" alt="The configuration's topology shows the 'otlphttp' destination is now paused.">
