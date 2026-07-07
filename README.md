# pawntest

`pawntest` is a Go CLI for discovering, compiling, listing, and running Pawn tests for SA-MP/open.mp-style projects.

The first-class runtime target is the pure-Go
[`pawnkit/goamx`](https://github.com/pawnkit/goamx) AMX backend. It parses AMX
metadata, executes 32-bit Pawn compiler output, dispatches built-in and mocked
natives, and reports assertion results without requiring a game server.

## Quick Start

```sh
# Build
go build -o pawntest ./cmd/pawntest

# Run tests
pawntest ./tests
pawntest ./tests --pawncc ./tools/pawncc -i include

# List tests without running them
pawntest --list ./tests

# Print environment diagnostics
pawntest doctor

# Print version
pawntest --version
```

Source tests must be named `<name>.test.pwn` or `<name>.test.inc`, and every
source or precompiled AMX test must include `<pawntest>`.

The `pawntest.inc` include is embedded in the binary. It is extracted to the
cache when missing or stale, and that cache include directory is passed to
`pawncc` automatically.

The default cache follows the platform user cache location: XDG cache on Linux,
`Library/Caches` on macOS, and LocalAppData on Windows.

If `pawncc` is not on `PATH`, interactive runs ask before downloading the
openmultiplayer compiler from GitHub releases into Pawntest's cache directory.
Release asset digests are verified when GitHub provides them.

Check the local toolchain and run a sample compile/run:

```sh
pawntest doctor
```

Run tests:

```sh
pawntest tests --format=plain
pawntest tests --format=json
pawntest tests --format=tap
pawntest tests --format=junit --output test-results.xml
```

Declared server natives can be configured with `MOCK_RETURN` without any flag.
To also permit calls that have no configured mock return, use:

```sh
pawntest tests --allow-unknown-natives
```

## Configuration

`pawntest` reads the first config file it finds in the working directory:
`pawntest.json`, `pawntest.yaml`, `pawntest.yml`, or `pawntest.toml`.

```toml
pawncc  = "./tools/pawncc"
include = ["include"]
tests   = ["tests/..."]
format  = "plain"
```

See [docs/config.md](docs/config.md) for the full field reference.

## Docs

- [Usage](docs/usage.md)
- [Writing tests](docs/writing-tests.md)
- [Assertions](docs/assertions.md)
- [Fixtures](docs/fixtures.md)
- [Mocking](docs/mocking.md)
- [Reports](docs/reports.md)
- [Configuration](docs/config.md)
- [Limitations](docs/limitations.md)
