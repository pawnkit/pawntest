package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *dialogState) assertVisible(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("dialog visibility assertion expects 5 arguments")
		}

		dialog, ok := state.dialogs[int(params[0])]
		expectedVisible := params[2] != 0
		actualVisible := ok && dialog.visible

		actualID := -1
		if actualVisible {
			actualID = dialog.id
		}

		if actualVisible != expectedVisible || expectedVisible && actualID != int(params[1]) {
			setFailure(result, params, 3, fmt.Sprintf("player %d dialog: expected visible=%t id=%d, got visible=%t id=%d", params[0], expectedVisible, params[1], actualVisible, actualID), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *dialogState) assertString(result *nativeState, label string, value func(testDialog) string) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("dialog string assertion expects 4 arguments")
		}

		dialog, ok := state.dialogs[int(params[0])]
		expected := readStringParam(ctx, params, 1)

		actual := ""
		if ok && dialog.visible {
			actual = value(dialog)
		}

		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("player %d dialog %s: expected %q, got %q", params[0], label, expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
