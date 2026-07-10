package runner

import (
	"math"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *vehicleState) changeVehicleColour(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	vehicle.colour1, vehicle.colour2 = int(params[1]), int(params[2])

	return 1, nil
}

func (state *vehicleState) getVehicleColour(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	if err := ctx.WriteCell(params[1], backend.Cell(vehicle.colour1)); err != nil {
		return 0, err
	}

	return 1, ctx.WriteCell(params[2], backend.Cell(vehicle.colour2))
}

func (state *vehicleState) setVehiclePlate(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	plate, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}

	vehicle.plate = truncateString(plate, 32)

	return 1, nil
}

func (state *vehicleState) getVehiclePlate(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 3 {
		return 0, nil
	}

	return 1, ctx.WriteString(params[1], truncateString(vehicle.plate, int(params[2])))
}

func (state *vehicleState) addVehicleComponent(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	vehicle.components[int(params[1])] = true

	return 1, nil
}

func (state *vehicleState) removeVehicleComponent(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	delete(vehicle.components, int(params[1]))

	return 1, nil
}

func (state *vehicleState) getVehicleComponentInSlot(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	for component := range vehicle.components {
		if componentSlot(component) == int(params[1]) {
			return backend.Cell(component), nil
		}
	}

	return 0, nil
}

func componentSlot(component int) int {
	if component < 1000 || component > 1193 {
		return -1
	}

	return component % 14
}

func (state *vehicleState) attachTrailer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	trailer, trailerOK := state.vehicle(params)
	if !trailerOK || len(params) < 2 {
		return 0, nil
	}

	vehicle, vehicleOK := state.vehicles[int(params[1])]
	if !vehicleOK {
		return 0, nil
	}

	vehicle.trailer = int(params[0])
	_ = trailer

	return 1, nil
}

func (state *vehicleState) detachTrailer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok {
		return 0, nil
	}

	vehicle.trailer = -1

	return 1, nil
}

func (state *vehicleState) isTrailerAttached(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if ok && vehicle.trailer >= 0 {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) getVehicleTrailer(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok {
		return -1, nil
	}

	return backend.Cell(vehicle.trailer), nil
}

func (state *vehicleState) setVehicleDamage(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	for index := range vehicle.damage {
		vehicle.damage[index] = int(params[index+1])
	}

	return 1, nil
}

func (state *vehicleState) getVehicleDamage(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	for index, value := range vehicle.damage {
		if err := ctx.WriteCell(params[index+1], backend.Cell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *vehicleState) repairVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok {
		return 0, nil
	}

	vehicle.health = 1000
	vehicle.damage = [4]int{}

	return 1, nil
}

func (state *vehicleState) respawnVehicle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok {
		return 0, nil
	}

	state.ejectVehicle(int(params[0]))

	vehicle.position = vehicle.spawn
	vehicle.angle = vehicle.spawnAngle
	vehicle.velocity = [3]float32{}
	vehicle.angular = [3]float32{}
	vehicle.health = 1000
	vehicle.damage = [4]int{}

	return 1, nil
}

func (state *vehicleState) getVehicleDistance(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	dx := vehicle.position[0] - cellFloat(params[1])
	dy := vehicle.position[1] - cellFloat(params[2])
	dz := vehicle.position[2] - cellFloat(params[3])
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))

	return floatCell(distance), nil
}

func (state *vehicleState) setVehicleParams(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 8 {
		return 0, nil
	}

	for index := range vehicle.params {
		vehicle.params[index] = int(params[index+1])
	}

	return 1, nil
}

func (state *vehicleState) getVehicleParams(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 8 {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[1:8], vehicle.params[:])
}

func (state *vehicleState) setVehicleDoors(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.setVehiclePart(params, func(vehicle *testVehicle) *[4]int { return &vehicle.doors })
}

func (state *vehicleState) getVehicleDoors(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getVehiclePart(ctx, params, func(vehicle *testVehicle) *[4]int { return &vehicle.doors })
}

func (state *vehicleState) setVehicleWindows(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.setVehiclePart(params, func(vehicle *testVehicle) *[4]int { return &vehicle.windows })
}

func (state *vehicleState) getVehicleWindows(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	return state.getVehiclePart(ctx, params, func(vehicle *testVehicle) *[4]int { return &vehicle.windows })
}

func (state *vehicleState) setVehiclePart(params []backend.Cell, field func(*testVehicle) *[4]int) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	for index := range *field(vehicle) {
		field(vehicle)[index] = int(params[index+1])
	}

	return 1, nil
}

func (state *vehicleState) getVehiclePart(ctx backend.NativeContext, params []backend.Cell, field func(*testVehicle) *[4]int) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 5 {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[1:5], field(vehicle)[:])
}

func writeVehicleInts(ctx backend.NativeContext, addresses []backend.Cell, values []int) (backend.Cell, error) {
	for index, value := range values {
		if err := ctx.WriteCell(addresses[index], backend.Cell(value)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *vehicleState) setVehicleSiren(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	vehicle.siren = params[1] != 0

	return 1, nil
}

func (state *vehicleState) getVehicleSiren(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if ok && vehicle.siren {
		return 1, nil
	}

	return 0, nil
}

func (state *vehicleState) setVehicleSirenEnabled(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if !ok || len(params) < 2 {
		return 0, nil
	}

	vehicle.sirenEnabled = params[1] != 0

	return 1, nil
}

func (state *vehicleState) getVehicleSirenEnabled(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	vehicle, ok := state.vehicle(params)
	if ok && vehicle.sirenEnabled {
		return 1, nil
	}

	return 0, nil
}
