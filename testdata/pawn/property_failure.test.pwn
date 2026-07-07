#include <pawntest>

forward value_must_be_ten(value);
public value_must_be_ten(value)
{
    ASSERT_EQ(value, 10);
    return 1;
}

TEST(property_shrinks_failure)
{
    FUZZ_INT(value_must_be_ten, 100, 10, 100);
}
