package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) openNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	state.nodes[int(params[0])] = &npcNode{open: true, point: -1}

	return 1, nil
}

func (state *npcState) closeNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	node, ok := state.nodes[int(params[0])]
	if !ok {
		return 0, nil
	}

	node.open = false

	return 1, nil
}

func (state *npcState) isNPCNodeOpen(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 && state.nodes[int(params[0])] != nil && state.nodes[int(params[0])].open {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) getNPCNodeType(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || state.nodes[int(params[0])] == nil {
		return -1, nil
	}

	return 0, nil
}

func (state *npcState) setNPCNodePoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || state.nodes[int(params[0])] == nil {
		return 0, nil
	}

	state.nodes[int(params[0])].point = int(params[1])

	return 1, nil
}

func (state *npcState) getNPCNodePointPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 || state.nodes[int(params[0])] == nil {
		return 0, nil
	}

	return writeFloatVector(ctx, params[1:4], [3]float32{})
}

func (state *npcState) getNPCNodePointCount(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || state.nodes[int(params[0])] == nil {
		return 0, nil
	}

	return 0, nil
}

func (state *npcState) getNPCNodeInfo(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 || state.nodes[int(params[0])] == nil {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[1:4], []int{0, 0, 0})
}

func (state *npcState) playNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.nodes[int(params[1])] == nil || !state.nodes[int(params[1])].open {
		return 0, nil
	}

	npc.playingNode, npc.nodeID, npc.nodePaused = int(params[1]), int(params[1]), false

	return 1, nil
}

func (state *npcState) stopNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.playingNode, npc.nodeID, npc.nodePaused = -1, -1, false

	return 1, nil
}

func (state *npcState) pauseNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || npc.playingNode < 0 {
		return 0, nil
	}

	npc.nodePaused = true

	return 1, nil
}

func (state *npcState) resumeNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || npc.playingNode < 0 {
		return 0, nil
	}

	npc.nodePaused = false

	return 1, nil
}

func (state *npcState) isNPCPlayingNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.playingNode >= 0 {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) isNPCNodePaused(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.playingNode >= 0 && npc.nodePaused {
		return 1, nil
	}

	return 0, nil
}

func (state *npcState) changeNPCNode(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || state.nodes[int(params[1])] == nil {
		return 0, nil
	}

	npc.playingNode, npc.nodeID = int(params[1]), int(params[1])

	return 1, nil
}

func (state *npcState) updateNPCNodePoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 2 || npc.playingNode < 0 {
		return 0, nil
	}

	state.nodes[npc.playingNode].point = int(params[1])

	return 1, nil
}
