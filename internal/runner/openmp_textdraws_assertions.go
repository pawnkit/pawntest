package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *textDrawState) assertValid(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		minimum := 3

		fileIndex := 1
		if player {
			minimum, fileIndex = 4, 2
		}

		if len(params) < minimum {
			return 0, errors.New("textdraw assertion has too few arguments")
		}

		if _, ok := state.draw(params, player); !ok {
			setFailure(result, params, fileIndex, "textdraw does not exist", ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *textDrawState) assertText(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		expectedIndex, fileIndex, minimum := 1, 2, 4
		if player {
			expectedIndex, fileIndex, minimum = 2, 3, 5
		}

		if len(params) < minimum {
			return 0, errors.New("textdraw text assertion has too few arguments")
		}

		draw, ok := state.draw(params, player)

		expected := readStringParam(ctx, params, expectedIndex)
		if !ok || draw.text != expected {
			actual := ""
			if ok {
				actual = draw.text
			}

			setFailure(result, params, fileIndex, fmt.Sprintf("textdraw: expected %q, got %q", expected, actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *textDrawState) assertVisible(result *nativeState, player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		visibleIndex, fileIndex, minimum := 2, 3, 5
		if len(params) < minimum {
			return 0, errors.New("textdraw visibility assertion has too few arguments")
		}

		drawParams := params[1:]
		if player {
			drawParams = params[:2]
		}

		draw, ok := state.draw(drawParams, player)
		actual := ok && draw.visible[int(params[0])]

		expected := params[visibleIndex] != 0
		if actual != expected {
			setFailure(result, params, fileIndex, fmt.Sprintf("textdraw visibility: expected %t, got %t", expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
