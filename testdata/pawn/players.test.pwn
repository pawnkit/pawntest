#include <pawntest>

#pragma rational Float

native IsPlayerConnected(playerid);
native GetPlayerName(playerid, name[], size);
native SetPlayerPos(playerid, Float:x, Float:y, Float:z);
native GetPlayerPos(playerid, &Float:x, &Float:y, &Float:z);
native SetPlayerMoney(playerid, money);
native GivePlayerMoney(playerid, money);
native GetPlayerMoney(playerid);
native SendClientMessage(playerid, color, const message[]);
native Kick(playerid);

TEST(player_scenario)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    ASSERT_PLAYER_CONNECTED(playerid);

    SetPlayerMoney(playerid, 1000);
    GivePlayerMoney(playerid, 250);
    ASSERT_PLAYER_MONEY(playerid, 1250);

    SetPlayerPos(playerid, 100.0, 200.0, 10.0);
    ASSERT_PLAYER_POS_NEAR(playerid, 100.0, 200.0, 10.0, 0.01);

    SendClientMessage(playerid, -1, "Welcome, Alice");
    ASSERT_PLAYER_MESSAGE_CONTAINS(playerid, "Welcome");

    new name[24];
    GetPlayerName(playerid, name, sizeof name);
    ASSERT_STR_EQ(name, "Alice");

    Kick(playerid);
    ASSERT_PLAYER_DISCONNECTED(playerid);
}
