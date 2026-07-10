package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *gangZoneState) createPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 || !state.validPlayer(params[0]) {
		return -1, nil
	}

	playerID := int(params[0])
	id := state.playerNext[playerID]

	state.playerNext[playerID]++
	if state.playerZones[playerID] == nil {
		state.playerZones[playerID] = map[int]*testGangZone{}
	}

	state.playerZones[playerID][id] = newTestGangZone(params[1:5])

	return backend.Cell(id), nil
}

func (state *gangZoneState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}

func (state *gangZoneState) playerZone(params []backend.Cell) (*testGangZone, bool) {
	if len(params) < 2 {
		return nil, false
	}

	zone, ok := state.playerZones[int(params[0])][int(params[1])]

	return zone, ok
}

func (state *gangZoneState) destroyPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerZone(params); !ok {
		return 0, nil
	}

	delete(state.playerZones[int(params[0])], int(params[1]))

	return 1, nil
}

func (state *gangZoneState) isValidPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.playerZone(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) showPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	zone.views[int(params[0])] = gangZoneView{visible: true, colour: int(params[2])}

	return 1, nil
}

func (state *gangZoneState) hidePlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.visible, view.flashing = false, false
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) flashPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.flashing, view.flash = true, int(params[2])
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) stopPlayerGangZoneFlash(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.flashing = false
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) isPlayerInPlayerGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if ok && state.playerInZone(int(params[0]), zone) {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) isPlayerGangZoneVisible(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if ok && zone.views[int(params[0])].visible {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) getPlayerGangZoneColour(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(zone.views[int(params[0])].colour), nil
}

func (state *gangZoneState) getPlayerGangZoneFlashColour(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(zone.views[int(params[0])].flash), nil
}

func (state *gangZoneState) isPlayerGangZoneFlashing(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if ok && zone.views[int(params[0])].flashing {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) getPlayerGangZonePosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok || len(params) < 6 {
		return 0, nil
	}

	return writeGangZoneBounds(ctx, params[2:6], zone.bounds)
}

func (state *gangZoneState) usePlayerGangZoneCheck(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZone(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	zone.check = params[2] != 0

	return 1, nil
}
