#include <pawntest>

forward value_stays_in_range(value);
public value_stays_in_range(value)
{
    ASSERT_GE(value, -100);
    ASSERT_LE(value, 100);
    return 1;
}

TEST(integer_property)
{
    FUZZ_INT(value_stays_in_range, 100, -100, 100);
}
