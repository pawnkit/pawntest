# pawntest

A Go CLI for discovering, compiling, and running Pawn tests for SA-MP and open.mp projects. Tests run on the pure-Go [`pawnkit/goamx`](https://github.com/pawnkit/goamx) backend, so no game server is required.

## Quick start

```sh
go build -o pawntest ./cmd/pawntest

./pawntest ./tests
./pawntest --list ./tests
./pawntest doctor
```

Pass a compiler or additional include directory when needed:

```sh
./pawntest ./tests --pawncc ./tools/pawncc -i include
```

Test source files must end in `.test.pwn` or `.test.inc`. Source and precompiled AMX tests must include `<pawntest>`.

## Output formats

```sh
./pawntest ./tests --format plain
./pawntest ./tests --format json
./pawntest ./tests --format tap
./pawntest ./tests --format junit --output test-results.xml
```

## Configuration

Pawntest loads the first `pawntest.json`, `pawntest.yaml`, `pawntest.yml`, or `pawntest.toml` file in the working directory.

```toml
pawncc  = "./tools/pawncc"
include = ["include"]
tests   = ["tests/..."]
format  = "plain"
```

See [Configuration](docs/config.md) for all options.

## Documentation

- [Usage](docs/usage.md)
- [Writing tests](docs/writing-tests.md)
- [Assertions](docs/assertions.md)
- [Fixtures](docs/fixtures.md)
- [Mocking](docs/mocking.md)
- [Reports](docs/reports.md)
- [Limitations](docs/limitations.md)
