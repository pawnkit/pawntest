package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *objectState) attachPlayerObject(kind int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		object, ok := state.playerObject(params)
		if !ok || len(params) < 9 {
			return 0, nil
		}

		object.attachment = objectAttachment{kind: kind, id: int(params[2]), offset: cellsToVector(params[3:6]), rotation: cellsToVector(params[6:9]), sync: len(params) < 10 || params[9] != 0}

		return 1, nil
	}
}

func (state *objectState) getPlayerObjectAttachedData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	values := []int{0, 0, 0}
	if object.attachment.kind > 0 {
		values[object.attachment.kind-1] = object.attachment.id
	}

	return writeVehicleInts(ctx, params[2:5], values)
}

func (state *objectState) getPlayerObjectAttachedOffset(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 8 {
		return 0, nil
	}

	if result, err := writeFloatVector(ctx, params[2:5], object.attachment.offset); err != nil || result == 0 {
		return result, err
	}

	return writeFloatVector(ctx, params[5:8], object.attachment.rotation)
}

func (state *objectState) getPlayerObjectSyncRotation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if ok && object.attachment.sync {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) setPlayerObjectMaterial(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 7 {
		return 0, nil
	}

	library, err := ctx.ReadString(params[4])
	if err != nil {
		return 0, err
	}

	name, err := ctx.ReadString(params[5])
	if err != nil {
		return 0, err
	}

	object.materials[int(params[2])] = objectMaterial{model: int(params[3]), library: library, name: name, colour: int(params[6])}

	return 1, nil
}

func (state *objectState) setPlayerObjectMaterialText(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	text, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	object.materials[int(params[3])] = objectMaterial{text: text, textMaterial: true}

	return 1, nil
}

func (state *objectState) isPlayerObjectMaterialSlotUsed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	if _, used := object.materials[int(params[2])]; used {
		return 1, nil
	}

	return 0, nil
}
