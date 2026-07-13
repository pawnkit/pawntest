#include <pawntest/provider>

new inventory[16][16];

forward PawntestProviderInit();
public PawntestProviderInit()
{
    PROVIDE_NATIVE(Inventory_Add, Provider_InventoryAdd, "i,i,i");
    PROVIDE_NATIVE(Inventory_Count, Provider_InventoryCount, "i,i");
    return 1;
}

forward Provider_InventoryAdd(playerid, itemid, amount);
public Provider_InventoryAdd(playerid, itemid, amount)
{
    inventory[playerid][itemid] += amount;
    return 1;
}

forward Provider_InventoryCount(playerid, itemid);
public Provider_InventoryCount(playerid, itemid)
{
    return inventory[playerid][itemid];
}
