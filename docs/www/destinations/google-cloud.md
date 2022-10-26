---
title: Google Cloud
category: 633dd7654359a20031653089
slug: google-cloud
hidden: false
---

## Supported Types

| Metrics | Logs | Traces |
| :------ | :--- | :----- |
| ✓       | ✓    | ✓      |

## Prerequisites

### Set up credentials

1. Enable billing in your GCP project.
2. Enable the Cloud Metrics and Cloud Trace APIs.
3. Ensure that your user GCP user has (at minimum) `roles/monitoring.metricWriter` and `roles/cloudtrace.agent`. You can learn about [metric-related](https://cloud.google.com/monitoring/access-control) and [trace-related](https://cloud.google.com/trace/docs/iam) IAM in the GCP documentation.
4. Obtain credentials.

```sh
gcloud auth application-default login
```



Alternatives

- You can run the collector as a service account, as long as it has the necessary roles. This is useful in production, because credentials for a user are short-lived.

- You can also run the collector on a GCE VM or as a GKE workload, which will use the service account associated with GCE/GKE.

## Configuration Table

| Parameter        | Type     | Default | Description                                                                                                                                                                                                                                                                                                                                                                             |
| :--------------- | :------- | :------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| project          | `string` |         | The Google Cloud Project ID to send logs, metrics, and traces to.                                                                                                                                                                                                                                                                                                                       |
| auth_type        | `enum`   | auto    | The method used for authenticating to Google Cloud. 'auto' will attempt to use the collector's environment, useful when running on Google Cloud or when you have set GOOGLE_APPLICATION_CREDENTIALS in the collector's environment. 'json' takes the json contents of a Google Service Account's credentials file. 'file' is the file path to a Google Service Account credential file. |
| credentials      | `string` |         | JSON value from a Google Service Account credential file.                                                                                                                                                                                                                                                                                                                               |
| credentials_file | `string` |         | Path to a Google Service Account credential file on the collector system. The collector's runtime user must have permission to read this file.                                                                                                                                                                                                                                          |