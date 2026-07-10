package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestTextDrawScenarioStoresProperties(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Title", "Updated"
	drawID := callScenarioNative(t, vm, "TextDrawCreate", floatCell(100), floatCell(50), 100)

	if drawID != 0 {
		t.Fatalf("textdraw ID = %d, want 0", drawID)
	}

	for _, call := range []struct {
		name   string
		params []backend.Cell
	}{
		{name: "TextDrawSetString", params: []backend.Cell{drawID, 200}},
		{name: "TextDrawLetterSize", params: []backend.Cell{drawID, floatCell(1), floatCell(2)}},
		{name: "TextDrawTextSize", params: []backend.Cell{drawID, floatCell(200), floatCell(40)}},
		{name: "TextDrawColor", params: []backend.Cell{drawID, 0xFF}},
		{name: "TextDrawUseBox", params: []backend.Cell{drawID, 1}},
		{name: "TextDrawFont", params: []backend.Cell{drawID, 2}},
		{name: "TextDrawSetSelectable", params: []backend.Cell{drawID, 1}},
		{name: "TextDrawSetPreviewModel", params: []backend.Cell{drawID, 411}},
		{name: "TextDrawSetPreviewRot", params: []backend.Cell{drawID, floatCell(1), floatCell(2), floatCell(3), floatCell(1.5)}},
	} {
		if result := callScenarioNative(t, vm, call.name, call.params...); result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	draw := textDrawScenarioState(t, registry).draws[int(drawID)]
	if draw.text != "Updated" || draw.position != [2]float32{100, 50} || draw.letterSize != [2]float32{1, 2} {
		t.Fatalf("unexpected textdraw: %#v", draw)
	}

	if draw.colour != 0xFF || !draw.box || draw.font != 2 || !draw.selectable || draw.previewModel != 411 {
		t.Fatalf("unexpected textdraw style: %#v", draw)
	}
}

func TestTextDrawScenarioTracksVisibilityAndSelection(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Title"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	drawID := callScenarioNative(t, vm, "TextDrawCreate", 0, 0, 200)

	if shown := callScenarioNative(t, vm, "TextDrawShowForPlayer", playerID, drawID); shown != 1 {
		t.Fatalf("TextDrawShowForPlayer returned %d", shown)
	}

	if visible := callScenarioNative(t, vm, "IsTextDrawVisibleForPlayer", playerID, drawID); visible != 1 {
		t.Fatalf("IsTextDrawVisibleForPlayer = %d", visible)
	}

	if selected := callScenarioNative(t, vm, "SelectTextDraw", playerID, 0xAA); selected != 1 {
		t.Fatalf("SelectTextDraw returned %d", selected)
	}

	selection := textDrawScenarioState(t, registry).selection[int(playerID)]
	if !selection.active || selection.colour != 0xAA {
		t.Fatalf("unexpected selection: %#v", selection)
	}

	callScenarioNative(t, vm, "CancelSelectTextDraw", playerID)

	if _, exists := textDrawScenarioState(t, registry).selection[int(playerID)]; exists {
		t.Fatal("textdraw selection remained active")
	}
}

func TestPlayerTextDrawScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "Alice", "Bob", "Private"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)

	firstDraw := callScenarioNative(t, vm, "CreatePlayerTextDraw", firstPlayer, 0, 0, 300)

	secondDraw := callScenarioNative(t, vm, "CreatePlayerTextDraw", secondPlayer, 0, 0, 300)
	if firstDraw != 0 || secondDraw != 0 {
		t.Fatalf("player textdraw IDs = %d and %d, want independent ID 0", firstDraw, secondDraw)
	}

	callScenarioNative(t, vm, "PlayerTextDrawShow", firstPlayer, firstDraw)

	if visible := callScenarioNative(t, vm, "IsPlayerTextDrawVisible", firstPlayer, firstDraw); visible != 1 {
		t.Fatalf("IsPlayerTextDrawVisible = %d", visible)
	}

	if owners := len(textDrawScenarioState(t, registry).playerDraws); owners != 2 {
		t.Fatalf("player textdraw owners = %d, want 2", owners)
	}
}

func TestTextDrawScenarioCloneIsolatesState(t *testing.T) {
	state := newTextDrawState()
	state.draws[0] = newTestTextDraw(1, 2, "Original")
	state.draws[0].visible[0] = true

	clone, ok := state.Clone().(*textDrawState)
	if !ok {
		t.Fatal("cloned scenario was not textdraw state")
	}

	clone.draws[0].text = "Changed"
	delete(clone.draws[0].visible, 0)

	if state.draws[0].text != "Original" || !state.draws[0].visible[0] {
		t.Fatal("textdraw clone shared mutable state")
	}
}

func textDrawScenarioState(t *testing.T, registry *scenarioRegistry) *textDrawState {
	t.Helper()

	state, ok := registry.modules[7].(*textDrawState)
	if !ok {
		t.Fatal("scenario registry did not contain textdraw state")
	}

	return state
}
