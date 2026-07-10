package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) navigationNatives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_npc_playback": state.assertPlayback(result), "__pt_npc_path_count": state.assertNavigationCount(result, "paths", func() int { return len(state.paths) }),
		"__pt_npc_record_count": state.assertNavigationCount(result, "records", func() int { return len(state.records) }),
		"NPC_StartPlayback":     state.startNPCPlayback, "NPC_StartPlaybackEx": state.startNPCPlaybackRecord,
		"NPC_StopPlayback": state.stopNPCPlayback, "NPC_PausePlayback": state.pauseNPCPlayback,
		"NPC_IsPlayingPlayback": state.isNPCPlayingPlayback, "NPC_IsPlaybackPaused": state.isNPCPlaybackPaused,
		"NPC_SetSurfingOffsets": state.setNPCVector(npcSurfingOffset), "NPC_GetSurfingOffsets": state.getNPCVector(npcSurfingOffset),
		"NPC_SetSurfingVehicle": state.setNPCInt(npcSurfingVehicle), "NPC_GetSurfingVehicle": state.getNPCInt(npcSurfingVehicle),
		"NPC_SetSurfingObject": state.setNPCInt(npcSurfingObject), "NPC_GetSurfingObject": state.getNPCInt(npcSurfingObject),
		"NPC_SetSurfingPlayerObject": state.setNPCInt(npcSurfingPlayerObject), "NPC_GetSurfingPlayerObject": state.getNPCInt(npcSurfingPlayerObject),
		"NPC_ResetSurfingData": state.resetNPCSurfing, "NPC_LoadRecord": state.loadNPCRecord,
		"NPC_UnloadRecord": state.unloadNPCRecord, "NPC_IsValidRecord": state.isValidNPCRecord,
		"NPC_GetRecordCount": state.getNPCRecordCount, "NPC_UnloadAllRecords": state.unloadAllNPCRecords,
		"NPC_CreatePath": state.createNPCPath, "NPC_DestroyPath": state.destroyNPCPath,
		"NPC_DestroyAllPath": state.destroyAllNPCPaths, "NPC_GetPathCount": state.getNPCPathCount,
		"NPC_AddPointToPath": state.addNPCPathPoint, "NPC_RemovePointFromPath": state.removeNPCPathPoint,
		"NPC_ClearPath": state.clearNPCPath, "NPC_GetPathPointCount": state.getNPCPathPointCount,
		"NPC_GetPathPoint": state.getNPCPathPoint, "NPC_IsValidPath": state.isValidNPCPath,
		"NPC_GetCurrentPathPointIndex": state.getNPCCurrentPathPoint, "NPC_MoveByPath": state.moveNPCByPath,
		"NPC_HasPathPointInRange": state.hasNPCPathPointInRange, "NPC_OpenNode": state.openNPCNode,
		"NPC_CloseNode": state.closeNPCNode, "NPC_IsNodeOpen": state.isNPCNodeOpen,
		"NPC_GetNodeType": state.getNPCNodeType, "NPC_SetNodePoint": state.setNPCNodePoint,
		"NPC_GetNodePointPosition": state.getNPCNodePointPosition, "NPC_GetNodePointCount": state.getNPCNodePointCount,
		"NPC_GetNodeInfo": state.getNPCNodeInfo, "NPC_PlayNode": state.playNPCNode,
		"NPC_StopPlayingNode": state.stopNPCNode, "NPC_PausePlayingNode": state.pauseNPCNode,
		"NPC_ResumePlayingNode": state.resumeNPCNode, "NPC_IsPlayingNode": state.isNPCPlayingNode,
		"NPC_IsPlayingNodePaused": state.isNPCNodePaused, "NPC_ChangeNode": state.changeNPCNode,
		"NPC_UpdateNodePoint": state.updateNPCNodePoint,
	}
}

func npcSurfingOffset(npc *testNPC) *[3]float32 { return &npc.surfingOffset }
func npcSurfingVehicle(npc *testNPC) *int       { return &npc.surfingVehicle }
func npcSurfingObject(npc *testNPC) *int        { return &npc.surfingObject }
func npcSurfingPlayerObject(npc *testNPC) *int  { return &npc.surfingPlayerObject }
