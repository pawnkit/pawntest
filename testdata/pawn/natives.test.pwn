#include <pawntest>

native external_native();

TEST(unknown_native_mock)
{
    MOCK_RETURN(external_native, 7);
    EXPECT_NATIVE_CALLS(external_native, 1);
    new result = external_native();
    ASSERT_EQ(result, 7);
}
