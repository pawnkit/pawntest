package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestNPCCoreScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Guard"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)

	if npcID != 0 {
		t.Fatalf("NPC ID = %d, want 0", npcID)
	}

	for _, call := range []struct {
		name   string
		params []backend.Cell
	}{
		{name: "NPC_Spawn", params: []backend.Cell{npcID}},
		{name: "NPC_SetPos", params: []backend.Cell{npcID, floatCell(10), floatCell(20), floatCell(30)}},
		{name: "NPC_SetSkin", params: []backend.Cell{npcID, 7}},
		{name: "NPC_SetHealth", params: []backend.Cell{npcID, floatCell(75)}},
		{name: "NPC_SetArmour", params: []backend.Cell{npcID, floatCell(50)}},
		{name: "NPC_SetVirtualWorld", params: []backend.Cell{npcID, 2}},
	} {
		if result := callScenarioNative(t, vm, call.name, call.params...); result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	npc := npcScenarioState(t, registry).npcs[int(npcID)]
	if npc.name != "Guard" || !npc.spawned || npc.skin != 7 || npc.health != 75 || npc.armour != 50 {
		t.Fatalf("unexpected NPC: %#v", npc)
	}

	if npc.position != [3]float32{10, 20, 30} || npc.world != 2 {
		t.Fatalf("unexpected NPC transform: %#v", npc)
	}
}

func TestNPCCoreScenarioTracksMovementAndStreaming(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Guard"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	npcID := callScenarioNative(t, vm, "NPC_Create", 200)
	callScenarioNative(t, vm, "NPC_Spawn", npcID)

	if streamed := callScenarioNative(t, vm, "NPC_IsStreamedIn", npcID, playerID); streamed != 1 {
		t.Fatalf("NPC_IsStreamedIn = %d", streamed)
	}

	if moved := callScenarioNative(t, vm, "NPC_MoveToPlayer", npcID, playerID); moved != 1 {
		t.Fatalf("NPC_MoveToPlayer returned %d", moved)
	}

	if moving := callScenarioNative(t, vm, "NPC_IsMovingToPlayer", npcID, playerID); moving != 1 {
		t.Fatalf("NPC_IsMovingToPlayer = %d", moving)
	}

	callScenarioNative(t, vm, "NPC_SetVirtualWorld", npcID, 2)

	if streamed := callScenarioNative(t, vm, "NPC_IsStreamedIn", npcID, playerID); streamed != 0 {
		t.Fatalf("NPC_IsStreamedIn = %d in different world", streamed)
	}
}

func TestNPCCoreScenarioUsesVehicles(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Driver"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)
	vehicleID := callScenarioNative(t, vm, "__pt_vehicle_create", 411, 0, 0, 0)

	if result := callScenarioNative(t, vm, "NPC_PutInVehicle", npcID, vehicleID, 0); result != 1 {
		t.Fatalf("NPC_PutInVehicle returned %d", result)
	}

	if vehicle := callScenarioNative(t, vm, "NPC_GetVehicleID", npcID); vehicle != vehicleID {
		t.Fatalf("NPC_GetVehicleID = %d, want %d", vehicle, vehicleID)
	}

	if seat := callScenarioNative(t, vm, "NPC_GetVehicleSeat", npcID); seat != 0 {
		t.Fatalf("NPC_GetVehicleSeat = %d", seat)
	}
}

func TestNPCCoreScenarioCloneIsolatesState(t *testing.T) {
	state := newNPCState()
	state.npcs[0] = &testNPC{name: "Original", position: [3]float32{1, 2, 3}}

	clone, ok := state.Clone().(*npcState)
	if !ok {
		t.Fatal("cloned scenario was not NPC state")
	}

	clone.npcs[0].name = "Changed"
	clone.npcs[0].position[0] = 99

	if state.npcs[0].name != "Original" || state.npcs[0].position[0] != 1 {
		t.Fatal("NPC clone shared state")
	}
}

func TestNPCCombatScenarioStoresWeaponsAndAim(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Guard"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	npcID := callScenarioNative(t, vm, "NPC_Create", 200)

	callScenarioNative(t, vm, "NPC_SetWeapon", npcID, 24)
	callScenarioNative(t, vm, "NPC_SetAmmo", npcID, 50)
	callScenarioNative(t, vm, "NPC_SetAmmoInClip", npcID, 7)
	callScenarioNative(t, vm, "NPC_AimAtPlayer", npcID, playerID, 1)
	callScenarioNative(t, vm, "NPC_SetWeaponAccuracy", npcID, 24, floatCell(0.75))

	npc := npcScenarioState(t, registry).npcs[int(npcID)]
	if npc.weapon != 24 || npc.ammo != 50 || npc.ammoInClip != 7 {
		t.Fatalf("unexpected NPC weapon state: %#v", npc)
	}

	if !npc.aiming || !npc.shooting || npc.aimPlayer != int(playerID) || npc.weaponAccuracy[24] != 0.75 {
		t.Fatalf("unexpected NPC aim state: %#v", npc)
	}
}

func TestNPCAnimationScenarioStoresAnimation(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "Guard", "PED", "WALK_player"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)

	result := callScenarioNative(t, vm, "NPC_ApplyAnimation", npcID, 200, 300, floatCell(4.1), 1, 0, 1, 0, 500)
	if result != 1 {
		t.Fatalf("NPC_ApplyAnimation returned %d", result)
	}

	npc := npcScenarioState(t, registry).npcs[int(npcID)]
	if !npc.hasAnimation || npc.animation.library != "PED" || npc.animation.name != "WALK_player" || !npc.animation.loop {
		t.Fatalf("unexpected NPC animation: %#v", npc.animation)
	}

	callScenarioNative(t, vm, "NPC_ClearAnimations", npcID)

	if npc.hasAnimation {
		t.Fatal("NPC animation remained active")
	}
}

func TestNPCCombatCloneIsolatesWeaponSettings(t *testing.T) {
	state := newNPCState()
	state.npcs[0] = &testNPC{weaponAccuracy: map[int]float32{24: 0.5}, weaponSkills: map[int]int{0: 100}}

	clone, ok := state.Clone().(*npcState)
	if !ok {
		t.Fatal("cloned scenario was not NPC state")
	}

	clone.npcs[0].weaponAccuracy[24] = 1
	clone.npcs[0].weaponSkills[0] = 999

	if state.npcs[0].weaponAccuracy[24] != 0.5 || state.npcs[0].weaponSkills[0] != 100 {
		t.Fatal("NPC combat clone shared mutable state")
	}
}

func TestNPCPlaybackAndRecordScenario(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Guard", "routes/guard.rec"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)
	recordID := callScenarioNative(t, vm, "NPC_LoadRecord", 200)

	if started := callScenarioNative(t, vm, "NPC_StartPlaybackEx", npcID, recordID); started != 1 {
		t.Fatalf("NPC_StartPlaybackEx returned %d", started)
	}

	if playing := callScenarioNative(t, vm, "NPC_IsPlayingPlayback", npcID); playing != 1 {
		t.Fatalf("NPC_IsPlayingPlayback = %d", playing)
	}

	if paused := callScenarioNative(t, vm, "NPC_PausePlayback", npcID, 1); paused != 1 {
		t.Fatalf("NPC_PausePlayback returned %d", paused)
	}

	state := npcScenarioState(t, registry)
	if state.records[int(recordID)] != "routes/guard.rec" || !state.npcs[int(npcID)].pausedPlayback {
		t.Fatalf("unexpected playback state: %#v", state.npcs[int(npcID)])
	}
}

func TestNPCPathScenarioStoresPoints(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Guard"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)
	pathID := callScenarioNative(t, vm, "NPC_CreatePath")
	callScenarioNative(t, vm, "NPC_AddPointToPath", pathID, floatCell(10), floatCell(20), floatCell(30), floatCell(0.5))
	callScenarioNative(t, vm, "NPC_AddPointToPath", pathID, floatCell(40), floatCell(50), floatCell(60), floatCell(1))

	if count := callScenarioNative(t, vm, "NPC_GetPathPointCount", pathID); count != 2 {
		t.Fatalf("NPC_GetPathPointCount = %d", count)
	}

	if moved := callScenarioNative(t, vm, "NPC_MoveByPath", npcID, pathID); moved != 1 {
		t.Fatalf("NPC_MoveByPath returned %d", moved)
	}

	npc := npcScenarioState(t, registry).npcs[int(npcID)]
	if npc.currentPath != int(pathID) || npc.moveTarget != [3]float32{10, 20, 30} {
		t.Fatalf("unexpected path movement: %#v", npc)
	}
}

func TestNPCNodeScenarioTracksPlayback(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Guard"
	npcID := callScenarioNative(t, vm, "NPC_Create", 100)
	callScenarioNative(t, vm, "NPC_OpenNode", 5)

	if played := callScenarioNative(t, vm, "NPC_PlayNode", npcID, 5); played != 1 {
		t.Fatalf("NPC_PlayNode returned %d", played)
	}

	if paused := callScenarioNative(t, vm, "NPC_PausePlayingNode", npcID); paused != 1 {
		t.Fatalf("NPC_PausePlayingNode returned %d", paused)
	}

	if result := callScenarioNative(t, vm, "NPC_IsPlayingNodePaused", npcID); result != 1 {
		t.Fatalf("NPC_IsPlayingNodePaused = %d", result)
	}
}

func TestNPCNavigationCloneIsolatesPaths(t *testing.T) {
	state := newNPCState()
	state.paths[0] = &npcPath{points: []npcPathPoint{{position: [3]float32{1, 2, 3}}}}
	state.records[0] = "record.rec"
	state.nodes[0] = &npcNode{open: true}

	clone, ok := state.Clone().(*npcState)
	if !ok {
		t.Fatal("cloned scenario was not NPC state")
	}

	clone.paths[0].points[0].position[0] = 99
	clone.records[0] = "changed.rec"
	clone.nodes[0].open = false

	if state.paths[0].points[0].position[0] != 1 || state.records[0] != "record.rec" || !state.nodes[0].open {
		t.Fatal("NPC navigation clone shared mutable state")
	}
}

func npcScenarioState(t *testing.T, registry *scenarioRegistry) *npcState {
	t.Helper()

	state, ok := registry.modules[14].(*npcState)
	if !ok {
		t.Fatal("scenario registry did not contain NPC state")
	}

	return state
}
