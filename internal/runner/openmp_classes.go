package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type classWeapon struct {
	weapon, ammo int
}

type playerClass struct {
	team, skin int
	position   [3]float32
	angle      float32
	weapons    [3]classWeapon
}

type classState struct {
	classes   []playerClass
	spawnInfo map[int]playerClass
	selected  map[int]int
	selecting map[int]bool
	players   *openMPState
}

func newClassState() *classState {
	return &classState{spawnInfo: map[int]playerClass{}, selected: map[int]int{}, selecting: map[int]bool{}}
}

func (state *classState) Clone() scenarioModule {
	clone := newClassState()

	clone.classes = append([]playerClass(nil), state.classes...)
	maps.Copy(clone.spawnInfo, state.spawnInfo)
	maps.Copy(clone.selected, state.selected)
	maps.Copy(clone.selecting, state.selecting)

	return clone
}

func (state *classState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *classState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_class_select":           state.selectClass,
		"__pt_class_count":            assertScenarioCount(result, "classes", func() int { return len(state.classes) }),
		"__pt_player_class":           state.assertPlayerClass(result),
		"__pt_player_selecting_class": state.assertSelecting(result),
		"AddPlayerClass":              state.addPlayerClass,
		"AddPlayerClassEx":            state.addPlayerClassEx,
		"SetSpawnInfo":                state.setSpawnInfo,
		"GetSpawnInfo":                state.getSpawnInfo,
		"SpawnPlayer":                 state.spawnPlayer,
		"ForceClassSelection":         state.forceClassSelection,
		"GetAvailableClasses":         state.getAvailableClasses,
		"GetPlayerClass":              state.getPlayerClass,
		"EditPlayerClass":             state.editPlayerClass,
	}
}

func (state *classState) addPlayerClass(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	class, ok := parsePlayerClass(params, false)
	if !ok {
		return -1, nil
	}

	state.classes = append(state.classes, class)

	return backend.Cell(len(state.classes) - 1), nil
}

func (state *classState) addPlayerClassEx(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	class, ok := parsePlayerClass(params, true)
	if !ok {
		return -1, nil
	}

	state.classes = append(state.classes, class)

	return backend.Cell(len(state.classes) - 1), nil
}

func parsePlayerClass(params []backend.Cell, hasTeam bool) (playerClass, bool) {
	minimum := 5

	offset := 0
	if hasTeam {
		minimum, offset = 6, 1
	}

	if len(params) < minimum {
		return playerClass{}, false
	}

	class := playerClass{team: -1, skin: int(params[offset]), position: cellsToVector(params[offset+1 : offset+4]), angle: cellFloat(params[offset+4])}
	if hasTeam {
		class.team = int(params[0])
	}

	for slot := range 3 {
		index := offset + 5 + slot*2
		if len(params) > index+1 {
			class.weapons[slot] = classWeapon{weapon: int(params[index]), ammo: int(params[index+1])}
		}
	}

	return class, true
}

func (state *classState) setSpawnInfo(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	class, ok := parsePlayerClass(params[1:], true)
	if !ok {
		return 0, nil
	}

	state.spawnInfo[int(params[0])] = class

	return 1, nil
}

func (state *classState) getSpawnInfo(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 13 {
		return 0, nil
	}

	class, ok := state.spawnInfo[int(params[0])]
	if !ok {
		return 0, nil
	}

	return writePlayerClass(ctx, params[1:13], class)
}

func (state *classState) getPlayerClass(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 13 || params[0] < 0 || int(params[0]) >= len(state.classes) {
		return 0, nil
	}

	return writePlayerClass(ctx, params[1:13], state.classes[int(params[0])])
}

func writePlayerClass(ctx backend.NativeContext, addresses []backend.Cell, class playerClass) (backend.Cell, error) {
	values := []backend.Cell{
		backend.Cell(class.team), backend.Cell(class.skin), floatCell(class.position[0]), floatCell(class.position[1]),
		floatCell(class.position[2]), floatCell(class.angle), backend.Cell(class.weapons[0].weapon), backend.Cell(class.weapons[0].ammo),
		backend.Cell(class.weapons[1].weapon), backend.Cell(class.weapons[1].ammo), backend.Cell(class.weapons[2].weapon), backend.Cell(class.weapons[2].ammo),
	}
	for index, value := range values {
		if err := ctx.WriteCell(addresses[index], value); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *classState) editPlayerClass(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 || params[0] < 0 || int(params[0]) >= len(state.classes) {
		return 0, nil
	}

	class, ok := parsePlayerClass(params[1:], true)
	if !ok {
		return 0, nil
	}

	state.classes[int(params[0])] = class

	return 1, nil
}

func (state *classState) getAvailableClasses(_ backend.NativeContext, _ []backend.Cell) (backend.Cell, error) {
	return backend.Cell(len(state.classes)), nil
}

func (state *classState) forceClassSelection(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	state.selecting[int(params[0])] = true
	state.players.players[int(params[0])].spawned = false

	return 1, nil
}

func (state *classState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}
