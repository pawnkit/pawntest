package runner

import (
	"math"
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *npcState) createNPCPath(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	id := state.nextPath
	state.nextPath++
	state.paths[id] = &npcPath{}

	return backend.Cell(id), nil
}

func (state *npcState) destroyNPCPath(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	if _, ok := state.paths[int(params[0])]; !ok {
		return 0, nil
	}

	delete(state.paths, int(params[0]))

	return 1, nil
}

func (state *npcState) destroyAllNPCPaths(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	clear(state.paths)

	return 1, nil
}

func (state *npcState) getNPCPathCount(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return backend.Cell(len(state.paths)), nil
}

func (state *npcState) addNPCPathPoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok {
		return 0, nil
	}

	stopRange := float32(0.2)
	if len(params) > 4 {
		stopRange = cellFloat(params[4])
	}

	path.points = append(path.points, npcPathPoint{position: cellsToVector(params[1:4]), stopRange: stopRange})

	return 1, nil
}

func (state *npcState) removeNPCPathPoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok || params[1] < 0 || int(params[1]) >= len(path.points) {
		return 0, nil
	}

	path.points = slices.Delete(path.points, int(params[1]), int(params[1])+1)

	return 1, nil
}

func (state *npcState) clearNPCPath(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok {
		return 0, nil
	}

	path.points = nil

	return 1, nil
}

func (state *npcState) getNPCPathPointCount(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok {
		return 0, nil
	}

	return backend.Cell(len(path.points)), nil
}

func (state *npcState) getNPCPathPoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 6 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok || params[1] < 0 || int(params[1]) >= len(path.points) {
		return 0, nil
	}

	point := path.points[int(params[1])]
	if result, err := writeFloatVector(ctx, params[2:5], point.position); err != nil || result == 0 {
		return result, err
	}

	return 1, ctx.WriteCell(params[5], floatCell(point.stopRange))
}

func (state *npcState) isValidNPCPath(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 {
		if _, ok := state.paths[int(params[0])]; ok {
			return 1, nil
		}
	}

	return 0, nil
}

func (state *npcState) moveNPCByPath(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	path, pathOK := state.paths[int(params[1])]
	if !pathOK || len(path.points) == 0 {
		return 0, nil
	}

	npc.currentPath, npc.currentPathPoint, npc.moving = int(params[1]), 0, true
	npc.moveTarget = path.points[0].position

	return 1, nil
}

func (state *npcState) getNPCCurrentPathPoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(npc.currentPathPoint), nil
}

func (state *npcState) hasNPCPathPointInRange(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return 0, nil
	}

	path, ok := state.paths[int(params[0])]
	if !ok {
		return 0, nil
	}

	target := cellsToVector(params[1:4])
	radius := cellFloat(params[4])

	for _, point := range path.points {
		dx, dy, dz := point.position[0]-target[0], point.position[1]-target[1], point.position[2]-target[2]
		if float32(math.Sqrt(float64(dx*dx+dy*dy+dz*dz))) <= radius {
			return 1, nil
		}
	}

	return 0, nil
}
