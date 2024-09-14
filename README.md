# analytics-server
Analytics server.

## Install

Install via Go toolchain.

```shell
$ go install github.com/dashanalytics/analytics-server/cmd/analytics-server
```

## Setup

### Server

| Key            | Description                                 |
|----------------|---------------------------------------------|
| `listen`       | Listen address and port.                    |
| `db`           | Redis database URL.                         |
| `access_token` | Access token for protecting sensitive data. |
| `cert`         | X.509 certificate file path.                |
| `key`          | X.509 private key file path.                |

When `cert` and `key` are not blank, the server will run as HTTPS.

### Provider-specific HTTP header specified in `header.key`

| HTTP Header     | Description                                           |
|-----------------|-------------------------------------------------------|
| `connecting_ip` | Source IP. Such as `CF-Connecting-IP` for Cloudflare. |

Example: `analytics-server.yaml`
```yaml
listen: ":443"
db: "redis://default:@localhost/0"

key: "/etc/letsencrypt/live/symboltics.com/privkey.pem"
cert: "/etc/letsencrypt/live/symboltics.com/fullchain.pem"

access_token: ""

header:
  key:
    connecting_ip: "CF-Connecting-IP"
```
