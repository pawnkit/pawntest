# Installation

Download the archive for your platform from GitHub Releases, extract it, and
place `pawntest` on `PATH`.

```sh
pawntest --version
pawntest doctor
```

Supported release targets:

| OS | amd64 | arm64 | Source compilation |
|---|---:|---:|---|
| Linux | Yes | Yes | amd64 only |
| macOS | Yes | Yes | Supported compiler required |
| Windows | Yes | Yes | amd64 only |

The official Linux Pawn compiler requires `libc6-i386` on Debian/Ubuntu or
`lib32-glibc` on Arch. Precompiled AMX tests do not require `pawncc`.

Build from source with the Go version declared in `go.mod`:

```sh
go install github.com/pawnkit/pawntest/cmd/pawntest@latest
```
