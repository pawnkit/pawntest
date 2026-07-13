#if defined _roleplay_test_support_included
    #endinput
#endif
#define _roleplay_test_support_included

#include <pawntest>

#pragma rational Float

#include <roleplay>

new testPlayer;

BEFORE_EACH()
{
    ResetRoleplayState();
}

FIXTURE_SETUP(spawned_player)
{
    testPlayer = TEST_CONNECT_PLAYER("Alice");
    TEST_SPAWN_PLAYER(testPlayer);
}

FIXTURE_TEARDOWN(spawned_player)
{
    TEST_DISCONNECT_PLAYER(testPlayer, 1);
}
