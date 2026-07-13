# Configuration

Pawntest loads the first config file found in the working directory:

```text
pawntest.json
pawntest.yaml
pawntest.yml
pawntest.toml
```

CLI arguments override config values.
Unknown fields are rejected.

```toml
pawncc = "./tools/pawncc"
include = ["include"]
tests = ["tests/..."]
providers = ["pawntest/inventory.provider.pwn"]
format = "plain"
```

## Fields

| Field | Default | Description |
|---|---|---|
| `pawncc` | `""` | Path to `pawncc`; defaults to `PATH`. |
| `include` | `[]` | Additional include directories. |
| `tests` | `[]` | Default test paths. |
| `providers` | `[]` | Pawn native provider source or AMX files. |
| `format` | `plain` | `plain`, `json`, `tap`, or `junit`. |
| `cache_dir` | OS cache | Cache directory. |
| `allow_unknown_natives` | `false` | Return zero for unconfigured natives. |
| `allow_empty` | `false` | Pass when no tests are found. |
| `isolation` | `test` | Use `test` or `suite` memory isolation. |
| `run` | `""` | Test name regular expression. |
| `tags` | `""` | Tag expression, such as `unit & !slow`. |
| `shuffle` | `false` | Shuffle tests. |
| `seed` | `1` | Shuffle seed. |
| `repeat` | `1` | Runs per test. |
| `max_instructions` | `1000000` | Instruction limit per invocation. |
| `jobs` | `1` | Concurrent test files. |
| `update_snapshots` | `false` | Update snapshots. |
| `coverage` | `false` | Collect coverage. |
| `coverage_output` | `""` | Coverage output path. |
| `coverage_format` | `lcov` | `lcov` or `json`. |
| `fuzz_seed` | `1` | Property-test seed. |
| `verbose` | `false` | Show durations and source paths. |
| `quiet` | `false` | Show failures and the summary only. |
