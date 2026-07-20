# Contributing

PawnKit is maintained by volunteers, so reviews may take a little time.

Small test-runner fixes, Pawn examples, and documentation corrections are
welcome. A minimal `.test.pwn` file is often enough to explain a problem.

Use the Go version from `go.mod` and run:

```sh
task check
```

Compiler integration changes also require:

```sh
PAWNTEST_PAWNCC=/path/to/pawncc \
go test ./internal/cli -tags integration -count=1
```

Add focused Go tests and Pawn corpus fixtures for behavior changes. Public CLI,
configuration, report, include, provider, and Go API changes must update the
compatibility documentation and changelog.
