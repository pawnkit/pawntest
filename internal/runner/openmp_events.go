package runner

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

type eventState struct {
	players              *openMPState
	vehicles             *vehicleState
	objects              *objectState
	actors               *actorState
	pickups              *pickupState
	checkpoints          *checkpointState
	textDraws            *textDrawState
	gangZones            *gangZoneState
	checkpointInside     map[int]bool
	raceCheckpointInside map[int]bool
	playerStreams        map[eventPair]bool
	vehicleStreams       map[eventPair]bool
	actorStreams         map[eventPair]bool
	gangZoneInside       map[eventPair]bool
	playerGangZoneInside map[eventPair]bool
}

type eventPair struct{ subject, viewer int }

func newEventState() *eventState {
	return &eventState{
		checkpointInside: map[int]bool{}, raceCheckpointInside: map[int]bool{},
		playerStreams: map[eventPair]bool{}, vehicleStreams: map[eventPair]bool{}, actorStreams: map[eventPair]bool{},
		gangZoneInside: map[eventPair]bool{}, playerGangZoneInside: map[eventPair]bool{},
	}
}

func (state *eventState) Clone() scenarioModule {
	clone := newEventState()
	maps.Copy(clone.checkpointInside, state.checkpointInside)
	maps.Copy(clone.raceCheckpointInside, state.raceCheckpointInside)
	maps.Copy(clone.playerStreams, state.playerStreams)
	maps.Copy(clone.vehicleStreams, state.vehicleStreams)
	maps.Copy(clone.actorStreams, state.actorStreams)
	maps.Copy(clone.gangZoneInside, state.gangZoneInside)
	maps.Copy(clone.playerGangZoneInside, state.playerGangZoneInside)

	return clone
}

func (state *eventState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.Players
	state.vehicles = context.scenarios.Vehicles
	state.objects = context.scenarios.Objects
	state.actors = context.scenarios.Actors
	state.pickups = context.scenarios.Pickups
	state.checkpoints = context.scenarios.Checkpoints
	state.textDraws = context.scenarios.TextDraws
	state.gangZones = context.scenarios.GangZones

	return registerScenarioNatives(vm, state.natives(), context.mocks, context.allowUnknown)
}

func (state *eventState) natives() map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_event_connect":          state.connect,
		"__pt_event_disconnect":       state.disconnect,
		"__pt_event_spawn":            state.spawn,
		"__pt_event_death":            state.death,
		"__pt_event_text":             state.text,
		"__pt_event_command":          state.command,
		"__pt_event_keys":             state.keys,
		"__pt_event_enter_vehicle":    state.enterVehicle,
		"__pt_event_exit_vehicle":     state.exitVehicle,
		"__pt_event_pickup":           state.pickup,
		"__pt_event_player_pickup":    state.playerPickup,
		"__pt_event_enter_checkpoint": state.enterCheckpoint,
		"__pt_event_leave_checkpoint": state.leaveCheckpoint,
		"__pt_event_enter_race_cp":    state.enterRaceCheckpoint,
		"__pt_event_leave_race_cp":    state.leaveRaceCheckpoint,
		"__pt_event_click_textdraw":   state.clickTextDraw,
		"__pt_event_click_player_td":  state.clickPlayerTextDraw,
		"__pt_event_damage_player":    state.damagePlayer,
		"__pt_event_damage_actor":     state.damageActor,
		"__pt_event_weapon_shot":      state.weaponShot,
		"__pt_event_object_moved":     state.objectMoved,
		"__pt_event_player_obj_moved": state.playerObjectMoved,
		"__pt_event_player_stream":    state.playerStream,
		"__pt_event_vehicle_stream":   state.vehicleStream,
		"__pt_event_actor_stream":     state.actorStream,
		"__pt_event_move_player":      state.movePlayer,
		"__pt_event_vehicle_damage":   state.damageVehicle,
		"__pt_event_vehicle_status":   state.vehicleDamageStatus,
		"__pt_event_vehicle_respawn":  state.respawnVehicle,
		"__pt_event_cancel_textdraw":  state.cancelTextDraw,
		"__pt_event_click_player":     state.clickPlayer,
		"__pt_event_click_map":        state.clickMap,
	}
}

func (state *eventState) connect(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 1 {
		return -1, errors.New("connect event expects a name")
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	playerID := state.players.addPlayer(name, false)
	if _, err := callEvent(ctx, "OnPlayerConnect", playerID); err != nil {
		delete(state.players.players, int(playerID))
		return -1, err
	}

	return playerID, nil
}

func (state *eventState) disconnect(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 2)
	if !ok {
		return 0, nil
	}

	result, err := callEvent(ctx, "OnPlayerDisconnect", params[0], params[1])
	if err != nil {
		return 0, err
	}

	if player.vehicle >= 0 {
		vehicleID := player.vehicle
		player.vehicle, player.seat = -1, -1

		state.vehicles.refreshOccupied(vehicleID)
	}

	player.connected, player.spawned = false, false

	return result, nil
}

func (state *eventState) spawn(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 1)
	if !ok {
		return 0, nil
	}

	player.spawned = true

	return callEvent(ctx, "OnPlayerSpawn", params[0])
}

func (state *eventState) death(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 3)
	if !ok || !player.spawned {
		return 0, nil
	}

	if player.vehicle >= 0 {
		vehicleID := player.vehicle
		player.vehicle, player.seat = -1, -1

		state.vehicles.refreshOccupied(vehicleID)
	}

	player.spawned = false

	return callEvent(ctx, "OnPlayerDeath", params[0], params[1], params[2])
}

func (state *eventState) text(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerText", params[0], params[1])
}

func (state *eventState) command(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerCommandText", params[0], params[1])
}

func (state *eventState) keys(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 4)
	if !ok {
		return 0, nil
	}

	oldKeys := player.keys[0]
	player.keys = [3]backend.Cell{params[1], params[2], params[3]}

	return callEvent(ctx, "OnPlayerKeyStateChange", params[0], params[1], oldKeys)
}

func (state *eventState) enterVehicle(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 3)
	if !ok || !player.spawned || player.vehicle >= 0 || state.vehicles.vehicles[int(params[1])] == nil {
		return 0, nil
	}

	passenger := params[2] != 0
	if _, err := callEvent(ctx, "OnPlayerEnterVehicle", params[0], params[1], boolCell(passenger)); err != nil {
		return 0, err
	}

	seat := backend.Cell(0)

	newPlayerState := backend.Cell(2)
	if passenger {
		seat, newPlayerState = 1, 3
	}

	if _, err := state.vehicles.putPlayerInVehicle(ctx, []backend.Cell{params[0], params[1], seat}); err != nil {
		return 0, err
	}

	player.spawned = true

	return callEvent(ctx, "OnPlayerStateChange", params[0], newPlayerState, 1)
}

func (state *eventState) exitVehicle(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 1)
	if !ok || player.vehicle < 0 {
		return 0, nil
	}

	vehicleID := backend.Cell(player.vehicle)

	oldPlayerState := backend.Cell(3)
	if player.seat == 0 {
		oldPlayerState = 2
	}

	if _, err := callEvent(ctx, "OnPlayerExitVehicle", params[0], vehicleID); err != nil {
		return 0, err
	}

	if _, err := state.vehicles.removePlayerFromVehicle(ctx, params[:1]); err != nil {
		return 0, err
	}

	return callEvent(ctx, "OnPlayerStateChange", params[0], 1, oldPlayerState)
}

func (state *eventState) pickup(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 2)
	if !ok {
		return 0, nil
	}

	pickup := state.pickups.pickups[int(params[1])]
	if pickup == nil || pickup.world != player.world || pickup.hidden[int(params[0])] {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerPickUpPickup", params[0], params[1])
}

func (state *eventState) playerPickup(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok || state.pickups.playerPickups[int(params[0])][int(params[1])] == nil {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerPickUpPlayerPickup", params[0], params[1])
}

func (state *eventState) enterCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 1); !ok {
		return 0, nil
	}

	playerID := int(params[0])
	if state.checkpointInside[playerID] || !state.checkpoints.playerInside(params, state.checkpoints.checkpoints) {
		return 0, nil
	}

	state.checkpointInside[playerID] = true

	return callEvent(ctx, "OnPlayerEnterCheckpoint", params[0])
}

func (state *eventState) leaveCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 1); !ok {
		return 0, nil
	}

	playerID := int(params[0])
	if !state.checkpointInside[playerID] || !state.activeCheckpoint(params, false) || state.checkpoints.playerInside(params, state.checkpoints.checkpoints) {
		return 0, nil
	}

	state.checkpointInside[playerID] = false

	return callEvent(ctx, "OnPlayerLeaveCheckpoint", params[0])
}

func (state *eventState) enterRaceCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 1); !ok {
		return 0, nil
	}

	playerID := int(params[0])
	if state.raceCheckpointInside[playerID] || !state.activeCheckpoint(params, true) || !state.insideRaceCheckpoint(params) {
		return 0, nil
	}

	state.raceCheckpointInside[playerID] = true

	return callEvent(ctx, "OnPlayerEnterRaceCheckpoint", params[0])
}

func (state *eventState) leaveRaceCheckpoint(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 1); !ok {
		return 0, nil
	}

	playerID := int(params[0])
	if !state.raceCheckpointInside[playerID] || !state.activeCheckpoint(params, true) || state.insideRaceCheckpoint(params) {
		return 0, nil
	}

	state.raceCheckpointInside[playerID] = false

	return callEvent(ctx, "OnPlayerLeaveRaceCheckpoint", params[0])
}

func (state *eventState) clickTextDraw(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok {
		return 0, nil
	}

	draw := state.textDraws.draws[int(params[1])]

	selection := state.textDraws.selection[int(params[0])]
	if draw == nil || !draw.visible[int(params[0])] || !draw.selectable || !selection.active {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerClickTextDraw", params[0], params[1])
}

func (state *eventState) clickPlayerTextDraw(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok {
		return 0, nil
	}

	draw := state.textDraws.playerDraws[int(params[0])][int(params[1])]

	selection := state.textDraws.selection[int(params[0])]
	if draw == nil || !draw.visible[int(params[0])] || !draw.selectable || !selection.active {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerClickPlayerTextDraw", params[0], params[1])
}

func (state *eventState) damagePlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	victim, ok := state.player(params, 5)
	if !ok || !victim.spawned {
		return 0, nil
	}

	amount := cellFloat(params[2])
	if amount < 0 {
		return 0, nil
	}

	remaining := amount
	if victim.armour > 0 {
		absorbed := min(victim.armour, remaining)
		victim.armour -= absorbed
		remaining -= absorbed
	}

	victim.health = max(0, victim.health-remaining)

	if issuer, issuerOK := state.players.player(params[1]); issuerOK && issuer.connected {
		if _, err := callEvent(ctx, "OnPlayerGiveDamage", params[1], params[0], params[2], params[3], params[4]); err != nil {
			return 0, err
		}
	}

	result, err := callEvent(ctx, "OnPlayerTakeDamage", params[0], params[1], params[2], params[3], params[4])
	if err != nil {
		return 0, err
	}

	if victim.health > 0 {
		return result, nil
	}

	if victim.vehicle >= 0 {
		vehicleID := victim.vehicle
		victim.vehicle, victim.seat = -1, -1

		state.vehicles.refreshOccupied(vehicleID)
	}

	victim.spawned = false

	return callEvent(ctx, "OnPlayerDeath", params[0], params[1], params[3])
}

func (state *eventState) damageActor(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 5); !ok {
		return 0, nil
	}

	actor := state.actors.actors[int(params[1])]

	amount := cellFloat(params[2])
	if actor == nil || actor.invulnerable || amount < 0 {
		return 0, nil
	}

	actor.health = max(0, actor.health-amount)

	return callEvent(ctx, "OnPlayerGiveDamageActor", params[0], params[1], params[2], params[3], params[4])
}

func (state *eventState) weaponShot(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 7); !ok {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerWeaponShot", params[0], params[1], params[2], params[3], params[4], params[5], params[6])
}

func (state *eventState) objectMoved(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 1 {
		return 0, nil
	}

	object := state.objects.objects[int(params[0])]
	if object == nil || !object.moving {
		return 0, nil
	}

	object.position, object.rotation, object.moving = object.targetPos, object.targetRot, false
	object.moveID++

	return callEvent(ctx, "OnObjectMoved", params[0])
}

func (state *eventState) playerObjectMoved(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 2); !ok {
		return 0, nil
	}

	object := state.objects.playerObjects[int(params[0])][int(params[1])]
	if object == nil || !object.moving {
		return 0, nil
	}

	object.position, object.rotation, object.moving = object.targetPos, object.targetRot, false
	object.moveID++

	return callEvent(ctx, "OnPlayerObjectMoved", params[0], params[1])
}

func (state *eventState) playerStream(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 3 || params[0] == params[1] {
		return 0, nil
	}

	if _, ok := state.player(params, 3); !ok {
		return 0, nil
	}

	viewer, ok := state.players.player(params[1])
	if !ok || !viewer.connected || params[2] != 0 && viewer.world != state.players.players[int(params[0])].world {
		return 0, nil
	}

	pair := eventPair{subject: int(params[0]), viewer: int(params[1])}

	return state.stream(ctx, state.playerStreams, pair, params[2] != 0, "OnPlayerStreamIn", "OnPlayerStreamOut", params[0], params[1])
}

func (state *eventState) vehicleStream(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.entityStream(ctx, params, state.vehicleStreams, "OnVehicleStreamIn", "OnVehicleStreamOut", func(id int) (int, bool) {
		vehicle, ok := state.vehicles.vehicles[id]
		if !ok {
			return 0, false
		}

		return vehicle.world, true
	})
}

func (state *eventState) actorStream(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.entityStream(ctx, params, state.actorStreams, "OnActorStreamIn", "OnActorStreamOut", func(id int) (int, bool) {
		actor, ok := state.actors.actors[id]
		if !ok {
			return 0, false
		}

		return actor.world, true
	})
}

func (state *eventState) entityStream(
	ctx backend.NativeContext,
	params []backend.Cell,
	streams map[eventPair]bool,
	inCallback string,
	outCallback string,
	entityWorld func(int) (int, bool),
) (backend.Cell, error) {
	if len(params) < 3 {
		return 0, nil
	}

	world, exists := entityWorld(int(params[0]))

	viewer, viewerOK := state.players.player(params[1])
	if !exists || !viewerOK || !viewer.connected || params[2] != 0 && world != viewer.world {
		return 0, nil
	}

	pair := eventPair{subject: int(params[0]), viewer: int(params[1])}

	return state.stream(ctx, streams, pair, params[2] != 0, inCallback, outCallback, params[0], params[1])
}

func (*eventState) stream(
	ctx backend.NativeContext,
	streams map[eventPair]bool,
	pair eventPair,
	in bool,
	inCallback string,
	outCallback string,
	params ...backend.Cell,
) (backend.Cell, error) {
	if streams[pair] == in {
		return 0, nil
	}

	streams[pair] = in

	callback := outCallback
	if in {
		callback = inCallback
	}

	return callEvent(ctx, callback, params...)
}

func (state *eventState) movePlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	player, ok := state.player(params, 4)
	if !ok {
		return 0, nil
	}

	player.x, player.y, player.z = cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])
	result := backend.Cell(1)

	var err error

	result, err = state.checkpointTransitions(ctx, params[0], result)
	if err != nil {
		return 0, err
	}

	return state.gangZoneTransitions(ctx, params[0], result)
}

func (state *eventState) checkpointTransitions(ctx backend.NativeContext, playerID, result backend.Cell) (backend.Cell, error) {
	id := int(playerID)

	inside := state.checkpoints.playerInside([]backend.Cell{playerID}, state.checkpoints.checkpoints)
	if checkpoint := state.checkpoints.checkpoints[id]; checkpoint.active && inside != state.checkpointInside[id] {
		state.checkpointInside[id] = inside

		callback := "OnPlayerLeaveCheckpoint"
		if inside {
			callback = "OnPlayerEnterCheckpoint"
		}

		var err error

		result, err = callEvent(ctx, callback, playerID)
		if err != nil {
			return 0, err
		}
	}

	raceInside := state.insideRaceCheckpoint([]backend.Cell{playerID})
	if checkpoint := state.checkpoints.races[id]; checkpoint.active && raceInside != state.raceCheckpointInside[id] {
		state.raceCheckpointInside[id] = raceInside

		callback := "OnPlayerLeaveRaceCheckpoint"
		if raceInside {
			callback = "OnPlayerEnterRaceCheckpoint"
		}

		var err error

		result, err = callEvent(ctx, callback, playerID)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (state *eventState) gangZoneTransitions(ctx backend.NativeContext, playerID, result backend.Cell) (backend.Cell, error) {
	zoneIDs := make([]int, 0, len(state.gangZones.zones))
	for zoneID := range state.gangZones.zones {
		zoneIDs = append(zoneIDs, zoneID)
	}

	slices.Sort(zoneIDs)

	for _, zoneID := range zoneIDs {
		zone := state.gangZones.zones[zoneID]
		if zone == nil || !zone.check {
			continue
		}

		pair := eventPair{subject: zoneID, viewer: int(playerID)}

		inside := state.gangZones.playerInZone(int(playerID), zone)
		if inside == state.gangZoneInside[pair] {
			continue
		}

		state.gangZoneInside[pair] = inside

		callback := "OnPlayerLeaveGangZone"
		if inside {
			callback = "OnPlayerEnterGangZone"
		}

		var err error

		result, err = callEvent(ctx, callback, playerID, backend.Cell(zoneID))
		if err != nil {
			return 0, err
		}
	}

	playerZoneIDs := make([]int, 0, len(state.gangZones.playerZones[int(playerID)]))
	for zoneID := range state.gangZones.playerZones[int(playerID)] {
		playerZoneIDs = append(playerZoneIDs, zoneID)
	}

	slices.Sort(playerZoneIDs)

	for _, zoneID := range playerZoneIDs {
		zone := state.gangZones.playerZones[int(playerID)][zoneID]
		if zone == nil || !zone.check {
			continue
		}

		pair := eventPair{subject: zoneID, viewer: int(playerID)}

		inside := state.gangZones.playerInZone(int(playerID), zone)
		if inside == state.playerGangZoneInside[pair] {
			continue
		}

		state.playerGangZoneInside[pair] = inside

		callback := "OnPlayerLeavePlayerGangZone"
		if inside {
			callback = "OnPlayerEnterPlayerGangZone"
		}

		var err error

		result, err = callEvent(ctx, callback, playerID, backend.Cell(zoneID))
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (state *eventState) damageVehicle(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 3 {
		return 0, nil
	}

	vehicle := state.vehicles.vehicles[int(params[0])]
	reporter, reporterOK := state.players.player(params[1])

	amount := cellFloat(params[2])
	if vehicle == nil || !reporterOK || !reporter.connected || vehicle.health <= 0 || amount < 0 {
		return 0, nil
	}

	vehicle.health = max(0, vehicle.health-amount)
	if vehicle.health > 0 {
		return 1, nil
	}

	state.vehicles.ejectVehicle(int(params[0]))

	return callEvent(ctx, "OnVehicleDeath", params[0], params[1])
}

func (state *eventState) vehicleDamageStatus(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 6 {
		return 0, nil
	}

	vehicle := state.vehicles.vehicles[int(params[0])]

	reporter, reporterOK := state.players.player(params[1])
	if vehicle == nil || !reporterOK || !reporter.connected {
		return 0, nil
	}

	for index := range vehicle.damage {
		vehicle.damage[index] = int(params[index+2])
	}

	return callEvent(ctx, "OnVehicleDamageStatusUpdate", params[0], params[1])
}

func (state *eventState) respawnVehicle(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 1 || state.vehicles.vehicles[int(params[0])] == nil {
		return 0, nil
	}

	if _, err := state.vehicles.respawnVehicle(ctx, params[:1]); err != nil {
		return 0, err
	}

	return callEvent(ctx, "OnVehicleSpawn", params[0])
}

func (state *eventState) cancelTextDraw(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 1); !ok || !state.textDraws.selection[int(params[0])].active {
		return 0, nil
	}

	delete(state.textDraws.selection, int(params[0]))

	return callEvent(ctx, "OnPlayerClickTextDraw", params[0], 0xFFFF)
}

func (state *eventState) clickPlayer(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 3); !ok {
		return 0, nil
	}

	clicked, clickedOK := state.players.player(params[1])
	if !clickedOK || !clicked.connected {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerClickPlayer", params[0], params[1], params[2])
}

func (state *eventState) clickMap(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.player(params, 4); !ok {
		return 0, nil
	}

	return callEvent(ctx, "OnPlayerClickMap", params[0], params[1], params[2], params[3])
}

func (state *eventState) player(params []backend.Cell, minimum int) (*testPlayer, bool) {
	if len(params) < minimum {
		return nil, false
	}

	player, ok := state.players.player(params[0])

	return player, ok && player.connected
}

func (state *eventState) activeCheckpoint(params []backend.Cell, race bool) bool {
	if _, ok := state.player(params, 1); !ok {
		return false
	}

	if race {
		return state.checkpoints.races[int(params[0])].active
	}

	return state.checkpoints.checkpoints[int(params[0])].active
}

func (state *eventState) insideRaceCheckpoint(params []backend.Cell) bool {
	value, ok := state.checkpoints.races[int(params[0])]
	if !ok {
		return false
	}

	return state.checkpoints.playerInside(params, map[int]checkpoint{int(params[0]): value.checkpoint})
}

func callEvent(ctx backend.NativeContext, name string, params ...backend.Cell) (backend.Cell, error) {
	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support event callbacks")
	}

	result, err := caller.CallPublic(name, params...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	return result, nil
}
