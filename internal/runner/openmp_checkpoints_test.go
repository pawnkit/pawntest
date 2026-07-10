package runner

import "testing"

func TestCheckpointScenarioTracksPlayerPosition(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	callScenarioNative(t, vm, "SetPlayerPos", playerID, floatCell(10), floatCell(20), floatCell(30))

	result := callScenarioNative(t, vm, "SetPlayerCheckpoint", playerID, floatCell(10), floatCell(20), floatCell(30), floatCell(2))
	if result != 1 {
		t.Fatalf("SetPlayerCheckpoint returned %d", result)
	}

	if active := callScenarioNative(t, vm, "IsPlayerCheckpointActive", playerID); active != 1 {
		t.Fatalf("IsPlayerCheckpointActive = %d", active)
	}

	if inside := callScenarioNative(t, vm, "IsPlayerInCheckpoint", playerID); inside != 1 {
		t.Fatalf("IsPlayerInCheckpoint = %d", inside)
	}

	callScenarioNative(t, vm, "SetPlayerPos", playerID, floatCell(20), floatCell(20), floatCell(30))

	if inside := callScenarioNative(t, vm, "IsPlayerInCheckpoint", playerID); inside != 0 {
		t.Fatalf("IsPlayerInCheckpoint = %d outside radius", inside)
	}

	if disabled := callScenarioNative(t, vm, "DisablePlayerCheckpoint", playerID); disabled != 1 {
		t.Fatalf("DisablePlayerCheckpoint returned %d", disabled)
	}

	if checkpointScenarioState(t, registry).checkpoints[int(playerID)].active {
		t.Fatal("checkpoint remained active")
	}
}

func TestRaceCheckpointScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)

	result := callScenarioNative(t, vm, "SetPlayerRaceCheckpoint", playerID, 1, floatCell(1), floatCell(2), floatCell(3), floatCell(4), floatCell(5), floatCell(6), floatCell(3))
	if result != 1 {
		t.Fatalf("SetPlayerRaceCheckpoint returned %d", result)
	}

	checkpoint := checkpointScenarioState(t, registry).races[int(playerID)]
	if checkpoint.checkpointType != 1 || checkpoint.position != [3]float32{1, 2, 3} || checkpoint.next != [3]float32{4, 5, 6} || checkpoint.radius != 3 {
		t.Fatalf("unexpected race checkpoint: %#v", checkpoint)
	}
}

func TestCheckpointScenarioCloneIsolatesState(t *testing.T) {
	state := newCheckpointState()
	state.checkpoints[0] = checkpoint{active: true, position: [3]float32{1, 2, 3}, radius: 4}
	state.races[0] = raceCheckpoint{checkpoint: checkpoint{active: true}, checkpointType: 1}

	clone, ok := state.Clone().(*checkpointState)
	if !ok {
		t.Fatal("cloned scenario was not checkpoint state")
	}

	value := clone.checkpoints[0]
	value.position[0] = 99
	clone.checkpoints[0] = value
	race := clone.races[0]
	race.checkpointType = 2
	clone.races[0] = race

	if state.checkpoints[0].position[0] != 1 || state.races[0].checkpointType != 1 {
		t.Fatal("checkpoint clone shared mutable state")
	}
}

func checkpointScenarioState(t *testing.T, registry *scenarioRegistry) *checkpointState {
	t.Helper()

	state, ok := registry.modules[5].(*checkpointState)
	if !ok {
		t.Fatal("scenario registry did not contain checkpoint state")
	}

	return state
}
