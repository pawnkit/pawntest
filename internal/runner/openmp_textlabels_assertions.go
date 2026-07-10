package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *textLabelState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("text label assertion expects 3 arguments")
		}

		if _, ok := state.label(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("text label %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *textLabelState) assertText(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("text label text assertion expects 4 arguments")
		}

		label, ok := state.label(params)

		expected := readStringParam(ctx, params, 1)
		if !ok || label.text != expected {
			actual := ""
			if ok {
				actual = label.text
			}

			setFailure(result, params, 2, fmt.Sprintf("text label %d: expected %q, got %q", params[0], expected, actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *textLabelState) assertPlayerValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player text label assertion expects 4 arguments")
		}

		if _, ok := state.playerLabel(params); !ok {
			setFailure(result, params, 2, fmt.Sprintf("player %d text label %d does not exist", params[0], params[1]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *textLabelState) assertPlayerText(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("player text label text assertion expects 5 arguments")
		}

		label, ok := state.playerLabel(params)

		expected := readStringParam(ctx, params, 2)
		if !ok || label.text != expected {
			actual := ""
			if ok {
				actual = label.text
			}

			setFailure(result, params, 3, fmt.Sprintf("player %d text label %d: expected %q, got %q", params[0], params[1], expected, actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}
