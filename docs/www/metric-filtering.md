---
title: Metric Filtering
category: 633dd7654359a20031653088
slug: metric-filtering
excerpt: Filter metrics with BindPlane OP
hidden: false
---
# Metric Filtering

Sources that collect metrics can be configured to filter out any number of metrics.
When a metric is filtered out, it will not be sent to any destination.

## Configuration

Metric filtering can be configured for a source when initially creating it, or an existing source
can be edited to change its filtering.
Once saved all agents collecting that source will be updated to use the updated filter settings.

The controls for filtering metrics are found in the `Advanced` section of the source configuration form.
Available metrics are organized in groups, allowing you to quickly enable or disable filtering for
multiple related metrics. The checkbox next to each metric indicates whether it will be sent to destinations.
That is, unchecking the box for a metric will filter it out.

<img src="https://storage.googleapis.com/bindplane-op-doc-images/guides/metric_filtering_example.png" width="700px" alt="Configuring F5 Big-IP filtering">

In the above image, all metrics in the `Virtual Server` group will be filtered out, as well as `bigip.pool.availability` and `bigip.pool.packet.count`.
