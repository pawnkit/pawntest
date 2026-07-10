package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *classState) selectClass(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 || !state.validPlayer(params[0]) || params[1] < 0 || int(params[1]) >= len(state.classes) {
		return 0, nil
	}

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support class callbacks")
	}

	playerID, classID := int(params[0]), int(params[1])
	state.selected[playerID] = classID
	state.spawnInfo[playerID] = state.classes[classID]
	state.selecting[playerID] = true

	result, err := caller.CallPublic("OnPlayerRequestClass", params[0], params[1])
	if err != nil {
		return 0, fmt.Errorf("class selection callback: %w", err)
	}

	return result, nil
}

func (state *classState) spawnPlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	info, ok := state.spawnInfo[int(params[0])]
	if !ok {
		return 0, nil
	}

	player := state.players.players[int(params[0])]
	player.team, player.skin = info.team, info.skin
	player.x, player.y, player.z = info.position[0], info.position[1], info.position[2]
	player.angle, player.spawned = info.angle, true

	player.weapons = map[int]int{}
	for _, weapon := range info.weapons {
		if weapon.weapon != 0 {
			player.weapons[weapon.weapon] = weapon.ammo
		}
	}

	state.selecting[int(params[0])] = false

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 1, nil
	}

	if _, err := caller.CallPublic("OnPlayerSpawn", params[0]); err != nil {
		return 0, fmt.Errorf("player spawn callback: %w", err)
	}

	return 1, nil
}
