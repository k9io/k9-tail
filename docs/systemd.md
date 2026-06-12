# Running as a Service

For production deployments, k9-tail should run as a systemd service so that it starts automatically at boot and restarts automatically if it crashes.

## Installing the service

A systemd unit file is included in the repository at `k9-tail.service`. Copy it to the systemd system directory:

```bash
sudo cp k9-tail.service /etc/systemd/system/k9-tail.service
```

## Enabling and starting

```bash
# Reload systemd to pick up the new unit file
sudo systemctl daemon-reload

# Enable the service to start at boot
sudo systemctl enable k9-tail

# Start the service now
sudo systemctl start k9-tail
```

## Checking status

```bash
sudo systemctl status k9-tail
```

To follow live log output:

```bash
sudo journalctl -u k9-tail -f
```

## Stopping and restarting

```bash
sudo systemctl stop k9-tail
sudo systemctl restart k9-tail
```

On a graceful stop (`systemctl stop`), k9-tail receives `SIGTERM`, flushes the current waldo position to disk, and exits cleanly. The next startup resumes from where it left off.

## Unit file reference

The included `k9-tail.service` unit file is configured as follows:

```ini
[Unit]
Description=Key9 Tail
After=network-online.target
Wants=network-online.target

[Service]
User=root
ExecStart=/opt/k9/bin/k9-tail
Restart=always
RestartSec=10
TimeoutStopSec=90
KillMode=process
OOMScoreAdjust=-900
SyslogIdentifier=k9-tail

[Install]
WantedBy=multi-user.target
```

### Key settings

| Setting | Value | Purpose |
|---|---|---|
| `After` / `Wants` | `network-online.target` | Ensures the network is available before k9-tail starts, so it can reach the Key9 API |
| `User` | `root` | Required to read `/var/log/auth.log` on most systems |
| `Restart` | `always` | Automatically restarts the process if it exits for any reason |
| `RestartSec` | `10` | Waits 10 seconds before each restart attempt |
| `TimeoutStopSec` | `90` | Allows up to 90 seconds for a graceful shutdown before forcibly killing the process |
| `OOMScoreAdjust` | `-900` | Makes k9-tail very unlikely to be killed by the kernel's OOM killer during memory pressure |
| `SyslogIdentifier` | `k9-tail` | Tags all log lines in journald with `k9-tail` for easy filtering |

## Viewing logs

All output from k9-tail is captured by journald and can be queried with `journalctl`:

```bash
# All logs since boot
sudo journalctl -u k9-tail -b

# Follow live output
sudo journalctl -u k9-tail -f

# Last 100 lines
sudo journalctl -u k9-tail -n 100
```
