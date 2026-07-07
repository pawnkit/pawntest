#include <pawntest>

#pragma rational Float

native Float:float(value);
native Float:floatadd(Float:a, Float:b);
native Float:floatmul(Float:a, Float:b);
native floatcmp(Float:a, Float:b);
native floatround(Float:value, mode = 0);
native Float:floatsin(Float:value, mode = 0);

TEST(float_helpers)
{
    new Float:sum = floatadd(1.5, 2.5);
    ASSERT_EQ(floatcmp(sum, 4.0), 0);
    ASSERT_EQ(floatcmp(floatmul(sum, 2.0), 8.0), 0);
    ASSERT_EQ(floatround(2.6), 3);
    ASSERT_EQ(floatcmp(floatsin(90.0, 1), 1.0), 0);
    ASSERT_EQ(floatcmp(float(5), 5.0), 0);
}
