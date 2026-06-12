# Building from Source

k9-tail is written in Go and compiles to a single static binary with no runtime dependencies.

## Prerequisites

* [Go](https://golang.org/dl/) 1.21 or later (the module requires Go 1.25.10+)
* Git

## Building for the current platform

```bash
# Clone the repository
git clone https://github.com/k9io/k9-tail.git
cd k9-tail

# Fetch dependencies
go mod tidy

# Build
./build.sh
```

The `build.sh` script produces a statically linked binary with debug symbols stripped:

```bash
env CGO_ENABLED=0 go build -ldflags "-s"
```

The resulting binary is named `k9-tail` in the project root.

### Install after building

```bash
sudo mkdir -p /opt/k9/bin
sudo cp k9-tail /opt/k9/bin/k9-tail
```

## Cross-compiling for all platforms

The `scripts/build-all` script cross-compiles k9-tail for every OS/architecture combination supported by Go, except plan9, iOS, Android, and JavaScript/WASM.

```bash
./scripts/build-all
```

Binaries are written to `bin/<os>/k9-tail.<arch>.gz` alongside a SHA256 checksum file at `bin/<os>/k9-tail.<arch>.gz-sha256.txt`.

### Supported platforms (examples)

| OS | Architectures |
|---|---|
| Linux | amd64, arm64, arm, i386, ppc64, ppc64le, mips, mips64, s390x, and more |
| FreeBSD | amd64, arm64, arm, i386 |
| NetBSD | amd64, arm64, arm, i386 |
| OpenBSD | amd64, arm64, arm, i386 |
| Solaris / illumos | amd64 |
| macOS (darwin) | amd64, arm64 |

Run `go tool dist list` to see the full list of targets recognized by your Go installation.

## Build flags explained

| Flag | Purpose |
|---|---|
| `CGO_ENABLED=0` | Disables C bindings, producing a fully static binary that runs without glibc or any shared library |
| `-ldflags "-s"` | Strips the symbol table from the binary, reducing file size |

## Dependencies

| Module | Purpose |
|---|---|
| `github.com/nxadm/tail` | File tailing with log-rotation support |
| `github.com/clarketm/json` | JSON marshaling |
| `gopkg.in/yaml.v2` | YAML configuration parsing |
