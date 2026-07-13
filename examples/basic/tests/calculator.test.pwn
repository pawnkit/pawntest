#include <pawntest>
#include <calculator>

TEST(adds_numbers)
{
    ASSERT_EQ(Add(20, 22), 42);
    ASSERT_EQ(Add(-5, 2), -3);
}

TEST(clamps_values)
{
    ASSERT_EQ(Clamp(-1, 0, 100), 0);
    ASSERT_EQ(Clamp(50, 0, 100), 50);
    ASSERT_EQ(Clamp(101, 0, 100), 100);
}
