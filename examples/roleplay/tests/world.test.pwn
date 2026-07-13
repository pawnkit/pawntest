#include <support>

native SetWeather(weatherid);
native SetWorldTime(hour);
native SetGameModeText(const text[]);
native SetSVarInt(const name[], value);

TEST(delivery_checkpoint)
{
    USE_FIXTURE(spawned_player);

    StartDelivery(testPlayer);
    ASSERT_CHECKPOINT_ACTIVE(testPlayer);

    TEST_MOVE_PLAYER(testPlayer, 50.0, 0.0, 3.0);
    ASSERT_PLAYER_MONEY(testPlayer, 1000);
    ASSERT_PLAYER_MESSAGE(testPlayer, "Delivery complete");
    ASSERT_CHECKPOINT_INACTIVE(testPlayer);
}

TEST(vehicle_lifecycle)
{
    USE_FIXTURE(spawned_player);
    new vehicleid = TEST_CREATE_VEHICLE(411, 0.0, 0.0, 3.0);

    TEST_ENTER_VEHICLE(testPlayer, vehicleid, false);
    ASSERT_PLAYER_MESSAGE(testPlayer, "Vehicle entered");
    TEST_DAMAGE_VEHICLE(vehicleid, testPlayer, 1000.0);
    ASSERT_PLAYER_MESSAGE(testPlayer, "Your vehicle was destroyed");

    TEST_RESPAWN_VEHICLE(vehicleid);
    ASSERT_VEHICLE_HEALTH(vehicleid, 1000.0);
}

TEST(world_entities)
{
    new playerid = TEST_CREATE_PLAYER("Builder");
    new objectid = TEST_CREATE_OBJECT(19379, 10.0, 20.0, 3.0);
    new actorid = TEST_CREATE_ACTOR(7, 15.0, 20.0, 3.0, 90.0);
    new pickupid = TEST_CREATE_PICKUP(1240, 1, 12.0, 20.0, 3.0, 0);
    new labelid = TEST_CREATE_TEXT_LABEL("Dealership", -1, 10.0, 20.0, 4.0, 50.0, 0);
    new zoneid = TEST_CREATE_GANGZONE(0.0, 0.0, 100.0, 100.0);

    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_OBJECT_VALID(objectid);
    ASSERT_ACTOR_VALID(actorid);
    ASSERT_PICKUP_VALID(pickupid);
    ASSERT_TEXT_LABEL_TEXT(labelid, "Dealership");
    ASSERT_GANGZONE_VALID(zoneid);
}

TEST(server_state)
{
    SetWeather(10);
    SetWorldTime(18);
    SetGameModeText("Roleplay");
    SetSVarInt("payday", 60);

    ASSERT_SERVER_WEATHER(10);
    ASSERT_SERVER_TIME(18);
    ASSERT_GAME_MODE_TEXT("Roleplay");
    ASSERT_SVAR_INT("payday", 60);
}

TEST_TAG(scn_delivery_checkpoint)
TEST_TAG(scn_vehicle_lifecycle)
TEST_TAG(scn_world_entities)
TEST_TAG(unit_server_state)
