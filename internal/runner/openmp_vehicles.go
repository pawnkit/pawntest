package runner

import (
	"errors"
	"fmt"
	"maps"
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testVehicle struct {
	model, colour1, colour2 int
	world, interior         int
	paintjob, respawnDelay  int
	plate                   string
	health                  float32
	position, spawn         [3]float32
	angle, spawnAngle       float32
	velocity, angular       [3]float32
	damage                  [4]int
	params                  [7]int
	doors, windows          [4]int
	components              map[int]bool
	trailer                 int
	siren, sirenEnabled     bool
	occupied                bool
	lastDriver              int
}

type vehicleState struct {
	next     int
	vehicles map[int]*testVehicle
	players  *openMPState
}

func newVehicleState() *vehicleState {
	return &vehicleState{next: 1, vehicles: map[int]*testVehicle{}}
}

func (state *vehicleState) Clone() scenarioModule {
	clone := newVehicleState()

	clone.next = state.next
	for id, vehicle := range state.vehicles {
		vehicleCopy := *vehicle
		vehicleCopy.components = maps.Clone(vehicle.components)
		clone.vehicles[id] = &vehicleCopy
	}

	return clone
}

func (state *vehicleState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	natives := state.natives(context.state)

	return registerScenarioNatives(vm, natives, context.mocks, context.allowUnknown)
}

func (state *vehicleState) natives(result *nativeState) map[string]backend.NativeFunc {
	natives := map[string]backend.NativeFunc{
		"__pt_vehicle_create":         state.createTestVehicle,
		"__pt_vehicle_valid":          state.assertValid(result),
		"__pt_vehicle_model":          state.assertModel(result),
		"__pt_vehicle_health":         state.assertHealth(result),
		"__pt_vehicle_pos_near":       state.assertPosition(result),
		"AddStaticVehicle":            state.createVehicle,
		"AddStaticVehicleEx":          state.createVehicle,
		"CreateVehicle":               state.createVehicle,
		"DestroyVehicle":              state.destroyVehicle,
		"IsValidVehicle":              state.isValidVehicle,
		"GetVehicleModel":             state.getVehicleInt(vehicleModel),
		"GetVehicleVirtualWorld":      state.getVehicleInt(vehicleWorld),
		"SetVehicleVirtualWorld":      state.setVehicleInt(vehicleWorld),
		"GetVehicleInterior":          state.getVehicleInt(vehicleInterior),
		"LinkVehicleToInterior":       state.setVehicleInt(vehicleInterior),
		"GetVehiclePaintjob":          state.getVehicleInt(vehiclePaintjob),
		"ChangeVehiclePaintjob":       state.setVehicleInt(vehiclePaintjob),
		"GetVehicleRespawnDelay":      state.getVehicleInt(vehicleRespawnDelay),
		"SetVehicleRespawnDelay":      state.setVehicleInt(vehicleRespawnDelay),
		"GetVehiclePos":               state.getVehicleVector(vehiclePosition),
		"SetVehiclePos":               state.setVehicleVector(vehiclePosition),
		"GetVehicleVelocity":          state.getVehicleVector(vehicleVelocity),
		"SetVehicleVelocity":          state.setVehicleVector(vehicleVelocity),
		"SetVehicleAngularVelocity":   state.setVehicleVector(vehicleAngularVelocity),
		"GetVehicleZAngle":            state.getVehicleFloat(vehicleAngle),
		"SetVehicleZAngle":            state.setVehicleFloat(vehicleAngle),
		"GetVehicleHealth":            state.getVehicleFloat(vehicleHealth),
		"SetVehicleHealth":            state.setVehicleFloat(vehicleHealth),
		"ChangeVehicleColor":          state.changeVehicleColour,
		"ChangeVehicleColours":        state.changeVehicleColour,
		"GetVehicleColor":             state.getVehicleColour,
		"GetVehicleColours":           state.getVehicleColour,
		"SetVehicleNumberPlate":       state.setVehiclePlate,
		"GetVehicleNumberPlate":       state.getVehiclePlate,
		"AddVehicleComponent":         state.addVehicleComponent,
		"RemoveVehicleComponent":      state.removeVehicleComponent,
		"GetVehicleComponentInSlot":   state.getVehicleComponentInSlot,
		"AttachTrailerToVehicle":      state.attachTrailer,
		"DetachTrailerFromVehicle":    state.detachTrailer,
		"IsTrailerAttachedToVehicle":  state.isTrailerAttached,
		"GetVehicleTrailer":           state.getVehicleTrailer,
		"GetVehicleDamageStatus":      state.getVehicleDamage,
		"UpdateVehicleDamageStatus":   state.setVehicleDamage,
		"RepairVehicle":               state.repairVehicle,
		"SetVehicleToRespawn":         state.respawnVehicle,
		"GetVehicleDistanceFromPoint": state.getVehicleDistance,
		"PutPlayerInVehicle":          state.putPlayerInVehicle,
		"RemovePlayerFromVehicle":     state.removePlayerFromVehicle,
		"GetPlayerVehicleID":          state.getPlayerVehicleID,
		"GetPlayerVehicleSeat":        state.getPlayerVehicleSeat,
		"IsPlayerInVehicle":           state.isPlayerInVehicle,
		"IsPlayerInAnyVehicle":        state.isPlayerInAnyVehicle,
		"GetVehicleDriver":            state.getVehicleDriver,
		"GetVehicleLastDriver":        state.getVehicleLastDriver,
		"GetVehicleOccupant":          state.getVehicleOccupant,
		"CountVehicleOccupants":       state.countVehicleOccupants,
		"IsVehicleOccupied":           state.isVehicleOccupied,
		"HasVehicleBeenOccupied":      state.hasVehicleBeenOccupied,
		"SetVehicleParamsEx":          state.setVehicleParams,
		"GetVehicleParamsEx":          state.getVehicleParams,
		"SetVehicleParamsCarDoors":    state.setVehicleDoors,
		"GetVehicleParamsCarDoors":    state.getVehicleDoors,
		"SetVehicleParamsCarWindows":  state.setVehicleWindows,
		"GetVehicleParamsCarWindows":  state.getVehicleWindows,
		"SetVehicleParamsSirenState":  state.setVehicleSiren,
		"GetVehicleParamsSirenState":  state.getVehicleSiren,
		"ToggleVehicleSirenEnabled":   state.setVehicleSirenEnabled,
		"IsVehicleSirenEnabled":       state.getVehicleSirenEnabled,
	}

	return natives
}

func (state *vehicleState) createTestVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return -1, nil
	}

	return state.addVehicle(int(params[0]), [3]float32{cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}, 0, -1, -1, -1), nil
}

func (state *vehicleState) createVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 {
		return -1, nil
	}

	respawnDelay := -1
	if len(params) > 7 {
		respawnDelay = int(params[7])
	}

	return state.addVehicle(int(params[0]), [3]float32{cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}, cellFloat(params[4]), int(params[5]), int(params[6]), respawnDelay), nil
}

func (state *vehicleState) addVehicle(model int, position [3]float32, angle float32, colour1, colour2, respawnDelay int) backend.Cell {
	id := state.next
	state.next++
	state.vehicles[id] = &testVehicle{
		model: model, colour1: colour1, colour2: colour2,
		position: position, spawn: position, angle: angle, spawnAngle: angle,
		health: 1000, paintjob: -1, respawnDelay: respawnDelay,
		components: map[int]bool{}, trailer: -1, lastDriver: -1,
	}

	return backend.Cell(id)
}

func (state *vehicleState) vehicle(params []backend.Cell) (*testVehicle, bool) {
	if len(params) == 0 {
		return nil, false
	}

	vehicle, ok := state.vehicles[int(params[0])]

	return vehicle, ok
}

func (state *vehicleState) destroyVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.vehicle(params); !ok {
		return 0, nil
	}

	state.ejectVehicle(int(params[0]))
	delete(state.vehicles, int(params[0]))

	return 1, nil
}

func (state *vehicleState) isValidVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.vehicle(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) setVehicleInt(field func(*testVehicle) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(vehicle) = int(params[1])

		return 1, nil
	}
}

func (state *vehicleState) getVehicleInt(field func(*testVehicle) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok {
			return 0, nil
		}

		return backend.Cell(*field(vehicle)), nil
	}
}

func (state *vehicleState) setVehicleFloat(field func(*testVehicle) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(vehicle) = cellFloat(params[1])

		return 1, nil
	}
}

func (state *vehicleState) getVehicleFloat(field func(*testVehicle) *float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return 1, ctx.WriteCell(params[1], floatCell(*field(vehicle)))
	}
}

func (state *vehicleState) setVehicleVector(field func(*testVehicle) *[3]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		*field(vehicle) = [3]float32{cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}

		return 1, nil
	}
}

func (state *vehicleState) getVehicleVector(field func(*testVehicle) *[3]float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		vehicle, ok := state.vehicle(params)
		if !ok || len(params) < 4 {
			return 0, nil
		}

		return writeFloatVector(ctx, params[1:4], *field(vehicle))
	}
}

func writeFloatVector(ctx backend.NativeContext, addresses []backend.Cell, values [3]float32) (backend.Cell, error) {
	for index, value := range values {
		if err := ctx.WriteCell(addresses[index], floatCell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func vehicleModel(vehicle *testVehicle) *int                  { return &vehicle.model }
func vehicleWorld(vehicle *testVehicle) *int                  { return &vehicle.world }
func vehicleInterior(vehicle *testVehicle) *int               { return &vehicle.interior }
func vehiclePaintjob(vehicle *testVehicle) *int               { return &vehicle.paintjob }
func vehicleRespawnDelay(vehicle *testVehicle) *int           { return &vehicle.respawnDelay }
func vehicleHealth(vehicle *testVehicle) *float32             { return &vehicle.health }
func vehicleAngle(vehicle *testVehicle) *float32              { return &vehicle.angle }
func vehiclePosition(vehicle *testVehicle) *[3]float32        { return &vehicle.position }
func vehicleVelocity(vehicle *testVehicle) *[3]float32        { return &vehicle.velocity }
func vehicleAngularVelocity(vehicle *testVehicle) *[3]float32 { return &vehicle.angular }

func (state *vehicleState) assertValid(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("vehicle assertion expects 3 arguments")
		}

		_, ok := state.vehicle(params)
		if !ok {
			setFailure(result, params, 1, fmt.Sprintf("vehicle %d does not exist", params[0]), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func (state *vehicleState) assertModel(result *nativeState) backend.NativeFunc {
	return state.assertVehicleInt(result, "model", func(vehicle *testVehicle) int { return vehicle.model })
}

func (state *vehicleState) assertHealth(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("vehicle health assertion expects 4 arguments")
		}

		vehicle, ok := state.vehicle(params)

		expected := cellFloat(params[1])
		if !ok || vehicle.health != expected {
			actual := float32(0)
			if ok {
				actual = vehicle.health
			}

			setFailure(result, params, 2, fmt.Sprintf("vehicle %d health: expected %g, got %g", params[0], expected, actual), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func (state *vehicleState) assertPosition(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			return 0, errors.New("vehicle position assertion expects 7 arguments")
		}

		vehicle, ok := state.vehicle(params)
		if !ok {
			setFailure(result, params, 5, fmt.Sprintf("vehicle %d does not exist", params[0]), ctx)
			return 0, nil
		}

		expected := [3]float32{cellFloat(params[1]), cellFloat(params[2]), cellFloat(params[3])}

		tolerance := float32(math.Abs(float64(cellFloat(params[4]))))
		for index := range expected {
			if absFloat(vehicle.position[index]-expected[index]) > tolerance {
				setFailure(result, params, 5, fmt.Sprintf("vehicle %d position: expected %v +/- %g, got %v", params[0], expected, tolerance, vehicle.position), ctx)
				return 0, nil
			}
		}

		return 1, nil
	}
}

func (state *vehicleState) assertVehicleInt(result *nativeState, label string, value func(*testVehicle) int) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, errors.New("vehicle integer assertion expects 4 arguments")
		}

		vehicle, ok := state.vehicle(params)

		actual := 0
		if ok {
			actual = value(vehicle)
		}

		if !ok || actual != int(params[1]) {
			setFailure(result, params, 2, fmt.Sprintf("vehicle %d %s: expected %d, got %d", params[0], label, params[1], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
