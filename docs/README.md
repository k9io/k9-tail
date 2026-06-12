# k9-tail

**k9-tail** is a lightweight daemon that monitors authentication log files and streams SSH events to the [Key9 Identity](https://k9.io) service in real time.

## What it does

k9-tail follows a system authentication log (typically `/var/log/auth.log`) and forwards SSH-related entries to Key9. This enables Key9 to:

* Identify which SSH public keys are being used for authentication
* Record successful and failed login attempts
* Establish geolocation data from client IP addresses
* Provide a full audit trail of SSH access across your fleet

Only `sshd` log entries are transmitted. Kernel audit log lines (containing `audit[` or `audit:`) are explicitly excluded.

## How it tracks progress

k9-tail uses a **waldo file** to record its current position in the authentication log. On restart, it resumes from where it left off, ensuring no events are duplicated or missed. If the log file is truncated (e.g., after log rotation), the waldo position is automatically reset to the beginning.

## Is k9-tail required?

k9-tail is optional but strongly recommended. Key9 SSH will function without it, but you will lose visibility into key usage, login history, and geolocation data.

## Quick links

* [Installation](installation.md)
* [Configuration](configuration.md)
* [Usage](usage.md)
* [Building from Source](building.md)
* [Running as a Service](systemd.md)
* [Key9 Slack Community](https://key9identity.slack.com/)
