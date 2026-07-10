package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestGangZoneScenarioTracksVisibility(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	zoneID := callScenarioNative(t, vm, "GangZoneCreate", floatCell(0), floatCell(0), floatCell(100), floatCell(100))

	if zoneID != 0 {
		t.Fatalf("gang zone ID = %d, want 0", zoneID)
	}

	if shown := callScenarioNative(t, vm, "GangZoneShowForPlayer", playerID, zoneID, 0xFF); shown != 1 {
		t.Fatalf("GangZoneShowForPlayer returned %d", shown)
	}

	if flashed := callScenarioNative(t, vm, "GangZoneFlashForPlayer", playerID, zoneID, 0xAA); flashed != 1 {
		t.Fatalf("GangZoneFlashForPlayer returned %d", flashed)
	}

	view := gangZoneScenarioState(t, registry).zones[int(zoneID)].views[int(playerID)]
	if !view.visible || !view.flashing || view.colour != 0xFF || view.flash != 0xAA {
		t.Fatalf("unexpected gang zone view: %#v", view)
	}
}

func TestGangZoneScenarioUsesPlayerPosition(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	zoneID := callScenarioNative(t, vm, "GangZoneCreate", 0, 0, floatCell(100), floatCell(100))

	callScenarioNative(t, vm, "SetPlayerPos", playerID, floatCell(50), floatCell(50), 0)

	if inside := callScenarioNative(t, vm, "IsPlayerInGangZone", playerID, zoneID); inside != 1 {
		t.Fatalf("IsPlayerInGangZone = %d inside bounds", inside)
	}

	callScenarioNative(t, vm, "SetPlayerPos", playerID, floatCell(150), floatCell(50), 0)

	if inside := callScenarioNative(t, vm, "IsPlayerInGangZone", playerID, zoneID); inside != 0 {
		t.Fatalf("IsPlayerInGangZone = %d outside bounds", inside)
	}
}

func TestPlayerGangZoneScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Bob"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)

	firstZone := callScenarioNative(t, vm, "CreatePlayerGangZone", firstPlayer, 0, 0, floatCell(10), floatCell(10))

	secondZone := callScenarioNative(t, vm, "CreatePlayerGangZone", secondPlayer, 0, 0, floatCell(20), floatCell(20))
	if firstZone != 0 || secondZone != 0 {
		t.Fatalf("player zone IDs = %d and %d, want independent ID 0", firstZone, secondZone)
	}

	callScenarioNative(t, vm, "PlayerGangZoneShow", firstPlayer, firstZone, 0xFF)

	if visible := callScenarioNative(t, vm, "IsPlayerGangZoneVisible", firstPlayer, firstZone); visible != 1 {
		t.Fatalf("IsPlayerGangZoneVisible = %d", visible)
	}

	if owners := len(gangZoneScenarioState(t, registry).playerZones); owners != 2 {
		t.Fatalf("player zone owners = %d, want 2", owners)
	}
}

func TestGangZoneScenarioCloneIsolatesState(t *testing.T) {
	state := newGangZoneState()
	state.zones[0] = newTestGangZone([]backend.Cell{0, 0, floatCell(10), floatCell(10)})
	state.zones[0].views[0] = gangZoneView{visible: true}

	clone, ok := state.Clone().(*gangZoneState)
	if !ok {
		t.Fatal("cloned scenario was not gang zone state")
	}

	clone.zones[0].bounds[0] = 5
	delete(clone.zones[0].views, 0)

	if state.zones[0].bounds[0] != 0 || !state.zones[0].views[0].visible {
		t.Fatal("gang zone clone shared mutable state")
	}
}

func gangZoneScenarioState(t *testing.T, registry *scenarioRegistry) *gangZoneState {
	t.Helper()

	state, ok := registry.modules[8].(*gangZoneState)
	if !ok {
		t.Fatal("scenario registry did not contain gang zone state")
	}

	return state
}
