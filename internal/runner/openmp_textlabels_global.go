package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *textLabelState) updateTextLabel(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	text, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	label.colour, label.text = int(params[1]), text

	return 1, nil
}

func (state *textLabelState) getTextLabelText(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return 1, ctx.WriteString(params[1], truncateString(label.text, int(params[2])))
}

func (state *textLabelState) setTextLabelInt(field func(*testTextLabel) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(label) = int(params[1])

		return 1, nil
	}
}

func (state *textLabelState) getTextLabelInt(field func(*testTextLabel) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(label)), nil
	}
}

func (state *textLabelState) setTextLabelFloat(field func(*testTextLabel) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(label) = cellFloat(params[1])

		return 1, nil
	}
}

func (state *textLabelState) getTextLabelFloat(field func(*testTextLabel) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if !ok {
			return 0, nil
		}

		return floatCell(*field(label)), nil
	}
}

func (state *textLabelState) setTextLabelBool(field func(*testTextLabel) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(label) = params[1] != 0

		return 1, nil
	}
}

func (state *textLabelState) getTextLabelBool(field func(*testTextLabel) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		label, ok := state.label(params)
		if ok && *field(label) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *textLabelState) getTextLabelPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[1:4], label.position)
}

func (state *textLabelState) attachTextLabelToPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 5 || state.players == nil {
		return 0, nil
	}

	if _, playerOK := state.players.player(params[1]); !playerOK {
		return 0, nil
	}

	label.parentPlayer, label.parentVehicle = int(params[1]), invalidScenarioID
	label.position = cellsToVector(params[2:5])

	return 1, nil
}

func (state *textLabelState) attachTextLabelToVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 5 || state.vehicles == nil {
		return 0, nil
	}

	if _, vehicleOK := state.vehicles.vehicles[int(params[1])]; !vehicleOK {
		return 0, nil
	}

	label.parentPlayer, label.parentVehicle = invalidScenarioID, int(params[1])
	label.position = cellsToVector(params[2:5])

	return 1, nil
}

func (state *textLabelState) getTextLabelAttachedData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	label, ok := state.label(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[1:3], []int{label.parentPlayer, label.parentVehicle})
}

func (state *textLabelState) isTextLabelStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[0])

	label, labelOK := state.labels[int(params[1])]
	if !playerOK || !labelOK || !player.connected || player.world != label.world {
		return 0, nil
	}

	position := state.textLabelPosition(label)
	if playerDistance(player, floatVectorCells(position)) <= label.drawDistance {
		return 1, nil
	}

	return 0, nil
}

func (state *textLabelState) textLabelPosition(label *testTextLabel) [3]float32 {
	if label.parentPlayer >= 0 && state.players != nil {
		if player, ok := state.players.players[label.parentPlayer]; ok {
			return addVectors([3]float32{player.x, player.y, player.z}, label.position)
		}
	}

	if label.parentVehicle >= 0 && state.vehicles != nil {
		if vehicle, ok := state.vehicles.vehicles[label.parentVehicle]; ok {
			return addVectors(vehicle.position, label.position)
		}
	}

	return label.position
}

func addVectors(left, right [3]float32) [3]float32 {
	return [3]float32{left[0] + right[0], left[1] + right[1], left[2] + right[2]}
}
