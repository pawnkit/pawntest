package runner

import (
	"errors"
	"fmt"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *serverState) assertInt(result *nativeState, label string, value func() int) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("server integer assertion expects 3 arguments")
		}

		actual := value()
		if actual != int(params[0]) {
			setFailure(result, params, 1, fmt.Sprintf("server %s: expected %d, got %d", label, params[0], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *serverState) assertGravity(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("server gravity assertion expects 4 arguments")
		}

		expected := cellFloat(params[0])

		tolerance := float32(math.Abs(float64(cellFloat(params[1]))))
		if absFloat(state.gravity-expected) > tolerance {
			setFailure(result, params, 2, fmt.Sprintf("server gravity: expected %g +/- %g, got %g", expected, tolerance, state.gravity), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *serverState) assertModeText(result *nativeState) backend.NativeFunc {
	return state.assertString(result, "game mode text", func() string { return state.gameModeText })
}

func (state *serverState) assertRule(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("server rule assertion expects 4 arguments")
		}

		name := readStringParam(ctx, params, 0)
		expected := readStringParam(ctx, params, 1)

		actual := state.rules[name]
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("server rule %q: expected %q, got %q", name, expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *serverState) assertString(result *nativeState, label string, value func() string) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("server string assertion expects 3 arguments")
		}

		expected := readStringParam(ctx, params, 0)

		actual := value()
		if actual != expected {
			setFailure(result, params, 1, fmt.Sprintf("server %s: expected %q, got %q", label, expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
