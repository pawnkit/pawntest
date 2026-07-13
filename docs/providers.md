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
`PROVIDER_CALL8` to invoke a public in the test AMX.
Callback arguments are cells or floats; string and array callback arguments are
not supported by the current ABI.

Signature kinds are `i` (cell), `f` (float), `s` (input string), `r` (cell
reference), `S:n` (output string), `a:n` (input array), and `A:n` (output
array). `n` is the zero-based argument containing the buffer or array length.

## Lifecycle

Providers may define these publics:

```pawn
PawntestProviderInit();
PawntestProviderBeforeTest();
PawntestProviderAfterTest();
PawntestProviderShutdown();
```

Lifecycle callbacks return nonzero on success. Use `PROVIDER_TEST_NAME` during
before/after callbacks to read the active test name.

Provider memory is restored before each test under test isolation. Suite
isolation preserves provider memory for the full test file.

Duplicate provider registrations fail. Explicit mocks override providers.
Custom Go natives override providers and mocks.

Providers are trusted dependency code. The provider bridge can read and write
declared test-AMX buffers. Bounds are enforced from registered signatures.
Precompiled providers must match `PAWNTEST_PROVIDER_ABI`.
