# Assertions

Pawntest provides fatal and non-fatal assertions.

```pawn
ASSERT(expression)
ASSERT_TRUE(expression)
ASSERT_FALSE(expression)
ASSERT_EQ(actual, expected)
ASSERT_NE(actual, expected)
ASSERT_GT(actual, expected)
ASSERT_GE(actual, expected)
ASSERT_LT(actual, expected)
ASSERT_LE(actual, expected)
ASSERT_STR_EQ(actual, expected)
ASSERT_STR_IEQ(actual, expected)
ASSERT_STR_CONTAINS(actual, expected)
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

Fatal assertions stop the current test. `CHECK`, `CHECK_EQ`, `CHECK_NE`, and `CHECK_STR_EQ` record a failure and continue.

Use `EXPECT_ERROR(callback, message)` to test runtime errors.
