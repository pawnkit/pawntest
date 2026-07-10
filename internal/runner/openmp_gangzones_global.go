package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *gangZoneState) destroyGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.zone(params); !ok {
		return 0, nil
	}

	delete(state.zones, int(params[0]))

	return 1, nil
}

func (state *gangZoneState) isValidGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.zone(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) showGangZoneForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.visible, view.colour = true, int(params[2])
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) showGangZoneForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	for id, player := range state.players.players {
		if player.connected {
			view := zone.views[id]
			view.visible, view.colour = true, int(params[1])
			zone.views[id] = view
		}
	}

	return 1, nil
}

func (state *gangZoneState) hideGangZoneForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.visible, view.flashing = false, false
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) hideGangZoneForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok {
		return 0, nil
	}

	clear(zone.views)

	return 1, nil
}

func (state *gangZoneState) flashGangZoneForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.flashing, view.flash = true, int(params[2])
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) flashGangZoneForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	for id, player := range state.players.players {
		if player.connected {
			view := zone.views[id]
			view.flashing, view.flash = true, int(params[1])
			zone.views[id] = view
		}
	}

	return 1, nil
}

func (state *gangZoneState) stopGangZoneFlashForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok {
		return 0, nil
	}

	view := zone.views[int(params[0])]
	view.flashing = false
	zone.views[int(params[0])] = view

	return 1, nil
}

func (state *gangZoneState) stopGangZoneFlashForAll(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok {
		return 0, nil
	}

	for id, view := range zone.views {
		view.flashing = false
		zone.views[id] = view
	}

	return 1, nil
}

func (state *gangZoneState) isPlayerInGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if ok && state.playerInZone(int(params[0]), zone) {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) isGangZoneVisibleForPlayer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if ok && zone.views[int(params[0])].visible {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) getGangZoneColour(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(zone.views[int(params[0])].colour), nil
}

func (state *gangZoneState) getGangZoneFlashColour(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if !ok {
		return 0, nil
	}

	return backend.Cell(zone.views[int(params[0])].flash), nil
}

func (state *gangZoneState) isGangZoneFlashing(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.playerZoneArgument(params)
	if ok && zone.views[int(params[0])].flashing {
		return 1, nil
	}

	return 0, nil
}

func (state *gangZoneState) getGangZonePosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	return writeGangZoneBounds(ctx, params[1:5], zone.bounds)
}

func (state *gangZoneState) useGangZoneCheck(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	zone, ok := state.zone(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	zone.check = params[1] != 0

	return 1, nil
}

func (state *gangZoneState) playerZoneArgument(params []backend.Cell) (*testGangZone, bool) {
	if len(params) < 2 || state.players == nil {
		return nil, false
	}

	if _, ok := state.players.player(params[0]); !ok {
		return nil, false
	}

	zone, ok := state.zones[int(params[1])]

	return zone, ok
}

func (state *gangZoneState) playerInZone(playerID int, zone *testGangZone) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.players[playerID]
	if !ok || !player.connected {
		return false
	}

	return player.x >= zone.bounds[0] && player.x <= zone.bounds[2] && player.y >= zone.bounds[1] && player.y <= zone.bounds[3]
}

func writeGangZoneBounds(ctx backend.NativeContext, addresses []backend.Cell, bounds [4]float32) (backend.Cell, error) {
	for index, value := range bounds {
		if err := ctx.WriteCell(addresses[index], floatCell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}
