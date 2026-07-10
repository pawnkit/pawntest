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
