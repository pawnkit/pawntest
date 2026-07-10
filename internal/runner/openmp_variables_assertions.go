package runner

import (
	"errors"
	"fmt"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *variableState) assertInt(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		offset, minimum := 0, 4
		if player {
			offset, minimum = 1, 5
		}

		if len(params) < minimum {
			return 0, errors.New("variable integer assertion has too few arguments")
		}

		value, ok, err := state.variable(ctx, params, player)
		if err != nil {
			return 0, err
		}

		actual := 0
		if ok && value.variableType == variableInt {
			actual = value.integer
		}

		if !ok || value.variableType != variableInt || actual != int(params[offset+1]) {
			setFailure(result, params, offset+2, fmt.Sprintf("integer variable: expected %d, got %d", params[offset+1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *variableState) assertFloat(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		offset, minimum := 0, 5
		if player {
			offset, minimum = 1, 6
		}

		if len(params) < minimum {
			return 0, errors.New("variable float assertion has too few arguments")
		}

		value, ok, err := state.variable(ctx, params, player)
		if err != nil {
			return 0, err
		}

		expected := cellFloat(params[offset+1])
		tolerance := float32(math.Abs(float64(cellFloat(params[offset+2]))))

		actual := float32(0)
		if ok && value.variableType == variableFloat {
			actual = value.floating
		}

		if !ok || value.variableType != variableFloat || absFloat(actual-expected) > tolerance {
			setFailure(result, params, offset+3, fmt.Sprintf("float variable: expected %g +/- %g, got %g", expected, tolerance, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *variableState) assertString(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		offset, minimum := 0, 4
		if player {
			offset, minimum = 1, 5
		}

		if len(params) < minimum {
			return 0, errors.New("variable string assertion has too few arguments")
		}

		value, ok, err := state.variable(ctx, params, player)
		if err != nil {
			return 0, err
		}

		expected := readStringParam(ctx, params, offset+1)

		actual := ""
		if ok && value.variableType == variableString {
			actual = value.text
		}

		if !ok || value.variableType != variableString || actual != expected {
			setFailure(result, params, offset+2, fmt.Sprintf("string variable: expected %q, got %q", expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
