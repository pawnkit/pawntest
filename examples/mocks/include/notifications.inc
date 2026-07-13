#if defined _notifications_included
    #endinput
#endif
#define _notifications_included

native Plugin_SendNotification(playerid, const title[], const body[]);

stock SendWelcomeNotification(playerid)
{
    return Plugin_SendNotification(playerid, "Welcome", "Your account is ready");
}
