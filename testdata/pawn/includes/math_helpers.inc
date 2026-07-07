#if defined _pawntest_math_helpers_included
    #endinput
#endif
#define _pawntest_math_helpers_included

stock helper_add(left, right)
{
    return left + right;
}

stock helper_clamp(value, minimum, maximum)
{
    if (value < minimum)
    {
        return minimum;
    }
    if (value > maximum)
    {
        return maximum;
    }
    return value;
}
