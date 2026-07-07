# Configuration

`pawntest` reads the first config file it finds in the current working
directory:

```text
pawntest.json
pawntest.yaml
pawntest.yml
pawntest.toml
```

CLI arguments take precedence over config values.

```json
{
  "pawncc": "./tools/pawncc",
  "include": ["include", "dependencies/samp-stdlib"],
  "tests": ["tests/..."],
  "format": "plain",
  "cache_dir": ".pawntest-cache",
  "allow_unknown_natives": true,
  "allow_empty": false
}
```

The same fields are available in YAML:

```yaml
pawncc: ./tools/pawncc
include:
  - include
  - dependencies/samp-stdlib
tests:
  - tests/...
format: plain
cache_dir: .pawntest-cache
allow_unknown_natives: true
allow_empty: false
```

Or TOML:

```toml
pawncc = "./tools/pawncc"
include = ["include", "dependencies/samp-stdlib"]
tests = ["tests/..."]
format = "plain"
cache_dir = ".pawntest-cache"
allow_unknown_natives = true
allow_empty = false
```

`cache_dir` is optional. When omitted, Pawntest uses the platform user cache
directory with a `pawntest` subdirectory.

## Full field reference

| Field | Type | Default | Description |
|---|---|---|---|
| `pawncc` | string | `""` | Path to pawncc. Falls back to PATH lookup. |
| `include` | string[] | `[]` | Additional Pawn include directories. |
| `tests` | string[] | `[]` | Default discovery paths (same as positional args). |
| `format` | string | `plain` | Report format: `plain`, `json`, `tap`, `junit`. |
| `cache_dir` | string | platform default | Compiled AMX and include cache location. |
| `allow_unknown_natives` | bool | `false` | Let unconfigured unknown natives return zero. |
| `allow_empty` | bool | `false` | Exit 0 when no tests are discovered. |
| `isolation` | string | `test` | Memory isolation: `test` or `suite`. |
| `run` | string | `""` | Regular expression filter on test names. |
| `tags` | string | `""` | Tag expression filter, e.g. `unit & !slow`. |
| `shuffle` | bool | `false` | Reproducibly shuffle test order. |
| `seed` | int | `1` | Shuffle/fuzz seed. |
| `repeat` | int | `1` | Run each selected test this many times. |
| `max_instructions` | int | `1000000` | Instruction budget per test invocation. |
| `jobs` | int | `1` | Concurrent files to compile and run. |
| `update_snapshots` | bool | `false` | Create or replace golden snapshots. |
| `coverage` | bool | `false` | Collect source-line coverage. |
| `coverage_output` | string | `""` | Coverage artifact path. |
| `coverage_format` | string | `lcov` | Coverage format: `lcov` or `json`. |
| `fuzz_seed` | int | `1` | Base seed for property-based tests. |
| `verbose` | bool | `false` | Include durations and source locations. |
| `quiet` | bool | `false` | Show only failures and the final summary. |
