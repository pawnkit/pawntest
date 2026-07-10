package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *openMPState) cameraNatives() map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"SetPlayerCameraPos":    state.setPlayerVector(playerCamera),
		"GetPlayerCameraPos":    state.getPlayerVector(playerCamera),
		"SetCameraBehindPlayer": state.setCameraBehindPlayer,
	}
}

func (state *openMPState) setCameraBehindPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	player.camera = [3]float32{player.x, player.y, player.z}

	return 1, nil
}
