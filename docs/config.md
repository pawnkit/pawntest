# Configuration

Pawntest reads the nearest `pawn.json` or `pawn.yaml`, then checks the working
directory for the first of these files:

```text
pawntest.json
pawntest.yaml
pawntest.yml
pawntest.toml
```

The local file overrides the project manifest. CLI arguments take precedence.
Unknown fields are rejected.

```toml
pawncc = "./tools/pawncc"
include = ["include"]
tests = ["tests/..."]
providers = ["pawntest/inventory.provider.pwn"]
format = "plain"
```

The same settings can live under `pawnkit.tool.pawntest` in the project
manifest. `pawnkit.tests.paths` supplies the default test paths, and
`pawnkit.includePaths` supplies include directories.

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
| `profile` | `false` | Count instructions by Pawn function. |
| `profile_output` | `profile.json` | Instruction profile output path. |
| `fuzz_seed` | `1` | Property-test seed. |
| `verbose` | `false` | Show durations and source paths. |
| `quiet` | `false` | Show failures and the summary only. |
| `native_plugin` | `""` | Legacy plugin library to load. |
| `plugin_architecture` | `x86` | Plugin architecture: `x86` or `x64`. |
| `plugin_worker_32` | `""` | Path to the x86 plugin worker. |
| `plugin_worker_64` | `""` | Path to the x64 plugin worker. |
