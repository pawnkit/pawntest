#include <support>

TEST(player_join)
{
    USE_FIXTURE(spawned_player);

    ASSERT_PLAYER_CONNECTED(testPlayer);
    ASSERT_PLAYER_MESSAGE_CONTAINS(testPlayer, "Welcome");
    ASSERT_PLAYER_POS_NEAR(testPlayer, 0.0, 0.0, 3.0, 0.01);
    ASSERT_EQ(Inventory_Count(testPlayer, 1), 1);
}

TEST(login_dialog)
{
    USE_FIXTURE(spawned_player);
    MOCK_RETURN(Plugin_Audit, 1);
    EXPECT_NATIVE_CALLS(Plugin_Audit, 1);
    EXPECT_NATIVE_STRING_ARG(Plugin_Audit, 0, 1, "login");

    TEST_PLAYER_COMMAND(testPlayer, "/login");
    ASSERT_DIALOG_VISIBLE(testPlayer, DIALOG_LOGIN);
    ASSERT_DIALOG_TITLE(testPlayer, "Account");

    TEST_RESPOND_DIALOG(testPlayer, true, 0, "secret");
    ASSERT_TRUE(gLoggedIn[testPlayer]);
    ASSERT_PLAYER_MONEY(testPlayer, 250);
    ASSERT_PLAYER_MESSAGE(testPlayer, "Login successful");
    ASSERT_DIALOG_HIDDEN(testPlayer);
}

TEST(provider_isolation)
{
    ASSERT_EQ(Inventory_Count(0, 1), 0);
}

TEST_TAG(scn_player_join)
TEST_TAG(scn_login_dialog)
TEST_TAG(int_provider_isolation)
