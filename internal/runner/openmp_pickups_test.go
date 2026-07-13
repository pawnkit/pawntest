package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestPickupScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	pickupID := callScenarioNative(t, vm, "CreatePickup", 1240, 2, floatCell(1), floatCell(2), floatCell(3), 4)

	if pickupID != 0 {
		t.Fatalf("pickup ID = %d, want 0", pickupID)
	}

	for _, call := range []struct {
		name   string
		params []backend.Cell
	}{
		{name: "SetPickupModel", params: []backend.Cell{pickupID, 1242, 1}},
		{name: "SetPickupType", params: []backend.Cell{pickupID, 14, 1}},
		{name: "SetPickupPos", params: []backend.Cell{pickupID, floatCell(4), floatCell(5), floatCell(6), 1}},
		{name: "SetPickupVirtualWorld", params: []backend.Cell{pickupID, 7}},
	} {
		if result := callScenarioNative(t, vm, call.name, call.params...); result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	pickup := pickupScenarioState(t, registry).pickups[int(pickupID)]
	if pickup.model != 1242 || pickup.pickupType != 14 || pickup.world != 7 || pickup.position != [3]float32{4, 5, 6} {
		t.Fatalf("unexpected pickup state: %#v", pickup)
	}
}

func TestPickupScenarioTracksVisibilityAndStreaming(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	pickupID := callScenarioNative(t, vm, "CreatePickup", 1240, 2, 0, 0, 0, 0)

	if streamed := callScenarioNative(t, vm, "IsPickupStreamedIn", playerID, pickupID); streamed != 1 {
		t.Fatalf("IsPickupStreamedIn = %d", streamed)
	}

	if hidden := callScenarioNative(t, vm, "HidePickupForPlayer", playerID, pickupID); hidden != 1 {
		t.Fatalf("HidePickupForPlayer = %d", hidden)
	}

	if streamed := callScenarioNative(t, vm, "IsPickupStreamedIn", playerID, pickupID); streamed != 0 {
		t.Fatalf("hidden pickup streamed = %d", streamed)
	}

	if shown := callScenarioNative(t, vm, "ShowPickupForPlayer", playerID, pickupID); shown != 1 {
		t.Fatalf("ShowPickupForPlayer = %d", shown)
	}
}

func TestPlayerPickupScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Bob"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)

	firstPickup := callScenarioNative(t, vm, "CreatePlayerPickup", firstPlayer, 1240, 2, 0, 0, 0, 0)

	secondPickup := callScenarioNative(t, vm, "CreatePlayerPickup", secondPlayer, 1242, 14, 0, 0, 0, 0)
	if firstPickup != 0 || secondPickup != 0 {
		t.Fatalf("player pickup IDs = %d and %d, want independent ID 0", firstPickup, secondPickup)
	}

	if model := callScenarioNative(t, vm, "GetPlayerPickupModel", firstPlayer, firstPickup); model != 1240 {
		t.Fatalf("first pickup model = %d", model)
	}

	if model := callScenarioNative(t, vm, "GetPlayerPickupModel", secondPlayer, secondPickup); model != 1242 {
		t.Fatalf("second pickup model = %d", model)
	}

	if owners := len(pickupScenarioState(t, registry).playerPickups); owners != 2 {
		t.Fatalf("player pickup owners = %d, want 2", owners)
	}
}

func TestPickupScenarioCloneIsolatesState(t *testing.T) {
	state := newPickupState()
	state.pickups[0] = newTestPickup(1240, 2, [3]float32{1, 2, 3}, 0)
	state.pickups[0].hidden[1] = true

	clone, ok := state.Clone().(*pickupState)
	if !ok {
		t.Fatal("cloned scenario was not pickup state")
	}

	clone.pickups[0].position[0] = 99
	delete(clone.pickups[0].hidden, 1)

	if state.pickups[0].position[0] != 1 || !state.pickups[0].hidden[1] {
		t.Fatal("pickup clone shared mutable state")
	}
}

func pickupScenarioState(t *testing.T, registry *scenarioRegistry) *pickupState {
	t.Helper()

	return registry.Pickups
}
