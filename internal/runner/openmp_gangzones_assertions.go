package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *gangZoneState) assertValid(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		minimum, fileIndex := 3, 1
		if player {
			minimum, fileIndex = 4, 2
		}

		if len(params) < minimum {
			return 0, errors.New("gang zone assertion has too few arguments")
		}

		var ok bool
		if player {
			_, ok = state.playerZone(params)
		} else {
			_, ok = state.zone(params)
		}

		if !ok {
			setFailure(result, params, fileIndex, "gang zone does not exist", ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *gangZoneState) assertVisible(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("gang zone visibility assertion expects 5 arguments")
		}

		var (
			zone *testGangZone
			ok   bool
		)

		if player {
			zone, ok = state.playerZone(params)
		} else {
			zone, ok = state.playerZoneArgument(params)
		}

		actual := ok && zone.views[int(params[0])].visible

		expected := params[2] != 0
		if actual != expected {
			setFailure(result, params, 3, fmt.Sprintf("gang zone visibility: expected %t, got %t", expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *gangZoneState) assertInside(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("gang zone containment assertion expects 5 arguments")
		}

		zone, ok := state.playerZoneArgument(params)
		actual := ok && state.playerInZone(int(params[0]), zone)

		expected := params[2] != 0
		if actual != expected {
			setFailure(result, params, 3, fmt.Sprintf("gang zone containment: expected %t, got %t", expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
