package runner

import (
	"testing"
)

func TestObjectScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	objectID := callScenarioNative(t, vm, "CreateObject", 19379, floatCell(1), floatCell(2), floatCell(3), 0, 0, 0, floatCell(200))

	if objectID != 1 {
		t.Fatalf("object ID = %d, want 1", objectID)
	}

	if result := callScenarioNative(t, vm, "SetObjectPos", objectID, floatCell(4), floatCell(5), floatCell(6)); result != 1 {
		t.Fatalf("SetObjectPos returned %d", result)
	}

	if result := callScenarioNative(t, vm, "MoveObject", objectID, floatCell(10), floatCell(20), floatCell(30), floatCell(2), 0, 0, 0); result <= 0 {
		t.Fatalf("MoveObject returned %d", result)
	}

	if result := callScenarioNative(t, vm, "AttachObjectToVehicle", objectID, 7, floatCell(1), floatCell(2), floatCell(3), 0, 0, 0); result != 1 {
		t.Fatalf("AttachObjectToVehicle returned %d", result)
	}

	state := objectScenarioState(t, registry)

	object := state.objects[int(objectID)]
	if object.model != 19379 || object.position != [3]float32{4, 5, 6} || object.targetPos != [3]float32{10, 20, 30} {
		t.Fatalf("unexpected object state: %#v", object)
	}

	if !object.moving || object.moveSpeed != 2 || object.attachment.kind != 1 || object.attachment.id != 7 {
		t.Fatalf("unexpected movement or attachment: %#v", object)
	}
}

func TestPlayerObjectScenarioIsScopedToPlayer(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "Alice", "Bob"
	firstPlayer := callScenarioNative(t, vm, "__pt_player_create", 100)
	secondPlayer := callScenarioNative(t, vm, "__pt_player_create", 200)

	firstObject := callScenarioNative(t, vm, "CreatePlayerObject", firstPlayer, 19379, 0, 0, 0, 0, 0, 0, 0)

	secondObject := callScenarioNative(t, vm, "CreatePlayerObject", secondPlayer, 19400, 0, 0, 0, 0, 0, 0, 0)
	if firstObject != 1 || secondObject != 1 {
		t.Fatalf("player object IDs = %d and %d, want independent ID 1", firstObject, secondObject)
	}

	if model := callScenarioNative(t, vm, "GetPlayerObjectModel", firstPlayer, firstObject); model != 19379 {
		t.Fatalf("first player object model = %d", model)
	}

	if model := callScenarioNative(t, vm, "GetPlayerObjectModel", secondPlayer, secondObject); model != 19400 {
		t.Fatalf("second player object model = %d", model)
	}

	state := objectScenarioState(t, registry)
	if len(state.playerObjects) != 2 {
		t.Fatalf("player object owners = %d, want 2", len(state.playerObjects))
	}
}

func TestObjectScenarioStoresMaterialsAndClones(t *testing.T) {
	state := newObjectState()
	objectID := int(state.addObject(19379, [3]float32{}, [3]float32{}, 0))
	state.objects[objectID].materials[0] = objectMaterial{model: 10765, library: "airport", name: "wall"}

	clone, ok := state.Clone().(*objectState)
	if !ok {
		t.Fatal("cloned scenario was not object state")
	}

	clone.objects[objectID].position[0] = 10
	delete(clone.objects[objectID].materials, 0)

	if state.objects[objectID].position[0] != 0 || len(state.objects[objectID].materials) != 1 {
		t.Fatal("object clone shared mutable state")
	}
}

func objectScenarioState(t *testing.T, registry *scenarioRegistry) *objectState {
	t.Helper()

	state, ok := registry.modules[2].(*objectState)
	if !ok {
		t.Fatal("scenario registry did not contain object state")
	}

	return state
}

func TestObjectScenarioRejectsMissingPlayer(t *testing.T) {
	vm, _ := registeredScenarios(t)

	if objectID := callScenarioNative(t, vm, "CreatePlayerObject", 99, 19379, 0, 0, 0, 0, 0, 0, 0); objectID != -1 {
		t.Fatalf("CreatePlayerObject returned %d for missing player", objectID)
	}
}
