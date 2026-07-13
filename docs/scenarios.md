# Scenario fidelity

Scenario natives model deterministic test state. They do not run an open.mp
server or reproduce networking, streaming distance, timers, persistence, or
plugin behavior unless documented.

| Scenario | Modeled behavior | Callback behavior |
|---|---|---|
| Players | Connection state, transforms, stats, equipment, camera, messages | `TEST_CREATE_PLAYER` does not call connect or spawn callbacks |
| Classes | Class data and selected spawn state | Selection calls `OnPlayerRequestClass`; spawn calls `OnPlayerSpawn` |
| Vehicles | Creation, transforms, health, damage, occupants, components | No automatic stream callbacks |
| Objects | Global/player objects, movement, materials, attachments | No automatic movement completion callback |
| Actors | State, animations, damage, streaming flags | No automatic stream callbacks |
| Pickups | Global/player pickup state and visibility | No automatic pickup callback |
| Checkpoints | Active state and geometric containment | No automatic enter/leave callbacks |
| Dialogs | Visible dialog state and response | Response calls `OnDialogResponse` |
| Menus | Menu state and selection | Selection and exit callbacks are called |
| Textdraws | State, visibility, selection, preview data | No automatic click callback |
| Text labels | State, attachment, visibility | No callbacks |
| Gang zones | Bounds, visibility, flashing, containment | No automatic enter/leave callbacks |
| Variables | Player and server variable storage | No callbacks |
| Server | Time, weather, gravity, mode text, messages | No lifecycle callbacks |
| NPCs | State and playback metadata | No network playback simulation |
| Database | In-memory SQLite operations | Callback behavior is modeled locally |
| HTTP | Configured requests and responses | Configured response calls its callback |

Call gamemode callbacks directly when a helper does not trigger them. Use strict
scenarios to fail on unused responses and incomplete configured behavior.
Unknown natives fail unless explicitly mocked, provided, or allowed with
`--allow-unknown-natives`; allowed calls are reported as warnings.
