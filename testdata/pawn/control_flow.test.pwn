#include <pawntest>

stock classify(value)
{
    switch (value)
    {
        case -2 .. -1:
        {
            return -1;
        }
        case 0:
        {
            return 0;
        }
        case 1 .. 3:
        {
            return 1;
        }
    }
    return 2;
}

TEST(control_flow_and_arrays)
{
    new values[5] = { -1, 0, 1, 2, 4 };
    new total = 0;

    for (new i = 0; i < sizeof(values); i++)
    {
        total += classify(values[i]);
    }

    ASSERT_EQ(total, 3);
}
