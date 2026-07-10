package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *textLabelState) createPlayerTextLabel(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 || !state.validPlayer(params[0]) {
		return -1, nil
	}

	text, err := ctx.ReadString(params[1])
	if err != nil {
		return -1, err
	}

	playerID := int(params[0])
	id := state.playerNext[playerID]

	state.playerNext[playerID]++
	if state.playerLabels[playerID] == nil {
		state.playerLabels[playerID] = map[int]*testTextLabel{}
	}

	label := newTestTextLabel(text, int(params[2]), cellsToVector(params[3:6]), cellFloat(params[6]), state.players.players[playerID].world)
	if len(params) > 7 {
		label.parentPlayer = int(params[7])
	}

	if len(params) > 8 {
		label.parentVehicle = int(params[8])
	}

	if len(params) > 9 {
		label.lineOfSight = params[9] != 0
	}

	state.playerLabels[playerID][id] = label

	return backend.Cell(id), nil
}

func (state *textLabelState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}

func (state *textLabelState) playerLabel(params []backend.Cell) (*testTextLabel, bool) {
	if len(params) < 2 {
		return nil, false
	}

	label, ok := state.playerLabels[int(params[0])][int(params[1])]

	return label, ok
}

func (state *textLabelState) deletePlayerTextLabel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerLabel(params); !ok {
		return 0, nil
	}

	delete(state.playerLabels[int(params[0])], int(params[1]))

	return 1, nil
}

func (state *textLabelState) isValidPlayerTextLabel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerLabel(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *textLabelState) updatePlayerTextLabel(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.playerLabel(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	text, err := ctx.ReadString(params[3])
	if err != nil {
		return 0, err
	}

	label.colour, label.text = int(params[2]), text

	return 1, nil
}

func (state *textLabelState) getPlayerTextLabelText(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.playerLabel(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return 1, ctx.WriteString(params[2], truncateString(label.text, int(params[3])))
}

func (state *textLabelState) setPlayerTextLabelFloat(field func(*testTextLabel) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.playerLabel(params)
		if !ok || len(params) < 3 {
			return 0, nil
		}

		*field(label) = cellFloat(params[2])

		return 1, nil
	}
}

func (state *textLabelState) getPlayerTextLabelFloat(field func(*testTextLabel) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.playerLabel(params)
		if !ok {
			return 0, nil
		}

		return floatCell(*field(label)), nil
	}
}

func (state *textLabelState) setPlayerTextLabelBool(field func(*testTextLabel) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.playerLabel(params)
		if !ok || len(params) < 3 {
			return 0, nil
		}

		*field(label) = params[2] != 0

		return 1, nil
	}
}

func (state *textLabelState) getPlayerTextLabelBool(field func(*testTextLabel) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.playerLabel(params)
		if ok && *field(label) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *textLabelState) getPlayerTextLabelInt(field func(*testTextLabel) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.playerLabel(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(label)), nil
	}
}

func (state *textLabelState) getPlayerTextLabelPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.playerLabel(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[2:5], label.position)
}

func (state *textLabelState) getPlayerTextLabelAttachedData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.playerLabel(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[2:4], []int{label.parentPlayer, label.parentVehicle})
}
