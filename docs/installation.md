# Installation

There are two ways to install k9-tail: download a pre-built binary, or compile from source.

## Option 1 — Pre-built binaries (recommended)

Pre-compiled static binaries are available for a wide range of operating systems and architectures, including Linux, FreeBSD, NetBSD, OpenBSD, Solaris/illumos, and macOS.

Binaries are published at:

```
https://github.com/k9io/k9-binaries/tree/main/k9-tail
```

Each binary is distributed as a gzip-compressed file alongside a SHA256 checksum. Always verify the checksum before installing.

### Example: installing on Linux amd64

```bash
# Download the binary and its checksum
curl -LO https://github.com/k9io/k9-binaries/raw/main/k9-tail/linux/k9-tail.amd64.gz
curl -LO https://github.com/k9io/k9-binaries/raw/main/k9-tail/linux/k9-tail.amd64.gz-sha256.txt

# Verify the checksum
sha256sum -c k9-tail.amd64.gz-sha256.txt

# Decompress and install
gunzip k9-tail.amd64.gz
sudo mkdir -p /opt/k9/bin
sudo mv k9-tail.amd64 /opt/k9/bin/k9-tail
sudo chmod +x /opt/k9/bin/k9-tail
```

## Option 2 — Build from source

See [Building from Source](building.md) for full instructions.

## Post-installation steps

After placing the binary at `/opt/k9/bin/k9-tail`:

1. Ensure the Key9 configuration file exists at `/opt/k9/etc/k9.yaml`. See [Configuration](configuration.md).
2. Create the directory for the waldo file if it does not exist:

```bash
sudo mkdir -p /var/lib/k9
```

3. (Optional) Install and enable the systemd service. See [Running as a Service](systemd.md).
