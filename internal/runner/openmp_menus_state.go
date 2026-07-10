package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *menuState) addMenuItem(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 3 || params[1] < 0 || params[1] >= backend.Cell(menu.columns) || params[1] >= 2 {
		return -1, nil
	}

	text, err := ctx.ReadString(params[2])
	if err != nil {
		return -1, err
	}

	column := int(params[1])
	row := len(menu.items[column])
	menu.items[column] = append(menu.items[column], text)

	return backend.Cell(row), nil
}

func (state *menuState) setMenuColumnHeader(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 3 || params[1] < 0 || params[1] >= backend.Cell(menu.columns) || params[1] >= 2 {
		return 0, nil
	}

	header, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	menu.headers[int(params[1])] = header

	return 1, nil
}

func (state *menuState) showMenuForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 2 || menu.disabled || !state.validPlayer(params[1]) {
		return 0, nil
	}

	for _, other := range state.menus {
		delete(other.visible, int(params[1]))
	}

	menu.visible[int(params[1])] = true

	return 1, nil
}

func (state *menuState) hideMenuForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	delete(menu.visible, int(params[1]))

	return 1, nil
}

func (state *menuState) disableMenu(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok {
		return 0, nil
	}

	menu.disabled = true
	clear(menu.visible)

	return 1, nil
}

func (state *menuState) disableMenuRow(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	menu.disabledRows[int(params[1])] = true

	return 1, nil
}

func (state *menuState) getPlayerMenu(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return -1, nil
	}

	for id, menu := range state.menus {
		if menu.visible[int(params[0])] {
			return backend.Cell(id), nil
		}
	}

	return -1, nil
}

func (state *menuState) isMenuDisabled(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if ok && menu.disabled {
		return 1, nil
	}

	return 0, nil
}

func (state *menuState) isMenuRowDisabled(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if ok && len(params) >= 2 && menu.disabledRows[int(params[1])] {
		return 1, nil
	}

	return 0, nil
}

func (state *menuState) getMenuColumns(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(menu.columns), nil
}

func (state *menuState) getMenuItems(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 2 || params[1] < 0 || params[1] >= 2 {
		return 0, nil
	}

	return backend.Cell(len(menu.items[int(params[1])])), nil
}

func (state *menuState) getMenuPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return writeMenuPair(ctx, params[1:3], menu.position)
}

func (state *menuState) getMenuWidths(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return writeMenuPair(ctx, params[1:3], menu.widths)
}

func writeMenuPair(ctx backend.NativeContext, addresses []backend.Cell, values [2]float32) (backend.Cell, error) {
	for index, value := range values {
		if err := ctx.WriteCell(addresses[index], floatCell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *menuState) getMenuColumnHeader(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 4 || params[1] < 0 || params[1] >= 2 {
		return 0, nil
	}

	return 1, ctx.WriteString(params[2], truncateString(menu.headers[int(params[1])], int(params[3])))
}

func (state *menuState) getMenuItem(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	menu, ok := state.menu(params)
	if !ok || len(params) < 5 || params[1] < 0 || params[1] >= 2 {
		return 0, nil
	}

	items := menu.items[int(params[1])]
	if params[2] < 0 || int(params[2]) >= len(items) {
		return 0, nil
	}

	return 1, ctx.WriteString(params[3], truncateString(items[int(params[2])], int(params[4])))
}

func (state *menuState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}
