package runner

import (
	"errors"
	"fmt"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *actorState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("actor assertion expects 3 arguments")
		}

		if _, ok := state.actor(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("actor %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *actorState) assertSkin(result *nativeState) backend.NativeFunc {
	return state.assertActorInt(result, "skin", func(actor *testActor) int { return actor.skin })
}

func (state *actorState) assertActorInt(result *nativeState, label string, value func(*testActor) int) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("actor integer assertion expects 4 arguments")
		}

		actor, ok := state.actor(params)

		actual := -1
		if ok {
			actual = value(actor)
		}

		if !ok || actual != int(params[1]) {
			setFailure(result, params, 2, fmt.Sprintf("actor %d %s: expected %d, got %d", params[0], label, params[1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *actorState) assertHealth(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("actor health assertion expects 4 arguments")
		}

		actor, ok := state.actor(params)

		expected := cellFloat(params[1])
		if !ok || actor.health != expected {
			actual := float32(0)
			if ok {
				actual = actor.health
			}

			setFailure(result, params, 2, fmt.Sprintf("actor %d health: expected %g, got %g", params[0], expected, actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *actorState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("actor position assertion expects 7 arguments")
		}

		actor, ok := state.actor(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("actor %d does not exist", params[0]), ctx)
			return 0, nil
		}

		expected := cellsToVector(params[1:4])

		tolerance := float32(math.Abs(float64(cellFloat(params[4]))))
		for index := range expected {
			if absFloat(actor.position[index]-expected[index]) > tolerance {
				setFailure(result, params, 5, fmt.Sprintf("actor %d position: expected %v +/- %g, got %v", params[0], expected, tolerance, actor.position), ctx)
				return 0, nil
			}
		}

		return 1, nil
	}
}
