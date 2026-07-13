#include <pawntest>

#pragma rational Float

#include <example_mode>

TEST(player_join_flow)
{
    new playerid = TEST_CONNECT_PLAYER("Alice");
    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_PLAYER_MESSAGE(playerid, "Welcome to the server");

    TEST_SPAWN_PLAYER(playerid);
    ASSERT_PLAYER_MONEY(playerid, 500);
    ASSERT_PLAYER_POS_NEAR(playerid, 100.0, 200.0, 10.0, 0.01);
}

TEST(player_command_flow)
{
    new playerid = TEST_CREATE_PLAYER("Alice");

    TEST_PLAYER_COMMAND(playerid, "/help");
    ASSERT_PLAYER_MESSAGE(playerid, "Command handled");
}
