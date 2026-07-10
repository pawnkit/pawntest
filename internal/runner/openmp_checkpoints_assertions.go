package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *checkpointState) assertCheckpointActive(result *nativeState) backend.NativeFunc {
	return state.assertCheckpointBool(result, "checkpoint active", func(playerID int) bool {
		return state.checkpoints[playerID].active
	})
}

func (state *checkpointState) assertCheckpointInside(result *nativeState) backend.NativeFunc {
	return state.assertCheckpointBool(result, "player in checkpoint", func(playerID int) bool {
		return state.playerInside([]backend.Cell{backend.Cell(playerID)}, state.checkpoints)
	})
}

func (state *checkpointState) assertRaceCheckpointActive(result *nativeState) backend.NativeFunc {
	return state.assertCheckpointBool(result, "race checkpoint active", func(playerID int) bool {
		return state.races[playerID].active
	})
}

func (state *checkpointState) assertRaceCheckpointInside(result *nativeState) backend.NativeFunc {
	return state.assertCheckpointBool(result, "player in race checkpoint", func(playerID int) bool {
		value, ok := state.races[playerID]
		if !ok {
			return false
		}

		return state.playerInside([]backend.Cell{backend.Cell(playerID)}, map[int]checkpoint{playerID: value.checkpoint})
	})
}

func (state *checkpointState) assertCheckpointBool(result *nativeState, label string, value func(int) bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("checkpoint assertion expects 4 arguments")
		}

		actual := value(int(params[0]))

		expected := params[1] != 0
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("player %d %s: expected %t, got %t", params[0], label, expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
