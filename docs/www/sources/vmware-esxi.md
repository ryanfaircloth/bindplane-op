---
title: "VMware ESXi"
slug: "vmware-esxi"
hidden: false
createdAt: "2022-08-02T13:48:45.045Z"
updatedAt: "2022-08-10T15:36:44.262Z"
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter     | Type     | Default   | Description                                                                                                                                     |
| :------------ | :------- | :-------- | :---------------------------------------------------------------------------------------------------------------------------------------------- |
| listen_port\* | `int`    | 5140      | The port to bind to and receive syslog. Collector must be running as root (Linux) or Administrator (windows) when binding to a port below 1024. |
| listen_ip     | `string` | "0.0.0.0" | The IP address to bind to and receive syslog.                                                                                                   |
| enable_tls    | `bool`   | false     | Whether or not to use TLS.                                                                                                                      |
| cert_file     | `string` |           | Path to the x509 PEM certificate.                                                                                                               |
| key_file      | `string` |           | Path to the x509 PEM private key.                                                                                                               |

<span style="color:red">\*_required field_</span>