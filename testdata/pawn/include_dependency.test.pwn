#include <pawntest>
#include "includes/math_helpers.inc"

TEST(include_dependency)
{
    ASSERT_EQ(helper_add(8, 13), 21);
    ASSERT_EQ(helper_clamp(42, 0, 10), 10);
    ASSERT_EQ(helper_clamp(-4, 0, 10), 0);
}
