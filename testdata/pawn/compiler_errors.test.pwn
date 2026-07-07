#include <pawntest>

// pawntest: expect-error 017

TEST(expected_compiler_error)
{
    return undefined_test_symbol;
}
