package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) shootNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 7 {
		return 0, nil
	}

	npc.weapon, npc.shooting = int(params[1]), true

	npc.aimPoint = cellsToVector(params[4:7])
	if !npc.infiniteAmmo && npc.ammo > 0 {
		npc.ammo--
	}

	return 1, nil
}

func (state *npcState) aimNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	npc.aimPoint, npc.aiming, npc.aimPlayer = cellsToVector(params[1:4]), true, -1
	if len(params) > 4 {
		npc.shooting = params[4] != 0
	}

	return 1, nil
}

func (state *npcState) aimNPCAtPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[1])
	if !playerOK {
		return 0, nil
	}

	npc.aimPoint = [3]float32{player.x, player.y, player.z}

	npc.aiming, npc.aimPlayer = true, int(params[1])
	if len(params) > 2 {
		npc.shooting = params[2] != 0
	}

	return 1, nil
}

func (state *npcState) stopNPCAim(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.aiming, npc.shooting, npc.aimPlayer = false, false, -1

	return 1, nil
}

func (state *npcState) isNPCAimingAtPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && len(params) >= 2 && npc.aiming && npc.aimPlayer == int(params[1]) {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) getNPCAimPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(npc.aimPlayer), nil
}

func (state *npcState) setNPCWeaponInt(field func(*testNPC) map[int]int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 3 {
			return 0, nil
		}

		field(npc)[int(params[1])] = int(params[2])

		return 1, nil
	}
}

func (state *npcState) getNPCWeaponInt(field func(*testNPC) map[int]int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return backend.Cell(field(npc)[int(params[1])]), nil
	}
}

func (state *npcState) setNPCWeaponFloat(field func(*testNPC) map[int]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 3 {
			return 0, nil
		}

		field(npc)[int(params[1])] = cellFloat(params[2])

		return 1, nil
	}
}

func (state *npcState) getNPCWeaponFloat(field func(*testNPC) map[int]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return floatCell(field(npc)[int(params[1])]), nil
	}
}
