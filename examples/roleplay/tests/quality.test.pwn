#include <support>

stock CheckFine(speed, limit, expected)
{
    ASSERT_EQ(CalculateSpeedingFine(speed, limit), expected);
    return 1;
}

forward FineNeverNegative(speed);
public FineNeverNegative(speed)
{
    ASSERT_GE(CalculateSpeedingFine(speed, 100), 0);
    return 1;
}

new runtimeZero;

forward TriggerRuntimeError();
public TriggerRuntimeError()
{
    return 10 / runtimeZero;
}

forward KnownFineBug();
public KnownFineBug()
{
    ASSERT_EQ(CalculateSpeedingFine(110, 100), 100);
    return 1;
}

TEST_CASE3(fine_below_limit, CheckFine, 90, 100, 0)
TEST_CASE3(fine_above_limit, CheckFine, 110, 100, 200)

TEST(fine_property)
{
    FUZZ_INT(FineNeverNegative, 100, 0, 200);
}

TEST(payday_timer)
{
    ASSERT_EQ(TEST_NOW(), 0);
    TEST_SCHEDULE(60000, RunPayday);
    TEST_ADVANCE_TIME(59999);
    ASSERT_EQ(gPaydayRuns, 0);
    TEST_ADVANCE_TIME(1);
    ASSERT_EQ(gPaydayRuns, 1);
}

TEST(mode_snapshot)
{
    ASSERT_SNAPSHOT("schema", "roleplay:account+inventory+delivery");
}

TEST(assertion_examples)
{
    new actual[] = {1, 2, 3};
    new expected[] = {1, 2, 3};

    ASSERT_ARRAY_EQ(actual, expected, sizeof actual);
    ASSERT_ARRAY_CONTAINS(actual, sizeof actual, 2);
    ASSERT_STR_CONTAINS("Roleplay account", "account");
    ASSERT_FLOAT_NEAR(1.0, 1.001, 0.01);
    CHECK_EQ(CalculateSpeedingFine(100, 100), 0);
}

TEST(expected_runtime_error)
{
    EXPECT_ERROR(TriggerRuntimeError, "divide by zero");
}

TEST(tracked_fine_bug)
{
    XFAIL(KnownFineBug, "replace the legacy fine table");
}

TEST(live_server_only)
{
    SKIP("requires a live networking test environment");
}

TEST_TAG(unit_fine_below_limit)
TEST_TAG(unit_fine_above_limit)
TEST_TAG(property_fine_property)
TEST_TAG(unit_payday_timer)
TEST_TAG(snapshot_mode_snapshot)
TEST_TAG(unit_assertion_examples)
TEST_TAG(unit_expected_runtime_error)
TEST_TAG(regression_tracked_fine_bug)
TEST_TAG(live_live_server_only)
