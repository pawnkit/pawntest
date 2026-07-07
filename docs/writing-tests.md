# Writing Tests

Name source tests `<name>.test.pwn` or `<name>.test.inc`, include
`pawntest.inc`, and declare tests with `TEST(name)`.

```pawn
#include <pawntest>

TEST(addition)
{
    ASSERT_EQ(2 + 2, 4);
}

TEST(strings)
{
    ASSERT_STR_EQ("hello", "hello");
	ASSERT_STR_CONTAINS("hello world", "world");
}

TEST(skip_example)
{
    SKIP("not implemented yet");
}
```

Lifecycle hooks use `BEFORE_EACH`, `AFTER_EACH`, `BEFORE_ALL`, and `AFTER_ALL`.

Assertions include `ASSERT`, `ASSERT_TRUE`, `ASSERT_FALSE`, `ASSERT_EQ`,
`ASSERT_NE`, `ASSERT_GT`, `ASSERT_GE`, `ASSERT_LT`, `ASSERT_LE`, string
equality/contains/prefix/suffix checks, `ASSERT_FLOAT_NEAR`, and
`ASSERT_ARRAY_EQ`. `CHECK`, `CHECK_EQ`, `CHECK_NE`, and `CHECK_STR_EQ` collect
multiple failures before the test returns. Use `EXPECT_ERROR(callback,
"message")` to verify a runtime error raised by a public callback.
Additional assertions cover ranges (`ASSERT_BETWEEN`), flags
(`ASSERT_HAS_FLAG`), array membership (`ASSERT_ARRAY_CONTAINS`),
case-insensitive strings (`ASSERT_STR_IEQ`), contextual equality
(`ASSERT_EQ_MSG`), and callbacks that must not fail (`ASSERT_NO_ERROR`).

Unknown server natives can be mocked with `MOCK_RETURN`:

```pawn
#include <pawntest>

native SendClientMessage(playerid, color, const message[]);

TEST(sends_message)
{
    MOCK_RETURN(SendClientMessage, 1);
	EXPECT_NATIVE_CALLS(SendClientMessage, 1);
	EXPECT_NATIVE_ARG(SendClientMessage, 0, 0, 0);
	EXPECT_NATIVE_STRING_ARG(SendClientMessage, 0, 2, "hello");
	EXPECT_NATIVE_ORDER(SendClientMessage);

    ASSERT_EQ(SendClientMessage(0, -1, "hello"), 1);
}
```

`EXPECT_NATIVE_CALLS_BETWEEN` supports a call-count range. Expectations are checked
after teardown, so the test does not need a final manual assertion. Return
queues, output parameters, string outputs, and callbacks are available through
`MOCK_RETURN_ONCE`, `MOCK_OUTPUT`, `MOCK_OUTPUT_STRING`, and `MOCK_CALLBACK`.
The expectation macros `EXPECT_NATIVE_CALLS`, `EXPECT_NATIVE_CALLED`,
`EXPECT_NATIVE_NOT_CALLED`, `EXPECT_NATIVE_ARG`, `EXPECT_NATIVE_STRING_ARG`,
and `EXPECT_NATIVE_ORDER` are preferred in new tests. Mismatches include the
recorded calls and arguments.

Named fixtures are reusable and opt-in. Their teardowns run automatically in
reverse use order, before the file-level `AFTER_EACH` hook:

```pawn
FIXTURE_SETUP(account)
{
    // Prepare account state.
}

FIXTURE_TEARDOWN(account)
{
    // Release account state.
}

TEST(loads_account)
{
    USE_FIXTURE(account);
}
```

Each test starts from the global-memory state left by `BEFORE_ALL` by
default. Its own setup, body, mocks, and virtual clock are isolated from every
other test. Select `--isolation=suite` only when a suite intentionally shares
mutable global state.

Use the deterministic virtual clock for timer-driven code:

```pawn
forward expire_session();
public expire_session() {}

TEST(session_expiry)
{
    TEST_SCHEDULE(1000, expire_session);
    TEST_ADVANCE_TIME(999);
    ASSERT_EQ(TEST_NOW(), 999);
    TEST_RUN_PENDING();
}
```

Callbacks scheduled for the same time execute in registration order. The
clock begins at zero for every test and never sleeps in real time.

Compiler diagnostics can be tested without producing a runnable AMX. The file
still uses the normal name and include contract:

```pawn
#include <pawntest>

// pawntest: expect-error 017
TEST(rejects_missing_symbol)
{
    return missing_symbol;
}
```

Use `expect-error` or `expect-warning` followed by the three-digit Pawn
diagnostic code. Missing expectations and unexpected compiler errors fail the
test; unrelated compiler warnings remain visible but do not fail it.

Generate independently reported parameter cases from a callback:

```pawn
stock assert_even(value)
{
    ASSERT_EQ(value % 2, 0);
    return 1;
}

TEST_CASE(even_two, assert_even, 2)
TEST_CASE(even_four, assert_even, 4)
```

`TEST_CASE2` and `TEST_CASE3` pass two or three values to the callback.

Each case becomes a normal public named `test_<case>`, so it has normal
filtering, fixtures, isolation, and timing. Keep the generated public name
within Pawn's 31-character symbol limit.

Known defects can have a distinct expected-failure result:

```pawn
forward known_failure();
public known_failure()
{
    ASSERT_EQ(1, 2);
}

TEST(known_defect)
{
    XFAIL(known_failure, "issue 42");
}
```

A failing callback reports `xfail` and does not fail the run. A callback that
passes reports `xpass` and fails the run, prompting removal of the stale
expectation.

Attach tags using a single metadata token. The first underscore separates the
tag from the test case name:

```pawn
TEST_CASE(even_two, assert_even, 2)
TEST_TAG(unit_even_two)
TEST_TAG(fast_even_two)
```

Run it with `pawntest --tags 'unit & !slow'`. Expressions support `!`, `&`,
`|`, and parentheses. Tags are compiled into the AMX and work without source.

Golden string snapshots are explicit and reviewable:

```pawn
TEST(serialized_player)
{
    ASSERT_SNAPSHOT("player", "{id: 7, name: Alice}");
}
```

Snapshots live in `__snapshots__/<test-file>.snap.json`, keyed by test and
snapshot name. Normal runs never rewrite them; use `--update-snapshots` to
accept values, then review and commit the JSON file.

Deterministic integer properties invoke a callback repeatedly:

```pawn
forward value_is_bounded(value);
public value_is_bounded(value)
{
    ASSERT_GE(value, -100);
    ASSERT_LE(value, 100);
    return 1;
}

TEST(integer_bounds)
{
    FUZZ_INT(value_is_bounded, 200, -100, 100);
}
```

Each test derives a stable seed from `--fuzz-seed` and its public name. A
failure is shrunk toward zero or the nearest range endpoint and reports both
the seed and minimized value. Property callbacks should avoid persistent side
effects because generated and shrinking calls execute within the same test.

## Player Scenarios

The built-in player model lets production code call common open.mp natives
without configuring each one as a generic mock:

```pawn
TEST(welcomes_player)
{
    new playerid = TEST_CREATE_PLAYER("Alice");

    SetPlayerMoney(playerid, 1000);
    SetPlayerPos(playerid, 100.0, 200.0, 10.0);
    SendClientMessage(playerid, -1, "Welcome, Alice");

    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_PLAYER_MONEY(playerid, 1000);
    ASSERT_PLAYER_POS_NEAR(playerid, 100.0, 200.0, 10.0, 0.01);
    ASSERT_PLAYER_MESSAGE_CONTAINS(playerid, "Welcome");
}
```

The initial module models `IsPlayerConnected`, player names, position, money,
messages, broadcast messages, and `Kick`. Explicit `MOCK_RETURN` configuration
or custom Go natives take precedence over the model.
