package runner

import (
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

func writeEntityIDs(ctx backend.NativeContext, params []backend.Cell, ids []int) (backend.Cell, error) {
	if len(params) < 2 || params[1] <= 0 {
		return 0, nil
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

func readAnimation(ctx backend.NativeContext, params []backend.Cell) (actorAnimation, error) {
	library, err := ctx.ReadString(params[0])
	if err != nil {
		return actorAnimation{}, err
	}

	name, err := ctx.ReadString(params[1])
	if err != nil {
		return actorAnimation{}, err
	}

	return actorAnimation{
		library: library, name: name, delta: cellFloat(params[2]),
		loop: params[3] != 0, lockX: params[4] != 0, lockY: params[5] != 0,
		freeze: params[6] != 0, time: int(params[7]),
	}, nil
}
