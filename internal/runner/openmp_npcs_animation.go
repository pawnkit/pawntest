package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) setNPCAnimation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 8 {
		return 0, nil
	}

	npc.animationID = int(params[1])
	npc.animation = actorAnimation{
		delta: cellFloat(params[2]), loop: params[3] != 0, lockX: params[4] != 0,
		lockY: params[5] != 0, freeze: params[6] != 0, time: int(params[7]),
	}
	npc.hasAnimation = true

	return 1, nil
}

func (state *npcState) applyNPCAnimation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 9 {
		return 0, nil
	}

	animation, err := readAnimation(ctx, params[1:9])
	if err != nil {
		return 0, err
	}

	npc.animation = animation
	npc.hasAnimation = true

	return 1, nil
}

func (state *npcState) clearNPCAnimation(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.animationID, npc.animation, npc.hasAnimation = 0, actorAnimation{}, false

	return 1, nil
}

func (state *npcState) getNPCAnimation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || !npc.hasAnimation || len(params) < 8 {
		return 0, nil
	}

	values := []backend.Cell{
		backend.Cell(npc.animationID), floatCell(npc.animation.delta), boolCell(npc.animation.loop),
		boolCell(npc.animation.lockX), boolCell(npc.animation.lockY), boolCell(npc.animation.freeze), backend.Cell(npc.animation.time),
	}
	for index, value := range values {
		if err := ctx.WriteCell(params[index+1], value); err != nil {
			return 0, err
		}
	}

	return 1, nil
}
