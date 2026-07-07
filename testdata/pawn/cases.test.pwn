#include <pawntest>

stock assert_even(value)
{
    ASSERT_EQ(value % 2, 0);
    return 1;
}

stock assert_sum(left, right, expected)
{
    ASSERT_EQ(left + right, expected);
    return 1;
}

forward known_failure();
public known_failure()
{
    ASSERT_EQ(1, 2);
}

TEST(expected_failure)
{
    XFAIL(known_failure, "tracked defect");
}

forward known_pass();
public known_pass()
{
    return 1;
}

TEST(unexpected_pass)
{
    XFAIL(known_pass, "stale tracked defect");
}

TEST_CASE(even_two, assert_even, 2)
TEST_CASE(even_four, assert_even, 4)
TEST_CASE3(sum_case, assert_sum, 2, 3, 5)
TEST_TAG(unit_even_two)
TEST_TAG(unit_even_four)
TEST_TAG(regression_expected_failure)

TEST(declarative_syntax)
{
    ASSERT_TRUE(1);
}
