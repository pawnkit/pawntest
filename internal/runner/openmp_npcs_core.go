package runner

import (
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *npcState) destroyNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.npc(params); !ok {
		return 0, nil
	}

	delete(state.npcs, int(params[0]))

	return 1, nil
}

func (state *npcState) isValidNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.npc(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isDeadNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.dead {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) spawnNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.spawned, npc.dead = true, false
	if npc.health <= 0 {
		npc.health = 100
	}

	return 1, nil
}

func (state *npcState) isSpawnedNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.spawned {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) getAllNPCs(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	ids := make([]int, 0, len(state.npcs))
	for id := range state.npcs {
		ids = append(ids, id)
	}

	return writeEntityIDs(ctx, params, ids)
}

func (state *npcState) setNPCVector(field func(*testNPC) *[3]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		*field(npc) = cellsToVector(params[1:4])

		return 1, nil
	}
}

func (state *npcState) getNPCVector(field func(*testNPC) *[3]float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		return writeFloatVector(ctx, params[1:4], *field(npc))
	}
}

func (state *npcState) setNPCFloat(field func(*testNPC) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(npc) = cellFloat(params[1])

		return 1, nil
	}
}

func (state *npcState) getNPCFloat(field func(*testNPC) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok {
			return 0, nil
		}

		return floatCell(*field(npc)), nil
	}
}

func (state *npcState) getNPCFloatOutput(field func(*testNPC) *float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return 1, ctx.WriteCell(params[1], floatCell(*field(npc)))
	}
}

func (state *npcState) setNPCInt(field func(*testNPC) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(npc) = int(params[1])

		return 1, nil
	}
}

func (state *npcState) getNPCInt(field func(*testNPC) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(npc)), nil
	}
}

func (state *npcState) setNPCInvulnerable(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.invulnerable = len(params) < 2 || params[1] != 0

	return 1, nil
}

func (state *npcState) isNPCInvulnerable(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.invulnerable {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) killNPC(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || npc.invulnerable {
		return 0, nil
	}

	npc.dead, npc.spawned, npc.moving, npc.health = true, false, false, 0

	return 1, nil
}

func (state *npcState) setNPCAngleToPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	dx := cellFloat(params[1]) - npc.position[0]
	dy := cellFloat(params[2]) - npc.position[1]
	npc.facing = float32(math.Atan2(float64(dy), float64(dx)) * 180 / math.Pi)

	return 1, nil
}

func (state *npcState) setNPCAngleToPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[1])
	if !playerOK {
		return 0, nil
	}

	dx, dy := player.x-npc.position[0], player.y-npc.position[1]
	npc.facing = float32(math.Atan2(float64(dy), float64(dx)) * 180 / math.Pi)

	return 1, nil
}
