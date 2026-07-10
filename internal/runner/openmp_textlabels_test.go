package runner

import "testing"

func TestTextLabelScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Welcome", "Updated"
	labelID := callScenarioNative(t, vm, "Create3DTextLabel", 100, -1, floatCell(1), floatCell(2), floatCell(3), floatCell(50), 4, 1)

	if labelID != 0 {
		t.Fatalf("text label ID = %d, want 0", labelID)
	}

	if updated := callScenarioNative(t, vm, "Update3DTextLabelText", labelID, 0xFF, 200); updated != 1 {
		t.Fatalf("Update3DTextLabelText returned %d", updated)
	}

	if result := callScenarioNative(t, vm, "Set3DTextLabelDrawDistance", labelID, floatCell(75)); result != 1 {
		t.Fatalf("Set3DTextLabelDrawDistance returned %d", result)
	}

	label := textLabelScenarioState(t, registry).labels[int(labelID)]
	if label.text != "Updated" || label.colour != 0xFF || label.position != [3]float32{1, 2, 3} {
		t.Fatalf("unexpected text label: %#v", label)
	}

	if label.drawDistance != 75 || !label.lineOfSight || label.world != 4 {
		t.Fatalf("unexpected text label settings: %#v", label)
	}
}

func TestTextLabelScenarioStreamingTracksAttachment(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Label"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	labelID := callScenarioNative(t, vm, "Create3DTextLabel", 200, -1, 0, 0, 0, floatCell(10), 0, 0)

	callScenarioNative(t, vm, "SetPlayerPos", playerID, floatCell(100), 0, 0)

	if streamed := callScenarioNative(t, vm, "Is3DTextLabelStreamedIn", playerID, labelID); streamed != 0 {
		t.Fatalf("distant label streamed = %d", streamed)
	}

	if attached := callScenarioNative(t, vm, "Attach3DTextLabelToPlayer", labelID, playerID, 0, 0, 0); attached != 1 {
		t.Fatalf("Attach3DTextLabelToPlayer returned %d", attached)
	}

	if streamed := callScenarioNative(t, vm, "Is3DTextLabelStreamedIn", playerID, labelID); streamed != 1 {
		t.Fatalf("attached label streamed = %d", streamed)
	}
}

func TestPlayerTextLabelScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "Alice", "Bob", "Private"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)

	firstLabel := callScenarioNative(t, vm, "CreatePlayer3DTextLabel", firstPlayer, 300, -1, 0, 0, 0, floatCell(50), -1, -1, 0)

	secondLabel := callScenarioNative(t, vm, "CreatePlayer3DTextLabel", secondPlayer, 300, -1, 0, 0, 0, floatCell(50), -1, -1, 0)
	if firstLabel != 0 || secondLabel != 0 {
		t.Fatalf("player label IDs = %d and %d, want independent ID 0", firstLabel, secondLabel)
	}

	state := textLabelScenarioState(t, registry)
	if len(state.playerLabels) != 2 || state.playerLabels[int(firstPlayer)][0].text != "Private" {
		t.Fatalf("unexpected player text labels: %#v", state.playerLabels)
	}
}

func TestTextLabelScenarioCloneIsolatesState(t *testing.T) {
	state := newTextLabelState()
	state.labels[0] = newTestTextLabel("Original", -1, [3]float32{1, 2, 3}, 50, 0)

	clone, ok := state.Clone().(*textLabelState)
	if !ok {
		t.Fatal("cloned scenario was not text label state")
	}

	clone.labels[0].text = "Changed"

	clone.labels[0].position[0] = 99
	if state.labels[0].text != "Original" || state.labels[0].position[0] != 1 {
		t.Fatal("text label clone shared mutable state")
	}
}

func textLabelScenarioState(t *testing.T, registry *scenarioRegistry) *textLabelState {
	t.Helper()

	state, ok := registry.modules[6].(*textLabelState)
	if !ok {
		t.Fatal("scenario registry did not contain text label state")
	}

	return state
}
