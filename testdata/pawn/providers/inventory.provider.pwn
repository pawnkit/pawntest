#include <pawntest/provider>

#define MAX_PROVIDER_PLAYERS (16)
#define MAX_PROVIDER_ITEMS (32)

new inventory[MAX_PROVIDER_PLAYERS][MAX_PROVIDER_ITEMS];
new before_count;

forward PawntestProviderInit();
forward PawntestProviderBeforeTest();
forward Provider_InventoryAdd(playerid, itemid, amount);
forward Provider_InventoryCount(playerid, itemid);
forward Provider_InventoryName(itemid, output[], size);
forward Provider_InventoryFind(name[]);
forward Provider_InventoryRead(&value);
forward Provider_InventoryLoad(playerid, callback[]);
forward Provider_InventoryDouble(input[], output[], size);
forward Provider_BeforeCount();

public PawntestProviderInit()
{
    PROVIDE_NATIVE(Inventory_Add, Provider_InventoryAdd, "i,i,i");
    PROVIDE_NATIVE(Inventory_Count, Provider_InventoryCount, "i,i");
    PROVIDE_NATIVE(Inventory_Name, Provider_InventoryName, "i,S:2,i");
    PROVIDE_NATIVE(Inventory_Find, Provider_InventoryFind, "s");
    PROVIDE_NATIVE(Inventory_Read, Provider_InventoryRead, "r");
    PROVIDE_NATIVE(Inventory_Load, Provider_InventoryLoad, "i,s");
    PROVIDE_NATIVE(Inventory_Double, Provider_InventoryDouble, "a:2,A:2,i");
    PROVIDE_NATIVE(Inventory_BeforeCount, Provider_BeforeCount, "");
    return 1;
}

public PawntestProviderBeforeTest()
{
    before_count++;
    return 1;
}

public Provider_InventoryAdd(playerid, itemid, amount)
{
    inventory[playerid][itemid] += amount;
    return 1;
}

public Provider_InventoryCount(playerid, itemid)
{
    return inventory[playerid][itemid];
}

public Provider_InventoryName(itemid, output[], size)
{
    #pragma unused itemid, output, size
    return PROVIDER_SET_STRING(1, "Phone", PROVIDER_ARG_CELL(2));
}

public Provider_InventoryFind(name[])
{
    #pragma unused name
    new value[16];
    PROVIDER_ARG_STRING(0, value);
    return value[0] == 'P' && value[1] == 'h' && value[2] == 'o' && value[3] == 'n' && value[4] == 'e' && value[5] == EOS;
}

public Provider_InventoryRead(&value)
{
    #pragma unused value
    return PROVIDER_SET_CELL(0, 42);
}

public Provider_InventoryLoad(playerid, callback[])
{
    #pragma unused callback
    new callback_name[32];
    PROVIDER_ARG_STRING(1, callback_name);
    return __pawntest_provider_call(callback_name, playerid, 1);
}

public Provider_InventoryDouble(input[], output[], size)
{
    #pragma unused input, output
    for (new index; index < size; index++)
    {
        PROVIDER_SET_ARRAY_CELL(1, index, PROVIDER_ARG_ARRAY_CELL(0, index) * 2);
    }
    return 1;
}

public Provider_BeforeCount()
{
    return before_count;
}
