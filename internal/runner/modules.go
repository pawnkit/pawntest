package runner

import (
	"errors"

	"github.com/pawnkit/pawntest/internal/backend"
)

type executionContext struct {
	state        *nativeState
	mocks        *mockState
	scheduler    *scheduler
	snapshots    *snapshotStore
	fixtures     *namedFixtureState
	scenarios    *scenarioRegistry
	publicName   string
	fuzzSeed     int64
	allowUnknown bool
	custom       map[string]backend.NativeFunc
	strict       bool
}

type nativeModule interface {
	Register(vm backend.VM, context *executionContext) error
}

type nativeModuleFunc func(vm backend.VM, context *executionContext) error

func (register nativeModuleFunc) Register(vm backend.VM, context *executionContext) error {
	return register(vm, context)
}

func defaultNativeModules() []nativeModule {
	return []nativeModule{
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerPawntestNatives(vm, context.state, context.mocks)
		}),
		nativeModuleFunc(func(vm backend.VM, _ *executionContext) error {
			return registerFloatNatives(vm)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerUnknownNativeMocks(vm, context.mocks, context.allowUnknown)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerMockControlNatives(vm, context.mocks)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerSchedulerNatives(vm, context.scheduler)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerSnapshotNative(vm, context.state, context.snapshots)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerFuzzNative(vm, context.state, context.fuzzSeed, context.publicName)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return registerNamedFixtureNative(vm, context.state, context.fixtures)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			return context.scenarios.Register(vm, context)
		}),
		nativeModuleFunc(func(vm backend.VM, context *executionContext) error {
			for name, native := range context.custom {
				if err := vm.RegisterNative(name, native); err != nil {
					return err
				}
			}

			return nil
		}),
	}
}

type scenarioModule interface {
	nativeModule
	Clone() scenarioModule
}

type scenarioRegistry struct {
	Players     *openMPState
	Vehicles    *vehicleState
	Objects     *objectState
	Actors      *actorState
	Pickups     *pickupState
	Checkpoints *checkpointState
	TextLabels  *textLabelState
	TextDraws   *textDrawState
	GangZones   *gangZoneState
	Dialogs     *dialogState
	Menus       *menuState
	Classes     *classState
	Variables   *variableState
	Server      *serverState
	NPCs        *npcState
	Database    *databaseState
	HTTP        *httpState
	modules     []scenarioModule
}

func newScenarioRegistry() *scenarioRegistry {
	registry := &scenarioRegistry{
		Players: newOpenMPState(), Vehicles: newVehicleState(), Objects: newObjectState(), Actors: newActorState(),
		Pickups: newPickupState(), Checkpoints: newCheckpointState(), TextLabels: newTextLabelState(), TextDraws: newTextDrawState(),
		GangZones: newGangZoneState(), Dialogs: newDialogState(), Menus: newMenuState(), Classes: newClassState(),
		Variables: newVariableState(), Server: newServerState(), NPCs: newNPCState(), Database: newDatabaseState(), HTTP: newHTTPState(),
	}
	registry.setModules()

	return registry
}

func (registry *scenarioRegistry) actorState() *actorState {
	return registry.Actors
}

func registerScenarioNatives(vm backend.VM, natives map[string]backend.NativeFunc, mocks *mockState, allowUnknown bool) error {
	for name, native := range natives {
		registered := native
		if !isPawntestNative(name) {
			registered = func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
				if mocks.configured(name) {
					return mockUnknownNative(name, mocks, allowUnknown)(ctx, params)
				}

				mocks.recordCall(name, ctx, params)

				return native(ctx, params)
			}
		}

		if err := vm.RegisterNative(name, registered); err != nil {
			return err
		}
	}

	return nil
}

func (registry *scenarioRegistry) playerState() *openMPState {
	return registry.Players
}

func (registry *scenarioRegistry) vehicleState() *vehicleState {
	return registry.Vehicles
}

func (registry *scenarioRegistry) Register(vm backend.VM, context *executionContext) error {
	for _, module := range registry.modules {
		if err := module.Register(vm, context); err != nil {
			return err
		}
	}

	return nil
}

func (registry *scenarioRegistry) Clone() *scenarioRegistry {
	clone := &scenarioRegistry{
		Players: registry.Players.clone(), Vehicles: cloneScenario[*vehicleState](registry.Vehicles),
		Objects: cloneScenario[*objectState](registry.Objects), Actors: cloneScenario[*actorState](registry.Actors),
		Pickups: cloneScenario[*pickupState](registry.Pickups), Checkpoints: cloneScenario[*checkpointState](registry.Checkpoints),
		TextLabels: cloneScenario[*textLabelState](registry.TextLabels), TextDraws: cloneScenario[*textDrawState](registry.TextDraws),
		GangZones: cloneScenario[*gangZoneState](registry.GangZones), Dialogs: cloneScenario[*dialogState](registry.Dialogs),
		Menus: cloneScenario[*menuState](registry.Menus), Classes: cloneScenario[*classState](registry.Classes),
		Variables: cloneScenario[*variableState](registry.Variables), Server: cloneScenario[*serverState](registry.Server),
		NPCs: cloneScenario[*npcState](registry.NPCs), Database: cloneScenario[*databaseState](registry.Database),
		HTTP: cloneScenario[*httpState](registry.HTTP),
	}
	clone.setModules()

	return clone
}

func (registry *scenarioRegistry) setModules() {
	registry.modules = []scenarioModule{
		registry.Players, registry.Vehicles, registry.Objects, registry.Actors,
		registry.Pickups, registry.Checkpoints, registry.TextLabels, registry.TextDraws,
		registry.GangZones, registry.Dialogs, registry.Menus, registry.Classes,
		registry.Variables, registry.Server, registry.NPCs, registry.Database, registry.HTTP,
	}
}

func cloneScenario[T scenarioModule](module T) T {
	clone, ok := module.Clone().(T)
	if !ok {
		panic("scenario clone returned a different module type")
	}

	return clone
}

func (registry *scenarioRegistry) Close() error {
	var closeErrors []error

	for _, module := range registry.modules {
		closer, ok := module.(interface{ Close() error })
		if ok {
			closeErrors = append(closeErrors, closer.Close())
		}
	}

	return errors.Join(closeErrors...)
}

func (registry *scenarioRegistry) StrictFailures() []string {
	failures := []string{}

	for _, module := range registry.modules {
		verifier, ok := module.(interface{ StrictFailures() []string })
		if ok {
			failures = append(failures, verifier.StrictFailures()...)
		}
	}

	return failures
}

func newExecutionContext(snapshots *snapshotStore, scenarios *scenarioRegistry, runner Runner) *executionContext {
	return &executionContext{
		mocks:        newMockState(),
		scheduler:    newScheduler(),
		snapshots:    snapshots,
		fixtures:     &namedFixtureState{},
		scenarios:    scenarios,
		fuzzSeed:     runner.FuzzSeed,
		allowUnknown: runner.AllowUnknownNatives,
		custom:       runner.Natives,
	}
}
