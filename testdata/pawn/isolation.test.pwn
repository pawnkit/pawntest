#include <pawntest>

new isolated_state;

BEFORE_ALL()
{
    isolated_state = 10;
}

BEFORE_EACH()
{
    isolated_state++;
}

TEST(isolation_first)
{
    ASSERT_EQ(isolated_state, 11);
}

TEST(isolation_second)
{
    ASSERT_EQ(isolated_state, 11);
}
