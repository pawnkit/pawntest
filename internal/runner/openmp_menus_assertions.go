package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *menuState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("menu assertion expects 3 arguments")
		}

		if _, ok := state.menu(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("menu %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *menuState) assertVisible(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("menu visibility assertion expects 5 arguments")
		}

		menu, ok := state.menus[int(params[1])]
		actual := ok && menu.visible[int(params[0])]

		expected := params[2] != 0
		if actual != expected {
			setFailure(result, params, 3, fmt.Sprintf("menu visibility: expected %t, got %t", expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *menuState) assertItems(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("menu item assertion expects 5 arguments")
		}

		menu, ok := state.menu(params)

		actual := 0
		if ok && params[1] >= 0 && params[1] < 2 {
			actual = len(menu.items[int(params[1])])
		}

		if !ok || actual != int(params[2]) {
			setFailure(result, params, 3, fmt.Sprintf("menu %d column %d items: expected %d, got %d", params[0], params[1], params[2], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
