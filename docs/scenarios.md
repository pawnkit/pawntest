# Scenario fidelity

Scenario natives model deterministic test state. They do not run an open.mp
server or reproduce networking, streaming distance, timers, persistence, or
plugin behavior unless documented.

| Scenario | Modeled behavior | Callback behavior |
|---|---|---|
| Players | Connection state, transforms, stats, equipment, camera, messages | Lifecycle, input, weapon, and damage helpers call callbacks |
| Classes | Class data and selected spawn state | Selection calls `OnPlayerRequestClass`; spawn calls `OnPlayerSpawn` |
| Vehicles | Creation, transforms, health, damage, occupants, components | Transition, damage, death, and respawn helpers call callbacks |
| Objects | Global/player objects, movement, materials, attachments | Virtual time or completion helpers apply targets and call callbacks |
| Actors | State, animations, damage, streaming flags | Damage and stream helpers call callbacks |
| Pickups | Global/player pickup state and visibility | Pickup helpers call callbacks |
| Checkpoints | Active state and geometric containment | Explicit events or `TEST_MOVE_PLAYER` call transition callbacks |
| Dialogs | Visible dialog state and response | Response calls `OnDialogResponse` |
| Menus | Menu state and selection | Selection and exit callbacks are called |
| Textdraws | State, visibility, selection, preview data | Click helpers require a visible selectable textdraw |
| Text labels | State, attachment, visibility | No callbacks |
| Gang zones | Bounds, visibility, flashing, containment | `TEST_MOVE_PLAYER` checks zones enabled with `UseGangZoneCheck` |
| Variables | Player and server variable storage | No callbacks |
| Server | Time, weather, gravity, mode text, messages | No lifecycle callbacks |
| NPCs | State and playback metadata | No network playback simulation |
| Database | In-memory SQLite operations | Callback behavior is modeled locally |
| HTTP | Configured requests and responses | Configured response calls its callback |

Call gamemode callbacks directly when a helper does not trigger them. Use strict
scenarios to fail on unused responses and incomplete configured behavior.
Unknown natives fail unless explicitly mocked, provided, or allowed with
`--allow-unknown-natives`; allowed calls are reported as warnings.

`TEST_CREATE_PLAYER` creates an already spawned player without callbacks. Use
`TEST_CONNECT_PLAYER` and `TEST_SPAWN_PLAYER` for lifecycle behavior.

## Event helpers

| Area | Helpers |
|---|---|
| Lifecycle | `TEST_CONNECT_PLAYER`, `TEST_SPAWN_PLAYER`, `TEST_KILL_PLAYER`, `TEST_DISCONNECT_PLAYER` |
| Input | `TEST_PLAYER_TEXT`, `TEST_PLAYER_COMMAND`, `TEST_PLAYER_KEYS` |
| Combat | `TEST_DAMAGE_PLAYER`, `TEST_DAMAGE_ACTOR`, `TEST_WEAPON_SHOT` |
| Vehicles | Enter, exit, health damage, damage status, and respawn helpers |
| World | `TEST_PICKUP`, `TEST_PLAYER_PICKUP`, checkpoint and race checkpoint enter/leave helpers |
| UI | Dialog, menu, textdraw, player, and map click helpers |
| Movement | `TEST_FINISH_OBJECT_MOVE`, `TEST_FINISH_PLAYER_OBJECT_MOVE` |
| Streaming | Player, vehicle, and actor stream-in/stream-out helpers |
| Position | `TEST_MOVE_PLAYER` applies position and area transitions |

`TEST_CONNECT_PLAYER` returns the player ID. Other event helpers return the
final callback result. Invalid transitions return `0` without a callback.
Damage consumes armour before health. Lethal player damage also calls
`OnPlayerDeath`.

Object movement completes when `TEST_ADVANCE_TIME` reaches the duration returned by
`MoveObject` or `MovePlayerObject`. Stopping or manually completing movement
cancels its pending completion.

Vehicle status damage is separate from vehicle health. Lethal health damage
calls `OnVehicleDeath`; respawning calls `OnVehicleSpawn`. Cancelling textdraw
selection calls `OnPlayerClickTextDraw` with `INVALID_TEXT_DRAW`.
