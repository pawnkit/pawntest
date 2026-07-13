package runner

import "testing"

func TestServerScenarioStoresSettings(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Freeroam"

	callScenarioNative(t, vm, "SetWeather", 10)
	callScenarioNative(t, vm, "SetWorldTime", 18)
	callScenarioNative(t, vm, "SetGravity", floatCell(0.01))
	callScenarioNative(t, vm, "SetGameModeText", 100)
	callScenarioNative(t, vm, "ToggleChatTextReplacement", 1)
	callScenarioNative(t, vm, "AllowAdminTeleport", 1)

	state := serverScenarioState(t, registry)
	if state.weather != 10 || state.worldTime != 18 || state.gravity != 0.01 || state.gameModeText != "Freeroam" {
		t.Fatalf("unexpected server settings: %#v", state)
	}

	if !state.chatReplacement || !state.adminTeleport {
		t.Fatalf("unexpected server toggles: %#v", state)
	}
}

func TestServerScenarioStoresRulesAndNicknameCharacters(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "language", "English", "Alice-Test"

	if result := callScenarioNative(t, vm, "AddServerRule", 100, 200); result != 1 {
		t.Fatalf("AddServerRule returned %d", result)
	}

	if valid := callScenarioNative(t, vm, "IsValidServerRule", 100); valid != 1 {
		t.Fatalf("IsValidServerRule = %d", valid)
	}

	if valid := callScenarioNative(t, vm, "IsValidNickName", 300); valid != 0 {
		t.Fatalf("IsValidNickName = %d before allowing hyphen", valid)
	}

	callScenarioNative(t, vm, "AllowNickNameCharacter", '-', 1)

	if valid := callScenarioNative(t, vm, "IsValidNickName", 300); valid != 1 {
		t.Fatalf("IsValidNickName = %d after allowing hyphen", valid)
	}

	if serverScenarioState(t, registry).rules["language"] != "English" {
		t.Fatalf("unexpected server rules: %#v", serverScenarioState(t, registry).rules)
	}
}

func TestServerScenarioReportsPoolSizes(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	vehicleID := callScenarioNative(t, vm, "__pt_vehicle_create", 411, 0, 0, 0)
	actorID := callScenarioNative(t, vm, "__pt_actor_create", 7, 0, 0, 0, 0)

	if size := callScenarioNative(t, vm, "GetPlayerPoolSize"); size != playerID {
		t.Fatalf("GetPlayerPoolSize = %d, want %d", size, playerID)
	}

	if size := callScenarioNative(t, vm, "GetVehiclePoolSize"); size != vehicleID {
		t.Fatalf("GetVehiclePoolSize = %d, want %d", size, vehicleID)
	}

	if size := callScenarioNative(t, vm, "GetActorPoolSize"); size != actorID {
		t.Fatalf("GetActorPoolSize = %d, want %d", size, actorID)
	}
}

func TestServerScenarioCloneIsolatesState(t *testing.T) {
	state := newServerState()
	state.rules["language"] = "English"

	clone, ok := state.Clone().(*serverState)
	if !ok {
		t.Fatal("cloned scenario was not server state")
	}

	clone.rules["language"] = "French"
	clone.nicknameCharacters['-'] = true

	if state.rules["language"] != "English" || state.nicknameCharacters['-'] {
		t.Fatal("server clone shared mutable state")
	}
}

func serverScenarioState(t *testing.T, registry *scenarioRegistry) *serverState {
	t.Helper()

	return registry.Server
}
