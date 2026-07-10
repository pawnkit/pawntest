package runner

import "github.com/pawnkit/pawntest/internal/backend"

type actorAnimation struct {
	library, name              string
	delta                      float32
	loop, lockX, lockY, freeze bool
	time                       int
}

type testActor struct {
	skin, world   int
	health, angle float32
	position      [3]float32
	spawnPosition [3]float32
	spawnSkin     int
	spawnAngle    float32
	invulnerable  bool
	animation     actorAnimation
	hasAnimation  bool
}

type actorState struct {
	next    int
	actors  map[int]*testActor
	players *openMPState
}

func newActorState() *actorState {
	return &actorState{actors: map[int]*testActor{}}
}

func (state *actorState) Clone() scenarioModule {
	clone := newActorState()

	clone.next = state.next
	for id, actor := range state.actors {
		actorCopy := *actor
		clone.actors[id] = &actorCopy
	}

	return clone
}

func (state *actorState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *actorState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_actor_create":    state.createActor,
		"__pt_actor_valid":     state.assertValid(result),
		"__pt_actor_skin":      state.assertSkin(result),
		"__pt_actor_health":    state.assertHealth(result),
		"__pt_actor_pos_near":  state.assertPosition(result),
		"CreateActor":          state.createActor,
		"DestroyActor":         state.destroyActor,
		"IsValidActor":         state.isValidActor,
		"IsActorStreamedIn":    state.isActorStreamedIn,
		"SetActorVirtualWorld": state.setActorInt(actorWorld),
		"GetActorVirtualWorld": state.getActorInt(actorWorld),
		"SetActorSkin":         state.setActorInt(actorSkin),
		"GetActorSkin":         state.getActorInt(actorSkin),
		"SetActorHealth":       state.setActorFloat(actorHealth),
		"GetActorHealth":       state.getActorFloat(actorHealth),
		"SetActorFacingAngle":  state.setActorFloat(actorAngle),
		"GetActorFacingAngle":  state.getActorFloat(actorAngle),
		"SetActorPos":          state.setActorPosition,
		"GetActorPos":          state.getActorPosition,
		"SetActorInvulnerable": state.setActorInvulnerable,
		"IsActorInvulnerable":  state.isActorInvulnerable,
		"ApplyActorAnimation":  state.applyActorAnimation,
		"ClearActorAnimations": state.clearActorAnimations,
		"GetActorAnimation":    state.getActorAnimation,
		"GetActorSpawnInfo":    state.getActorSpawnInfo,
		"GetActors":            state.getActors,
	}
}

func (state *actorState) createActor(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return -1, nil
	}

	id := state.next
	state.next++
	position := cellsToVector(params[1:4])
	angle := cellFloat(params[4])
	state.actors[id] = &testActor{
		skin: int(params[0]), health: 100, position: position, angle: angle,
		spawnSkin: int(params[0]), spawnPosition: position, spawnAngle: angle,
	}

	return backend.Cell(id), nil
}

func (state *actorState) actor(params []backend.Cell) (*testActor, bool) {
	if len(params) == 0 {
		return nil, false
	}

	actor, ok := state.actors[int(params[0])]

	return actor, ok
}

func (state *actorState) destroyActor(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.actor(params); !ok {
		return 0, nil
	}

	delete(state.actors, int(params[0]))

	return 1, nil
}

func (state *actorState) isValidActor(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.actor(params); ok {
		return 1, nil
	}

	return 0, nil
}

func (state *actorState) isActorStreamedIn(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || len(params) < 2 || state.players == nil {
		return 0, nil
	}

	player, playerOK := state.players.player(params[1])
	if playerOK && player.connected && player.world == actor.world {
		return 1, nil
	}

	return 0, nil
}

func (state *actorState) setActorInt(field func(*testActor) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		actor, ok := state.actor(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(actor) = int(params[1])

		return 1, nil
	}
}

func (state *actorState) getActorInt(field func(*testActor) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		actor, ok := state.actor(params)
		if !ok {
			return -1, nil
		}

		return backend.Cell(*field(actor)), nil
	}
}

func (state *actorState) setActorFloat(field func(*testActor) *float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		actor, ok := state.actor(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		*field(actor) = cellFloat(params[1])

		return 1, nil
	}
}

func (state *actorState) getActorFloat(field func(*testActor) *float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		actor, ok := state.actor(params)
		if !ok || len(params) < 2 {
			return 0, nil
		}

		return 1, ctx.WriteCell(params[1], floatCell(*field(actor)))
	}
}

func actorWorld(actor *testActor) *int      { return &actor.world }
func actorSkin(actor *testActor) *int       { return &actor.skin }
func actorHealth(actor *testActor) *float32 { return &actor.health }
func actorAngle(actor *testActor) *float32  { return &actor.angle }

func (state *actorState) setActorPosition(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	actor.position = cellsToVector(params[1:4])

	return 1, nil
}

func (state *actorState) getActorPosition(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeFloatVector(ctx, params[1:4], actor.position)
}

func (state *actorState) setActorInvulnerable(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if !ok {
		return 0, nil
	}

	actor.invulnerable = len(params) < 2 || params[1] != 0

	return 1, nil
}

func (state *actorState) isActorInvulnerable(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	actor, ok := state.actor(params)
	if ok && actor.invulnerable {
		return 1, nil
	}

	return 0, nil
}
