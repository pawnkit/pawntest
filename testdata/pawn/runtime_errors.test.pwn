#include <pawntest>

TEST(divide_by_zero_errors)
{
    new zero = 0;
    new value = 1 / zero;
    ASSERT_TRUE(value == 0);
}
