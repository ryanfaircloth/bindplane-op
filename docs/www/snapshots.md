---
title: Snapshots
excerpt: Utilize Snapshots in BindPlane OP
category: 633dd7654359a20031653088
slug: snapshots
hidden: false
---

# Snapshots

Snapshots provide a way to view logs, metrics, and traces recently collected by an agent.

## Viewing

Snapshots can be viewed by clicking the `View Recent Telemetry` button found below the details table on an agent page.

### Logs

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/snapshot_logs_example.png" width="1000px" alt="PostgreSQL logs in the snapshot view">

### Metrics

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/snapshot_metrics_example.png" width="1000px" alt="PostgreSQL metrics in the snapshot view">

If no metrics are available, it's possible none have been collected yet. For example, if an agent has been running
for 30 seconds and collecting metrics with a 60 second interval, it won't have any recent metrics to show.

## Limitations

For recent telemetry to be available for an agent, there are two requirements:

1. The agent must be connected to BindPlane
2. The agent must be using a managed configuration
