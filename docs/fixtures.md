# Fixtures

Use hooks for shared setup and cleanup:

```pawn
BEFORE_ALL() {}
BEFORE_EACH() {}
AFTER_EACH() {}
AFTER_ALL() {}
```

They run in this order:

```text
BEFORE_ALL
  BEFORE_EACH
  test
  AFTER_EACH
AFTER_ALL
```

Cleanup hooks still run after a failure.

Use named fixtures for reusable, opt-in setup:

```pawn
FIXTURE_SETUP(database) {}
FIXTURE_TEARDOWN(database) {}

TEST(example)
{
    USE_FIXTURE(database);
}
```

Named fixtures are cleaned up in reverse order.
