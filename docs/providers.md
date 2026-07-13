# Native providers

Pawn providers implement plugin natives in a separate AMX. Add provider source
or AMX files with `--provider` or the `providers` config field.

```toml
providers = ["pawntest/inventory.provider.pwn"]
```

## Registration

Include `<pawntest/provider>` and register each native during initialization:

```pawn
#include <pawntest/provider>

new inventory[16][32];

forward PawntestProviderInit();
public PawntestProviderInit()
{
    PROVIDE_NATIVE(Inventory_Add, Provider_InventoryAdd);
    PROVIDE_NATIVE(Inventory_Count, Provider_InventoryCount);
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
```

Native arguments are passed to the handler unchanged. Cells and floats can be
used directly. Strings, arrays, and references point into the test AMX and use
the provider bridge:

```pawn
PROVIDER_ARG_STRING(argument, output);
PROVIDER_ARG_ARRAY_CELL(argument, offset);
PROVIDER_SET_CELL(argument, value);
PROVIDER_SET_ARRAY_CELL(argument, offset, value);
PROVIDER_SET_STRING(argument, value, size);
```

Argument indexes are zero-based. Use `PROVIDER_CALL0` through
`PROVIDER_CALL3` to invoke a public in the test AMX.

## Lifecycle

Providers may define these publics:

```pawn
PawntestProviderInit();
PawntestProviderBeforeTest();
PawntestProviderAfterTest();
PawntestProviderShutdown();
```

Provider memory is restored before each test under test isolation. Suite
isolation preserves provider memory for the full test file.

Duplicate provider registrations fail. Explicit mocks override providers.
Custom Go natives override providers and mocks.
