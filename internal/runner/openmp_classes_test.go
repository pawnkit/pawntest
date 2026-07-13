package runner

import (
	"slices"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestClassScenarioStoresClasses(t *testing.T) {
	vm, registry := registeredScenarios(t)
	classID := callScenarioNative(t, vm, "AddPlayerClassEx", 2, 7, floatCell(10), floatCell(20), floatCell(30), floatCell(90), 24, 50, 31, 100, 0, 0)

	if classID != 0 {
		t.Fatalf("class ID = %d, want 0", classID)
	}

	if count := callScenarioNative(t, vm, "GetAvailableClasses"); count != 1 {
		t.Fatalf("GetAvailableClasses = %d", count)
	}

	class := classScenarioState(t, registry).classes[0]
	if class.team != 2 || class.skin != 7 || class.position != [3]float32{10, 20, 30} || class.angle != 90 {
		t.Fatalf("unexpected class: %#v", class)
	}

	if class.weapons[0] != (classWeapon{weapon: 24, ammo: 50}) || class.weapons[1] != (classWeapon{weapon: 31, ammo: 100}) {
		t.Fatalf("unexpected class weapons: %#v", class.weapons)
	}
}

func TestClassScenarioSelectsAndSpawnsPlayer(t *testing.T) {
	base := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{100: "Alice"}}
	vm := &dialogMockVM{mockVM: base}
	registry := newScenarioRegistry()

	context := &executionContext{state: &nativeState{status: Pass}, mocks: newMockState(), scenarios: registry}
	if err := registry.Register(vm, context); err != nil {
		t.Fatal(err)
	}

	playerID, err := base.natives["__pt_player_create"](vm, []backend.Cell{100})
	if err != nil {
		t.Fatal(err)
	}

	classID, err := base.natives["AddPlayerClassEx"](vm, []backend.Cell{2, 7, floatCell(10), floatCell(20), floatCell(30), floatCell(90), 24, 50})
	if err != nil {
		t.Fatal(err)
	}

	selected, err := base.natives["__pt_class_select"](vm, []backend.Cell{playerID, classID})
	if err != nil {
		t.Fatal(err)
	}

	if selected != 1 || vm.callback != "OnPlayerRequestClass" || !slices.Equal(vm.args, []backend.Cell{playerID, classID}) {
		t.Fatalf("selected=%d callback=%q args=%v", selected, vm.callback, vm.args)
	}

	spawned, err := base.natives["SpawnPlayer"](vm, []backend.Cell{playerID})
	if err != nil {
		t.Fatal(err)
	}

	player := registry.playerState().players[int(playerID)]
	if spawned != 1 || !player.spawned || player.skin != 7 || player.team != 2 || player.x != 10 || player.weapons[24] != 50 {
		t.Fatalf("unexpected spawned player: %#v", player)
	}

	if vm.callback != "OnPlayerSpawn" {
		t.Fatalf("spawn callback = %q", vm.callback)
	}
}

func TestClassScenarioForceSelection(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)

	if result := callScenarioNative(t, vm, "ForceClassSelection", playerID); result != 1 {
		t.Fatalf("ForceClassSelection returned %d", result)
	}

	if !classScenarioState(t, registry).selecting[int(playerID)] || registry.playerState().players[int(playerID)].spawned {
		t.Fatal("player did not enter class selection")
	}
}

func TestClassScenarioCloneIsolatesState(t *testing.T) {
	state := newClassState()
	state.classes = []playerClass{{skin: 7, weapons: [3]classWeapon{{weapon: 24, ammo: 50}}}}
	state.selected[0] = 0

	clone, ok := state.Clone().(*classState)
	if !ok {
		t.Fatal("cloned scenario was not class state")
	}

	clone.classes[0].skin = 10
	clone.selected[0] = 1

	if state.classes[0].skin != 7 || state.selected[0] != 0 {
		t.Fatal("class clone shared state")
	}
}

func classScenarioState(t *testing.T, registry *scenarioRegistry) *classState {
	t.Helper()

	return registry.Classes
}
