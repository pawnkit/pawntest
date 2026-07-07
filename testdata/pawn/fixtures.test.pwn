#include <pawntest>

new fixture_state;
new named_fixture_state;

FIXTURE_SETUP(account)
{
    named_fixture_state = 42;
}

FIXTURE_TEARDOWN(account)
{
    named_fixture_state = 0;
}

BEFORE_ALL()
{
    fixture_state = 10;
}

AFTER_ALL()
{
    ASSERT_EQ(fixture_state, 12);
}

BEFORE_EACH()
{
    ASSERT_EQ(fixture_state, 10);
    fixture_state = 11;
}

AFTER_EACH()
{
    ASSERT_EQ(fixture_state, 12);
}

TEST(fixtures_run_in_order)
{
    ASSERT_EQ(fixture_state, 11);
    fixture_state = 12;
}

TEST(named_fixture_is_available)
{
    USE_FIXTURE(account);
    ASSERT_EQ(named_fixture_state, 42);
    fixture_state = 12;
}
