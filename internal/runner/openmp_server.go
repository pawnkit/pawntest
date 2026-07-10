package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type serverState struct {
	weather, worldTime             int
	gravity, restartTime           float32
	gameModeText                   string
	rules                          map[string]string
	nicknameCharacters             map[rune]bool
	chatRadius, markerRadius       float32
	nameTagDistance                float32
	interiorExits, nameTagLOS      bool
	showNameTags, pedAnimations    bool
	chatReplacement, stuntBonus    bool
	adminTeleport, interiorWeapons bool
	allAnimations, zoneNames       bool
	markersMode                    int
	exited                         bool
	players                        *openMPState
	vehicles                       *vehicleState
	actors                         *actorState
	scheduler                      *scheduler
}

func newServerState() *serverState {
	return &serverState{
		gravity: 0.008, restartTime: 12, rules: map[string]string{}, nicknameCharacters: defaultNicknameCharacters(),
		showNameTags: true, nameTagLOS: true, interiorExits: true, interiorWeapons: true,
	}
}

func defaultNicknameCharacters() map[rune]bool {
	allowed := map[rune]bool{'_': true, '[': true, ']': true, '(': true, ')': true, '$': true, '@': true, '.': true, '=': true}
	for value := '0'; value <= '9'; value++ {
		allowed[value] = true
	}

	for value := 'A'; value <= 'Z'; value++ {
		allowed[value], allowed[value+('a'-'A')] = true, true
	}

	return allowed
}

func (state *serverState) Clone() scenarioModule {
	clone := *state
	clone.rules = maps.Clone(state.rules)
	clone.nicknameCharacters = maps.Clone(state.nicknameCharacters)
	clone.players, clone.vehicles, clone.actors, clone.scheduler = nil, nil, nil, nil

	return &clone
}

func (state *serverState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()
	state.vehicles = context.scenarios.vehicleState()
	state.actors = context.scenarios.actorState()
	state.scheduler = context.scheduler

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *serverState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_server_weather": state.assertInt(result, "weather", func() int { return state.weather }),
		"__pt_server_time":    state.assertInt(result, "world time", func() int { return state.worldTime }),
		"__pt_server_gravity": state.assertGravity(result), "__pt_server_mode_text": state.assertModeText(result),
		"__pt_server_rule": state.assertRule(result), "SetWeather": state.setInt(&state.weather),
		"GetWeather": state.getInt(&state.weather), "SetWorldTime": state.setInt(&state.worldTime),
		"GetWorldTime": state.getInt(&state.worldTime), "SetGravity": state.setFloat(&state.gravity),
		"GetGravity": state.getFloat(&state.gravity), "SetModeRestartTime": state.setFloat(&state.restartTime),
		"GetModeRestartTime": state.getFloat(&state.restartTime), "SetGameModeText": state.setGameModeText,
		"GameModeExit": state.gameModeExit, "AddServerRule": state.setServerRule,
		"SetServerRule": state.setServerRule, "RemoveServerRule": state.removeServerRule,
		"IsValidServerRule": state.isValidServerRule, "AllowNickNameCharacter": state.allowNicknameCharacter,
		"IsNickNameCharacterAllowed": state.isNicknameCharacterAllowed, "IsValidNickName": state.isValidNickname,
		"DisableInteriorEnterExits": state.disableInteriorExits, "DisableNameTagLOS": state.disableNameTagLOS,
		"LimitGlobalChatRadius": state.setFloat(&state.chatRadius), "LimitPlayerMarkerRadius": state.setFloat(&state.markerRadius),
		"SetNameTagDrawDistance": state.setFloat(&state.nameTagDistance), "ShowNameTags": state.setBool(&state.showNameTags),
		"ShowPlayerMarkers": state.setInt(&state.markersMode), "UsePlayerPedAnims": state.enablePedAnimations,
		"ToggleChatTextReplacement": state.setBool(&state.chatReplacement), "ChatTextReplacementToggled": state.getBool(&state.chatReplacement),
		"EnableStuntBonusForAll": state.setBool(&state.stuntBonus), "AllowAdminTeleport": state.setBool(&state.adminTeleport),
		"IsAdminTeleportAllowed": state.getBool(&state.adminTeleport), "AllowInteriorWeapons": state.setBool(&state.interiorWeapons),
		"AreInteriorWeaponsAllowed": state.getBool(&state.interiorWeapons), "EnableAllAnimations": state.setBool(&state.allAnimations),
		"AreAllAnimationsEnabled": state.getBool(&state.allAnimations), "EnableZoneNames": state.setBool(&state.zoneNames),
		"GetMaxPlayers": state.getMaxPlayers, "GetPlayerPoolSize": state.getPlayerPoolSize,
		"GetVehiclePoolSize": state.getVehiclePoolSize, "GetActorPoolSize": state.getActorPoolSize,
		"GetTickCount": state.getTickCount, "GetServerTickRate": state.getServerTickRate,
		"VectorSize": vectorSize, "GetWeaponSlot": getWeaponSlot,
	}
}

func (state *serverState) setInt(field *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) == 0 {
			return 0, nil
		}

		*field = int(params[0])

		return 1, nil
	}
}

func (state *serverState) getInt(field *int) backend.NativeFunc {
	return func(backend.NativeContext, []backend.Cell) (backend.Cell, error) { return backend.Cell(*field), nil }
}

func (state *serverState) setFloat(field *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) == 0 {
			return 0, nil
		}

		*field = cellFloat(params[0])

		return 1, nil
	}
}

func (state *serverState) getFloat(field *float32) backend.NativeFunc {
	return func(backend.NativeContext, []backend.Cell) (backend.Cell, error) { return floatCell(*field), nil }
}

func (state *serverState) setBool(field *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) == 0 {
			return 0, nil
		}

		*field = params[0] != 0

		return 1, nil
	}
}

func (state *serverState) getBool(field *bool) backend.NativeFunc {
	return func(backend.NativeContext, []backend.Cell) (backend.Cell, error) { return boolCell(*field), nil }
}
