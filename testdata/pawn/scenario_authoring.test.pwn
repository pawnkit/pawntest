#include <pawntest>

#pragma rational Float

TEST(scenario_authoring)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    new vehicleid = TEST_CREATE_VEHICLE(411, 10.0, 20.0, 30.0);
    new objectid = TEST_CREATE_OBJECT(19300, 1.0, 2.0, 3.0);
    new actorid = TEST_CREATE_ACTOR(7, 4.0, 5.0, 6.0, 90.0);
    new pickupid = TEST_CREATE_PICKUP(1240, 1, 7.0, 8.0, 9.0, 0);

    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_VEHICLE_VALID(vehicleid);
    ASSERT_OBJECT_VALID(objectid);
    ASSERT_ACTOR_VALID(actorid);
    ASSERT_PICKUP_VALID(pickupid);
}
