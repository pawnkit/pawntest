# Fixtures

Fixture hooks are optional. Prefer the declarative form:

```pawn
BEFORE_ALL() {}
BEFORE_EACH() {}
AFTER_EACH() {}
AFTER_ALL() {}
```

Execution order:

```text
BEFORE_ALL
  BEFORE_EACH
  test body
  AFTER_EACH
AFTER_ALL
```

Fixture hooks are not listed as tests. A failing setup prevents the body from
running, but per-test and suite cleanup still run. Failures from the body,
named teardown, file teardown, and mock verification are aggregated in phase
order rather than replacing one another.

Named fixtures use `FIXTURE_SETUP(name)`, `FIXTURE_TEARDOWN(name)`, and
`USE_FIXTURE(name)`. Acquired named fixtures always tear down in reverse order,
including when a later setup step fails.
