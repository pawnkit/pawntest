package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *pickupState) createPlayerPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 6 || !state.validPlayer(params[0]) {
		return -1, nil
	}

	world := 0
	if len(params) > 6 {
		world = int(params[6])
	}

	playerID := int(params[0])
	id := state.playerNext[playerID]

	state.playerNext[playerID]++
	if state.playerPickups[playerID] == nil {
		state.playerPickups[playerID] = map[int]*testPickup{}
	}

	state.playerPickups[playerID][id] = newTestPickup(int(params[1]), int(params[2]), cellsToVector(params[3:6]), world)

	return backend.Cell(id), nil
}

func (state *pickupState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}

func (state *pickupState) playerPickup(params []backend.Cell) (*testPickup, bool) {
	if len(params) < 2 {
		return nil, false
	}

	pickup, ok := state.playerPickups[int(params[0])][int(params[1])]

	return pickup, ok
}

func (state *pickupState) destroyPlayerPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerPickup(params); !ok {
		return 0, nil
	}

	delete(state.playerPickups[int(params[0])], int(params[1]))

	return 1, nil
}

func (state *pickupState) isValidPlayerPickup(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerPickup(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *pickupState) isPlayerPickupStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickup(params)
	if !ok || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[0])
	if playerOK && player.connected && player.world == pickup.world {
		return 1, nil
	}

	return 0, nil
}

func (state *pickupState) setPlayerPickupInt(field func(*testPickup) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		pickup, ok := state.playerPickup(params)
		if !ok || len(params) < 3 {
			return 0, nil
		}

		*field(pickup) = int(params[2])

		return 1, nil
	}
}

func (state *pickupState) getPlayerPickupInt(field func(*testPickup) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		pickup, ok := state.playerPickup(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(pickup)), nil
	}
}

func (state *pickupState) setPlayerPickupPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickup(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	pickup.position = cellsToVector(params[2:5])

	return 1, nil
}

func (state *pickupState) getPlayerPickupPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	pickup, ok := state.playerPickup(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[2:5], pickup.position)
}
