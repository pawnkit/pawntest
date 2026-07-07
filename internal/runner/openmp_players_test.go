package runner

import (
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
	original := registry.modules[0].(*openMPState)
	original.players[0] = &testPlayer{name: "Alice", connected: true, messages: []string{"hello"}}

	clone := registry.Clone()
	clonedPlayer := clone.modules[0].(*openMPState).players[0]
	clonedPlayer.name = "Bob"
	clonedPlayer.messages[0] = "changed"

	if original.players[0].name != "Alice" || original.players[0].messages[0] != "hello" {
		t.Fatal("scenario clone shared mutable player state")
	}
}
