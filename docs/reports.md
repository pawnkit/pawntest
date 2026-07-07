# Reports

Supported formats:

```sh
pawntest tests --format=plain
pawntest tests --format=json
pawntest tests --format=tap
pawntest tests --format=junit --output test-results.xml
```

Plain output is intended for local development. JSON, TAP, and JUnit are
intended for CI systems.

Plain results are grouped by source file. Failures use multiline diagnostics,
source locations, and an exact `--run` command. `--verbose` adds durations and
absolute paths; `--quiet` prints only failures and the final summary. Statuses
and summary counts are colored independently, so a successful run remains
visibly successful when tests are skipped.

JSON output contains a summary and result list. Durations are reported in
milliseconds:

```json
{
  "summary": {
    "total": 1,
    "passed": 1,
    "failed": 0,
    "skipped": 0,
    "errored": 0
  },
  "results": [
    {
      "name": "test_addition",
      "file": "tests/math.test.pwn",
      "status": "pass",
      "duration_ms": 0
    }
  ]
}
```

JUnit output uses seconds for the testcase `time` attribute, as expected by CI
test report parsers.

Statuses:

```text
pass
fail
skip
error
xfail
xpass
```

Failures come from assertion natives. Errors come from runtime, load, or fixture
execution problems. An `xfail` is an expected known failure; an `xpass` is an
unexpected pass and fails the run.
