package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *objectState) attachObject(kind int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		object, ok := state.object(params)
		if !ok || len(params) < 8 {
			return 0, nil
		}

		object.attachment = objectAttachment{kind: kind, id: int(params[1]), offset: cellsToVector(params[2:5]), rotation: cellsToVector(params[5:8]), sync: len(params) < 9 || params[8] != 0}

		return 1, nil
	}
}

func (state *objectState) getObjectAttachedData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	values := []int{0, 0, 0}
	if object.attachment.kind > 0 {
		values[object.attachment.kind-1] = object.attachment.id
	}

	return writeVehicleInts(ctx, params[1:4], values)
}

func (state *objectState) getObjectAttachedOffset(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 7 {
		return 0, nil
	}

	if result, err := writeFloatVector(ctx, params[1:4], object.attachment.offset); err != nil || result == 0 {
		return result, err
	}

	return writeFloatVector(ctx, params[4:7], object.attachment.rotation)
}

func (state *objectState) getObjectSyncRotation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if ok && object.attachment.sync {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) setObjectMaterial(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 6 {
		return 0, nil
	}

	library, err := ctx.ReadString(params[3])
	if err != nil {
		return 0, err
	}

	name, err := ctx.ReadString(params[4])
	if err != nil {
		return 0, err
	}

	object.materials[int(params[1])] = objectMaterial{model: int(params[2]), library: library, name: name, colour: int(params[5])}

	return 1, nil
}

func (state *objectState) setObjectMaterialText(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	text, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	object.materials[int(params[2])] = objectMaterial{text: text, textMaterial: true}

	return 1, nil
}

func (state *objectState) isObjectMaterialSlotUsed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	if _, used := object.materials[int(params[1])]; used {
		return 1, nil
	}

	return 0, nil
}
