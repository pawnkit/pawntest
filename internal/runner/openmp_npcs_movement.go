package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) moveNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	npc.moveTarget = cellsToVector(params[1:4])
	npc.moving, npc.movingToPlayer, npc.movePlayer = true, false, -1

	return 1, nil
}

func (state *npcState) moveNPCToPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[1])
	if !playerOK {
		return 0, nil
	}

	npc.moveTarget = [3]float32{player.x, player.y, player.z}
	npc.moving, npc.movingToPlayer, npc.movePlayer = true, true, int(params[1])

	return 1, nil
}

func (state *npcState) stopNPCMove(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.moving, npc.movingToPlayer, npc.movePlayer = false, false, -1

	return 1, nil
}

func (state *npcState) isNPCMoving(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.moving {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isNPCMovingToPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && len(params) >= 2 && npc.movingToPlayer && npc.movePlayer == int(params[1]) {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isNPCStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.players == nil || !npc.spawned {
		return 0, nil
	}

	player, playerOK := state.players.player(params[1])
	if playerOK && player.connected && player.world == npc.world {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isNPCAnyStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || state.players == nil || !npc.spawned {
		return 0, nil
	}

	for _, player := range state.players.players {
		if player.connected && player.world == npc.world {
			return 1, nil
		}
	}

	return 0, nil
}

func (state *npcState) putNPCInVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 3 || state.vehicles == nil {
		return 0, nil
	}

	if _, vehicleOK := state.vehicles.vehicles[int(params[1])]; !vehicleOK {
		return 0, nil
	}

	npc.vehicle, npc.seat = int(params[1]), int(params[2])

	return 1, nil
}

func (state *npcState) removeNPCFromVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || npc.vehicle < 0 {
		return 0, nil
	}

	npc.vehicle, npc.seat = -1, -1

	return 1, nil
}

func (state *npcState) getNPCVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || npc.vehicle < 0 {
		return -1, nil
	}

	return backend.Cell(npc.vehicle), nil
}

func (state *npcState) getNPCVehicleSeat(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(npc.seat), nil
}
