#include <pawntest>
#include <inventory>

TEST(grants_starter_items)
{
    GrantStarterItems(0);

    ASSERT_EQ(Inventory_Count(0, 1), 2);
    ASSERT_EQ(Inventory_Count(0, 5), 1);
}

TEST(provider_state_is_isolated)
{
    ASSERT_EQ(Inventory_Count(0, 1), 0);
    ASSERT_EQ(Inventory_Count(0, 5), 0);
}
