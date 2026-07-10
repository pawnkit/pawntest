package runner

import (
	"errors"
	"fmt"
	"maps"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

type objectMaterial struct {
	model, colour int
	library, name string
	text          string
	textMaterial  bool
}

type objectAttachment struct {
	kind, id int
	offset   [3]float32
	rotation [3]float32
	sync     bool
}

type testObject struct {
	model                   int
	position, rotation      [3]float32
	targetPos, targetRot    [3]float32
	drawDistance, moveSpeed float32
	moving, noCamera        bool
	attachment              objectAttachment
	materials               map[int]objectMaterial
}

type objectState struct {
	next          int
	objects       map[int]*testObject
	playerNext    map[int]int
	playerObjects map[int]map[int]*testObject
	players       *openMPState
}

func newObjectState() *objectState {
	return &objectState{
		next: 1, objects: map[int]*testObject{},
		playerNext: map[int]int{}, playerObjects: map[int]map[int]*testObject{},
	}
}

func (state *objectState) Clone() scenarioModule {
	clone := newObjectState()
	clone.next = state.next
	maps.Copy(clone.playerNext, state.playerNext)

	for id, object := range state.objects {
		clone.objects[id] = cloneObject(object)
	}

	for playerID, objects := range state.playerObjects {
		clone.playerObjects[playerID] = map[int]*testObject{}
		for id, object := range objects {
			clone.playerObjects[playerID][id] = cloneObject(object)
		}
	}

	return clone
}

func cloneObject(object *testObject) *testObject {
	clone := *object
	clone.materials = maps.Clone(object.materials)

	return &clone
}

func (state *objectState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *objectState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_object_create":         state.createTestObject,
		"__pt_object_valid":          state.assertValid(result),
		"__pt_object_model":          state.assertModel(result),
		"__pt_object_pos_near":       state.assertPosition(result),
		"CreateObject":               state.createObject,
		"DestroyObject":              state.destroyObject,
		"IsValidObject":              state.isValidObject,
		"GetObjectModel":             state.getObjectModel,
		"SetObjectPos":               state.setObjectPosition,
		"GetObjectPos":               state.getObjectPosition,
		"SetObjectRot":               state.setObjectRotation,
		"GetObjectRot":               state.getObjectRotation,
		"MoveObject":                 state.moveObject,
		"StopObject":                 state.stopObject,
		"IsObjectMoving":             state.isObjectMoving,
		"SetObjectMoveSpeed":         state.setObjectMoveSpeed,
		"GetObjectMoveSpeed":         state.getObjectMoveSpeed,
		"GetObjectMovingTargetPos":   state.getObjectTargetPosition,
		"GetObjectTarget":            state.getObjectTargetPosition,
		"GetObjectMovingTargetRot":   state.getObjectTargetRotation,
		"GetObjectDrawDistance":      state.getObjectDrawDistance,
		"SetObjectNoCameraCol":       state.setObjectNoCameraCollision,
		"SetObjectNoCameraCollision": state.setObjectNoCameraCollision,
		"IsObjectNoCameraCol":        state.isObjectNoCameraCollision,
		"AttachObjectToVehicle":      state.attachObject(1),
		"AttachObjectToObject":       state.attachObject(2),
		"AttachObjectToPlayer":       state.attachObject(3),
		"GetObjectAttachedData":      state.getObjectAttachedData,
		"GetObjectAttachedOffset":    state.getObjectAttachedOffset,
		"GetObjectSyncRotation":      state.getObjectSyncRotation,
		"SetObjectMaterial":          state.setObjectMaterial,
		"SetObjectMaterialText":      state.setObjectMaterialText,
		"IsObjectMaterialSlotUsed":   state.isObjectMaterialSlotUsed,

		"__pt_player_object_create":      state.createTestPlayerObject,
		"__pt_player_object_valid":       state.assertPlayerObjectValid(result),
		"CreatePlayerObject":             state.createPlayerObject,
		"DestroyPlayerObject":            state.destroyPlayerObject,
		"IsValidPlayerObject":            state.isValidPlayerObject,
		"GetPlayerObjectModel":           state.getPlayerObjectModel,
		"SetPlayerObjectPos":             state.setPlayerObjectPosition,
		"GetPlayerObjectPos":             state.getPlayerObjectPosition,
		"SetPlayerObjectRot":             state.setPlayerObjectRotation,
		"GetPlayerObjectRot":             state.getPlayerObjectRotation,
		"MovePlayerObject":               state.movePlayerObject,
		"StopPlayerObject":               state.stopPlayerObject,
		"IsPlayerObjectMoving":           state.isPlayerObjectMoving,
		"SetPlayerObjectMoveSpeed":       state.setPlayerObjectMoveSpeed,
		"GetPlayerObjectMoveSpeed":       state.getPlayerObjectMoveSpeed,
		"GetPlayerObjectMovingTargetPos": state.getPlayerObjectTargetPosition,
		"GetPlayerObjectTarget":          state.getPlayerObjectTargetPosition,
		"GetPlayerObjectMovingTargetRot": state.getPlayerObjectTargetRotation,
		"GetPlayerObjectDrawDistance":    state.getPlayerObjectDrawDistance,
		"SetPlayerObjectNoCameraCol":     state.setPlayerObjectNoCameraCollision,
		"IsPlayerObjectNoCameraCol":      state.isPlayerObjectNoCameraCollision,
		"AttachPlayerObjectToVehicle":    state.attachPlayerObject(1),
		"AttachPlayerObjectToObject":     state.attachPlayerObject(2),
		"AttachPlayerObjectToPlayer":     state.attachPlayerObject(3),
		"GetPlayerObjectAttachedData":    state.getPlayerObjectAttachedData,
		"GetPlayerObjectAttachedOffset":  state.getPlayerObjectAttachedOffset,
		"GetPlayerObjectSyncRotation":    state.getPlayerObjectSyncRotation,
		"SetPlayerObjectMaterial":        state.setPlayerObjectMaterial,
		"SetPlayerObjectMaterialText":    state.setPlayerObjectMaterialText,
		"IsPlayerObjectMaterialSlotUsed": state.isPlayerObjectMaterialSlotUsed,
	}
}

func newTestObject(model int, position, rotation [3]float32, drawDistance float32) *testObject {
	return &testObject{model: model, position: position, rotation: rotation, targetPos: position, targetRot: rotation, drawDistance: drawDistance, materials: map[int]objectMaterial{}}
}

func (state *objectState) createTestObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return -1, nil
	}

	return state.addObject(int(params[0]), cellsToVector(params[1:4]), [3]float32{}, 0), nil
}

func (state *objectState) createObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 {
		return -1, nil
	}

	drawDistance := float32(0)
	if len(params) > 7 {
		drawDistance = cellFloat(params[7])
	}

	return state.addObject(int(params[0]), cellsToVector(params[1:4]), cellsToVector(params[4:7]), drawDistance), nil
}

func (state *objectState) addObject(model int, position, rotation [3]float32, drawDistance float32) backend.Cell {
	id := state.next
	state.next++
	state.objects[id] = newTestObject(model, position, rotation, drawDistance)

	return backend.Cell(id)
}

func (state *objectState) object(params []backend.Cell) (*testObject, bool) {
	if len(params) == 0 {
		return nil, false
	}

	object, ok := state.objects[int(params[0])]

	return object, ok
}

func (state *objectState) destroyObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.object(params); !ok {
		return 0, nil
	}

	delete(state.objects, int(params[0]))

	return 1, nil
}

func (state *objectState) isValidObject(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.object(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *objectState) getObjectModel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	object, ok := state.object(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(object.model), nil
}

func cellsToVector(cells []backend.Cell) [3]float32 {
	return [3]float32{cellFloat(cells[0]), cellFloat(cells[1]), cellFloat(cells[2])}
}

func (state *objectState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("object assertion expects 3 arguments")
		}

		if _, ok := state.object(params); !ok {
			setFailure(result, params, 1, fmt.Sprintf("object %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *objectState) assertModel(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("object model assertion expects 4 arguments")
		}

		object, ok := state.object(params)
		if !ok || object.model != int(params[1]) {
			actual := -1
			if ok {
				actual = object.model
			}

			setFailure(result, params, 2, fmt.Sprintf("object %d model: expected %d, got %d", params[0], params[1], actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *objectState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("object position assertion expects 7 arguments")
		}

		object, ok := state.object(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("object %d does not exist", params[0]), ctx)
			return 0, nil
		}

		expected := cellsToVector(params[1:4])

		tolerance := float32(math.Abs(float64(cellFloat(params[4]))))
		for index := range expected {
			if absFloat(object.position[index]-expected[index]) > tolerance {
				setFailure(result, params, 5, fmt.Sprintf("object %d position: expected %v +/- %g, got %v", params[0], expected, tolerance, object.position), ctx)
				return 0, nil
			}
		}

		return 1, nil
	}
}
