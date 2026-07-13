package runner

import (
	"math"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestPlayerModelHandlesZeroPlayersWithoutMockConfiguration(t *testing.T) {
	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}

	state := newOpenMPState()
	if err := registerOpenMPPlayerNatives(vm, &nativeState{status: Pass}, state, newMockState(), false); err != nil {
		t.Fatal(err)
	}

	connected, err := vm.natives["IsPlayerConnected"](vm, []backend.Cell{0})
	if err != nil {
		t.Fatal(err)
	}

	if connected != 0 {
		t.Fatalf("IsPlayerConnected(0) = %d, want 0", connected)
	}
}

func TestPlayerModelStoresCoreState(t *testing.T) {
	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{100: "Alice"}}

	state := newOpenMPState()
	if err := registerOpenMPPlayerNatives(vm, &nativeState{status: Pass}, state, newMockState(), false); err != nil {
		t.Fatal(err)
	}

	playerID, err := vm.natives["__pt_player_create"](vm, []backend.Cell{100})
	if err != nil {
		t.Fatal(err)
	}

	calls := []struct {
		name   string
		params []backend.Cell
	}{
		{name: "SetPlayerHealth", params: []backend.Cell{playerID, floatCell(75.5)}},
		{name: "SetPlayerSkin", params: []backend.Cell{playerID, 42}},
		{name: "SetPlayerInterior", params: []backend.Cell{playerID, 3}},
		{name: "SetPlayerVirtualWorld", params: []backend.Cell{playerID, 7}},
		{name: "SetPlayerFacingAngle", params: []backend.Cell{playerID, floatCell(90)}},
		{name: "SetPlayerVelocity", params: []backend.Cell{playerID, floatCell(1), floatCell(2), floatCell(3)}},
		{name: "GivePlayerWeapon", params: []backend.Cell{playerID, 24, 50}},
	}
	for _, call := range calls {
		result, callErr := vm.natives[call.name](vm, call.params)
		if callErr != nil {
			t.Fatalf("%s: %v", call.name, callErr)
		}

		if result != 1 {
			t.Fatalf("%s returned %d", call.name, result)
		}
	}

	player := state.players[int(playerID)]
	if player.health != 75.5 || player.skin != 42 || player.interior != 3 || player.world != 7 {
		t.Fatalf("unexpected player state: %#v", player)
	}

	if player.angle != 90 || player.velocity != [3]float32{1, 2, 3} {
		t.Fatalf("unexpected player movement: angle=%g velocity=%v", player.angle, player.velocity)
	}

	if player.armedWeapon != 24 || player.weapons[24] != 50 {
		t.Fatalf("unexpected weapons: armed=%d weapons=%v", player.armedWeapon, player.weapons)
	}

	health, err := vm.natives["GetPlayerGravity"](vm, []backend.Cell{playerID})
	if err != nil {
		t.Fatal(err)
	}

	if math.Float32frombits(uint32(health)) != 0.008 {
		t.Fatalf("gravity = %g", math.Float32frombits(uint32(health)))
	}
}

func TestPlayerAssertionRejectsMalformedArguments(t *testing.T) {
	vm := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}
	if err := registerOpenMPPlayerNatives(vm, &nativeState{status: Pass}, newOpenMPState(), newMockState(), false); err != nil {
		t.Fatal(err)
	}

	if _, err := vm.natives["__pt_player_money"](vm, nil); err == nil {
		t.Fatal("expected malformed player assertion error")
	}
}

func TestScenarioRegistryCloneIsolatesPlayerState(t *testing.T) {
	registry := newScenarioRegistry()

	original := registry.Players
	original.players[0] = &testPlayer{name: "Alice", connected: true, messages: []string{"hello"}}

	clone := registry.Clone()

	clonedState := clone.Players
	clonedPlayer := clonedState.players[0]
	clonedPlayer.name = "Bob"
	clonedPlayer.messages[0] = "changed"

	if original.players[0].name != "Alice" || original.players[0].messages[0] != "hello" {
		t.Fatal("scenario clone shared mutable player state")
	}
}
