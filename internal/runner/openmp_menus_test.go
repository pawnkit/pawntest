package runner

import (
	"slices"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestMenuScenarioStoresItems(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "Shop", "Item", "Health"
	menuID := callScenarioNative(t, vm, "CreateMenu", 100, 2, floatCell(10), floatCell(20), floatCell(100), floatCell(80))

	if menuID != 0 {
		t.Fatalf("menu ID = %d, want 0", menuID)
	}

	if result := callScenarioNative(t, vm, "SetMenuColumnHeader", menuID, 0, 200); result != 1 {
		t.Fatalf("SetMenuColumnHeader returned %d", result)
	}

	if row := callScenarioNative(t, vm, "AddMenuItem", menuID, 0, 300); row != 0 {
		t.Fatalf("AddMenuItem row = %d, want 0", row)
	}

	menu := menuScenarioState(t, registry).menus[int(menuID)]
	if menu.title != "Shop" || menu.columns != 2 || menu.headers[0] != "Item" {
		t.Fatalf("unexpected menu: %#v", menu)
	}

	if !slices.Equal(menu.items[0], []string{"Health"}) || menu.widths != [2]float32{100, 80} {
		t.Fatalf("unexpected menu items: %#v", menu)
	}
}

func TestMenuScenarioTracksVisibilityAndDisabledRows(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Menu"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	menuID := callScenarioNative(t, vm, "CreateMenu", 200, 1, 0, 0, floatCell(100), 0)

	if shown := callScenarioNative(t, vm, "ShowMenuForPlayer", menuID, playerID); shown != 1 {
		t.Fatalf("ShowMenuForPlayer returned %d", shown)
	}

	if current := callScenarioNative(t, vm, "GetPlayerMenu", playerID); current != menuID {
		t.Fatalf("GetPlayerMenu = %d, want %d", current, menuID)
	}

	callScenarioNative(t, vm, "DisableMenuRow", menuID, 2)

	if disabled := callScenarioNative(t, vm, "IsMenuRowDisabled", menuID, 2); disabled != 1 {
		t.Fatalf("IsMenuRowDisabled = %d", disabled)
	}
}

func TestMenuScenarioInvokesSelectionCallback(t *testing.T) {
	base := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{100: "Alice", 200: "Menu"}}
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

	menuID, err := base.natives["CreateMenu"](vm, []backend.Cell{200, 1, 0, 0, floatCell(100), 0})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := base.natives["ShowMenuForPlayer"](vm, []backend.Cell{menuID, playerID}); err != nil {
		t.Fatal(err)
	}

	selected, err := base.natives["__pt_menu_select"](vm, []backend.Cell{playerID, 3})
	if err != nil {
		t.Fatal(err)
	}

	if selected != 1 || vm.callback != "OnPlayerSelectedMenuRow" || !slices.Equal(vm.args, []backend.Cell{playerID, 3}) {
		t.Fatalf("selected=%d callback=%q args=%v", selected, vm.callback, vm.args)
	}
}

func TestMenuScenarioCloneIsolatesState(t *testing.T) {
	state := newMenuState()
	state.menus[0] = &testMenu{title: "Original", items: [2][]string{{"Item"}}, disabledRows: map[int]bool{}, visible: map[int]bool{0: true}}

	clone, ok := state.Clone().(*menuState)
	if !ok {
		t.Fatal("cloned scenario was not menu state")
	}

	clone.menus[0].title = "Changed"
	clone.menus[0].items[0][0] = "Changed"
	delete(clone.menus[0].visible, 0)

	if state.menus[0].title != "Original" || state.menus[0].items[0][0] != "Item" || !state.menus[0].visible[0] {
		t.Fatal("menu clone shared mutable state")
	}
}

func menuScenarioState(t *testing.T, registry *scenarioRegistry) *menuState {
	t.Helper()

	return registry.Menus
}
