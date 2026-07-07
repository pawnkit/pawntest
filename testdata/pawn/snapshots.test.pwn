#include <pawntest>

TEST(string_snapshot)
{
    ASSERT_SNAPSHOT("greeting", "hello from Pawn");
}
