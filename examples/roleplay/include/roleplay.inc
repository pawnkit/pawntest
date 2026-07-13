#if defined _roleplay_included
    #endinput
#endif
#define _roleplay_included

#define DIALOG_LOGIN (100)
#define HTTP_GET (1)

native GivePlayerMoney(playerid, amount);
native SendClientMessage(playerid, colour, const message[]);
native SetPlayerPos(playerid, Float:x, Float:y, Float:z);
native SetPlayerCheckpoint(playerid, Float:x, Float:y, Float:z, Float:radius);
native DisablePlayerCheckpoint(playerid);
native ShowPlayerDialog(playerid, dialogid, style, const title[], const body[], const button1[], const button2[]);
native bool:HTTP(index, method, const url[], const data[], const callback[]);
native Plugin_Audit(playerid, const event[]);

native DB:DB_Open(const name[]);
native DB_Close(DB:database);
native DBResult:DB_ExecuteQuery(DB:database, const query[]);
native DB_FreeResultSet(DBResult:result);
native DB_GetRowCount(DBResult:result);

#include <inventory>

new bool:gLoggedIn[16];
new gPaydayRuns;

stock ResetRoleplayState()
{
    for (new playerid; playerid < sizeof gLoggedIn; playerid++)
    {
        gLoggedIn[playerid] = false;
    }
    gPaydayRuns = 0;
}

stock CalculateSpeedingFine(speed, limit)
{
    if (speed <= limit)
    {
        return 0;
    }
    return (speed - limit) * 20;
}

stock OpenLoginDialog(playerid)
{
    return ShowPlayerDialog(playerid, DIALOG_LOGIN, 1, "Account", "Enter any password", "Login", "Leave");
}

stock StartDelivery(playerid)
{
    return SetPlayerCheckpoint(playerid, 50.0, 0.0, 3.0, 2.0);
}

stock bool:RequestServiceStatus(playerid)
{
    return HTTP(playerid, HTTP_GET, "api.example.test/status", "", "OnServiceStatus");
}

stock DB:OpenAccountDatabase()
{
    new DB:database = DB_Open(":memory:");
    new DBResult:result = DB_ExecuteQuery(database, "CREATE TABLE accounts (name TEXT)");
    DB_FreeResultSet(result);
    return database;
}

forward OnPlayerConnect(playerid);
public OnPlayerConnect(playerid)
{
    gLoggedIn[playerid] = false;
    SendClientMessage(playerid, -1, "Welcome. Use /login to continue.");
    return 1;
}

forward OnPlayerDisconnect(playerid, reason);
public OnPlayerDisconnect(playerid, reason)
{
    #pragma unused reason
    gLoggedIn[playerid] = false;
    return 1;
}

forward OnPlayerSpawn(playerid);
public OnPlayerSpawn(playerid)
{
    SetPlayerPos(playerid, 0.0, 0.0, 3.0);
    Inventory_Add(playerid, 1, 1);
    return 1;
}

forward OnPlayerCommandText(playerid, cmdtext[]);
public OnPlayerCommandText(playerid, cmdtext[])
{
    #pragma unused cmdtext
    OpenLoginDialog(playerid);
    return 1;
}

forward OnDialogResponse(playerid, dialogid, response, listitem, inputtext[]);
public OnDialogResponse(playerid, dialogid, response, listitem, inputtext[])
{
    #pragma unused listitem, inputtext
    if (dialogid != DIALOG_LOGIN || !response)
    {
        return 0;
    }

    gLoggedIn[playerid] = true;
    GivePlayerMoney(playerid, 250);
    Plugin_Audit(playerid, "login");
    SendClientMessage(playerid, -1, "Login successful");
    return 1;
}

forward OnPlayerEnterCheckpoint(playerid);
public OnPlayerEnterCheckpoint(playerid)
{
    GivePlayerMoney(playerid, 1000);
    DisablePlayerCheckpoint(playerid);
    SendClientMessage(playerid, -1, "Delivery complete");
    return 1;
}

forward OnPlayerEnterVehicle(playerid, vehicleid, ispassenger);
public OnPlayerEnterVehicle(playerid, vehicleid, ispassenger)
{
    #pragma unused vehicleid, ispassenger
    SendClientMessage(playerid, -1, "Vehicle entered");
    return 1;
}

forward OnPlayerStateChange(playerid, newstate, oldstate);
public OnPlayerStateChange(playerid, newstate, oldstate)
{
    #pragma unused playerid, newstate, oldstate
    return 1;
}

forward OnVehicleDeath(vehicleid, killerid);
public OnVehicleDeath(vehicleid, killerid)
{
    #pragma unused vehicleid
    SendClientMessage(killerid, -1, "Your vehicle was destroyed");
    return 1;
}

forward OnVehicleSpawn(vehicleid);
public OnVehicleSpawn(vehicleid)
{
    #pragma unused vehicleid
    return 1;
}

forward OnServiceStatus(index, responseCode, data[]);
public OnServiceStatus(index, responseCode, data[])
{
    #pragma unused data
    if (responseCode == 200)
    {
        SendClientMessage(index, -1, "Service online");
    }
    return 1;
}

forward RunPayday();
public RunPayday()
{
    gPaydayRuns++;
    return 1;
}
