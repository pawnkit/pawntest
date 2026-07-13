#if defined _example_mode_included
    #endinput
#endif
#define _example_mode_included

native GivePlayerMoney(playerid, amount);
native SendClientMessage(playerid, colour, const message[]);
native SetPlayerPos(playerid, Float:x, Float:y, Float:z);

forward OnPlayerConnect(playerid);
public OnPlayerConnect(playerid)
{
    SendClientMessage(playerid, -1, "Welcome to the server");
    return 1;
}

forward OnPlayerSpawn(playerid);
public OnPlayerSpawn(playerid)
{
    GivePlayerMoney(playerid, 500);
    SetPlayerPos(playerid, 100.0, 200.0, 10.0);
    return 1;
}

forward OnPlayerCommandText(playerid, cmdtext[]);
public OnPlayerCommandText(playerid, cmdtext[])
{
    #pragma unused cmdtext
    SendClientMessage(playerid, -1, "Command handled");
    return 1;
}
