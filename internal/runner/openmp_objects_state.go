package runner

import (
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *objectState) setObjectPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return mutateObjectVector(state.object, params, func(object *testObject) *[3]float32 { return &object.position })
}

func (state *objectState) getObjectPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return readObjectVector(ctx, state.object, params, func(object *testObject) *[3]float32 { return &object.position })
}

func (state *objectState) setObjectRotation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return mutateObjectVector(state.object, params, func(object *testObject) *[3]float32 { return &object.rotation })
}

func (state *objectState) getObjectRotation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return readObjectVector(ctx, state.object, params, func(object *testObject) *[3]float32 { return &object.rotation })
}

func mutateObjectVector(find func([]backend.Cell) (*testObject, bool), params []backend.Cell, field func(*testObject) *[3]float32) (backend.Cell, error) {
	object, ok := find(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	*field(object) = cellsToVector(params[1:4])

	return 1, nil
}

func readObjectVector(ctx backend.NativeContext, find func([]backend.Cell) (*testObject, bool), params []backend.Cell, field func(*testObject) *[3]float32) (backend.Cell, error) {
	object, ok := find(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[1:4], *field(object))
}

func (state *objectState) moveObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	object.targetPos = cellsToVector(params[1:4])
	object.moveSpeed = cellFloat(params[4])

	object.moving = true

	object.moveID++
	if len(params) >= 8 {
		object.targetRot = cellsToVector(params[5:8])
	}

	duration := objectMoveDuration(object)
	state.scheduleObjectMove(int(params[0]), object, duration)

	return duration, nil
}

func (state *objectState) stopObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok {
		return 0, nil
	}

	object.moving = false
	object.moveID++

	return 1, nil
}

func (state *objectState) scheduleObjectMove(id int, object *testObject, duration backend.Cell) {
	if state.scheduler == nil || duration <= 0 {
		return
	}

	moveID := object.moveID

	state.scheduler.schedule(int64(duration), "OnObjectMoved", []backend.Cell{backend.Cell(id)}, func() bool {
		current := state.objects[id]
		if current == nil || !current.moving || current.moveID != moveID {
			return false
		}

		current.position, current.rotation, current.moving = current.targetPos, current.targetRot, false

		return true
	})
}

func (state *objectState) isObjectMoving(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if ok && object.moving {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) setObjectMoveSpeed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	object.moveSpeed = cellFloat(params[1])

	return 1, nil
}

func (state *objectState) getObjectMoveSpeed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok {
		return 0, nil
	}

	return floatCell(object.moveSpeed), nil
}

func (state *objectState) getObjectTargetPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return readObjectVector(ctx, state.object, params, func(object *testObject) *[3]float32 { return &object.targetPos })
}

func (state *objectState) getObjectTargetRotation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return readObjectVector(ctx, state.object, params, func(object *testObject) *[3]float32 { return &object.targetRot })
}

func (state *objectState) getObjectDrawDistance(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok {
		return 0, nil
	}

	return floatCell(object.drawDistance), nil
}

func (state *objectState) setObjectNoCameraCollision(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok {
		return 0, nil
	}

	object.noCamera = true

	return 1, nil
}

func (state *objectState) isObjectNoCameraCollision(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if ok && object.noCamera {
		return 1, nil
	}

	return 0, nil
}

func objectDistance(from, to [3]float32) float32 {
	dx, dy, dz := from[0]-to[0], from[1]-to[1], from[2]-to[2]

	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func objectMoveDuration(object *testObject) backend.Cell {
	if object.moveSpeed <= 0 {
		return 0
	}

	return backend.Cell(math.Ceil(float64(objectDistance(object.position, object.targetPos) / object.moveSpeed * 1000)))
}
