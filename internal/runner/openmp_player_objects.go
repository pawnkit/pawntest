package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *objectState) createTestPlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 || !state.validPlayer(params[0]) {
		return -1, nil
	}

	return state.addPlayerObject(int(params[0]), int(params[1]), cellsToVector(params[2:5]), [3]float32{}, 0), nil
}

func (state *objectState) createPlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 8 || !state.validPlayer(params[0]) {
		return -1, nil
	}

	drawDistance := float32(0)
	if len(params) > 8 {
		drawDistance = cellFloat(params[8])
	}

	return state.addPlayerObject(int(params[0]), int(params[1]), cellsToVector(params[2:5]), cellsToVector(params[5:8]), drawDistance), nil
}

func (state *objectState) addPlayerObject(playerID, model int, position, rotation [3]float32, drawDistance float32) backend.Cell {
	next := state.playerNext[playerID]
	if next == 0 {
		next = 1
	}

	state.playerNext[playerID] = next + 1
	if state.playerObjects[playerID] == nil {
		state.playerObjects[playerID] = map[int]*testObject{}
	}

	state.playerObjects[playerID][next] = newTestObject(model, position, rotation, drawDistance)

	return backend.Cell(next)
}

func (state *objectState) playerObject(params []backend.Cell) (*testObject, bool) {
	if len(params) < 2 {
		return nil, false
	}

	object, ok := state.playerObjects[int(params[0])][int(params[1])]

	return object, ok
}

func (state *objectState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}

func (state *objectState) destroyPlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerObject(params); !ok {
		return 0, nil
	}

	delete(state.playerObjects[int(params[0])], int(params[1]))

	return 1, nil
}

func (state *objectState) isValidPlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerObject(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) getPlayerObjectModel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(object.model), nil
}

func (state *objectState) setPlayerObjectPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.setPlayerObjectVector(params, func(object *testObject) *[3]float32 { return &object.position })
}

func (state *objectState) getPlayerObjectPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getPlayerObjectVector(ctx, params, func(object *testObject) *[3]float32 { return &object.position })
}

func (state *objectState) setPlayerObjectRotation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.setPlayerObjectVector(params, func(object *testObject) *[3]float32 { return &object.rotation })
}

func (state *objectState) getPlayerObjectRotation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getPlayerObjectVector(ctx, params, func(object *testObject) *[3]float32 { return &object.rotation })
}

func (state *objectState) setPlayerObjectVector(params []backend.Cell, field func(*testObject) *[3]float32) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	*field(object) = cellsToVector(params[2:5])

	return 1, nil
}

func (state *objectState) getPlayerObjectVector(ctx backend.NativeContext, params []backend.Cell, field func(*testObject) *[3]float32) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[2:5], *field(object))
}

func (state *objectState) movePlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 6 {
		return 0, nil
	}

	object.targetPos = cellsToVector(params[2:5])
	object.moveSpeed = cellFloat(params[5])

	object.moving = true

	object.moveID++
	if len(params) >= 9 {
		object.targetRot = cellsToVector(params[6:9])
	}

	duration := objectMoveDuration(object)
	state.schedulePlayerObjectMove(int(params[0]), int(params[1]), object, duration)

	return duration, nil
}

func (state *objectState) stopPlayerObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok {
		return 0, nil
	}

	object.moving = false
	object.moveID++

	return 1, nil
}

func (state *objectState) schedulePlayerObjectMove(playerID, objectID int, object *testObject, duration backend.Cell) {
	if state.scheduler == nil || duration <= 0 {
		return
	}

	moveID := object.moveID

	state.scheduler.schedule(int64(duration), "OnPlayerObjectMoved", []backend.Cell{backend.Cell(playerID), backend.Cell(objectID)}, func() bool {
		current := state.playerObjects[playerID][objectID]
		if current == nil || !current.moving || current.moveID != moveID {
			return false
		}

		current.position, current.rotation, current.moving = current.targetPos, current.targetRot, false

		return true
	})
}

func (state *objectState) isPlayerObjectMoving(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if ok && object.moving {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) setPlayerObjectMoveSpeed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	object.moveSpeed = cellFloat(params[2])

	return 1, nil
}

func (state *objectState) getPlayerObjectMoveSpeed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok {
		return 0, nil
	}

	return floatCell(object.moveSpeed), nil
}

func (state *objectState) getPlayerObjectTargetPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getPlayerObjectVector(ctx, params, func(object *testObject) *[3]float32 { return &object.targetPos })
}

func (state *objectState) getPlayerObjectTargetRotation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getPlayerObjectVector(ctx, params, func(object *testObject) *[3]float32 { return &object.targetRot })
}

func (state *objectState) getPlayerObjectDrawDistance(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok {
		return 0, nil
	}

	return floatCell(object.drawDistance), nil
}

func (state *objectState) setPlayerObjectNoCameraCollision(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if !ok {
		return 0, nil
	}

	object.noCamera = true

	return 1, nil
}

func (state *objectState) isPlayerObjectNoCameraCollision(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.playerObject(params)
	if ok && object.noCamera {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) assertPlayerObjectValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player object assertion expects 4 arguments")
		}

		if _, ok := state.playerObject(params); !ok {
			setFailure(result, params, 2, fmt.Sprintf("player %d object %d does not exist", params[0], params[1]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
