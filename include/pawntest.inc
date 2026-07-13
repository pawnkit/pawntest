#if defined _pawntest_included
    #endinput
#endif
#define _pawntest_included

#include <pawntest/core>
#include <pawntest/assertions>
#include <pawntest/mocks>
#include <pawntest/scenarios>

forward __pawntest_marker();
public __pawntest_marker()
{
    return 1;
}

#define TEST(%0) \
    forward test_%0(); \
    public test_%0()

#define BEFORE_EACH() \
    forward test_setup(); \
    public test_setup()

#define AFTER_EACH() \
    forward test_teardown(); \
    public test_teardown()

#define BEFORE_ALL() \
    forward test_suite_setup(); \
    public test_suite_setup()

#define AFTER_ALL() \
    forward test_suite_teardown(); \
    public test_suite_teardown()

#define FIXTURE_SETUP(%0) \
    forward fixture_%0_setup(); \
    public fixture_%0_setup()

#define FIXTURE_TEARDOWN(%0) \
    forward fixture_%0_teardown(); \
    public fixture_%0_teardown()

#define USE_FIXTURE(%0) \
    if (!__pawntest_use_fixture("fixture_" #%0 "_setup", "fixture_" #%0 "_teardown", __file, __line)) return 0

#include <pawntest/scenario_assertions>

#define ASSERT(%0) \
    if (!__pawntest_assert((%0), #%0, __file, __line)) return 0

#define ASSERT_TRUE(%0) \
    ASSERT(%0)

#define ASSERT_FALSE(%0) \
    if (!__pawntest_assert(!(%0), "!(" #%0 ")", __file, __line)) return 0

#define ASSERT_EQ(%0,%1) \
    if (!__pawntest_assert_eq((%0), (%1), #%0, #%1, __file, __line)) return 0

#define ASSERT_NE(%0,%1) \
    if (!__pawntest_assert_ne((%0), (%1), #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_EQ(%0,%1) \
    if (!__pawntest_assert_str_eq(%0, %1, #%0, #%1, __file, __line)) return 0

#define ASSERT_GT(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 1, #%0, #%1, __file, __line)) return 0

#define ASSERT_GE(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 2, #%0, #%1, __file, __line)) return 0

#define ASSERT_LT(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 3, #%0, #%1, __file, __line)) return 0

#define ASSERT_LE(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 4, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_CONTAINS(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 1, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_PREFIX(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 2, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_SUFFIX(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 3, #%0, #%1, __file, __line)) return 0

#define ASSERT_FLOAT_NEAR(%0,%1,%2) \
    if (!__pawntest_assert_float_near((%0), (%1), (%2), #%0, #%1, __file, __line)) return 0

#define ASSERT_ARRAY_EQ(%0,%1,%2) \
    if (!__pawntest_assert_array_eq(%0, %1, (%2), #%0, #%1, __file, __line)) return 0

#define ASSERT_BETWEEN(%0,%1,%2) \
    if (!__pawntest_assert_between((%0), (%1), (%2), #%0, __file, __line)) return 0

#define ASSERT_HAS_FLAG(%0,%1) \
    if (!__pawntest_assert_has_flag((%0), (%1), #%0, __file, __line)) return 0

#define ASSERT_ARRAY_CONTAINS(%0,%1,%2) \
    if (!__pt_assert_array_contains(%0, (%1), (%2), #%0, __file, __line)) return 0

#define ASSERT_STR_IEQ(%0,%1) \
    if (!__pawntest_assert_str_ieq(%0, %1, #%0, #%1, __file, __line)) return 0

#define ASSERT_EQ_MSG(%0,%1,%2) \
    if (!__pawntest_assert_eq_message((%0), (%1), #%0, #%1, (%2), __file, __line)) return 0

#define EXPECT_ERROR(%0,%1) \
    if (!__pawntest_expect_error(#%0, (%1), __file, __line)) return 0

#define ASSERT_NO_ERROR(%0) \
    if (!__pawntest_expect_no_error(#%0, __file, __line)) return 0

#define XFAIL(%0,%1) \
    return __pawntest_xfail(#%0, (%1), __file, __line)

#define TEST_CASE(%0,%1,%2) \
    forward test_%0(); \
    public test_%0() { return %1(%2); }

#define TEST_CASE2(%0,%1,%2,%3) \
    forward test_%0(); \
    public test_%0() { return %1(%2, %3); }

#define TEST_CASE3(%0,%1,%2,%3,%4) \
    forward test_%0(); \
    public test_%0() { return %1(%2, %3, %4); }

#define TEST_TAG(%0) \
    forward ptt_%0(); \
    public ptt_%0() { return 1; }

#define ASSERT_SNAPSHOT(%0,%1) \
    if (!__pawntest_snapshot((%0), (%1), __file, __line)) return 0

#define FUZZ_INT(%0,%1,%2,%3) \
    return __pawntest_fuzz_int(#%0, (%1), (%2), (%3), __file, __line)

#define CHECK(%0) \
    __pawntest_assert((%0), #%0, __file, __line)

#define CHECK_EQ(%0,%1) \
    __pawntest_assert_eq((%0), (%1), #%0, #%1, __file, __line)

#define CHECK_NE(%0,%1) \
    __pawntest_assert_ne((%0), (%1), #%0, #%1, __file, __line)

#define CHECK_STR_EQ(%0,%1) \
    __pawntest_assert_str_eq(%0, %1, #%0, #%1, __file, __line)

#define FAIL(%0) \
    return __pawntest_fail(%0, __file, __line)

#define SKIP(%0) \
    return __pawntest_skip(%0, __file, __line)

#define MOCK_RETURN(%0,%1) \
    __pawntest_mock_return(#%0, (%1))

#define MOCK_RETURN_ONCE(%0,%1) \
    __pawntest_mock_return_once(#%0, (%1))

#define MOCK_OUTPUT(%0,%1,%2) \
    __pawntest_mock_output(#%0, (%1), (%2))

#define MOCK_OUTPUT_STRING(%0,%1,%2,%3) \
    __pawntest_mock_output_string(#%0, (%1), (%2), (%3))

#define MOCK_CALLBACK(%0,%1) \
    __pawntest_mock_callback(#%0, #%1)

#define MOCK_RESET(%0) \
    __pawntest_mock_reset(#%0)

#define EXPECT_NATIVE_CALLS(%0,%1) \
    __pawntest_expect_calls(#%0, (%1), (%1), __file, __line)

#define EXPECT_NATIVE_CALLS_BETWEEN(%0,%1,%2) \
    __pawntest_expect_calls(#%0, (%1), (%2), __file, __line)

#define EXPECT_NATIVE_CALLED(%0) \
    __pawntest_expect_calls(#%0, 1, 2147483647, __file, __line)

#define EXPECT_NATIVE_NOT_CALLED(%0) \
    __pawntest_expect_calls(#%0, 0, 0, __file, __line)

#define EXPECT_NATIVE_ARG(%0,%1,%2,%3) \
    __pawntest_expect_arg(#%0, (%1), (%2), (%3), __file, __line)

#define EXPECT_NATIVE_STRING_ARG(%0,%1,%2,%3) \
    __pawntest_expect_string_arg(#%0, (%1), (%2), (%3), __file, __line)

#define EXPECT_NATIVE_ORDER(%0) \
    __pawntest_expect_order(#%0, __file, __line)

#define TEST_SCHEDULE(%0,%1) \
    __pawntest_schedule((%0), #%1)

#define TEST_ADVANCE_TIME(%0) \
    __pawntest_advance_time((%0))

#define TEST_RUN_PENDING() \
    __pawntest_run_pending()

#define TEST_NOW() \
    __pawntest_now()
