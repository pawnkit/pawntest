package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestVehicleScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)

	vehicleID := callScenarioNative(t, vm, "CreateVehicle", 411, floatCell(10), floatCell(20), floatCell(30), floatCell(90), 1, 2, 60)
	if vehicleID != 1 {
		t.Fatalf("vehicle ID = %d, want 1", vehicleID)
	}

	for _, call := range []struct {
		name   string
		params []backend.Cell
	}{
		{name: "SetVehicleHealth", params: []backend.Cell{vehicleID, floatCell(750)}},
		{name: "SetVehiclePos", params: []backend.Cell{vehicleID, floatCell(40), floatCell(50), floatCell(60)}},
		{name: "ChangeVehicleColor", params: []backend.Cell{vehicleID, 3, 4}},
		{name: "ChangeVehiclePaintjob", params: []backend.Cell{vehicleID, 2}},
		{name: "SetVehicleVirtualWorld", params: []backend.Cell{vehicleID, 7}},
		{name: "AddVehicleComponent", params: []backend.Cell{vehicleID, 1010}},
	} {
		if result := callScenarioNative(t, vm, call.name, call.params...); result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	state := vehicleScenarioState(t, registry)

	vehicle := state.vehicles[int(vehicleID)]
	if vehicle.model != 411 || vehicle.health != 750 || vehicle.position != [3]float32{40, 50, 60} {
		t.Fatalf("unexpected vehicle state: %#v", vehicle)
	}

	if vehicle.colour1 != 3 || vehicle.colour2 != 4 || vehicle.paintjob != 2 || vehicle.world != 7 {
		t.Fatalf("unexpected vehicle appearance: %#v", vehicle)
	}

	if !vehicle.components[1010] {
		t.Fatal("vehicle component was not stored")
	}
}

func TestVehicleScenarioTracksOccupants(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = "Alice"

	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)

	vehicleID := callScenarioNative(t, vm, "__pt_vehicle_create", 411, 0, 0, 0)
	if result := callScenarioNative(t, vm, "PutPlayerInVehicle", playerID, vehicleID, 0); result != 1 {
		t.Fatalf("PutPlayerInVehicle returned %d", result)
	}

	if got := callScenarioNative(t, vm, "GetPlayerVehicleID", playerID); got != vehicleID {
		t.Fatalf("GetPlayerVehicleID = %d, want %d", got, vehicleID)
	}

	if got := callScenarioNative(t, vm, "GetVehicleDriver", vehicleID); got != playerID {
		t.Fatalf("GetVehicleDriver = %d, want %d", got, playerID)
	}

	if got := callScenarioNative(t, vm, "CountVehicleOccupants", vehicleID); got != 1 {
		t.Fatalf("CountVehicleOccupants = %d, want 1", got)
	}

	if result := callScenarioNative(t, vm, "SetVehicleToRespawn", vehicleID); result != 1 {
		t.Fatalf("SetVehicleToRespawn returned %d", result)
	}

	if got := callScenarioNative(t, vm, "IsPlayerInAnyVehicle", playerID); got != 0 {
		t.Fatalf("IsPlayerInAnyVehicle = %d after respawn", got)
	}

	vehicle := vehicleScenarioState(t, registry).vehicles[int(vehicleID)]
	if vehicle.position != vehicle.spawn || vehicle.health != 1000 {
		t.Fatalf("vehicle did not reset: %#v", vehicle)
	}
}

func TestVehicleScenarioCloneIsolatesState(t *testing.T) {
	state := newVehicleState()
	vehicleID := int(state.addVehicle(411, [3]float32{1, 2, 3}, 90, 1, 2, 60))
	state.vehicles[vehicleID].components[1010] = true

	clone, ok := state.Clone().(*vehicleState)
	if !ok {
		t.Fatal("cloned scenario was not vehicle state")
	}

	clone.vehicles[vehicleID].position[0] = 99
	delete(clone.vehicles[vehicleID].components, 1010)

	if state.vehicles[vehicleID].position[0] != 1 || !state.vehicles[vehicleID].components[1010] {
		t.Fatal("vehicle clone shared mutable state")
	}
}

func registeredScenarios(t *testing.T) (*mockVM, *scenarioRegistry) {
	t.Helper()

	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}
	registry := newScenarioRegistry()

	context := &executionContext{state: &nativeState{status: Pass}, mocks: newMockState(), scenarios: registry}
	if err := registry.Register(vm, context); err != nil {
		t.Fatal(err)
	}

	return vm, registry
}

func callScenarioNative(t *testing.T, vm *mockVM, name string, params ...backend.Cell) backend.Cell {
	t.Helper()

	native, ok := vm.natives[name]
	if !ok {
		t.Fatalf("native %s was not registered", name)
	}

	result, err := native(vm, params)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}

	return result
}

func vehicleScenarioState(t *testing.T, registry *scenarioRegistry) *vehicleState {
	t.Helper()

	return registry.Vehicles
}
