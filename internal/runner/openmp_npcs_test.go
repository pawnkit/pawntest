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

func npcScenarioState(t *testing.T, registry *scenarioRegistry) *npcState {
	t.Helper()

	state, ok := registry.modules[14].(*npcState)
	if !ok {
		t.Fatal("scenario registry did not contain NPC state")
	}

	return state
}
