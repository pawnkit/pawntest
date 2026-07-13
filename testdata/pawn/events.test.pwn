#include <pawntest>

#pragma rational Float

native GivePlayerMoney(playerid, money);
native SendClientMessage(playerid, colour, const message[]);
native GetPlayerState(playerid);
native MoveObject(objectid, Float:x, Float:y, Float:z, Float:speed, Float:rx = -1000.0, Float:ry = -1000.0, Float:rz = -1000.0);
native SetPlayerCheckpoint(playerid, Float:x, Float:y, Float:z, Float:radius);
native UseGangZoneCheck(zoneid, bool:enable);
native SelectTextDraw(playerid, colour);

forward OnPlayerConnect(playerid);
forward OnPlayerSpawn(playerid);
forward OnPlayerCommandText(playerid, cmdtext[]);
forward OnPlayerKeyStateChange(playerid, newkeys, oldkeys);
forward OnPlayerDeath(playerid, killerid, reason);
forward OnPlayerDisconnect(playerid, reason);
forward OnPlayerEnterVehicle(playerid, vehicleid, ispassenger);
forward OnPlayerExitVehicle(playerid, vehicleid);
forward OnPlayerStateChange(playerid, newstate, oldstate);
forward OnPlayerGiveDamage(playerid, damagedid, Float:amount, weaponid, bodypart);
forward OnPlayerTakeDamage(playerid, issuerid, Float:amount, weaponid, bodypart);
forward OnPlayerGiveDamageActor(playerid, actorid, Float:amount, weaponid, bodypart);
forward OnPlayerWeaponShot(playerid, weaponid, hittype, hitid, Float:x, Float:y, Float:z);
forward OnObjectMoved(objectid);
forward OnPlayerStreamIn(playerid, forplayerid);
forward OnPlayerStreamOut(playerid, forplayerid);
forward OnPlayerEnterCheckpoint(playerid);
forward OnPlayerLeaveCheckpoint(playerid);
forward OnPlayerEnterGangZone(playerid, zoneid);
forward OnPlayerLeaveGangZone(playerid, zoneid);
forward OnVehicleDamageStatusUpdate(vehicleid, playerid);
forward OnVehicleDeath(vehicleid, killerid);
forward OnVehicleSpawn(vehicleid);
forward OnPlayerClickTextDraw(playerid, textid);
forward OnPlayerClickPlayer(playerid, clickedplayerid, source);
forward OnPlayerClickMap(playerid, Float:x, Float:y, Float:z);

new bool:object_moved;
new bool:vehicle_spawned;

public OnPlayerConnect(playerid)
{
    SendClientMessage(playerid, -1, "connected");
    return 1;
}

public OnPlayerSpawn(playerid)
{
    GivePlayerMoney(playerid, 500);
    return 1;
}

public OnPlayerCommandText(playerid, cmdtext[])
{
    SendClientMessage(playerid, -1, cmdtext);
    return 1;
}

public OnPlayerKeyStateChange(playerid, newkeys, oldkeys)
{
    #pragma unused newkeys, oldkeys
    SendClientMessage(playerid, -1, "keys changed");
    return 1;
}

public OnPlayerDeath(playerid, killerid, reason)
{
    #pragma unused killerid, reason
    SendClientMessage(playerid, -1, "died");
    return 1;
}

public OnPlayerDisconnect(playerid, reason)
{
    #pragma unused reason
    SendClientMessage(playerid, -1, "disconnected");
    return 1;
}

public OnPlayerEnterVehicle(playerid, vehicleid, ispassenger)
{
    #pragma unused playerid, vehicleid, ispassenger
    return 1;
}

public OnPlayerExitVehicle(playerid, vehicleid)
{
    #pragma unused playerid, vehicleid
    return 1;
}

public OnPlayerStateChange(playerid, newstate, oldstate)
{
    #pragma unused playerid, newstate, oldstate
    return 1;
}

public OnPlayerGiveDamage(playerid, damagedid, Float:amount, weaponid, bodypart)
{
    #pragma unused damagedid, amount, weaponid, bodypart
    SendClientMessage(playerid, -1, "gave damage");
    return 1;
}

public OnPlayerTakeDamage(playerid, issuerid, Float:amount, weaponid, bodypart)
{
    #pragma unused issuerid, amount, weaponid, bodypart
    SendClientMessage(playerid, -1, "took damage");
    return 1;
}

public OnPlayerGiveDamageActor(playerid, actorid, Float:amount, weaponid, bodypart)
{
    #pragma unused actorid, amount, weaponid, bodypart
    SendClientMessage(playerid, -1, "damaged actor");
    return 1;
}

public OnPlayerWeaponShot(playerid, weaponid, hittype, hitid, Float:x, Float:y, Float:z)
{
    #pragma unused weaponid, hittype, hitid, x, y, z
    SendClientMessage(playerid, -1, "fired weapon");
    return 1;
}

public OnObjectMoved(objectid)
{
    #pragma unused objectid
    object_moved = true;
    return 1;
}

public OnPlayerStreamIn(playerid, forplayerid)
{
    #pragma unused playerid
    SendClientMessage(forplayerid, -1, "streamed in");
    return 1;
}

public OnPlayerStreamOut(playerid, forplayerid)
{
    #pragma unused playerid
    SendClientMessage(forplayerid, -1, "streamed out");
    return 1;
}

public OnPlayerEnterCheckpoint(playerid)
{
    SendClientMessage(playerid, -1, "entered checkpoint");
    return 1;
}

public OnPlayerLeaveCheckpoint(playerid)
{
    SendClientMessage(playerid, -1, "left checkpoint");
    return 1;
}

public OnPlayerEnterGangZone(playerid, zoneid)
{
    #pragma unused zoneid
    SendClientMessage(playerid, -1, "entered gang zone");
    return 1;
}

public OnPlayerLeaveGangZone(playerid, zoneid)
{
    #pragma unused zoneid
    SendClientMessage(playerid, -1, "left gang zone");
    return 1;
}

public OnVehicleDamageStatusUpdate(vehicleid, playerid)
{
    #pragma unused vehicleid
    SendClientMessage(playerid, -1, "vehicle status changed");
    return 1;
}

public OnVehicleDeath(vehicleid, killerid)
{
    #pragma unused vehicleid
    SendClientMessage(killerid, -1, "vehicle destroyed");
    return 1;
}

public OnVehicleSpawn(vehicleid)
{
    #pragma unused vehicleid
    vehicle_spawned = true;
    return 1;
}

public OnPlayerClickTextDraw(playerid, textid)
{
    #pragma unused textid
    SendClientMessage(playerid, -1, "selection cancelled");
    return 1;
}

public OnPlayerClickPlayer(playerid, clickedplayerid, source)
{
    #pragma unused clickedplayerid, source
    SendClientMessage(playerid, -1, "player clicked");
    return 1;
}

public OnPlayerClickMap(playerid, Float:x, Float:y, Float:z)
{
    #pragma unused x, y, z
    SendClientMessage(playerid, -1, "map clicked");
    return 1;
}

TEST(player_event_flow)
{
    new playerid = TEST_CONNECT_PLAYER("Alice");
    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_PLAYER_MESSAGE(playerid, "connected");

    ASSERT_EQ(TEST_SPAWN_PLAYER(playerid), 1);
    ASSERT_PLAYER_MONEY(playerid, 500);

    TEST_PLAYER_COMMAND(playerid, "/help");
    ASSERT_PLAYER_MESSAGE(playerid, "/help");

    TEST_PLAYER_KEYS(playerid, 4, 0, 0);
    ASSERT_PLAYER_MESSAGE(playerid, "keys changed");

    TEST_KILL_PLAYER(playerid, -1, 54);
    ASSERT_PLAYER_MESSAGE(playerid, "died");

    TEST_DISCONNECT_PLAYER(playerid, 1);
    ASSERT_PLAYER_DISCONNECTED(playerid);
}

TEST(vehicle_event_flow)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    new vehicleid = TEST_CREATE_VEHICLE(411, 0.0, 0.0, 3.0);

    TEST_ENTER_VEHICLE(playerid, vehicleid, false);
    ASSERT_EQ(GetPlayerState(playerid), 2);

    TEST_EXIT_VEHICLE(playerid);
    ASSERT_EQ(GetPlayerState(playerid), 1);

    TEST_VEHICLE_DAMAGE_STATUS(vehicleid, playerid, 1, 2, 3, 4);
    ASSERT_PLAYER_MESSAGE(playerid, "vehicle status changed");

    TEST_DAMAGE_VEHICLE(vehicleid, playerid, 1000.0);
    ASSERT_PLAYER_MESSAGE(playerid, "vehicle destroyed");

    vehicle_spawned = false;
    TEST_RESPAWN_VEHICLE(vehicleid);
    ASSERT_VEHICLE_HEALTH(vehicleid, 1000.0);
    ASSERT_TRUE(vehicle_spawned);
}

TEST(combat_event_flow)
{
    new attackerid = TEST_CREATE_PLAYER("Attacker");
    new victimid = TEST_CREATE_PLAYER("Victim");
    new actorid = TEST_CREATE_ACTOR(7, 0.0, 0.0, 3.0, 0.0);

    TEST_DAMAGE_PLAYER(victimid, attackerid, 25.0, 24, 3);
    ASSERT_PLAYER_MESSAGE(attackerid, "gave damage");
    ASSERT_PLAYER_MESSAGE(victimid, "took damage");

    TEST_DAMAGE_ACTOR(attackerid, actorid, 30.0, 24, 3);
    ASSERT_ACTOR_HEALTH(actorid, 70.0);
    ASSERT_PLAYER_MESSAGE(attackerid, "damaged actor");

    TEST_WEAPON_SHOT(attackerid, 24, 1, victimid, 0.0, 0.0, 0.0);
    ASSERT_PLAYER_MESSAGE(attackerid, "fired weapon");
}

TEST(movement_and_stream_events)
{
    new subjectid = TEST_CREATE_PLAYER("Subject");
    new viewerid = TEST_CREATE_PLAYER("Viewer");
    new objectid = TEST_CREATE_OBJECT(19379, 0.0, 0.0, 0.0);

    MoveObject(objectid, 10.0, 20.0, 30.0, 2.0);
    TEST_FINISH_OBJECT_MOVE(objectid);
    ASSERT_OBJECT_POS_NEAR(objectid, 10.0, 20.0, 30.0, 0.01);
    ASSERT_TRUE(object_moved);

    new timed_objectid = TEST_CREATE_OBJECT(19379, 0.0, 0.0, 0.0);
    object_moved = false;
    new duration = MoveObject(timed_objectid, 4.0, 0.0, 0.0, 2.0);
    TEST_ADVANCE_TIME(duration);
    ASSERT_OBJECT_POS_NEAR(timed_objectid, 4.0, 0.0, 0.0, 0.01);
    ASSERT_TRUE(object_moved);

    TEST_STREAM_PLAYER_IN(subjectid, viewerid);
    TEST_STREAM_PLAYER_OUT(subjectid, viewerid);
    ASSERT_PLAYER_MESSAGE(viewerid, "streamed in");
    ASSERT_PLAYER_MESSAGE(viewerid, "streamed out");

    SetPlayerCheckpoint(subjectid, 50.0, 0.0, 0.0, 2.0);
    new zoneid = TEST_CREATE_GANGZONE(48.0, -2.0, 52.0, 2.0);
    UseGangZoneCheck(zoneid, true);

    TEST_MOVE_PLAYER(subjectid, 50.0, 0.0, 0.0);
    ASSERT_PLAYER_MESSAGE(subjectid, "entered checkpoint");
    ASSERT_PLAYER_MESSAGE(subjectid, "entered gang zone");

    TEST_MOVE_PLAYER(subjectid, 60.0, 0.0, 0.0);
    ASSERT_PLAYER_MESSAGE(subjectid, "left checkpoint");
    ASSERT_PLAYER_MESSAGE(subjectid, "left gang zone");

    SelectTextDraw(subjectid, -1);
    TEST_CANCEL_TEXTDRAW_SELECTION(subjectid);
    TEST_CLICK_PLAYER(subjectid, viewerid, 0);
    TEST_CLICK_MAP(subjectid, 1.0, 2.0, 3.0);
    ASSERT_PLAYER_MESSAGE(subjectid, "selection cancelled");
    ASSERT_PLAYER_MESSAGE(subjectid, "player clicked");
    ASSERT_PLAYER_MESSAGE(subjectid, "map clicked");
}
