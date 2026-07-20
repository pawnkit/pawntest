# Usage

Run a directory, source file, or precompiled AMX file:

```sh
pawntest ./tests
pawntest ./tests/...
pawntest tests/math.test.pwn
pawntest tests/math_test.amx
```

Source files end in `.test.pwn` or `.test.inc` and include `<pawntest>`.

Pass a compiler, include directory, or define when the project needs one:

```sh
pawntest ./tests --pawncc ./tools/pawncc -i include -DTESTING
```

Interactive runs can install the openmultiplayer compiler when `pawncc` is missing. Its Linux binary is 32-bit; Debian and Ubuntu need `libc6-i386`, while Arch uses `lib32-glibc`.

## Common options

```text
--list                  List tests.
--run expression        Filter test names.
--tags expression       Filter tags.
--recursive             Scan directories recursively.
--fail-fast             Stop after the first failure.
--shuffle               Shuffle tests.
--seed number           Set the shuffle seed.
--repeat number         Repeat each test.
-j, --jobs number       Run files concurrently.
--isolation test|suite  Set memory isolation.
--watch                 Rerun after changes.
--coverage              Collect coverage.
--format format         Use plain, JSON, TAP, or JUnit output.
--output file           Write output to a file.
-v, --verbose           Show durations and source paths.
-q, --quiet             Show failures and the summary only.
--allow-empty           Pass when no tests are found.
--allow-unknown-natives Return zero for unconfigured natives.
--provider path         Load a Pawn native provider source or AMX file.
```

Allowed unknown-native calls appear as warnings in every report format. Run `pawntest --help` for the complete option list.

## Coverage and profiling

`--coverage` collects line coverage. `--profile` writes deterministic JSON instruction counts by Pawn function; choose its path with `--profile-output`. Both modes can run together.

## Legacy plugins

Native plugins are opt-in. Provide the plugin and a worker with the same architecture:

```sh
pawntest --native-plugin ./plugins/fixture.so \
  --plugin-architecture x64 \
  --plugin-worker-64 ./bin/pawn-plugin-host-x64
```

The worker loads the native library in a separate process. The current bridge supports one registered native with up to 32 scalar cell arguments. Strings, references, arrays, callbacks, and worker pooling are not supported yet.

Install the worker and its matching `pawn-plugin-loader` from the
[`pawn-plugin-host` releases](https://github.com/pawnkit/pawn-plugin-host/releases).
Legacy loading currently works on Linux only.

## Diagnostics and cache

```sh
pawntest doctor
pawntest cache clean
```

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Tests passed |
| `1` | A test failed or errored |
| `2` | Usage, discovery, compilation, loading, or output error |
| `3` | Internal CLI error |
