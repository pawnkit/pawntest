package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *textDrawState) showGlobalTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	draw, ok := state.draw(params[1:], false)
	if !ok {
		return 0, nil
	}

	draw.visible[int(params[0])] = true

	return 1, nil
}

func (state *textDrawState) hideGlobalTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	draw, ok := state.draw(params[1:], false)
	if !ok {
		return 0, nil
	}

	delete(draw.visible, int(params[0]))

	return 1, nil
}

func (state *textDrawState) showGlobalTextDrawForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	draw, ok := state.draw(params, false)
	if !ok || state.players == nil {
		return 0, nil
	}

	for id, player := range state.players.players {
		if player.connected {
			draw.visible[id] = true
		}
	}

	return 1, nil
}

func (state *textDrawState) hideGlobalTextDrawForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	draw, ok := state.draw(params, false)
	if !ok {
		return 0, nil
	}

	clear(draw.visible)

	return 1, nil
}

func (state *textDrawState) isGlobalTextDrawVisible(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	draw, ok := state.draw(params[1:], false)
	if ok && draw.visible[int(params[0])] {
		return 1, nil
	}

	return 0, nil
}

func (state *textDrawState) setGlobalTextDrawStringForPlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 3 || !state.validPlayer(params[1]) {
		return 0, nil
	}

	draw, ok := state.draw(params, false)
	if !ok {
		return 0, nil
	}

	text, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	draw.playerText[int(params[1])] = text

	return 1, nil
}

func (state *textDrawState) showPlayerTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	draw, ok := state.draw(params, true)
	if !ok {
		return 0, nil
	}

	draw.visible[int(params[0])] = true

	return 1, nil
}

func (state *textDrawState) hidePlayerTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	draw, ok := state.draw(params, true)
	if !ok {
		return 0, nil
	}

	delete(draw.visible, int(params[0]))

	return 1, nil
}

func (state *textDrawState) isPlayerTextDrawVisible(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	draw, ok := state.draw(params, true)
	if ok && draw.visible[int(params[0])] {
		return 1, nil
	}

	return 0, nil
}

func (state *textDrawState) selectTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	state.selection[int(params[0])] = textDrawSelection{active: true, colour: int(params[1])}

	return 1, nil
}

func (state *textDrawState) cancelSelectTextDraw(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	delete(state.selection, int(params[0]))

	return 1, nil
}
