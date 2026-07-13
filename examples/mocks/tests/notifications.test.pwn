#include <pawntest>
#include <notifications>

TEST(sends_welcome_notification)
{
    MOCK_RETURN(Plugin_SendNotification, 1);
    EXPECT_NATIVE_CALLS(Plugin_SendNotification, 1);
    EXPECT_NATIVE_ARG(Plugin_SendNotification, 0, 0, 7);
    EXPECT_NATIVE_STRING_ARG(Plugin_SendNotification, 0, 1, "Welcome");
    EXPECT_NATIVE_STRING_ARG(Plugin_SendNotification, 0, 2, "Your account is ready");

    ASSERT_EQ(SendWelcomeNotification(7), 1);
}

TEST(handles_plugin_failure)
{
    MOCK_RETURN(Plugin_SendNotification, 0);

    ASSERT_EQ(SendWelcomeNotification(7), 0);
    EXPECT_NATIVE_CALLS(Plugin_SendNotification, 1);
}
