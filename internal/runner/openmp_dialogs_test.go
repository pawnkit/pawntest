package runner

import (
	"slices"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestDialogScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Title", "Choose an option"
	vm.strings[300], vm.strings[400] = "OK", "Cancel"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)

	result := callScenarioNative(t, vm, "ShowPlayerDialog", playerID, 42, 2, 100, 200, 300, 400)
	if result != 1 {
		t.Fatalf("ShowPlayerDialog returned %d", result)
	}

	dialog := dialogScenarioState(t, registry).dialogs[int(playerID)]
	if dialog.id != 42 || dialog.style != 2 || dialog.title != "Title" || dialog.body != "Choose an option" {
		t.Fatalf("unexpected dialog: %#v", dialog)
	}

	if dialog.button1 != "OK" || dialog.button2 != "Cancel" || !dialog.visible {
		t.Fatalf("unexpected dialog buttons: %#v", dialog)
	}

	if id := callScenarioNative(t, vm, "GetPlayerDialogID", playerID); id != 42 {
		t.Fatalf("GetPlayerDialogID = %d", id)
	}
}

func TestDialogScenarioHidesDialog(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Text"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	callScenarioNative(t, vm, "ShowPlayerDialog", playerID, 5, 0, 200, 200, 200, 200)

	if hidden := callScenarioNative(t, vm, "HidePlayerDialog", playerID); hidden != 1 {
		t.Fatalf("HidePlayerDialog returned %d", hidden)
	}

	if id := callScenarioNative(t, vm, "GetPlayerDialogID", playerID); id != -1 {
		t.Fatalf("GetPlayerDialogID = %d after hide", id)
	}
}

func TestDialogScenarioInvokesResponse(t *testing.T) {
	base := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{100: "Alice", 200: "Text", 300: "input"}}
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

	if _, err := base.natives["ShowPlayerDialog"](vm, []backend.Cell{playerID, 9, 1, 200, 200, 200, 200}); err != nil {
		t.Fatal(err)
	}

	response, err := base.natives["__pt_dialog_respond"](vm, []backend.Cell{playerID, 1, 3, 300})
	if err != nil {
		t.Fatal(err)
	}

	if response != 1 || vm.callback != "OnDialogResponse" {
		t.Fatalf("response=%d callback=%q", response, vm.callback)
	}

	want := []backend.Cell{playerID, 9, 1, 3, 300}
	if !slices.Equal(vm.args, want) {
		t.Fatalf("callback args = %v, want %v", vm.args, want)
	}
}

func TestDialogScenarioCloneIsolatesState(t *testing.T) {
	state := newDialogState()
	state.dialogs[0] = testDialog{id: 1, title: "Original", visible: true}

	clone, ok := state.Clone().(*dialogState)
	if !ok {
		t.Fatal("cloned scenario was not dialog state")
	}

	value := clone.dialogs[0]
	value.title = "Changed"
	clone.dialogs[0] = value

	if state.dialogs[0].title != "Original" {
		t.Fatal("dialog clone shared state")
	}
}

type dialogMockVM struct {
	*mockVM
	callback string
	args     []backend.Cell
}

func (vm *dialogMockVM) CallPublic(name string, args ...backend.Cell) (backend.Cell, error) {
	vm.callback = name

	vm.args = append([]backend.Cell(nil), args...)

	return 1, nil
}

func dialogScenarioState(t *testing.T, registry *scenarioRegistry) *dialogState {
	t.Helper()

	return registry.Dialogs
}
