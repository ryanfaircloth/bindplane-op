---
title: "TCP"
category: 633dd7654359a20031653089
slug: "tcp"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter            | Type     | Default   | Description                                                                                                                        |
| :------------------- | :------- | :-------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| listen_port\*        | `int`    |           | Port to listen on.                                                                                                                 |
| listen_ip            | `string` | "0.0.0.0" | IP Address to listen on.                                                                                                           |
| log_type             | `string` | tcp       | Arbitrary for attribute 'log_type'. Useful for filtering between many tcp sources.                                                 |
| enable_tls           | `bool`   | false     | Whether or not to use TLS.                                                                                                         |
| tls_certificate_path | `string` |           | Path to the TLS cert to use for TLS required connections.                                                                          |
| tls_private_key_path | `string` |           | Path to the TLS key to use for TLS required connections.                                                                           |
| tls_min_version      | `enum`   | "1.2"     | The minimum TLS version to support. 1.0 and 1.1 should not be considered secure. Valid values include: `1.3`, `1.2`, `1.1`, `1.0`. |

<span style="color:red">\*_required field_</span>