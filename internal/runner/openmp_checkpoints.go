package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type checkpoint struct {
	active   bool
	position [3]float32
	radius   float32
}

type raceCheckpoint struct {
	checkpoint
	checkpointType int
	next           [3]float32
}

type checkpointState struct {
	checkpoints map[int]checkpoint
	races       map[int]raceCheckpoint
	players     *openMPState
}

func newCheckpointState() *checkpointState {
	return &checkpointState{checkpoints: map[int]checkpoint{}, races: map[int]raceCheckpoint{}}
}

func (state *checkpointState) Clone() scenarioModule {
	clone := newCheckpointState()
	maps.Copy(clone.checkpoints, state.checkpoints)
	maps.Copy(clone.races, state.races)

	return clone
}

func (state *checkpointState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *checkpointState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_checkpoint_active":       state.assertCheckpointActive(result),
		"__pt_checkpoint_inside":       state.assertCheckpointInside(result),
		"__pt_race_checkpoint_active":  state.assertRaceCheckpointActive(result),
		"__pt_race_checkpoint_inside":  state.assertRaceCheckpointInside(result),
		"SetPlayerCheckpoint":          state.setPlayerCheckpoint,
		"DisablePlayerCheckpoint":      state.disablePlayerCheckpoint,
		"IsPlayerInCheckpoint":         state.isPlayerInCheckpoint,
		"IsPlayerCheckpointActive":     state.isPlayerCheckpointActive,
		"GetPlayerCheckpoint":          state.getPlayerCheckpoint,
		"SetPlayerRaceCheckpoint":      state.setPlayerRaceCheckpoint,
		"DisablePlayerRaceCheckpoint":  state.disablePlayerRaceCheckpoint,
		"IsPlayerInRaceCheckpoint":     state.isPlayerInRaceCheckpoint,
		"IsPlayerRaceCheckpointActive": state.isPlayerRaceCheckpointActive,
		"GetPlayerRaceCheckpoint":      state.getPlayerRaceCheckpoint,
	}
}

func (state *checkpointState) setPlayerCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	state.checkpoints[int(params[0])] = checkpoint{active: true, position: cellsToVector(params[1:4]), radius: cellFloat(params[4])}

	return 1, nil
}

func (state *checkpointState) disablePlayerCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	value := state.checkpoints[int(params[0])]
	value.active = false
	state.checkpoints[int(params[0])] = value

	return 1, nil
}

func (state *checkpointState) isPlayerInCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if state.playerInside(params, state.checkpoints) {
		return 1, nil
	}

	return 0, nil
}

func (state *checkpointState) isPlayerCheckpointActive(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 && state.checkpoints[int(params[0])].active {
		return 1, nil
	}

	return 0, nil
}

func (state *checkpointState) getPlayerCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return 0, nil
	}

	value, ok := state.checkpoints[int(params[0])]
	if !ok {
		return 0, nil
	}

	if result, err := writeFloatVector(ctx, params[1:4], value.position); err != nil || result == 0 {
		return result, err
	}

	return 1, ctx.WriteCell(params[4], floatCell(value.radius))
}

func (state *checkpointState) setPlayerRaceCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 9 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	state.races[int(params[0])] = raceCheckpoint{
		checkpoint:     checkpoint{active: true, position: cellsToVector(params[2:5]), radius: cellFloat(params[8])},
		checkpointType: int(params[1]), next: cellsToVector(params[5:8]),
	}

	return 1, nil
}

func (state *checkpointState) disablePlayerRaceCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	value := state.races[int(params[0])]
	value.active = false
	state.races[int(params[0])] = value

	return 1, nil
}

func (state *checkpointState) isPlayerInRaceCheckpoint(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	checkpoints := make(map[int]checkpoint, len(state.races))
	for playerID, value := range state.races {
		checkpoints[playerID] = value.checkpoint
	}

	if state.playerInside(params, checkpoints) {
		return 1, nil
	}

	return 0, nil
}

func (state *checkpointState) isPlayerRaceCheckpointActive(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) > 0 && state.races[int(params[0])].active {
		return 1, nil
	}

	return 0, nil
}

func (state *checkpointState) getPlayerRaceCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 8 {
		return 0, nil
	}

	value, ok := state.races[int(params[0])]
	if !ok {
		return 0, nil
	}

	if result, err := writeFloatVector(ctx, params[1:4], value.position); err != nil || result == 0 {
		return result, err
	}

	if result, err := writeFloatVector(ctx, params[4:7], value.next); err != nil || result == 0 {
		return result, err
	}

	return 1, ctx.WriteCell(params[7], floatCell(value.radius))
}

func (state *checkpointState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}

func (state *checkpointState) playerInside(params []backend.Cell, checkpoints map[int]checkpoint) bool {
	if len(params) == 0 || state.players == nil {
		return false
	}

	player, playerOK := state.players.player(params[0])

	value, checkpointOK := checkpoints[int(params[0])]
	if !playerOK || !checkpointOK || !value.active {
		return false
	}

	return playerDistance(player, floatVectorCells(value.position)) <= value.radius
}

func floatVectorCells(vector [3]float32) []backend.Cell {
	return []backend.Cell{floatCell(vector[0]), floatCell(vector[1]), floatCell(vector[2])}
}
