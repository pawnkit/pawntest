# Roleplay gamemode example

This project shows how pawntest fits an existing gamemode. Production code is
in `include/roleplay.inc`; tests share `tests/support.inc` and do not copy the
gamemode callbacks.

```text
include/roleplay.inc                 gamemode logic and callbacks
include/inventory.inc                plugin native declarations
provider/inventory.provider.pwn      test implementation of plugin natives
tests/support.inc                    shared fixtures and setup
tests/account.test.pwn               lifecycle, dialogs, mocks, provider state
tests/world.test.pwn                 entities, checkpoints, vehicles, server state
tests/integrations.test.pwn          HTTP, SQLite, strict scenario checks
tests/quality.test.pwn               cases, tags, property tests, time, snapshots
tests/compile_contract.test.pwn      expected compiler diagnostics
```

## Feature map

| Capability | Where to look |
|---|---|
| Config, discovery, compiler includes, providers | `pawntest.toml` |
| Production callbacks and reusable logic | `include/roleplay.inc` |
| Setup hooks and named fixtures | `tests/support.inc` |
| Assertions and non-fatal checks | `tests/quality.test.pwn` |
| Player lifecycle, commands, dialogs | `tests/account.test.pwn` |
| Entities, server state, checkpoints, vehicles | `tests/world.test.pwn` |
| Native mocks and argument expectations | `tests/account.test.pwn` |
| Pawn native providers and isolation | `provider/` and `tests/account.test.pwn` |
| HTTP, SQLite, strict resource checks | `tests/integrations.test.pwn` |
| Cases, tags, property tests, snapshots, virtual time | `tests/quality.test.pwn` |
| Expected errors, skips, tracked failures | `tests/quality.test.pwn` |
| Expected compiler diagnostics | `tests/compile_contract.test.pwn` |

Run everything:

```sh
cd examples/roleplay
pawntest
```

Common development commands:

```sh
pawntest --list
pawntest --run 'login|delivery'
pawntest --tags unit
pawntest --shuffle --seed 42
pawntest --repeat 10
pawntest --coverage --coverage-output coverage.lcov
pawntest --format json
pawntest --format tap
pawntest --format junit --output test-results.xml
pawntest --update-snapshots
pawntest --watch
```

The default `test` isolation restores gamemode and provider memory before every
test. Use `isolation = "suite"` only when tests intentionally share state.

Pawntest models deterministic server state and callback transitions. It does
not start an open.mp server or simulate networking, physics, persistence, or
plugin internals. Use mocks for interactions and Pawn providers when tests need
a reusable native implementation.
