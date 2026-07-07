#include <pawntest>

#pragma rational Float

native send_message(playerid, const message[]);

new runtime_zero;

forward trigger_divide_by_zero();
public trigger_divide_by_zero()
{
    return 10 / runtime_zero;
}

TEST(rich_assertions)
{
    new actual[] = {1, 2, 3};
    new expected[] = {1, 2, 3};

    ASSERT_GT(5, 4);
    ASSERT_GE(5, 5);
    ASSERT_LT(4, 5);
    ASSERT_LE(5, 5);
    ASSERT_STR_CONTAINS("hello world", "world");
    ASSERT_STR_PREFIX("hello world", "hello");
    ASSERT_STR_SUFFIX("hello world", "world");
    ASSERT_FLOAT_NEAR(0.3333, 0.3334, 0.001);
    ASSERT_ARRAY_EQ(actual, expected, sizeof actual);
    ASSERT_ARRAY_CONTAINS(actual, sizeof actual, 2);
    ASSERT_BETWEEN(5, 1, 10);
    ASSERT_HAS_FLAG(5, 4);
    ASSERT_STR_IEQ("Alice", "alice");
    ASSERT_EQ_MSG(2 + 2, 4, "math remains stable");
    EXPECT_ERROR(trigger_divide_by_zero, "divide by zero");
}

TEST(mock_expectations)
{
    MOCK_RETURN(send_message, 1);
    EXPECT_NATIVE_CALLS(send_message, 1);
    EXPECT_NATIVE_ARG(send_message, 0, 0, 7);
    EXPECT_NATIVE_STRING_ARG(send_message, 0, 1, "hello");
    EXPECT_NATIVE_ORDER(send_message);

    new result = send_message(7, "hello");
    ASSERT_EQ(result, 1);
}
