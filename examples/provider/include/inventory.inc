#if defined _inventory_included
    #endinput
#endif
#define _inventory_included

native Inventory_Add(playerid, itemid, amount);
native Inventory_Count(playerid, itemid);

stock GrantStarterItems(playerid)
{
    Inventory_Add(playerid, 1, 2);
    Inventory_Add(playerid, 5, 1);
}
