package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) startNPCPlayback(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	npc.playbackName, npc.playbackRecord = name, -1
	npc.playingPlayback, npc.pausedPlayback = true, false

	return 1, nil
}

func (state *npcState) startNPCPlaybackRecord(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	if _, exists := state.records[int(params[1])]; !exists {
		return 0, nil
	}

	npc.playbackName, npc.playbackRecord = "", int(params[1])
	npc.playingPlayback, npc.pausedPlayback = true, false

	return 1, nil
}

func (state *npcState) stopNPCPlayback(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.playingPlayback, npc.pausedPlayback = false, false

	return 1, nil
}

func (state *npcState) pauseNPCPlayback(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || !npc.playingPlayback {
		return 0, nil
	}

	npc.pausedPlayback = len(params) < 2 || params[1] != 0

	return 1, nil
}

func (state *npcState) isNPCPlayingPlayback(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.playingPlayback {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isNPCPlaybackPaused(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.playingPlayback && npc.pausedPlayback {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) resetNPCSurfing(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.surfingOffset = [3]float32{}
	npc.surfingVehicle, npc.surfingObject, npc.surfingPlayerObject = -1, -1, -1

	return 1, nil
}

func (state *npcState) loadNPCRecord(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return -1, nil
	}

	path, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	id := state.nextRecord
	state.nextRecord++
	state.records[id] = path

	return backend.Cell(id), nil
}

func (state *npcState) unloadNPCRecord(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	if _, ok := state.records[int(params[0])]; !ok {
		return 0, nil
	}

	delete(state.records, int(params[0]))

	return 1, nil
}

func (state *npcState) isValidNPCRecord(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 {
		if _, ok := state.records[int(params[0])]; ok {
			return 1, nil
		}
	}

	return 0, nil
}

func (state *npcState) getNPCRecordCount(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return backend.Cell(len(state.records)), nil
}

func (state *npcState) unloadAllNPCRecords(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	clear(state.records)

	return 1, nil
}
