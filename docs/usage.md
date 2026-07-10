# Usage

Run a directory, source file, or precompiled AMX file:

```sh
pawntest ./tests
pawntest ./tests/...
pawntest tests/math.test.pwn
pawntest tests/math_test.amx
```

Source files must end in `.test.pwn` or `.test.inc`. Every test file must include `<pawntest>`.

Use a specific compiler or include directory when needed:

```sh
pawntest ./tests --pawncc ./tools/pawncc -i include -DTESTING
```

If `pawncc` is missing, interactive runs can install the openmultiplayer compiler. On Linux, its 32-bit binary requires `libc6-i386` on Debian/Ubuntu or `lib32-glibc` on Arch.

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
```

Run `pawntest --help` for all options.

## Diagnostics

Check the compiler, cache, and a sample test:

```sh
pawntest doctor
```

## Exit codes

```text
0  Tests passed.
1  A test failed or errored.
2  Usage, discovery, compile, load, or output error.
3  Internal CLI error.
```
