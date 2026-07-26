# Reports

Pawntest supports plain, JSON, TAP, and JUnit output:

```sh
pawntest tests --format plain
pawntest tests --format json
pawntest tests --format tap
pawntest tests --format junit --output test-results.xml
```

Use plain output locally and structured formats in CI. `--verbose` adds durations
and absolute source paths. `--quiet` shows failures and the summary only.

JSON results include `source` and assertion/runtime `file` locations. JUnit uses
the source as `classname` and `file`. TAP emits source and warning diagnostics.

## Runtime fidelity

Every test run records the runtime it used. Normal pawntest runs use
`platform-simulation`; runs with `--native-plugin` use
`native-plugin-integration`.

JSON stores the full record in `runtime`. TAP writes it as leading comments,
and JUnit writes it as `pawnkit.*` properties. Plain output only prints the
tier when a native plugin is active.

The record describes what ran, not how trustworthy the result is. Platform
simulation does not cover network behavior or real server scheduling. Plugin
workers limit crashes and hangs, but they are not security sandboxes.

## Statuses

| Status | Meaning |
|---|---|
| `pass` | Test passed. |
| `fail` | Assertion failed. |
| `skip` | Test was skipped. |
| `error` | Runtime, load, or fixture error. |
| `xfail` | Expected failure. |
| `xpass` | Expected failure passed; the run fails. |

JSON durations use milliseconds. JUnit durations use seconds.
