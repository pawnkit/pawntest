package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *classState) assertPlayerClass(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player class assertion expects 4 arguments")
		}

		actual, ok := state.selected[int(params[0])]
		if !ok || actual != int(params[1]) {
			setFailure(result, params, 2, fmt.Sprintf("player %d class: expected %d, got %d", params[0], params[1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *classState) assertSelecting(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("class selection assertion expects 4 arguments")
		}

		actual := state.selecting[int(params[0])]

		expected := params[1] != 0
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("player %d selecting class: expected %t, got %t", params[0], expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
