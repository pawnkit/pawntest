package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *openMPState) equipmentNatives() map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"ResetPlayerMoney":     state.resetPlayerMoney,
		"GivePlayerWeapon":     state.givePlayerWeapon,
		"RemovePlayerWeapon":   state.removePlayerWeapon,
		"ResetPlayerWeapons":   state.resetPlayerWeapons,
		"SetPlayerAmmo":        state.setPlayerAmmo,
		"SetPlayerArmedWeapon": state.setPlayerArmedWeapon,
		"GetPlayerWeapon":      state.getPlayerWeapon,
		"GetPlayerAmmo":        state.getPlayerAmmo,
		"GetPlayerWeaponData":  state.getPlayerWeaponData,
	}
}

func (state *openMPState) resetPlayerMoney(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	player.money = 0

	return 1, nil
}

func (state *openMPState) givePlayerWeapon(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	weapon := int(params[1])
	player.weapons[weapon] += int(params[2])
	player.armedWeapon = weapon

	return 1, nil
}

func (state *openMPState) removePlayerWeapon(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	weapon := int(params[1])
	delete(player.weapons, weapon)

	if player.armedWeapon == weapon {
		player.armedWeapon = 0
	}

	return 1, nil
}

func (state *openMPState) resetPlayerWeapons(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	player.weapons = map[int]int{}
	player.armedWeapon = 0

	return 1, nil
}

func (state *openMPState) setPlayerAmmo(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	player.weapons[int(params[1])] = int(params[2])

	return 1, nil
}

func (state *openMPState) setPlayerArmedWeapon(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	player.armedWeapon = int(params[1])

	return 1, nil
}

func (state *openMPState) getPlayerWeapon(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(player.armedWeapon), nil
}

func (state *openMPState) getPlayerAmmo(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(player.weapons[player.armedWeapon]), nil
}

func (state *openMPState) getPlayerWeaponData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	weapon := weaponForSlot(player, int(params[1]))
	if err := ctx.WriteCell(params[2], backend.Cell(weapon)); err != nil {
		return 0, err
	}

	return 1, ctx.WriteCell(params[3], backend.Cell(player.weapons[weapon]))
}

func weaponForSlot(player *testPlayer, slot int) int {
	for weapon := range player.weapons {
		if weaponSlot(weapon) == slot {
			return weapon
		}
	}

	return 0
}

func weaponSlot(weapon int) int {
	switch {
	case weapon == 0 || weapon == 1:
		return 0
	case weapon <= 9:
		return 1
	case weapon <= 15:
		return 10
	case weapon <= 18:
		return 8
	case weapon <= 24:
		return 2
	case weapon <= 27:
		return 3
	case weapon <= 29:
		return 4
	case weapon <= 32:
		return 5
	case weapon <= 34:
		return 6
	case weapon <= 38:
		return 7
	case weapon <= 40:
		return 8
	case weapon <= 43:
		return 9
	case weapon <= 46:
		return 11
	}

	return 12
}
