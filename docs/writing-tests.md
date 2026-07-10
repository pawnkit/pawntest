# Writing Tests

Name files `<name>.test.pwn` or `<name>.test.inc`, include `<pawntest>`, and use `TEST`:

```pawn
#include <pawntest>

TEST(addition)
{
    ASSERT_EQ(2 + 2, 4);
}

TEST(not_ready)
{
    SKIP("not implemented");
}
```

See [Assertions](assertions.md), [Fixtures](fixtures.md), and [Mocking](mocking.md) for the core test APIs.

## Isolation

Each test starts from the state created by `BEFORE_ALL`. Test setup, mocks, memory, and time are isolated by default. Use `--isolation suite` only when tests must share state.

## Timers

Use the virtual clock for timer-based code:

```pawn
TEST(session_expiry)
{
    TEST_SCHEDULE(1000, expire_session);
    TEST_ADVANCE_TIME(1000);
    TEST_RUN_PENDING();
}
```

The clock starts at zero for each test and does not wait in real time.

## Compiler diagnostics

Declare expected compiler errors or warnings in the source:

```pawn
// pawntest: expect-error 017
TEST(missing_symbol)
{
    return missing_symbol;
}
```

Use a three-digit Pawn diagnostic code.

## Test cases

Generate tests from a callback:

```pawn
stock assert_even(value)
{
    ASSERT_EQ(value % 2, 0);
    return 1;
}

TEST_CASE(even_two, assert_even, 2)
TEST_CASE(even_four, assert_even, 4)
```

Use `TEST_CASE2` or `TEST_CASE3` for more arguments. Keep generated names within Pawn's 31-character symbol limit.

## Expected failures

Mark a known failure with `XFAIL`:

```pawn
TEST(known_defect)
{
    XFAIL(known_failure, "issue 42");
}
```

A passing callback reports `xpass` and fails the run.

## Tags

Attach tags to a generated test name:

```pawn
TEST_CASE(even_two, assert_even, 2)
TEST_TAG(unit_even_two)
TEST_TAG(fast_even_two)
```

Filter them with `pawntest --tags 'unit & !slow'`. Expressions support `!`, `&`, `|`, and parentheses.

## Snapshots

Compare strings with a stored snapshot:

```pawn
ASSERT_SNAPSHOT("player", "{id: 7, name: Alice}");
```

Snapshots are stored in `__snapshots__/<test-file>.snap.json`. Update them with `--update-snapshots`.

## Property tests

Run a callback with generated integers:

```pawn
TEST(integer_bounds)
{
    FUZZ_INT(value_is_bounded, 200, -100, 100);
}
```

Runs are deterministic from `--fuzz-seed`. Failures report the seed and a reduced value.

## Player scenarios

The built-in player model supports common open.mp behavior:

```pawn
TEST(welcomes_player)
{
    new playerid = TEST_CREATE_PLAYER("Alice");

    SetPlayerMoney(playerid, 1000);
    ASSERT_PLAYER_CONNECTED(playerid);
    ASSERT_PLAYER_MONEY(playerid, 1000);
}
```

Other server or plugin natives require mocks or custom Go natives.

## Vehicle scenarios

Create and inspect vehicles without a server:

```pawn
TEST(vehicle_damage)
{
    new vehicleid = TEST_CREATE_VEHICLE(411, 10.0, 20.0, 30.0);

    SetVehicleHealth(vehicleid, 750.0);
    ASSERT_VEHICLE_VALID(vehicleid);
    ASSERT_VEHICLE_MODEL(vehicleid, 411);
    ASSERT_VEHICLE_HEALTH(vehicleid, 750.0);
}
```

The vehicle model tracks transforms, health, appearance, components, damage, trailers, parameters, respawns, and occupants.

## Object scenarios

Create global or player-scoped objects:

```pawn
TEST(moving_gate)
{
    new objectid = TEST_CREATE_OBJECT(19379, 0.0, 0.0, 5.0);

    MoveObject(objectid, 0.0, 0.0, 10.0, 2.0);
    ASSERT_OBJECT_VALID(objectid);
    ASSERT_OBJECT_MODEL(objectid, 19379);
}
```

Use `TEST_CREATE_PLAYER_OBJECT(playerid, modelid, x, y, z)` and `ASSERT_PLAYER_OBJECT_VALID` for player objects. The model tracks transforms, movement, materials, camera collision, and attachments.

## Actor scenarios

Create actors and test their state:

```pawn
TEST(actor_damage)
{
    new actorid = TEST_CREATE_ACTOR(7, 10.0, 20.0, 30.0, 90.0);

    SetActorHealth(actorid, 75.0);
    ASSERT_ACTOR_VALID(actorid);
    ASSERT_ACTOR_SKIN(actorid, 7);
    ASSERT_ACTOR_HEALTH(actorid, 75.0);
}
```

The actor model tracks transforms, health, skins, virtual worlds, invulnerability, animations, and player streaming.

## Pickup scenarios

Create global or player-scoped pickups:

```pawn
TEST(health_pickup)
{
    new pickupid = TEST_CREATE_PICKUP(1240, 2, 10.0, 20.0, 30.0, 0);

    ASSERT_PICKUP_VALID(pickupid);
    ASSERT_PICKUP_MODEL(pickupid, 1240);
}
```

Use `TEST_CREATE_PLAYER_PICKUP(playerid, model, type, x, y, z, world)` for player pickups. The model tracks models, types, positions, virtual worlds, visibility, and streaming.

## Checkpoint scenarios

Checkpoint state follows the player's modeled position:

```pawn
TEST(reaches_checkpoint)
{
    new playerid = TEST_CREATE_PLAYER("Alice");

    SetPlayerCheckpoint(playerid, 10.0, 20.0, 30.0, 2.0);
    SetPlayerPos(playerid, 10.0, 20.0, 30.0);

    ASSERT_CHECKPOINT_ACTIVE(playerid);
    ASSERT_PLAYER_IN_CHECKPOINT(playerid);
}
```

Equivalent `ASSERT_RACE_CHECKPOINT_*` helpers cover race checkpoints.

## Text label scenarios

Create global or player-scoped 3D text labels:

```pawn
TEST(welcome_label)
{
    new labelid = TEST_CREATE_TEXT_LABEL("Welcome", -1, 10.0, 20.0, 30.0, 50.0, 0);

    ASSERT_TEXT_LABEL_VALID(labelid);
    ASSERT_TEXT_LABEL_TEXT(labelid, "Welcome");
}
```

Use `TEST_CREATE_PLAYER_TEXT_LABEL` for private labels. The model tracks text, colour, position, draw distance, virtual world, line of sight, streaming, and attachments.

## Textdraw scenarios

Create and display global or player textdraws:

```pawn
TEST(shows_title)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    new textid = TEST_CREATE_TEXTDRAW(100.0, 50.0, "Title");

    TextDrawShowForPlayer(playerid, textid);
    ASSERT_TEXTDRAW_TEXT(textid, "Title");
    ASSERT_TEXTDRAW_VISIBLE(playerid, textid);
}
```

Use `TEST_CREATE_PLAYER_TEXTDRAW` for private textdraws. The model tracks styling, text, position, visibility, selection, and model previews.

## Gang zone scenarios

Create global or player-scoped gang zones:

```pawn
TEST(zone_entry)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    new zoneid = TEST_CREATE_GANGZONE(0.0, 0.0, 100.0, 100.0);

    GangZoneShowForPlayer(playerid, zoneid, -1);
    SetPlayerPos(playerid, 50.0, 50.0, 0.0);

    ASSERT_GANGZONE_VISIBLE(playerid, zoneid);
    ASSERT_PLAYER_IN_GANGZONE(playerid, zoneid);
}
```

Use `TEST_CREATE_PLAYER_GANGZONE` for private zones. The model tracks bounds, visibility, colours, flashing, and containment.

## Dialog scenarios

Show a dialog and simulate its response:

```pawn
TEST(confirms_action)
{
    new playerid = TEST_CREATE_PLAYER("Alice");

    ShowPlayerDialog(playerid, 10, DIALOG_STYLE_MSGBOX, "Confirm", "Continue?", "Yes", "No");
    ASSERT_DIALOG_VISIBLE(playerid, 10);
    ASSERT_DIALOG_TITLE(playerid, "Confirm");

    TEST_RESPOND_DIALOG(playerid, true, 0, "");
    ASSERT_DIALOG_HIDDEN(playerid);
}
```

`TEST_RESPOND_DIALOG` calls `OnDialogResponse` with the modeled dialog ID and response values.

## Menu scenarios

Build a menu and simulate a row selection:

```pawn
TEST(selects_item)
{
    new playerid = TEST_CREATE_PLAYER("Alice");
    new Menu:menuid = CreateMenu("Shop", 1, 10.0, 20.0, 100.0);

    AddMenuItem(menuid, 0, "Health");
    ShowMenuForPlayer(menuid, playerid);

    ASSERT_MENU_VISIBLE(playerid, menuid);
    ASSERT_MENU_ITEMS(menuid, 0, 1);
    TEST_SELECT_MENU_ROW(playerid, 0);
}
```

Selection and exit helpers call `OnPlayerSelectedMenuRow` and `OnPlayerExitedMenu`.

## Class scenarios

Add a class, select it, and spawn the player:

```pawn
TEST(selects_class)
{
    AddPlayerClass(7, 10.0, 20.0, 30.0, 90.0);
    new playerid = TEST_CREATE_PLAYER("Alice");

    TEST_SELECT_CLASS(playerid, 0);
    ASSERT_CLASS_COUNT(1);
    ASSERT_PLAYER_CLASS(playerid, 0);

    SpawnPlayer(playerid);
    ASSERT_PLAYER_NOT_SELECTING_CLASS(playerid);
}
```

Class selection calls `OnPlayerRequestClass`. Spawning applies the selected spawn data and calls `OnPlayerSpawn`.
