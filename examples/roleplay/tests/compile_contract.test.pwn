#include <pawntest>

// pawntest: expect-error 017

TEST(rejects_unknown_account_field)
{
    return missing_account_field;
}
