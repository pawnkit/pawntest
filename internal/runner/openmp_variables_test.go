package runner

import (
	"slices"
	"testing"
)

func TestServerVariableScenarioStoresTypes(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "count", "name"
	vm.strings[300], vm.strings[400] = "ratio", "Alice"

	callScenarioNative(t, vm, "SetSVarInt", 100, 42)
	callScenarioNative(t, vm, "SetSVarString", 200, 400)
	callScenarioNative(t, vm, "SetSVarFloat", 300, floatCell(1.5))

	state := variableScenarioState(t, registry)
	if state.server.values["count"] != (testVariable{variableType: variableInt, integer: 42}) {
		t.Fatalf("unexpected integer variable: %#v", state.server.values["count"])
	}

	if state.server.values["name"].text != "Alice" || state.server.values["ratio"].floating != 1.5 {
		t.Fatalf("unexpected variables: %#v", state.server.values)
	}

	if upper := callScenarioNative(t, vm, "GetSVarsUpperIndex"); upper != 2 {
		t.Fatalf("GetSVarsUpperIndex = %d, want 2", upper)
	}

	if variableType := callScenarioNative(t, vm, "GetSVarType", 200); variableType != variableString {
		t.Fatalf("GetSVarType = %d, want string", variableType)
	}
}

func TestServerVariableScenarioPreservesEnumerationOrder(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "first", "second", "third"
	callScenarioNative(t, vm, "SetSVarInt", 100, 1)
	callScenarioNative(t, vm, "SetSVarInt", 200, 2)
	callScenarioNative(t, vm, "SetSVarInt", 300, 3)
	callScenarioNative(t, vm, "DeleteSVar", 200)

	order := variableScenarioState(t, registry).server.order
	if !slices.Equal(order, []string{"first", "third"}) {
		t.Fatalf("variable order = %v", order)
	}
}

func TestPlayerVariableScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200], vm.strings[300] = "Alice", "Bob", "score"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)
	callScenarioNative(t, vm, "SetPVarInt", firstPlayer, 300, 10)
	callScenarioNative(t, vm, "SetPVarInt", secondPlayer, 300, 20)

	state := variableScenarioState(t, registry)
	if state.players[int(firstPlayer)].values["score"].integer != 10 || state.players[int(secondPlayer)].values["score"].integer != 20 {
		t.Fatalf("player variables were not isolated: %#v", state.players)
	}
}

func TestVariableScenarioCloneIsolatesState(t *testing.T) {
	state := newVariableState()
	setVariable(&state.server, "name", testVariable{variableType: variableString, text: "Original"})

	clone, ok := state.Clone().(*variableState)
	if !ok {
		t.Fatal("cloned scenario was not variable state")
	}

	setVariable(&clone.server, "name", testVariable{variableType: variableString, text: "Changed"})
	setVariable(&clone.server, "new", testVariable{variableType: variableInt, integer: 1})

	if state.server.values["name"].text != "Original" || len(state.server.order) != 1 {
		t.Fatal("variable clone shared mutable state")
	}
}

func variableScenarioState(t *testing.T, registry *scenarioRegistry) *variableState {
	t.Helper()

	return registry.Variables
}
