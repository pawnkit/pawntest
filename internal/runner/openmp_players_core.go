package runner

import (
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *openMPState) coreNatives() map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"SetPlayerName":              state.setPlayerName,
		"SetPlayerHealth":            state.setPlayerFloat(playerHealth),
		"GetPlayerHealth":            state.getPlayerFloat(playerHealth),
		"SetPlayerArmour":            state.setPlayerFloat(playerArmour),
		"GetPlayerArmour":            state.getPlayerFloat(playerArmour),
		"SetPlayerSkin":              state.setPlayerInt(playerSkin),
		"GetPlayerSkin":              state.getPlayerInt(playerSkin),
		"SetPlayerTeam":              state.setPlayerInt(playerTeam),
		"GetPlayerTeam":              state.getPlayerInt(playerTeam),
		"SetPlayerScore":             state.setPlayerInt(playerScore),
		"GetPlayerScore":             state.getPlayerInt(playerScore),
		"SetPlayerColor":             state.setPlayerInt(playerColour),
		"GetPlayerColor":             state.getPlayerInt(playerColour),
		"SetPlayerInterior":          state.setPlayerInt(playerInterior),
		"GetPlayerInterior":          state.getPlayerInt(playerInterior),
		"SetPlayerVirtualWorld":      state.setPlayerInt(playerWorld),
		"GetPlayerVirtualWorld":      state.getPlayerInt(playerWorld),
		"SetPlayerWantedLevel":       state.setPlayerInt(playerWanted),
		"GetPlayerWantedLevel":       state.getPlayerInt(playerWanted),
		"SetPlayerWeather":           state.setPlayerInt(playerWeather),
		"GetPlayerWeather":           state.getPlayerInt(playerWeather),
		"SetPlayerDrunkLevel":        state.setPlayerInt(playerDrunk),
		"GetPlayerDrunkLevel":        state.getPlayerInt(playerDrunk),
		"SetPlayerFightingStyle":     state.setPlayerInt(playerFightingStyle),
		"GetPlayerFightingStyle":     state.getPlayerInt(playerFightingStyle),
		"SetPlayerSpecialAction":     state.setPlayerInt(playerAction),
		"GetPlayerSpecialAction":     state.getPlayerInt(playerAction),
		"SetPlayerGravity":           state.setPlayerFloat(playerGravity),
		"GetPlayerGravity":           state.getPlayerFloatValue(playerGravity),
		"SetPlayerFacingAngle":       state.setPlayerFloat(playerAngle),
		"GetPlayerFacingAngle":       state.getPlayerFloat(playerAngle),
		"SetPlayerVelocity":          state.setPlayerVector(playerVelocity),
		"GetPlayerVelocity":          state.getPlayerVector(playerVelocity),
		"SetPlayerTime":              state.setPlayerTime,
		"GetPlayerTime":              state.getPlayerTime,
		"TogglePlayerClock":          state.setPlayerBool(playerClock),
		"PlayerHasClockEnabled":      state.getPlayerBool(playerClock),
		"TogglePlayerControllable":   state.setPlayerBool(playerControllable),
		"IsPlayerControllable":       state.getPlayerBool(playerControllable),
		"TogglePlayerGhostMode":      state.setPlayerBool(playerGhost),
		"GetPlayerGhostMode":         state.getPlayerBool(playerGhost),
		"IsPlayerSpawned":            state.getPlayerBool(playerSpawned),
		"IsPlayerNPC":                state.falseForPlayer,
		"GetPlayerState":             state.getPlayerState,
		"GetPlayerKeys":              state.getPlayerKeys,
		"GetPlayerDistanceFromPoint": state.getPlayerDistanceFromPoint,
		"IsPlayerInRangeOfPoint":     state.isPlayerInRangeOfPoint,
	}
}

func (state *openMPState) setPlayerName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	player.name = name

	return 1, nil
}

func (state *openMPState) setPlayerInt(field func(*testPlayer) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(player) = int(params[1])

		return 1, nil
	}
}

func (state *openMPState) getPlayerInt(field func(*testPlayer) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok {
			return 0, nil
		}

		return backend.Cell(*field(player)), nil
	}
}

func (state *openMPState) setPlayerFloat(field func(*testPlayer) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(player) = cellFloat(params[1])

		return 1, nil
	}
}

func (state *openMPState) getPlayerFloat(field func(*testPlayer) *float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return 1, ctx.WriteCell(params[1], floatCell(*field(player)))
	}
}

func (state *openMPState) getPlayerFloatValue(field func(*testPlayer) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok {
			return 0, nil
		}

		return floatCell(*field(player)), nil
	}
}

func (state *openMPState) setPlayerBool(field func(*testPlayer) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(player) = params[1] != 0

		return 1, nil
	}
}

func (state *openMPState) getPlayerBool(field func(*testPlayer) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if ok && *field(player) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *openMPState) setPlayerVector(field func(*testPlayer) *[3]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		*field(player) = [3]float32{cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}

		return 1, nil
	}
}

func (state *openMPState) getPlayerVector(field func(*testPlayer) *[3]float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		player, ok := state.paramPlayer(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		return writeFloatVector(ctx, params[1:4], *field(player))
	}
}

func (state *openMPState) setPlayerTime(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	player.hour, player.minute = int(params[1]), int(params[2])

	return 1, nil
}

func (state *openMPState) getPlayerTime(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	if err := ctx.WriteCell(params[1], backend.Cell(player.hour)); err != nil {
		return 0, err
	}

	return 1, ctx.WriteCell(params[2], backend.Cell(player.minute))
}

func (state *openMPState) getPlayerKeys(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	for index, value := range player.keys {
		if err := ctx.WriteCell(params[index+1], value); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *openMPState) getPlayerState(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || !player.connected {
		return 0, nil
	}

	if player.spawned {
		return 1, nil
	}

	return 0, nil
}

func (state *openMPState) falseForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.paramPlayer(params); !ok {
		return 0, nil
	}

	return 0, nil
}

func (state *openMPState) getPlayerDistanceFromPoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	distance := playerDistance(player, params[1:4])

	return floatCell(distance), nil
}

func (state *openMPState) isPlayerInRangeOfPoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	if playerDistance(player, params[2:5]) <= cellFloat(params[1]) {
		return 1, nil
	}

	return 0, nil
}

func playerDistance(player *testPlayer, point []backend.Cell) float32 {
	dx := player.x - cellFloat(point[0])
	dy := player.y - cellFloat(point[1])
	dz := player.z - cellFloat(point[2])

	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func playerHealth(player *testPlayer) *float32      { return &player.health }
func playerArmour(player *testPlayer) *float32      { return &player.armour }
func playerGravity(player *testPlayer) *float32     { return &player.gravity }
func playerAngle(player *testPlayer) *float32       { return &player.angle }
func playerSkin(player *testPlayer) *int            { return &player.skin }
func playerTeam(player *testPlayer) *int            { return &player.team }
func playerScore(player *testPlayer) *int           { return &player.score }
func playerColour(player *testPlayer) *int          { return &player.colour }
func playerInterior(player *testPlayer) *int        { return &player.interior }
func playerWorld(player *testPlayer) *int           { return &player.world }
func playerWanted(player *testPlayer) *int          { return &player.wanted }
func playerWeather(player *testPlayer) *int         { return &player.weather }
func playerDrunk(player *testPlayer) *int           { return &player.drunk }
func playerFightingStyle(player *testPlayer) *int   { return &player.fightingStyle }
func playerAction(player *testPlayer) *int          { return &player.action }
func playerClock(player *testPlayer) *bool          { return &player.clock }
func playerControllable(player *testPlayer) *bool   { return &player.controllable }
func playerGhost(player *testPlayer) *bool          { return &player.ghost }
func playerSpawned(player *testPlayer) *bool        { return &player.spawned }
func playerVelocity(player *testPlayer) *[3]float32 { return &player.velocity }
func playerCamera(player *testPlayer) *[3]float32   { return &player.camera }
