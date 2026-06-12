# Usage

## Running k9-tail

Once installed and configured, k9-tail can be run directly from the command line:

```bash
/opt/k9/bin/k9-tail
```

The process runs in the foreground. For production deployments, it is recommended to run k9-tail as a systemd service so it starts automatically and restarts on failure. See [Running as a Service](systemd.md).

## Command-line flags

k9-tail has a single optional flag:

| Flag | Default | Description |
|---|---|---|
| `-debug` | `false` | Enable verbose debug logging |

### Debug mode

```bash
/opt/k9/bin/k9-tail -debug
```

Debug mode logs additional information to standard output, including:

* The `tail_file` path being followed
* The current waldo (seek) position on startup
* Each `sshd` log line as it is read, along with its byte offset and line number

This is useful for troubleshooting configuration problems or verifying that log lines are being picked up and transmitted correctly.

## What k9-tail processes

k9-tail reads each new line appended to `tail_file` and applies the following filter before transmitting:

* The line **must** contain the string `sshd`
* The line **must not** contain ` audit[` or ` audit:`

Lines that do not match this filter are silently skipped.

Matching lines are sent to Key9 as a JSON POST request:

```json
{
  "log": "<original log line>",
  "host": "<system hostname>"
}
```

## Retry behavior

If a log line cannot be delivered (non-200 response or connection error), k9-tail retries with **exponential backoff**:

* Initial retry delay: 2 seconds
* Each subsequent retry doubles the delay, capped at 60 seconds
* Maximum retries: 10

After 10 failed attempts, the log line is discarded and a warning is written to the log. k9-tail continues processing subsequent lines.

## Signal handling

k9-tail handles the following signals:

| Signal | Behavior |
|---|---|
| `SIGTERM` | Flush waldo position to disk and exit cleanly |
| `SIGINT` (Ctrl-C) | Flush waldo position to disk and exit cleanly |
| `SIGABRT` | Flush waldo position to disk and exit cleanly |

On a graceful shutdown, any pending waldo position is written to disk before the process exits. This ensures that no already-processed log lines are re-transmitted on the next startup.

## Verifying operation

To confirm k9-tail is running and transmitting logs, you can:

1. Run with `-debug` and watch for lines being printed as SSH events occur.
2. Inspect the waldo file — its contents will update every 5 seconds while the tail file has new activity:

```bash
cat /var/lib/k9/waldo
```

The waldo file contains a byte offset, for example `{2227597 0}`. A changing value confirms that k9-tail is actively reading the log.

3. Trigger a test SSH login to the system and verify it appears in the Key9 dashboard.
