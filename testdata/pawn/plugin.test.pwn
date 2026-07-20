#include <pawntest>

native PawnKitFixture(left, right);

TEST(isolated_plugin_native)
{
    ASSERT_EQ(PawnKitFixture(20, 21), 42);
}
