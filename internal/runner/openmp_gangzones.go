package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type gangZoneView struct {
	visible, flashing bool
	colour, flash     int
}

type testGangZone struct {
	bounds [4]float32
	check  bool
	views  map[int]gangZoneView
}

type gangZoneState struct {
	next        int
	zones       map[int]*testGangZone
	playerNext  map[int]int
	playerZones map[int]map[int]*testGangZone
	players     *openMPState
}

func newGangZoneState() *gangZoneState {
	return &gangZoneState{zones: map[int]*testGangZone{}, playerNext: map[int]int{}, playerZones: map[int]map[int]*testGangZone{}}
}

func (state *gangZoneState) Clone() scenarioModule {
	clone := newGangZoneState()
	clone.next = state.next
	maps.Copy(clone.playerNext, state.playerNext)

	for id, zone := range state.zones {
		clone.zones[id] = cloneGangZone(zone)
	}

	for playerID, zones := range state.playerZones {
		clone.playerZones[playerID] = map[int]*testGangZone{}
		for id, zone := range zones {
			clone.playerZones[playerID][id] = cloneGangZone(zone)
		}
	}

	return clone
}

func cloneGangZone(zone *testGangZone) *testGangZone {
	clone := *zone
	clone.views = maps.Clone(zone.views)

	return &clone
}

func (state *gangZoneState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *gangZoneState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_gangzone_create":            state.createGangZone,
		"__pt_gangzone_valid":             state.assertValid(result, false),
		"__pt_gangzone_visible":           state.assertVisible(result, false),
		"__pt_gangzone_inside":            state.assertInside(result),
		"GangZoneCreate":                  state.createGangZone,
		"GangZoneDestroy":                 state.destroyGangZone,
		"IsValidGangZone":                 state.isValidGangZone,
		"GangZoneShowForPlayer":           state.showGangZoneForPlayer,
		"GangZoneShowForAll":              state.showGangZoneForAll,
		"GangZoneHideForPlayer":           state.hideGangZoneForPlayer,
		"GangZoneHideForAll":              state.hideGangZoneForAll,
		"GangZoneFlashForPlayer":          state.flashGangZoneForPlayer,
		"GangZoneFlashForAll":             state.flashGangZoneForAll,
		"GangZoneStopFlashForPlayer":      state.stopGangZoneFlashForPlayer,
		"GangZoneStopFlashForAll":         state.stopGangZoneFlashForAll,
		"IsPlayerInGangZone":              state.isPlayerInGangZone,
		"IsGangZoneVisibleForPlayer":      state.isGangZoneVisibleForPlayer,
		"GangZoneGetColourForPlayer":      state.getGangZoneColour,
		"GangZoneGetColorForPlayer":       state.getGangZoneColour,
		"GangZoneGetFlashColourForPlayer": state.getGangZoneFlashColour,
		"GangZoneGetFlashColorForPlayer":  state.getGangZoneFlashColour,
		"IsGangZoneFlashingForPlayer":     state.isGangZoneFlashing,
		"GangZoneGetPos":                  state.getGangZonePosition,
		"UseGangZoneCheck":                state.useGangZoneCheck,

		"__pt_player_gangzone_create":  state.createPlayerGangZone,
		"__pt_player_gangzone_valid":   state.assertValid(result, true),
		"__pt_player_gangzone_visible": state.assertVisible(result, true),
		"CreatePlayerGangZone":         state.createPlayerGangZone,
		"PlayerGangZoneDestroy":        state.destroyPlayerGangZone,
		"IsValidPlayerGangZone":        state.isValidPlayerGangZone,
		"PlayerGangZoneShow":           state.showPlayerGangZone,
		"PlayerGangZoneHide":           state.hidePlayerGangZone,
		"PlayerGangZoneFlash":          state.flashPlayerGangZone,
		"PlayerGangZoneStopFlash":      state.stopPlayerGangZoneFlash,
		"IsPlayerInPlayerGangZone":     state.isPlayerInPlayerGangZone,
		"IsPlayerGangZoneVisible":      state.isPlayerGangZoneVisible,
		"PlayerGangZoneGetColour":      state.getPlayerGangZoneColour,
		"PlayerGangZoneGetColor":       state.getPlayerGangZoneColour,
		"PlayerGangZoneGetFlashColour": state.getPlayerGangZoneFlashColour,
		"PlayerGangZoneGetFlashColor":  state.getPlayerGangZoneFlashColour,
		"IsPlayerGangZoneFlashing":     state.isPlayerGangZoneFlashing,
		"PlayerGangZoneGetPos":         state.getPlayerGangZonePosition,
		"UsePlayerGangZoneCheck":       state.usePlayerGangZoneCheck,
	}
}

func newTestGangZone(params []backend.Cell) *testGangZone {
	return &testGangZone{bounds: [4]float32{cellFloat(params[0]), cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}, views: map[int]gangZoneView{}}
}

func (state *gangZoneState) createGangZone(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return -1, nil
	}

	id := state.next
	state.next++
	state.zones[id] = newTestGangZone(params)

	return backend.Cell(id), nil
}

func (state *gangZoneState) zone(params []backend.Cell) (*testGangZone, bool) {
	if len(params) == 0 {
		return nil, false
	}

	zone, ok := state.zones[int(params[0])]

	return zone, ok
}
