#if defined _calculator_included
    #endinput
#endif
#define _calculator_included

stock Add(left, right)
{
    return left + right;
}

stock Clamp(value, minimum, maximum)
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
