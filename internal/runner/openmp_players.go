package runner

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testPlayer struct {
	name                       string
	connected, spawned         bool
	controllable, ghost, clock bool
	money, score, team, colour int
	skin, interior, world      int
	wanted, weather, drunk     int
	fightingStyle, action      int
	hour, minute               int
	health, armour, gravity    float32
	x, y, z, angle             float32
	velocity                   [3]float32
	camera                     [3]float32
	keys                       [3]backend.Cell
	weapons                    map[int]int
	armedWeapon                int
	vehicle, seat              int
	messages                   []string
}

type openMPState struct {
	nextPlayer int
	players    map[int]*testPlayer
}

func newOpenMPState() *openMPState {
	return &openMPState{players: map[int]*testPlayer{}}
}

func (state *openMPState) clone() *openMPState {
	clone := newOpenMPState()

	clone.nextPlayer = state.nextPlayer
	for id, player := range state.players {
		playerCopy := *player
		playerCopy.messages = append([]string(nil), player.messages...)
		playerCopy.weapons = make(map[int]int, len(player.weapons))
		maps.Copy(playerCopy.weapons, player.weapons)
		clone.players[id] = &playerCopy
	}

	return clone
}

func (state *openMPState) Clone() scenarioModule { return state.clone() }

func (state *openMPState) Register(vm backend.VM, context *executionContext) error {
	return registerOpenMPPlayerNatives(vm, context.state, state, context.mocks, context.allowUnknown)
}

func registerOpenMPPlayerNatives(vm backend.VM, nativeState *nativeState, state *openMPState, mocks *mockState, allowUnknown bool) error {
	natives := map[string]backend.NativeFunc{
		"__pt_player_create":     state.createPlayer,
		"__pt_player_connected":  state.assertConnected(nativeState),
		"__pt_player_money":      state.assertMoney(nativeState),
		"__pt_player_message":    state.assertMessage(nativeState),
		"__pt_player_pos_near":   state.assertPosition(nativeState),
		"IsPlayerConnected":      state.isPlayerConnected,
		"GetPlayerName":          state.getPlayerName,
		"SetPlayerPos":           state.setPlayerPos,
		"GetPlayerPos":           state.getPlayerPos,
		"SetPlayerMoney":         state.setPlayerMoney,
		"GivePlayerMoney":        state.givePlayerMoney,
		"GetPlayerMoney":         state.getPlayerMoney,
		"SendClientMessage":      state.sendClientMessage,
		"SendClientMessageToAll": state.sendClientMessageToAll,
		"Kick":                   state.kick,
	}
	maps.Copy(natives, state.coreNatives())
	maps.Copy(natives, state.equipmentNatives())
	maps.Copy(natives, state.cameraNatives())

	return registerScenarioNatives(vm, natives, mocks, allowUnknown)
}

func (state *openMPState) createPlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 1 {
		return -1, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	id := state.nextPlayer
	state.nextPlayer++
	state.players[id] = &testPlayer{
		name:         name,
		connected:    true,
		spawned:      true,
		controllable: true,
		health:       100,
		gravity:      0.008,
		weapons:      map[int]int{},
		vehicle:      -1,
		seat:         -1,
	}

	return backend.Cell(id), nil
}

func (state *openMPState) player(id backend.Cell) (*testPlayer, bool) {
	player, ok := state.players[int(id)]
	return player, ok
}

func (state *openMPState) isPlayerConnected(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if ok && player.connected {
		return 1, nil
	}

	return 0, nil
}

func (state *openMPState) getPlayerName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return 1, ctx.WriteString(params[1], truncateString(player.name, int(params[2])))
}

func (state *openMPState) setPlayerPos(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	player.x, player.y, player.z = cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])

	return 1, nil
}

func (state *openMPState) getPlayerPos(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	for index, value := range []float32{player.x, player.y, player.z} {
		if err := ctx.WriteCell(params[index+1], floatCell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *openMPState) setPlayerMoney(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	player.money = int(params[1])

	return 1, nil
}

func (state *openMPState) givePlayerMoney(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	player.money += int(params[1])

	return 1, nil
}

func (state *openMPState) getPlayerMoney(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(player.money), nil
}

func (state *openMPState) sendClientMessage(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	message, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}

	player.messages = append(player.messages, message)

	return 1, nil
}

func (state *openMPState) sendClientMessageToAll(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}

	message, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	for _, player := range state.players {
		if player.connected {
			player.messages = append(player.messages, message)
		}
	}

	return 1, nil
}

func (state *openMPState) kick(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.paramPlayer(params)
	if !ok {
		return 0, nil
	}

	player.connected = false

	return 1, nil
}

func (state *openMPState) assertConnected(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player connected assertion expects 4 arguments")
		}

		player, exists := state.paramPlayer(params)
		actual := exists && player.connected

		expected := len(params) >= 2 && params[1] != 0
		if actual != expected {
			setFailure(result, params, 2, fmt.Sprintf("player %d connected state: expected %t, got %t", params[0], expected, actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *openMPState) assertMoney(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("player money assertion expects 4 arguments")
		}

		player, ok := state.paramPlayer(params)
		if !ok || player.money != int(params[1]) {
			actual := 0
			if ok {
				actual = player.money
			}

			setFailure(result, params, 2, fmt.Sprintf("player %d money: expected %d, got %d", params[0], params[1], actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *openMPState) assertMessage(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("player message assertion expects 5 arguments")
		}

		player, ok := state.paramPlayer(params)
		expected := readStringParam(ctx, params, 1)
		contains := len(params) >= 3 && params[2] != 0

		if ok {
			for _, message := range player.messages {
				if message == expected || contains && strings.Contains(message, expected) {
					return 1, nil
				}
			}
		}

		messages := []string(nil)
		if ok {
			messages = player.messages
		}

		setFailure(result, params, 3, fmt.Sprintf("player %d message: expected %q; recorded %q", params[0], expected, messages), ctx)

		return 0, nil
	}
}

func (state *openMPState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("player position assertion expects 7 arguments")
		}

		player, ok := state.paramPlayer(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("player %d does not exist", params[0]), ctx)
			return 0, nil
		}

		x, y, z, tolerance := cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3]), float32(math.Abs(float64(cellFloat(params[4]))))
		if absFloat(player.x-x) > tolerance || absFloat(player.y-y) > tolerance || absFloat(player.z-z) > tolerance {
			setFailure(result, params, 5, fmt.Sprintf("player %d position: expected (%g, %g, %g) +/- %g, got (%g, %g, %g)", params[0], x, y, z, tolerance, player.x, player.y, player.z), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *openMPState) paramPlayer(params []backend.Cell) (*testPlayer, bool) {
	if len(params) < 1 {
		return nil, false
	}

	return state.player(params[0])
}

func cellFloat(value backend.Cell) float32 { return math.Float32frombits(uint32(value)) }
func floatCell(value float32) backend.Cell { return backend.Cell(int32(math.Float32bits(value))) }
func absFloat(value float32) float32 {
	if value < 0 {
		return -value
	}

	return value
}
