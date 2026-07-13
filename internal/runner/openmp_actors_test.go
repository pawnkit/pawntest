package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestActorScenarioStoresState(t *testing.T) {
	vm, registry := registeredScenarios(t)
	actorID := callScenarioNative(t, vm, "CreateActor", 7, floatCell(1), floatCell(2), floatCell(3), floatCell(90))

	if actorID != 0 {
		t.Fatalf("actor ID = %d, want 0", actorID)
	}

	for _, call := range []struct {
		name   string
		params []backend.Cell
	}{
		{name: "SetActorSkin", params: []backend.Cell{actorID, 10}},
		{name: "SetActorHealth", params: []backend.Cell{actorID, floatCell(75)}},
		{name: "SetActorPos", params: []backend.Cell{actorID, floatCell(4), floatCell(5), floatCell(6)}},
		{name: "SetActorFacingAngle", params: []backend.Cell{actorID, floatCell(180)}},
		{name: "SetActorVirtualWorld", params: []backend.Cell{actorID, 2}},
		{name: "SetActorInvulnerable", params: []backend.Cell{actorID, 1}},
	} {
		if result := callScenarioNative(t, vm, call.name, call.params...); result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	actor := actorScenarioState(t, registry).actors[int(actorID)]
	if actor.skin != 10 || actor.health != 75 || actor.position != [3]float32{4, 5, 6} {
		t.Fatalf("unexpected actor state: %#v", actor)
	}

	if actor.angle != 180 || actor.world != 2 || !actor.invulnerable {
		t.Fatalf("unexpected actor settings: %#v", actor)
	}
}

func TestActorScenarioTracksAnimation(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100], vm.strings[200] = "PED", "WALK_player"
	actorID := callScenarioNative(t, vm, "CreateActor", 7, 0, 0, 0, 0)

	result := callScenarioNative(t, vm, "ApplyActorAnimation", actorID, 100, 200, floatCell(4.1), 1, 0, 1, 0, 500)
	if result != 1 {
		t.Fatalf("ApplyActorAnimation returned %d", result)
	}

	actor := actorScenarioState(t, registry).actors[int(actorID)]
	if !actor.hasAnimation || actor.animation.library != "PED" || actor.animation.name != "WALK_player" {
		t.Fatalf("unexpected actor animation: %#v", actor.animation)
	}

	if !actor.animation.loop || actor.animation.lockX || !actor.animation.lockY || actor.animation.time != 500 {
		t.Fatalf("unexpected actor animation flags: %#v", actor.animation)
	}

	if cleared := callScenarioNative(t, vm, "ClearActorAnimations", actorID); cleared != 1 || actor.hasAnimation {
		t.Fatalf("ClearActorAnimations returned %d, animation=%t", cleared, actor.hasAnimation)
	}
}

func TestActorScenarioStreamingUsesPlayerWorld(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = "Alice"
	playerID := callScenarioNative(t, vm, "__pt_player_create", 100)
	actorID := callScenarioNative(t, vm, "CreateActor", 7, 0, 0, 0, 0)

	if streamed := callScenarioNative(t, vm, "IsActorStreamedIn", actorID, playerID); streamed != 1 {
		t.Fatalf("IsActorStreamedIn = %d in matching world", streamed)
	}

	callScenarioNative(t, vm, "SetActorVirtualWorld", actorID, 2)

	if streamed := callScenarioNative(t, vm, "IsActorStreamedIn", actorID, playerID); streamed != 0 {
		t.Fatalf("IsActorStreamedIn = %d in different world", streamed)
	}
}

func TestActorScenarioCloneIsolatesState(t *testing.T) {
	state := newActorState()
	state.actors[0] = &testActor{skin: 7, health: 100, position: [3]float32{1, 2, 3}}

	clone, ok := state.Clone().(*actorState)
	if !ok {
		t.Fatal("cloned scenario was not actor state")
	}

	clone.actors[0].skin = 10

	clone.actors[0].position[0] = 99
	if state.actors[0].skin != 7 || state.actors[0].position[0] != 1 {
		t.Fatal("actor clone shared mutable state")
	}
}

func actorScenarioState(t *testing.T, registry *scenarioRegistry) *actorState {
	t.Helper()

	return registry.Actors
}
