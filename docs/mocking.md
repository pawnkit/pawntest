# Mocking

Declared natives that are not part of `pawntest.inc` are registered as mocks.
Configure their return values with `MOCK_RETURN`; calls and arguments are
recorded automatically.

```pawn
native SendClientMessage(playerid, color, const message[]);

TEST(sends_message)
{
    MOCK_RETURN(SendClientMessage, 1);
    EXPECT_NATIVE_CALLS(SendClientMessage, 1);
    EXPECT_NATIVE_ARG(SendClientMessage, 0, 0, 0);
    EXPECT_NATIVE_ARG(SendClientMessage, 0, 1, -1);

    ASSERT_EQ(SendClientMessage(0, -1, "hello"), 1);
}
```

One-shot return values are consumed in order before the default return:

```pawn
MOCK_RETURN(fetch_value, 30);
MOCK_RETURN_ONCE(fetch_value, 10);
MOCK_RETURN_ONCE(fetch_value, 20);
```

Mocks can populate scalar and string output parameters or invoke a Pawn public:

```pawn
MOCK_OUTPUT(read_value, 0, 42);
MOCK_OUTPUT_STRING(read_name, 0, name_buffer_size, "Alice");
MOCK_CALLBACK(start_request, on_request_complete);
```

Argument indexes are zero-based. String output includes the destination capacity
and is truncated safely. Callback arguments are the arguments passed to the
mocked native. Global changes made by a callback persist after it returns.

Calling an unknown native without first configuring it is an error by default.
Use `--allow-unknown-natives` to permit unconfigured calls; they return `0` and
are still recorded.

Use the declarative `EXPECT_NATIVE_*` macros.
Expectation failures include every recorded scalar call, so manual inspection
is unnecessary. Output-only and callback-only behavior
also counts as explicit mock configuration.

The built-in player scenario models common open.mp natives when no explicit
mock is configured. Explicit `MOCK_*` behavior wins over the model, and custom
Go natives win over both.
