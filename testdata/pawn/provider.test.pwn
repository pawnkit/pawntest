#include <pawntest>

native Inventory_Add(playerid, itemid, amount);
native Inventory_Count(playerid, itemid);
native Inventory_Name(itemid, output[], size = sizeof(output));
native Inventory_Find(const name[]);
native Inventory_Read(&value);
native Inventory_Load(playerid, const callback[]);
native Inventory_Double(const input[], output[], size);
native Inventory_BeforeCount();

new loaded_player = -1;
new bool:loaded_success;

forward OnInventoryLoaded(playerid, bool:success);
public OnInventoryLoaded(playerid, bool:success)
{
    loaded_player = playerid;
    loaded_success = success;
    return 1;
}

TEST(provider_keeps_state)
{
    Inventory_Add(3, 7, 2);
    Inventory_Add(3, 7, 4);
    ASSERT_EQ(Inventory_Count(3, 7), 6);
}

TEST(provider_state_is_isolated)
{
    ASSERT_EQ(Inventory_Count(3, 7), 0);
    ASSERT_EQ(Inventory_BeforeCount(), 1);
}

TEST(provider_marshals_values)
{
    new name[16];
    new value;
    Inventory_Name(7, name);
    Inventory_Read(value);

    ASSERT_STR_EQ(name, "Phone");
    ASSERT_TRUE(Inventory_Find(name));
    ASSERT_EQ(value, 42);

    new input[] = {2, 4, 8};
    new output[3];
    new expected[] = {4, 8, 16};
    Inventory_Double(input, output, sizeof(input));
    ASSERT_ARRAY_EQ(output, expected, sizeof(expected));
}

TEST(provider_calls_test_public)
{
    Inventory_Load(9, "OnInventoryLoaded");
    ASSERT_EQ(loaded_player, 9);
    ASSERT_TRUE(loaded_success);
}

TEST(mock_overrides_provider)
{
    MOCK_RETURN(Inventory_Count, 99);
    ASSERT_EQ(Inventory_Count(3, 7), 99);
}
