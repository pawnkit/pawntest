# pawntest

A Go CLI for discovering, compiling, and running Pawn tests for SA-MP and open.mp projects. Tests run on the pure-Go [`pawnkit/goamx`](https://github.com/pawnkit/goamx) backend, so no game server is required.

See [Installation](docs/install.md) for release binaries and supported platforms.

## Quick start

```sh
pawntest ./tests
pawntest --list ./tests
pawntest doctor
```

Pass a compiler or additional include directory when needed:

```sh
pawntest ./tests --pawncc ./tools/pawncc -i include
```

Test source files must end in `.test.pwn` or `.test.inc`. Source and precompiled AMX tests must include `<pawntest>`.

## Output formats

```sh
pawntest ./tests --format plain
pawntest ./tests --format json
pawntest ./tests --format tap
pawntest ./tests --format junit --output test-results.xml
```

## Configuration

Pawntest reads `pawn.json` or `pawn.yaml` from the project root. A local
`pawntest.json`, `pawntest.yaml`, `pawntest.yml`, or `pawntest.toml` can override
the project settings.

```toml
pawncc  = "./tools/pawncc"
include = ["include"]
tests   = ["tests/..."]
format  = "plain"
```

See [Configuration](docs/config.md) for all options.

## Documentation

- [Runnable examples](examples/README.md)
- [Usage](docs/usage.md)
- [Installation](docs/install.md)
- [Writing tests](docs/writing-tests.md)
- [Assertions](docs/assertions.md)
- [Fixtures](docs/fixtures.md)
- [Mocking](docs/mocking.md)
- [Native providers](docs/providers.md)
- [Scenario fidelity](docs/scenarios.md)
- [Compatibility](docs/compatibility.md)
- [Go API](docs/go-api.md)
- [Reports](docs/reports.md)
- [Limitations](docs/limitations.md)

## Contributing

Small test cases, runner fixes, and documentation corrections are welcome. See
[CONTRIBUTING.md](CONTRIBUTING.md).
