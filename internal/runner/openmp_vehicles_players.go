package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *vehicleState) putPlayerInVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, playerOK := state.player(params)
	if !playerOK || len(params) < 3 {
		return 0, nil
	}

	vehicle, vehicleOK := state.vehicles[int(params[1])]
	if !vehicleOK {
		return 0, nil
	}

	player.vehicle, player.seat = int(params[1]), int(params[2])

	vehicle.occupied = true
	if player.seat == 0 {
		vehicle.lastDriver = int(params[0])
	}

	return 1, nil
}

func (state *vehicleState) removePlayerFromVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params)
	if !ok || player.vehicle < 0 {
		return 0, nil
	}

	vehicleID := player.vehicle
	player.vehicle, player.seat = -1, -1

	state.refreshOccupied(vehicleID)

	return 1, nil
}

func (state *vehicleState) getPlayerVehicleID(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params)
	if !ok || player.vehicle < 0 {
		return 0, nil
	}

	return backend.Cell(player.vehicle), nil
}

func (state *vehicleState) getPlayerVehicleSeat(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params)
	if !ok || player.seat < 0 {
		return -1, nil
	}

	return backend.Cell(player.seat), nil
}

func (state *vehicleState) isPlayerInVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params)
	if ok && len(params) >= 2 && player.vehicle == int(params[1]) {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) isPlayerInAnyVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params)
	if ok && player.vehicle >= 0 {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) getVehicleDriver(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicleID, ok := state.vehicleID(params)
	if !ok {
		return -1, nil
	}

	for id, player := range state.players.players {
		if player.vehicle == vehicleID && player.seat == 0 {
			return backend.Cell(id), nil
		}
	}

	return -1, nil
}

func (state *vehicleState) getVehicleLastDriver(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(vehicle.lastDriver), nil
}

func (state *vehicleState) getVehicleOccupant(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicleID, ok := state.vehicleID(params)
	if !ok || len(params) < 2 {
		return -1, nil
	}

	for id, player := range state.players.players {
		if player.vehicle == vehicleID && player.seat == int(params[1]) {
			return backend.Cell(id), nil
		}
	}

	return -1, nil
}

func (state *vehicleState) countVehicleOccupants(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicleID, ok := state.vehicleID(params)
	if !ok {
		return 0, nil
	}

	count := 0

	for _, player := range state.players.players {
		if player.vehicle == vehicleID {
			count++
		}
	}

	return backend.Cell(count), nil
}

func (state *vehicleState) isVehicleOccupied(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if ok && vehicle.occupied {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) hasVehicleBeenOccupied(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if ok && vehicle.lastDriver >= 0 {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) player(params []backend.Cell) (*testPlayer, bool) {
	if state.players == nil || len(params) == 0 {
		return nil, false
	}

	return state.players.player(params[0])
}

func (state *vehicleState) vehicleID(params []backend.Cell) (int, bool) {
	if _, ok := state.vehicle(params); !ok {
		return 0, false
	}

	return int(params[0]), true
}

func (state *vehicleState) ejectVehicle(vehicleID int) {
	if state.players == nil {
		return
	}

	for _, player := range state.players.players {
		if player.vehicle == vehicleID {
			player.vehicle, player.seat = -1, -1
		}
	}

	state.refreshOccupied(vehicleID)
}

func (state *vehicleState) refreshOccupied(vehicleID int) {
	vehicle, ok := state.vehicles[vehicleID]
	if !ok {
		return
	}

	vehicle.occupied = false

	if state.players == nil {
		return
	}

	for _, player := range state.players.players {
		if player.vehicle == vehicleID {
			vehicle.occupied = true

			return
		}
	}
}
