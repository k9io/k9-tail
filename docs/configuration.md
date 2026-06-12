# Configuration

k9-tail reads its configuration from a YAML file located at:

```
/opt/k9/etc/k9.yaml
```

This path is fixed and cannot be changed via command-line flags. The same file is shared with other Key9 components. A template is available in the [k9-ssh repository](https://github.com/k9io/k9-ssh/blob/main/etc/k9.yaml).

## Configuration file structure

```yaml
authentication:
  api_key: "your-api-key"
  company_uuid: "your-company-uuid"

tail:
  tail_file: "/var/log/auth.log"
  waldo_file: "/var/lib/k9/waldo"
  client_logging_url: "https://api.k9.io/client/logging"
```

## Field reference

### `authentication` section

| Field | Type | Required | Description |
|---|---|---|---|
| `api_key` | string | Yes | API key issued by Key9 for your account |
| `company_uuid` | string | Yes | UUID identifying your organization in Key9 |

These two values are combined to form the `API_KEY` HTTP header sent with every log submission:

```
API_KEY: <company_uuid>:<api_key>
```

### `tail` section

| Field | Type | Required | Description |
|---|---|---|---|
| `tail_file` | string | Yes | Path to the authentication log file to monitor |
| `waldo_file` | string | Yes | Path where k9-tail stores its file position state |
| `client_logging_url` | string | Yes | HTTPS endpoint for the Key9 client logging API |

#### `tail_file`

The authentication log to follow. On most Debian/Ubuntu systems this is `/var/log/auth.log`. On RHEL/CentOS systems it may be `/var/log/secure`.

#### `waldo_file`

k9-tail writes a small state file here containing the current byte offset of `tail_file`. This prevents duplicate log submissions if the process is restarted.

* The file is created automatically if it does not exist.
* It is updated every **5 seconds** during normal operation and immediately upon a graceful shutdown.
* File permissions are set to `0600` (owner read/write only).
* If `tail_file` is truncated (e.g., log rotation), the waldo position is reset to zero automatically.

#### `client_logging_url`

The Key9 API endpoint that receives log events. This value **must** use `https://` — k9-tail will refuse to start if an HTTP URL is provided.

## Validation

On startup, k9-tail validates that all five fields are present and that `client_logging_url` uses HTTPS. If any check fails, the process exits immediately with a descriptive error message.
