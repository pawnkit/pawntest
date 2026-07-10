package runner

import (
	"math"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *serverState) setGameModeText(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	value, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	state.gameModeText = value

	return 1, nil
}

func (state *serverState) gameModeExit(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	state.exited = true

	return 1, nil
}

func (state *serverState) setServerRule(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	value, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	state.rules[name] = value

	return 1, nil
}

func (state *serverState) removeServerRule(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	if _, ok := state.rules[name]; !ok {
		return 0, nil
	}

	delete(state.rules, name)

	return 1, nil
}

func (state *serverState) isValidServerRule(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	if _, ok := state.rules[name]; ok {
		return 1, nil
	}

	return 0, nil
}

func (state *serverState) allowNicknameCharacter(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	state.nicknameCharacters[rune(params[0])] = params[1] != 0

	return 1, nil
}

func (state *serverState) isNicknameCharacterAllowed(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 && state.nicknameCharacters[rune(params[0])] {
		return 1, nil
	}

	return 0, nil
}

func (state *serverState) isValidNickname(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}

	if name == "" || len(name) > 24 || strings.TrimSpace(name) != name {
		return 0, nil
	}

	for _, value := range name {
		if !state.nicknameCharacters[value] {
			return 0, nil
		}
	}

	return 1, nil
}

func (state *serverState) disableInteriorExits(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	state.interiorExits = false

	return 1, nil
}

func (state *serverState) disableNameTagLOS(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	state.nameTagLOS = false

	return 1, nil
}

func (state *serverState) enablePedAnimations(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	state.pedAnimations = true

	return 1, nil
}

func (state *serverState) getMaxPlayers(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return 1000, nil
}

func (state *serverState) getPlayerPoolSize(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return highestID(state.players.players), nil
}

func (state *serverState) getVehiclePoolSize(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return highestID(state.vehicles.vehicles), nil
}

func (state *serverState) getActorPoolSize(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return highestID(state.actors.actors), nil
}

func highestID[T any](values map[int]T) backend.Cell {
	highest := -1
	for id := range values {
		if id > highest {
			highest = id
		}
	}

	return backend.Cell(highest)
}

func (state *serverState) getTickCount(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return backend.Cell(state.scheduler.now), nil
}

func (state *serverState) getServerTickRate(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
	return 50, nil
}

func vectorSize(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 3 {
		return 0, nil
	}

	x, y, z := cellFloat(params[0]), cellFloat(params[1]), cellFloat(params[2])

	return floatCell(float32(math.Sqrt(float64(x*x + y*y + z*z)))), nil
}

func getWeaponSlot(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return -1, nil
	}

	return backend.Cell(weaponSlot(int(params[0]))), nil
}
