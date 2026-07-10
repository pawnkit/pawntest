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

func (state *npcState) assertWeapon(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("NPC weapon assertion expects 4 arguments")
		}

		npc, ok := state.npc(params)

		actual := 0
		if ok {
			actual = npc.weapon
		}

		if !ok || actual != int(params[1]) {
			setFailure(result, params, 2, fmt.Sprintf("NPC %d weapon: expected %d, got %d", params[0], params[1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *npcState) assertAiming(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("NPC aiming assertion expects 4 arguments")
		}

		npc, ok := state.npc(params)
		actual := ok && npc.aiming

		expected := params[1] != 0
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("NPC %d aiming: expected %t, got %t", params[0], expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *npcState) assertAnimation(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("NPC animation assertion expects 5 arguments")
		}

		npc, ok := state.npc(params)
		expectedLibrary := readStringParam(ctx, params, 1)

		expectedName := readStringParam(ctx, params, 2)
		if !ok || !npc.hasAnimation || npc.animation.library != expectedLibrary || npc.animation.name != expectedName {
			setFailure(result, params, 3, fmt.Sprintf("NPC %d animation: expected %s/%s, got %s/%s", params[0], expectedLibrary, expectedName, npc.animation.library, npc.animation.name), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
