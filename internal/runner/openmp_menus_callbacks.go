package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *menuState) selectMenuRow(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, errors.New("menu selection expects 2 arguments")
	}

	menuID, menu := state.visibleMenu(int(params[0]))
	if menu == nil || menu.disabledRows[int(params[1])] {
		return 0, nil
	}

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support menu callbacks")
	}

	delete(menu.visible, int(params[0]))

	result, err := caller.CallPublic("OnPlayerSelectedMenuRow", params[0], params[1])
	if err != nil {
		return 0, fmt.Errorf("menu %d selection callback: %w", menuID, err)
	}

	return result, nil
}

func (state *menuState) exitMenu(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, errors.New("menu exit expects a player")
	}

	menuID, menu := state.visibleMenu(int(params[0]))
	if menu == nil {
		return 0, nil
	}

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support menu callbacks")
	}

	delete(menu.visible, int(params[0]))

	result, err := caller.CallPublic("OnPlayerExitedMenu", params[0])
	if err != nil {
		return 0, fmt.Errorf("menu %d exit callback: %w", menuID, err)
	}

	return result, nil
}

func (state *menuState) visibleMenu(playerID int) (int, *testMenu) {
	for id, menu := range state.menus {
		if menu.visible[playerID] {
			return id, menu
		}
	}

	return -1, nil
}
