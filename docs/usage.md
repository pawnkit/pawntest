# Usage

`pawntest` discovers Pawn test sources and AMX files, then lists or runs public
functions whose names begin with `test_`.

```sh
pawntest ./tests
pawntest ./tests/...
pawntest tests/math.test.pwn
pawntest tests/helpers.test.inc
pawntest tests/math_test.amx
```

Pawn source tests must be named `<name>.test.pwn` or `<name>.test.inc`. Source
tests and precompiled AMX tests must include `<pawntest>`; the include emits a
reserved marker public that `pawntest` checks before listing or running tests.

Pawn source files are compiled with `pawncc` before loading:

```sh
pawntest ./tests --pawncc ./tools/pawncc -i include -DTESTING
```

`pawntest.inc` is embedded in the binary. Before compiling, `pawntest` extracts
it to `<cache-dir>/include/pawntest.inc` if the cached copy is missing or stale,
then automatically passes that include directory to `pawncc`.

By default, `<cache-dir>` is the OS user cache directory plus `pawntest`: XDG
cache on Linux, `Library/Caches` on macOS, and LocalAppData on Windows. Use
`--cache-dir` to override it.

If `pawncc` is not found in `PATH` and the run is interactive, `pawntest` asks
before downloading the openmultiplayer compiler from GitHub releases and stores
it under the Pawntest cache directory:

```text
<cache-dir>/tools/openmp-compiler/
```

When GitHub provides a release asset digest, Pawntest verifies the downloaded
compiler archive before installing it.

On Linux, the official openmultiplayer compiler release is a 32-bit glibc
binary. Pawntest wraps the downloaded compiler so the bundled `libpawnc.so` is
on the library path, but the system runtime loader must still be visible where
`pawntest` is running, including `/lib/ld-linux.so.2`. On Arch/CachyOS install
`lib32-glibc`; on Debian/Ubuntu install `libc6-i386`.

## Sub-commands

### `pawntest doctor`

Print environment diagnostics and run a sample compile/run:

```sh
pawntest doctor
pawntest doctor --pawncc ./tools/pawncc --cache-dir .cache
```

Output includes the active platform, config file, cache location, embedded
include hash, resolved pawncc path and version, and the result of a tiny
sample test compile/run.

## Flags

```sh
--list                  list tests without running them
--run expression        only include test publics matching a Go regular expression
--tags expression       filter tags, for example 'unit & !slow'
--shuffle               reproducibly shuffle tests (default seed: 1)
--seed number           choose the shuffle seed
--repeat number         run every selected test repeatedly
--max-instructions n    cap instructions per setup, test, or teardown call
-j, --jobs number       compile and run test files concurrently
--isolation test|suite  isolate global memory per test or share it per suite
--update-snapshots      create or replace golden string snapshots
--coverage              collect source-line coverage
--coverage-format       write lcov or json coverage
--coverage-output file  choose the coverage artifact path
--watch                 rerun after source, include, or config changes
--watch-interval time   choose the polling interval
--fuzz-seed number      base seed for deterministic property tests
--recursive             recursively scan input directories
--format plain|json|tap|junit
--color auto|always|never
--output file           write the report to a file
-v, --verbose           include durations and absolute source locations
-q, --quiet             show failures and the final summary only
--cache-dir dir         store generated includes and AMX files in dir
--no-cache              compile to a source-named AMX instead of cache-keyed AMX
--count=1               force recompilation
--allow-empty           treat no discovered tests as success
--allow-unknown-natives allow unconfigured unknown natives to return zero
-V, --version           print version and exit
```

`--fail-fast` runs files serially, stops after the first failed file, and stops
within that file after the first failed test, even when `--jobs` is greater than
one.
Parallel runs preserve discovery order in reports. Invalid regular expressions,
non-positive job counts, and instruction-budget overruns are reported as
errors. Runtime errors use compiler debug metadata to report the Pawn source
file, line, and public name when that metadata is present.

The Go library runner also accepts a `Natives` map. This lets an embedding
application register realistic domain natives with Pawn memory and callback
access instead of configuring every interaction as a generic mock.

With `--repeat`, reports label every result as `[attempt N/total]`; `--list`
continues to list each discovered test only once. Expected failures appear as
`xfail`. Unexpected passes appear as `xpass` and fail the command.

Coverage requires compiler debug metadata, which Pawntest enables by default.
LCOV includes executable zero-count lines and excludes Pawntest's own metadata;
JSON maps source files to line/count pairs. Watch mode uses content hashes,
keeps running after failures, and monitors Pawn source, includes, and config.

## Exit codes

```text
0 = all tests passed, skipped tests allowed
1 = one or more tests failed or errored
2 = usage, discovery, compile, load, or output error
3 = internal CLI setup error
```
