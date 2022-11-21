---
title: "AWS Cloudwatch"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "aws-cloudwatch"
hidden: false
---

## Prerequisites

While installing AWS CLI is not required in order to collect logs using the AWS Cloudwatch source type, it still provides an easier means of authentication. With AWS Cloudwatch, users are required to provide some form of authentication. This can be either profile credentials or environment variables that provide access keys for user accounts. Setting these credentials up can prove tedious and confusing, luckily AWS CLI can generate the user's profile / credentials with the `aws configure` command. Here is the [AWS CLI Getting Started](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) Guide which outlines a few prerequisites.

Some prerequisites:

1. Creating an IAM user account
   - Required Permissions
     - `logs:GetLogEvents`
     - `logs:DescribeLogGroups`
     - `logs:DescribeLogStreams`
   - User does not require console access
2. Create an access key ID and secret access key

---

**NOTE**
In the [AWS CLI getting started guide](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html), it will instruct you to install for your current user or all users.
The observIQ OTEL Collector runs [as root](https://github.com/observIQ/observiq-otel-collector/blob/main/docs/installation-linux.md#configuring-the-collector) by default, meaning the AWS CLI and credentials should be installed under the collector system's root account.

---

### Credentials

### Credential and Config Files

AWS Authentication utilizes a user profile specified in the user's home directory at `.aws/credentials`. Each profiles' credentials should include at minimum
the profile name, access key, and secret access key.

```
[profile_name]
aws_access_key_id=******
aws_secret_access_key=******
```

In addition to the credentials file, there is also a `.aws/config`. This includes less sensitvie configuration options such as the region, output format, etc. A typical entry in the `config` file should look as such.

```
[profile_name]
region=us-west-2
output=json
```

More information on [AWS Configuration and Credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

### Environment Variables

Alternatively, [AWS Environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) can be specified to override a credentials file.
You can modify the collector's environment variables by configuring a Systemd override.
Run `sudo systemctl edit observiq-otel-collector` and add your access key, secret key, and region.

```
[Service]
Environment=AWS_ACCESS_KEY_ID=******
Environment=AWS_SECRET_ACCESS_KEY=******
Environment=AWS_DEFAULT_REGION=us-west-2
```

Then restart the collector after making environment changes.

## Setup

1. Once logged in, Select the `Configs` tab at the top of the Bindplane Home page.
2. Select a pre-existing config or create a new one.
3. Add a new source and select `AWS Cloudwatch`.
4. After configuring credentials using the AWS CLI on the collector system, using the default values in the source form should enable the collector to collect logs from Cloudwatch. If credentials were configured using environment variables, you will need to leave the `Profile` field blank
5. Click save configuration.
6. Add Destination type of your choosing.
7. Apply Configuration to desired Agent.
8. Voila. Logs should be collecting

## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter              | Type                      | Default      | Description                                                                                                       |
| :--------------------- | :------------------------ | :----------- | :---------------------------------------------------------------------------------------------------------------- |
| region\*               | `enum`                    | us-east-1    | The AWS recognized region string.                                                                                 |
| profile \*             | `string`                  | "default"    | The AWS profile used to authenticate, if none is specified the default is chosen from the list of profiles.       |
| credential_type        | `enum`                    | profile      | Determines whether to pull credentials from a `credentials` file or use environment variables for authentication. |
| discovery_type         | `enum`                    | AutoDiscover | Configuration for Log Groups, by default all Log Groups and Log Streams will be collected.                        |
| limit                  | `int`                     | 50           | Limits the number of discovered log groups.                                                                       |
| prefix                 | `string`                  | ""           | A prefix for log groups to limit the number of log groups discovered.                                             |
| names                  | `strings`                 | []           | A list of full log stream names to filter the discovered log groups to collect from.                              |
| prefixes               | `strings`                 | []           | A list of prefixes to filter the discovered log groups to collect from.                                           |
| named_groups           | `awsCloudwatchNamedField` | []           | Configuration for Log Groups, by default all Log Groups and Log Streams will be collected.                        |
| imds_endpoint          | `string`                  | ""           | A way of specifying a custom URL to be used by the EC2 IMDS client to validate the session.                       |
| poll_interval          | `int`                     | 1            | The duration waiting in between requests (minutes).                                                               |
| max_events_per_request | `int`                     | 50           | The maximum number of events to process per request to Cloudwatch.                                                |

<span style="color:red">\*_required field_</span>

## Discovery Type

### Default Settings

When starting with an AWS Cloudwatch set to its default values, you should see log collection from all logs groups with no filtering of log streams. The default polling interval for Cloudwatch is 1 minute, so there may be a delay before seeing any logs coming through.

### AutoDiscover

When using Discovery Type `AutoDiscover` there are some optional parameters that can be added to / filter the amount of logs collected.

- `limit`: limits the number of discovered log groups(default = 50).
- `prefix`: Prefix for log groups to limit the number of log groups discovered
  - `prefix: /aws/eks/`
  - if omitted, all log groups up to the limit will be collected from.
- `names`: A list of full log stream names to filter the discovered log groups to collect from.
  - `names: [kube-apiserver-ea9c831555adca1815ae04b87661klasdj]`
- `prefixes`: A list of log stream prefixes to filter the discovered log groups to collect from.
  - `prefixes: [kube-api-controller]`

### Named

This Discovery Type filters logs by listing only the desired log groups to collect from and omitting any other log groups. When selecting this Discovery Type, at least one log group is required otherwise no logs would be collected. When listing log groups the `ID` field of each log group instance should match the full name of the group.

Additionally, `Named` also provides `prefixes` and `names` parameters for each listed log group that filters out the listed log streams. These parameters should be listed underneath each log group's ID's as they are unique to each individual log group.
