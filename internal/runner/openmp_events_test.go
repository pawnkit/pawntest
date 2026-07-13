package runner

import (
	"reflect"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

type eventCall struct {
	name string
	args []backend.Cell
}

type eventMockVM struct {
	*mockVM
	calls []eventCall
}

func (vm *eventMockVM) CallPublic(name string, args ...backend.Cell) (backend.Cell, error) {
	vm.calls = append(vm.calls, eventCall{name: name, args: append([]backend.Cell(nil), args...)})
	return 1, nil
}

func registeredEventScenarios(t *testing.T) (*eventMockVM, *scenarioRegistry) {
	t.Helper()

	vm := &eventMockVM{mockVM: &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}}
	registry := newScenarioRegistry()

	context := &executionContext{state: &nativeState{status: Pass}, mocks: newMockState(), scenarios: registry}
	if err := registry.Register(vm, context); err != nil {
		t.Fatal(err)
	}

	return vm, registry
}

func callEventNative(t *testing.T, vm *eventMockVM, name string, params ...backend.Cell) backend.Cell {
	t.Helper()

	native := vm.natives[name]
	if native == nil {
		t.Fatalf("native %s was not registered", name)
	}

	result, err := native(vm, params)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}

	return result
}

func TestPlayerLifecycleEventsUpdateStateAndInvokeCallbacks(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	vm.strings[100] = "Alice"

	playerID := callEventNative(t, vm, "__pt_event_connect", 100)

	player := registry.Players.players[int(playerID)]
	if player == nil || !player.connected || player.spawned {
		t.Fatalf("connected player state = %+v", player)
	}

	callEventNative(t, vm, "__pt_event_spawn", playerID)
	callEventNative(t, vm, "__pt_event_death", playerID, -1, 54)
	callEventNative(t, vm, "__pt_event_disconnect", playerID, 1)

	if player.connected || player.spawned {
		t.Fatalf("disconnected player state = %+v", player)
	}

	got := []string{}
	for _, call := range vm.calls {
		got = append(got, call.name)
	}

	want := []string{"OnPlayerConnect", "OnPlayerSpawn", "OnPlayerDeath", "OnPlayerDisconnect"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestInputEventsExposeUpdatedKeyState(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	vm.strings[100] = "/help"
	playerID := registry.Players.addPlayer("Alice", true)

	callEventNative(t, vm, "__pt_event_keys", playerID, 4, -1, 1)
	callEventNative(t, vm, "__pt_event_command", playerID, 100)

	if got := registry.Players.players[int(playerID)].keys; got != [3]backend.Cell{4, -1, 1} {
		t.Fatalf("keys = %v", got)
	}

	if vm.calls[0].name != "OnPlayerKeyStateChange" || !reflect.DeepEqual(vm.calls[0].args, []backend.Cell{playerID, 4, 0}) {
		t.Fatalf("key callback = %+v", vm.calls[0])
	}

	if vm.calls[1].name != "OnPlayerCommandText" {
		t.Fatalf("command callback = %+v", vm.calls[1])
	}
}

func TestVehicleEventsUpdatePlayerState(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	vehicleID := registry.Vehicles.addVehicle(411, [3]float32{}, 0, -1, -1, -1)

	callEventNative(t, vm, "__pt_event_enter_vehicle", playerID, vehicleID, 0)

	if got := callEventNative(t, vm, "GetPlayerState", playerID); got != 2 {
		t.Fatalf("driver state = %d", got)
	}

	if registry.Players.players[int(playerID)].vehicle != int(vehicleID) {
		t.Fatal("player was not placed in the vehicle")
	}

	callEventNative(t, vm, "__pt_event_exit_vehicle", playerID)

	if got := callEventNative(t, vm, "GetPlayerState", playerID); got != 1 {
		t.Fatalf("on-foot state = %d", got)
	}

	want := []string{"OnPlayerEnterVehicle", "OnPlayerStateChange", "OnPlayerExitVehicle", "OnPlayerStateChange"}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestWorldEventsValidateScenarioState(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	player := registry.Players.players[int(playerID)]

	registry.Checkpoints.checkpoints[int(playerID)] = checkpoint{active: true, radius: 5}
	if result := callEventNative(t, vm, "__pt_event_enter_checkpoint", playerID); result != 1 {
		t.Fatalf("enter checkpoint = %d", result)
	}

	if result := callEventNative(t, vm, "__pt_event_enter_checkpoint", playerID); result != 0 {
		t.Fatalf("repeated enter checkpoint = %d", result)
	}

	player.x = 10

	if result := callEventNative(t, vm, "__pt_event_leave_checkpoint", playerID); result != 1 {
		t.Fatalf("leave checkpoint = %d", result)
	}

	pickupID := registry.Pickups.next
	registry.Pickups.next++

	registry.Pickups.pickups[pickupID] = newTestPickup(1240, 1, [3]float32{}, 0)
	if result := callEventNative(t, vm, "__pt_event_pickup", playerID, backend.Cell(pickupID)); result != 1 {
		t.Fatalf("pickup = %d", result)
	}

	want := []string{"OnPlayerEnterCheckpoint", "OnPlayerLeaveCheckpoint", "OnPlayerPickUpPickup"}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestTextDrawClickRequiresSelection(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	draw := newTestTextDraw(10, 20, "Click")
	draw.selectable = true
	draw.visible[int(playerID)] = true
	registry.TextDraws.draws[0] = draw

	if result := callEventNative(t, vm, "__pt_event_click_textdraw", playerID, 0); result != 0 {
		t.Fatalf("click without selection = %d", result)
	}

	registry.TextDraws.selection[int(playerID)] = textDrawSelection{active: true}
	if result := callEventNative(t, vm, "__pt_event_click_textdraw", playerID, 0); result != 1 {
		t.Fatalf("click with selection = %d", result)
	}

	if len(vm.calls) != 1 || vm.calls[0].name != "OnPlayerClickTextDraw" {
		t.Fatalf("callbacks = %+v", vm.calls)
	}
}

func TestEventScenarioCloneIsolatesTransitions(t *testing.T) {
	state := newEventState()
	state.checkpointInside[3] = true
	state.playerStreams[eventPair{subject: 1, viewer: 2}] = true
	state.gangZoneInside[eventPair{subject: 4, viewer: 3}] = true

	clone, ok := state.Clone().(*eventState)
	if !ok {
		t.Fatal("cloned scenario was not event state")
	}

	clone.checkpointInside[3] = false
	clone.playerStreams[eventPair{subject: 1, viewer: 2}] = false
	clone.gangZoneInside[eventPair{subject: 4, viewer: 3}] = false

	if !state.checkpointInside[3] {
		t.Fatal("event clone shared transition state")
	}

	if !state.playerStreams[eventPair{subject: 1, viewer: 2}] {
		t.Fatal("event clone shared stream state")
	}

	if !state.gangZoneInside[eventPair{subject: 4, viewer: 3}] {
		t.Fatal("event clone shared gang-zone state")
	}
}

func TestPlayerDamageAppliesArmourAndDeath(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	attackerID := registry.Players.addPlayer("Attacker", true)
	victimID := registry.Players.addPlayer("Victim", true)
	victim := registry.Players.players[int(victimID)]
	victim.armour = 20

	callEventNative(t, vm, "__pt_event_damage_player", victimID, attackerID, floatCell(30), 24, 3)

	if victim.armour != 0 || victim.health != 90 || !victim.spawned {
		t.Fatalf("damaged player = %+v", victim)
	}

	callEventNative(t, vm, "__pt_event_damage_player", victimID, attackerID, floatCell(100), 24, 3)

	if victim.health != 0 || victim.spawned {
		t.Fatalf("dead player = %+v", victim)
	}

	want := []string{
		"OnPlayerGiveDamage", "OnPlayerTakeDamage",
		"OnPlayerGiveDamage", "OnPlayerTakeDamage", "OnPlayerDeath",
	}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestActorDamageHonoursInvulnerability(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	actorID := callEventNative(t, vm, "CreateActor", 7, 0, 0, 0, 0)
	actor := registry.Actors.actors[int(actorID)]

	callEventNative(t, vm, "__pt_event_damage_actor", playerID, actorID, floatCell(25), 24, 3)

	if actor.health != 75 || len(vm.calls) != 1 || vm.calls[0].name != "OnPlayerGiveDamageActor" {
		t.Fatalf("damaged actor = %+v; callbacks = %+v", actor, vm.calls)
	}

	actor.invulnerable = true

	if result := callEventNative(t, vm, "__pt_event_damage_actor", playerID, actorID, floatCell(25), 24, 3); result != 0 {
		t.Fatalf("invulnerable damage = %d", result)
	}

	if actor.health != 75 || len(vm.calls) != 1 {
		t.Fatalf("invulnerable actor changed = %+v; callbacks = %+v", actor, vm.calls)
	}
}

func TestObjectMovementEventsCompleteTargets(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	objectID := registry.Objects.addObject(19379, [3]float32{}, [3]float32{}, 0)
	callEventNative(t, vm, "MoveObject", objectID, floatCell(10), floatCell(20), floatCell(30), floatCell(2))

	if result := callEventNative(t, vm, "__pt_event_object_moved", objectID); result != 1 {
		t.Fatalf("object moved = %d", result)
	}

	object := registry.Objects.objects[int(objectID)]
	if object.moving || object.position != [3]float32{10, 20, 30} {
		t.Fatalf("completed object = %+v", object)
	}

	if len(vm.calls) != 1 || vm.calls[0].name != "OnObjectMoved" {
		t.Fatalf("callbacks = %+v", vm.calls)
	}
}

func TestObjectMovementCompletesWithVirtualTime(t *testing.T) {
	vm := &eventMockVM{mockVM: &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}}
	scheduler := newScheduler()
	objects := newObjectState()
	objects.scheduler = scheduler
	objectID := objects.addObject(19379, [3]float32{}, [3]float32{}, 0)

	duration, err := objects.moveObject(vm, []backend.Cell{objectID, floatCell(10), 0, 0, floatCell(2)})
	if err != nil || duration != 5000 {
		t.Fatalf("MoveObject duration = %d, err = %v", duration, err)
	}

	if _, err := nativeAdvanceTime(scheduler)(vm, []backend.Cell{4999}); err != nil {
		t.Fatal(err)
	}

	if len(vm.calls) != 0 || !objects.objects[int(objectID)].moving {
		t.Fatal("object completed before its duration")
	}

	if _, err := nativeAdvanceTime(scheduler)(vm, []backend.Cell{1}); err != nil {
		t.Fatal(err)
	}

	object := objects.objects[int(objectID)]
	if object.moving || object.position != [3]float32{10, 0, 0} {
		t.Fatalf("completed object = %+v", object)
	}

	if len(vm.calls) != 1 || vm.calls[0].name != "OnObjectMoved" {
		t.Fatalf("callbacks = %+v", vm.calls)
	}
}

func TestStoppedObjectCancelsScheduledCompletion(t *testing.T) {
	vm := &eventMockVM{mockVM: &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{}}}
	scheduler := newScheduler()
	objects := newObjectState()
	objects.scheduler = scheduler
	objectID := objects.addObject(19379, [3]float32{}, [3]float32{}, 0)

	if _, err := objects.moveObject(vm, []backend.Cell{objectID, floatCell(10), 0, 0, floatCell(2)}); err != nil {
		t.Fatal(err)
	}

	if _, err := objects.stopObject(vm, []backend.Cell{objectID}); err != nil {
		t.Fatal(err)
	}

	if _, err := nativeAdvanceTime(scheduler)(vm, []backend.Cell{5000}); err != nil {
		t.Fatal(err)
	}

	if len(vm.calls) != 0 || objects.objects[int(objectID)].position != [3]float32{} {
		t.Fatalf("stopped object completed: object=%+v callbacks=%+v", objects.objects[int(objectID)], vm.calls)
	}
}

func TestPlayerObjectMovementEventCompletesTarget(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	objectID := registry.Objects.addPlayerObject(int(playerID), 19379, [3]float32{}, [3]float32{}, 0)
	object := registry.Objects.playerObjects[int(playerID)][int(objectID)]
	object.targetPos = [3]float32{4, 5, 6}
	object.moving = true

	if result := callEventNative(t, vm, "__pt_event_player_obj_moved", playerID, objectID); result != 1 {
		t.Fatalf("player object moved = %d", result)
	}

	if object.moving || object.position != [3]float32{4, 5, 6} {
		t.Fatalf("completed player object = %+v", object)
	}

	if len(vm.calls) != 1 || vm.calls[0].name != "OnPlayerObjectMoved" {
		t.Fatalf("callbacks = %+v", vm.calls)
	}
}

func TestStreamEventsRequireTransitionsAndMatchingWorlds(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	subjectID := registry.Players.addPlayer("Subject", true)
	viewerID := registry.Players.addPlayer("Viewer", true)

	if result := callEventNative(t, vm, "__pt_event_player_stream", subjectID, viewerID, 1); result != 1 {
		t.Fatalf("stream in = %d", result)
	}

	if result := callEventNative(t, vm, "__pt_event_player_stream", subjectID, viewerID, 1); result != 0 {
		t.Fatalf("repeated stream in = %d", result)
	}

	registry.Players.players[int(subjectID)].world = 2
	if result := callEventNative(t, vm, "__pt_event_player_stream", subjectID, viewerID, 0); result != 1 {
		t.Fatalf("stream out after world change = %d", result)
	}

	if result := callEventNative(t, vm, "__pt_event_player_stream", subjectID, viewerID, 1); result != 0 {
		t.Fatalf("cross-world stream in = %d", result)
	}

	want := []string{"OnPlayerStreamIn", "OnPlayerStreamOut"}

	got := []string{vm.calls[0].name, vm.calls[1].name}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestEntityStreamEventsInvokeTypedCallbacks(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	viewerID := registry.Players.addPlayer("Viewer", true)
	vehicleID := registry.Vehicles.addVehicle(411, [3]float32{}, 0, -1, -1, -1)
	actorID := callEventNative(t, vm, "CreateActor", 7, 0, 0, 0, 0)

	callEventNative(t, vm, "__pt_event_vehicle_stream", vehicleID, viewerID, 1)
	callEventNative(t, vm, "__pt_event_vehicle_stream", vehicleID, viewerID, 0)
	callEventNative(t, vm, "__pt_event_actor_stream", actorID, viewerID, 1)
	callEventNative(t, vm, "__pt_event_actor_stream", actorID, viewerID, 0)

	want := []string{"OnVehicleStreamIn", "OnVehicleStreamOut", "OnActorStreamIn", "OnActorStreamOut"}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestMovePlayerEmitsAreaTransitions(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	registry.Checkpoints.checkpoints[int(playerID)] = checkpoint{active: true, position: [3]float32{10, 0, 0}, radius: 2}
	registry.GangZones.zones[7] = &testGangZone{bounds: [4]float32{8, -2, 12, 2}, check: true, views: map[int]gangZoneView{}}

	if result := callEventNative(t, vm, "__pt_event_move_player", playerID, floatCell(10), 0, 0); result != 1 {
		t.Fatalf("move into areas = %d", result)
	}

	if result := callEventNative(t, vm, "__pt_event_move_player", playerID, floatCell(20), 0, 0); result != 1 {
		t.Fatalf("move out of areas = %d", result)
	}

	want := []string{
		"OnPlayerEnterCheckpoint", "OnPlayerEnterGangZone",
		"OnPlayerLeaveCheckpoint", "OnPlayerLeaveGangZone",
	}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestVehicleDamageDeathAndRespawnEvents(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	reporterID := registry.Players.addPlayer("Reporter", true)
	vehicleID := registry.Vehicles.addVehicle(411, [3]float32{1, 2, 3}, 90, -1, -1, -1)
	callEventNative(t, vm, "PutPlayerInVehicle", reporterID, vehicleID, 0)

	callEventNative(t, vm, "__pt_event_vehicle_status", vehicleID, reporterID, 1, 2, 3, 4)

	vehicle := registry.Vehicles.vehicles[int(vehicleID)]
	if vehicle.damage != [4]int{1, 2, 3, 4} {
		t.Fatalf("vehicle damage status = %v", vehicle.damage)
	}

	callEventNative(t, vm, "__pt_event_vehicle_damage", vehicleID, reporterID, floatCell(1000))

	if vehicle.health != 0 || registry.Players.players[int(reporterID)].vehicle != -1 {
		t.Fatalf("destroyed vehicle = %+v", vehicle)
	}

	callEventNative(t, vm, "__pt_event_vehicle_respawn", vehicleID)

	if vehicle.health != 1000 || vehicle.position != [3]float32{1, 2, 3} {
		t.Fatalf("respawned vehicle = %+v", vehicle)
	}

	want := []string{"OnVehicleDamageStatusUpdate", "OnVehicleDeath", "OnVehicleSpawn"}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestClickAndSelectionEvents(t *testing.T) {
	vm, registry := registeredEventScenarios(t)
	playerID := registry.Players.addPlayer("Alice", true)
	clickedID := registry.Players.addPlayer("Bob", true)
	registry.TextDraws.selection[int(playerID)] = textDrawSelection{active: true}

	callEventNative(t, vm, "__pt_event_cancel_textdraw", playerID)
	callEventNative(t, vm, "__pt_event_click_player", playerID, clickedID, 0)
	callEventNative(t, vm, "__pt_event_click_map", playerID, floatCell(1), floatCell(2), floatCell(3))

	if registry.TextDraws.selection[int(playerID)].active {
		t.Fatal("textdraw selection remained active")
	}

	want := []string{"OnPlayerClickTextDraw", "OnPlayerClickPlayer", "OnPlayerClickMap"}

	got := make([]string, len(vm.calls))
	for index, call := range vm.calls {
		got[index] = call.name
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}

	if !reflect.DeepEqual(vm.calls[0].args, []backend.Cell{playerID, 0xFFFF}) {
		t.Fatalf("cancel callback args = %v", vm.calls[0].args)
	}
}
