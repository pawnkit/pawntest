# Assertions

`pawntest.inc` provides assertion macros backed by Go native callbacks:

```pawn
ASSERT(expression)
ASSERT_TRUE(expression)
ASSERT_FALSE(expression)
ASSERT_EQ(actual, expected)
ASSERT_NE(actual, expected)
ASSERT_STR_EQ(actual, expected)
ASSERT_STR_IEQ(actual, expected)
ASSERT_BETWEEN(actual, minimum, maximum)
ASSERT_HAS_FLAG(actual, flag)
ASSERT_ARRAY_EQ(actual, expected, length)
ASSERT_ARRAY_CONTAINS(values, length, expected)
ASSERT_FLOAT_NEAR(actual, expected, tolerance)
ASSERT_EQ_MSG(actual, expected, context)
ASSERT_NO_ERROR(callback)
FAIL(message)
SKIP(reason)
```

Assertions return from the current Pawn test public when they fail. The runner
records the message, file, and line supplied by the include macros.
`CHECK`, `CHECK_EQ`, `CHECK_NE`, and `CHECK_STR_EQ` record non-fatal failures and
allow the current public to continue.
