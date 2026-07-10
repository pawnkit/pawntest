package runner

import (
	"errors"
	"fmt"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *pickupState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("pickup assertion expects 3 arguments")
		}

		if _, ok := state.pickup(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("pickup %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *pickupState) assertModel(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("pickup model assertion expects 4 arguments")
		}

		pickup, ok := state.pickup(params)

		actual := -1
		if ok {
			actual = pickup.model
		}

		if !ok || actual != int(params[1]) {
			setFailure(result, params, 2, fmt.Sprintf("pickup %d model: expected %d, got %d", params[0], params[1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *pickupState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("pickup position assertion expects 7 arguments")
		}

		pickup, ok := state.pickup(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("pickup %d does not exist", params[0]), ctx)
			return 0, nil
		}

		expected := cellsToVector(params[1:4])

		tolerance := float32(math.Abs(float64(cellFloat(params[4]))))
		for index := range expected {
			if absFloat(pickup.position[index]-expected[index]) > tolerance {
				setFailure(result, params, 5, fmt.Sprintf("pickup %d position: expected %v +/- %g, got %v", params[0], expected, tolerance, pickup.position), ctx)
				return 0, nil
			}
		}

		return 1, nil
	}
}

func (state *pickupState) assertPlayerPickupValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player pickup assertion expects 4 arguments")
		}

		if _, ok := state.playerPickup(params); !ok {
			setFailure(result, params, 2, fmt.Sprintf("player %d pickup %d does not exist", params[0], params[1]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
