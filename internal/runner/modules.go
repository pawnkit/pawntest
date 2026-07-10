package runner

import "github.com/pawnkit/pawntest/internal/backend"

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
	modules []scenarioModule
}

func newScenarioRegistry() *scenarioRegistry {
	return &scenarioRegistry{modules: []scenarioModule{newOpenMPState()}}
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
	clone := &scenarioRegistry{modules: make([]scenarioModule, 0, len(registry.modules))}
	for _, module := range registry.modules {
		clone.modules = append(clone.modules, module.Clone())
	}

	return clone
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
