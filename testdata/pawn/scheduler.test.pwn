#include <pawntest>

new callback_count;
new callback_time;

forward scheduled_callback();
public scheduled_callback()
{
    callback_count++;
    callback_time = TEST_NOW();
}

TEST(virtual_time)
{
    callback_count = 0;
    callback_time = -1;

    ASSERT_EQ(TEST_NOW(), 0);
    TEST_SCHEDULE(100, scheduled_callback);
    TEST_ADVANCE_TIME(99);
    ASSERT_EQ(callback_count, 0);
    TEST_ADVANCE_TIME(1);
    ASSERT_EQ(callback_count, 1);
    ASSERT_EQ(callback_time, 100);
}

TEST(pending_callback_order)
{
    callback_count = 0;
    TEST_SCHEDULE(20, scheduled_callback);
    TEST_SCHEDULE(10, scheduled_callback);
    TEST_RUN_PENDING();

    ASSERT_EQ(callback_count, 2);
    ASSERT_EQ(callback_time, 20);
    ASSERT_EQ(TEST_NOW(), 20);
}
