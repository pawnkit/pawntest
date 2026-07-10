package runner

import (
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *actorState) applyActorAnimation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || len(params) < 9 {
		return 0, nil
	}

	library, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	name, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	actor.animation = actorAnimation{
		library: library, name: name, delta: cellFloat(params[3]),
		loop: params[4] != 0, lockX: params[5] != 0, lockY: params[6] != 0,
		freeze: params[7] != 0, time: int(params[8]),
	}
	actor.hasAnimation = true

	return 1, nil
}

func (state *actorState) clearActorAnimations(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok {
		return 0, nil
	}

	actor.animation = actorAnimation{}
	actor.hasAnimation = false

	return 1, nil
}

func (state *actorState) getActorAnimation(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || !actor.hasAnimation || len(params) < 11 {
		return 0, nil
	}

	if err := ctx.WriteString(params[1], truncateString(actor.animation.library, int(params[2]))); err != nil {
		return 0, err
	}

	if err := ctx.WriteString(params[3], truncateString(actor.animation.name, int(params[4]))); err != nil {
		return 0, err
	}

	values := []backend.Cell{
		floatCell(actor.animation.delta), boolCell(actor.animation.loop), boolCell(actor.animation.lockX),
		boolCell(actor.animation.lockY), boolCell(actor.animation.freeze), backend.Cell(actor.animation.time),
	}
	for index, value := range values {
		if err := ctx.WriteCell(params[index+5], value); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func boolCell(value bool) backend.Cell {
	if value {
		return 1
	}

	return 0
}

func (state *actorState) getActorSpawnInfo(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || len(params) < 6 {
		return 0, nil
	}

	if err := ctx.WriteCell(params[1], backend.Cell(actor.spawnSkin)); err != nil {
		return 0, err
	}

	if result, err := writeFloatVector(ctx, params[2:5], actor.spawnPosition); err != nil || result == 0 {
		return result, err
	}

	return 1, ctx.WriteCell(params[5], floatCell(actor.spawnAngle))
}

func (state *actorState) getActors(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || params[1] <= 0 {
		return 0, nil
	}

	ids := make([]int, 0, len(state.actors))
	for id := range state.actors {
		ids = append(ids, id)
	}

	slices.Sort(ids)

	limit := min(len(ids), int(params[1]))
	for index, id := range ids[:limit] {
		if err := ctx.WriteCell(params[0]+backend.Cell(index*4), backend.Cell(id)); err != nil {
			return 0, err
		}
	}

	return backend.Cell(limit), nil
}
