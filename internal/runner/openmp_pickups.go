package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testPickup struct {
	model, pickupType, world int
	position                 [3]float32
	hidden                   map[int]bool
}

type pickupState struct {
	next          int
	pickups       map[int]*testPickup
	playerNext    map[int]int
	playerPickups map[int]map[int]*testPickup
	players       *openMPState
}

func newPickupState() *pickupState {
	return &pickupState{
		pickups: map[int]*testPickup{}, playerNext: map[int]int{},
		playerPickups: map[int]map[int]*testPickup{},
	}
}

func (state *pickupState) Clone() scenarioModule {
	clone := newPickupState()
	clone.next = state.next
	maps.Copy(clone.playerNext, state.playerNext)

	for id, pickup := range state.pickups {
		clone.pickups[id] = clonePickup(pickup)
	}

	for playerID, pickups := range state.playerPickups {
		clone.playerPickups[playerID] = map[int]*testPickup{}
		for id, pickup := range pickups {
			clone.playerPickups[playerID][id] = clonePickup(pickup)
		}
	}

	return clone
}

func clonePickup(pickup *testPickup) *testPickup {
	clone := *pickup
	clone.hidden = maps.Clone(pickup.hidden)

	return &clone
}

func (state *pickupState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *pickupState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_pickup_create":      state.createPickup,
		"__pt_pickup_valid":       state.assertValid(result),
		"__pt_pickup_model":       state.assertModel(result),
		"__pt_pickup_pos_near":    state.assertPosition(result),
		"AddStaticPickup":         state.createPickup,
		"CreatePickup":            state.createPickup,
		"DestroyPickup":           state.destroyPickup,
		"IsValidPickup":           state.isValidPickup,
		"IsPickupStreamedIn":      state.isPickupStreamedIn,
		"GetPickupPos":            state.getPickupPosition,
		"SetPickupPos":            state.setPickupPosition,
		"GetPickupModel":          state.getPickupInt(pickupModel),
		"SetPickupModel":          state.setPickupInt(pickupModel),
		"GetPickupType":           state.getPickupInt(pickupType),
		"SetPickupType":           state.setPickupInt(pickupType),
		"GetPickupVirtualWorld":   state.getPickupInt(pickupWorld),
		"SetPickupVirtualWorld":   state.setPickupInt(pickupWorld),
		"ShowPickupForPlayer":     state.showPickupForPlayer,
		"HidePickupForPlayer":     state.hidePickupForPlayer,
		"IsPickupHiddenForPlayer": state.isPickupHiddenForPlayer,

		"__pt_player_pickup_create":   state.createPlayerPickup,
		"__pt_player_pickup_valid":    state.assertPlayerPickupValid(result),
		"CreatePlayerPickup":          state.createPlayerPickup,
		"DestroyPlayerPickup":         state.destroyPlayerPickup,
		"IsValidPlayerPickup":         state.isValidPlayerPickup,
		"IsPlayerPickupStreamedIn":    state.isPlayerPickupStreamedIn,
		"GetPlayerPickupPos":          state.getPlayerPickupPosition,
		"SetPlayerPickupPos":          state.setPlayerPickupPosition,
		"GetPlayerPickupModel":        state.getPlayerPickupInt(pickupModel),
		"SetPlayerPickupModel":        state.setPlayerPickupInt(pickupModel),
		"GetPlayerPickupType":         state.getPlayerPickupInt(pickupType),
		"SetPlayerPickupType":         state.setPlayerPickupInt(pickupType),
		"GetPlayerPickupVirtualWorld": state.getPlayerPickupInt(pickupWorld),
		"SetPlayerPickupVirtualWorld": state.setPlayerPickupInt(pickupWorld),
	}
}

func (state *pickupState) createPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return -1, nil
	}

	world := 0
	if len(params) > 5 {
		world = int(params[5])
	}

	id := state.next
	state.next++
	state.pickups[id] = newTestPickup(int(params[0]), int(params[1]), cellsToVector(params[2:5]), world)

	return backend.Cell(id), nil
}

func newTestPickup(model, pickupType int, position [3]float32, world int) *testPickup {
	return &testPickup{model: model, pickupType: pickupType, position: position, world: world, hidden: map[int]bool{}}
}

func (state *pickupState) pickup(params []backend.Cell) (*testPickup, bool) {
	if len(params) == 0 {
		return nil, false
	}

	pickup, ok := state.pickups[int(params[0])]

	return pickup, ok
}

func (state *pickupState) destroyPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.pickup(params); !ok {
		return 0, nil
	}

	delete(state.pickups, int(params[0]))

	return 1, nil
}

func (state *pickupState) isValidPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.pickup(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *pickupState) setPickupInt(field func(*testPickup) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		pickup, ok := state.pickup(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(pickup) = int(params[1])

		return 1, nil
	}
}

func (state *pickupState) getPickupInt(field func(*testPickup) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		pickup, ok := state.pickup(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(pickup)), nil
	}
}

func pickupModel(pickup *testPickup) *int { return &pickup.model }
func pickupType(pickup *testPickup) *int  { return &pickup.pickupType }
func pickupWorld(pickup *testPickup) *int { return &pickup.world }

func (state *pickupState) setPickupPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.pickup(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	pickup.position = cellsToVector(params[1:4])

	return 1, nil
}

func (state *pickupState) getPickupPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.pickup(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[1:4], pickup.position)
}

func (state *pickupState) isPickupStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[0])

	pickup, pickupOK := state.pickups[int(params[1])]
	if playerOK && pickupOK && player.connected && player.world == pickup.world && !pickup.hidden[int(params[0])] {
		return 1, nil
	}

	return 0, nil
}

func (state *pickupState) showPickupForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickupArgument(params)
	if !ok {
		return 0, nil
	}

	delete(pickup.hidden, int(params[0]))

	return 1, nil
}

func (state *pickupState) hidePickupForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickupArgument(params)
	if !ok {
		return 0, nil
	}

	pickup.hidden[int(params[0])] = true

	return 1, nil
}

func (state *pickupState) isPickupHiddenForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickupArgument(params)
	if ok && pickup.hidden[int(params[0])] {
		return 1, nil
	}

	return 0, nil
}

func (state *pickupState) playerPickupArgument(params []backend.Cell) (*testPickup, bool) {
	if len(params) < 2 {
		return nil, false
	}

	pickup, ok := state.pickups[int(params[1])]

	return pickup, ok
}
