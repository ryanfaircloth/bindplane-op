---
title: "Syslog"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "syslog"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    |         | ✓    |        |
| Windows  |         | ✓    |        |
| macOS    |         | ✓    |        |

## Configuration Table

| Parameter                       | Type     | Default   | Description                                                                                                                                                                                                    |
| :------------------------------ | :------- | :-------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| protocol\*                      | `enum`   | "rfc3164" | The RFC protocol to use when parsing incoming syslog. Valid values are `rfc3164` or `rfc5424`.                                                                                                                 |
| connection_type                 | `enum`   | udp       | The transport protocol to use. Valid values are `udp` or `tcp`.                                                                                                                                                |
| data_flow                       | `enum`   | high      | Enable high flow or reduced low flow.                                                                                                                                                                          |
| listen_port\*                   | `int`    | 5140      | The port to bind to and receive syslog. Collector must be running as root (Linux) or Administrator (windows) when binding to a port below 1024.                                                                |
| listen_ip\*                     | `string` | "0.0.0.0" | The IP address to bind to and receive syslog.                                                                                                                                                                  |
| timezone                        | `enum`   | UTC       | RFC3164 only. The timezone to use when parsing timestamps.                                                                                                                                                     |
| enable_octet_counting           | `bool`   | false     | Whether or not to parse using a trailer character. This is a special character that will be the termination character for syslog messages. This is only applicable for `tcp` and `rfc5424` configurations.     |
| non_transparent_framing_trailer | `enum`   | LF        | Whether or not to enable octet counting on syslog framing. This framing allows for the transmission of all characters inside a syslog message. This is only applicable for `tcp` and `rfc5424` configurations. |
| enable_mutual_tls                      | `bool`   | false     | Whether or not to use mutual TLS.                                                                                                                                                                                     |
| cert_file                       | `string` |           | Path to the TLS cert to use for TLS required connections.                                                                                                                                                      |
| key_file                        | `string` |           | Path to the TLS key to use for TLS required connections.                                                                                                                                                       |
| ca_file                         | `string` |           | When set, enforces mutual TLS authentication and verifies client certificates.                                                                                                                                 |
| tls_min_version                 | `enum`   | "1.2"     | The minimum TLS version to support. 1.0 and 1.1 should not be considered secure.                                                                                                                               |

<span style="color:red">\*_required field_</span>
