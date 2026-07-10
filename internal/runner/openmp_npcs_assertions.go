package runner

import (
	"errors"
	"fmt"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *npcState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("NPC assertion expects 3 arguments")
		}

		if _, ok := state.npc(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("NPC %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *npcState) assertSpawned(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("NPC spawn assertion expects 4 arguments")
		}

		npc, ok := state.npc(params)
		actual := ok && npc.spawned

		expected := params[1] != 0
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("NPC %d spawned: expected %t, got %t", params[0], expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *npcState) assertHealth(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("NPC health assertion expects 5 arguments")
		}

		npc, ok := state.npc(params)
		expected := cellFloat(params[1])
		tolerance := float32(math.Abs(float64(cellFloat(params[2]))))

		actual := float32(0)
		if ok {
			actual = npc.health
		}

		if !ok || absFloat(actual-expected) > tolerance {
			setFailure(result, params, 3, fmt.Sprintf("NPC %d health: expected %g +/- %g, got %g", params[0], expected, tolerance, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *npcState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("NPC position assertion expects 7 arguments")
		}

		npc, ok := state.npc(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("NPC %d does not exist", params[0]), ctx)
			return 0, nil
		}

		expected := cellsToVector(params[1:4])

		tolerance := float32(math.Abs(float64(cellFloat(params[4]))))
		for index := range expected {
			if absFloat(npc.position[index]-expected[index]) > tolerance {
				setFailure(result, params, 5, fmt.Sprintf("NPC %d position: expected %v +/- %g, got %v", params[0], expected, tolerance, npc.position), ctx)
				return 0, nil
			}
		}

		return 1, nil
	}
}
