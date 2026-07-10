# Mocking

Pawntest records calls to declared server natives. Configure a return value before calling one:

```pawn
native SendClientMessage(playerid, color, const message[]);

TEST(sends_message)
{
    MOCK_RETURN(SendClientMessage, 1);

    ASSERT_EQ(SendClientMessage(0, -1, "hello"), 1);
    EXPECT_NATIVE_CALLS(SendClientMessage, 1);
    EXPECT_NATIVE_ARG(SendClientMessage, 0, 0, 0);
    EXPECT_NATIVE_STRING_ARG(SendClientMessage, 0, 2, "hello");
}
```

Other mock controls:

```pawn
MOCK_RETURN_ONCE(fetch_value, 10);
MOCK_OUTPUT(read_value, 0, 42);
MOCK_OUTPUT_STRING(read_name, 0, buffer_size, "Alice");
MOCK_CALLBACK(start_request, on_complete);
```

Useful expectations include:

```pawn
EXPECT_NATIVE_CALLED(name);
EXPECT_NATIVE_NOT_CALLED(name);
EXPECT_NATIVE_CALLS(name, count);
EXPECT_NATIVE_CALLS_BETWEEN(name, minimum, maximum);
EXPECT_NATIVE_ORDER(name);
```

Argument and call indexes start at zero. Expectations are checked after teardown.

Unconfigured native calls fail by default. Use `--allow-unknown-natives` to return `0` instead.

Explicit mocks override built-in player scenarios. Custom Go natives override both.
