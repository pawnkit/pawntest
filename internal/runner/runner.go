package runner

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	goamx "github.com/pawnkit/goamx/vm"
	"github.com/pawnkit/pawntest/internal/backend"
	"github.com/pawnkit/pawntest/internal/discovery"
)

const markerPublic = "__pawntest_marker"

var (
	ErrMissingPawntestInclude = errors.New("missing pawntest include marker")
	ErrNoTestsFound           = errors.New("no tests found")
)

type Runner struct {
	Backend             backend.Backend
	Run                 string
	FailFast            bool
	AllowUnknownNatives bool
	Isolation           string
	Shuffle             bool
	Seed                int64
	Repeat              int
	MaxInstructions     int
	Natives             map[string]backend.NativeFunc
	TagExpression       string
	SourcePath          string
	UpdateSnapshots     bool
	Coverage            *Coverage
	FuzzSeed            int64
	Providers           []string
}

type suiteRunContext struct {
	fixtures  map[string]backend.Public
	snapshots *snapshotStore
	scenarios *scenarioRegistry
	baseline  []byte
	path      string
	strict    bool
	providers *providerSet
}

func (r Runner) List(path string) ([]backend.Public, error) {
	vm, publics, err := r.loadVM(path)
	if err != nil {
		return nil, err
	}
	defer vm.Close()

	return r.selectTests(publics)
}

func (r Runner) RunFile(path string) (Suite, error) {
	vm, publics, err := r.loadVM(path)
	if err != nil {
		return Suite{}, err
	}
	defer vm.Close()

	providers, err := loadProviders(r.backend(), r.Providers, r.MaxInstructions)
	if err != nil {
		return Suite{}, err
	}
	defer providers.Close()

	if r.Coverage != nil {
		r.Coverage.instrument(vm)
	}

	tests, err := r.selectTests(publics)
	if err != nil {
		return Suite{}, err
	}

	if limiter, ok := vm.(backend.InstructionLimiter); ok && r.MaxInstructions > 0 {
		limiter.SetInstructionLimit(r.MaxInstructions)
	}

	runs := r.testRuns(tests)

	source := r.SourcePath
	if source == "" {
		source = path
	}

	sc := suiteRunContext{
		fixtures:  findFixtures(publics),
		snapshots: newSnapshotStore(source, r.UpdateSnapshots),
		scenarios: newScenarioRegistry(),
		path:      path,
		strict:    hasPublic(publics, "__pawntest_strict_scenarios"),
		providers: providers,
	}
	defer sc.scenarios.Close()

	suiteContext := newExecutionContext(sc.snapshots, sc.scenarios, r)
	suiteContext.strict = sc.strict
	suiteContext.providers = providers
	suite := Suite{}

	if fixture, ok := sc.fixtures["test_suite_setup"]; ok {
		if failed := r.runSuiteSetup(vm, fixture, suiteContext, runs, sc); failed != nil {
			suite.Results = append(suite.Results, failed...)
			return suite, nil
		}
	}

	if snapshotter, ok := vm.(backend.MemorySnapshotter); ok && r.Isolation != "suite" {
		sc.baseline = snapshotter.SnapshotMemory()
	}
	providers.snapshot()

	for _, run := range runs {
		result, err := r.runTest(vm, run, sc)
		if err != nil {
			return suite, err
		}

		suite.Results = append(suite.Results, result)
		if r.FailFast && shouldStop(result.Status) {
			break
		}
	}

	if cleanup := r.runNamedTeardowns(vm, suiteContext, "test_suite_teardown", path); cleanup != nil {
		suite.Results = append(suite.Results, *cleanup)
	}

	if fixture, ok := sc.fixtures["test_suite_teardown"]; ok {
		state, err := r.exec(vm, fixture, suiteContext)
		if err != nil || state.status != Pass {
			suite.Results = append(suite.Results, resultFromState("test_suite_teardown", path, state, err))
		}
	}

	if sc.strict {
		for _, message := range sc.scenarios.StrictFailures() {
			suite.Results = append(suite.Results, Result{Name: "strict scenarios", File: path, Status: Fail, Message: message})
		}
	}

	if len(suite.Results) == 0 {
		return suite, ErrNoTestsFound
	}

	return suite, nil
}

func (r Runner) backend() backend.Backend {
	if r.Backend != nil {
		return r.Backend
	}

	return backend.NewGoAMXBackend()
}

func (r Runner) loadVM(path string) (backend.VM, []backend.Public, error) {
	vm, err := r.backend().LoadFile(path)
	if err != nil {
		return nil, nil, err
	}

	publics, err := vm.Publics()
	if err != nil {
		return nil, nil, errors.Join(err, vm.Close())
	}

	if err := requirePawntestInclude(publics, path); err != nil {
		return nil, nil, errors.Join(err, vm.Close())
	}

	return vm, publics, nil
}

func (r Runner) runSuiteSetup(vm backend.VM, fixture backend.Public, context *executionContext, runs []testRun, sc suiteRunContext) []Result {
	state, err := r.exec(vm, fixture, context)
	if err == nil && state.status == Pass {
		return nil
	}

	message := fixtureMessage(state, err)

	results := make([]Result, 0, len(runs)+2)
	for _, run := range runs {
		results = append(results, Result{Name: run.name, File: sc.path, Status: Error, Message: message})
	}

	if cleanup := r.runNamedTeardowns(vm, context, "test_suite_teardown", sc.path); cleanup != nil {
		results = append(results, *cleanup)
	}

	if teardown, ok := sc.fixtures["test_suite_teardown"]; ok {
		teardownState, teardownErr := r.exec(vm, teardown, context)
		if teardownErr != nil || teardownState.status != Pass {
			results = append(results, resultFromState("test_suite_teardown", sc.path, teardownState, teardownErr))
		}
	}

	return results
}

func (r Runner) runTest(vm backend.VM, run testRun, sc suiteRunContext) (Result, error) {
	sc.snapshots.test = run.public.Name
	if sc.baseline != nil {
		snapshotter, ok := vm.(backend.MemorySnapshotter)
		if !ok {
			return Result{}, errors.New("runtime does not support memory snapshots")
		}

		if err := snapshotter.RestoreMemory(sc.baseline); err != nil {
			return Result{}, err
		}
	}

	scenarios := sc.scenarios
	if r.Isolation != "suite" {
		scenarios = sc.scenarios.Clone()
		defer scenarios.Close()
	}

	context := newExecutionContext(sc.snapshots, scenarios, r)
	context.strict = sc.strict
	context.providers = sc.providers
	start := time.Now()
	result := Result{Name: run.name, File: sc.path, Status: Pass}
	setupPassed := true

	if err := sc.providers.beforeTest(r.Isolation, run.name); err != nil {
		setupPassed = false
		result = mergePhase(result, "provider setup", Result{Name: run.name, File: sc.path, Status: Error, Message: err.Error()})
	}

	if fixture, ok := sc.fixtures["test_setup"]; ok {
		state, err := r.exec(vm, fixture, context)
		if err != nil || state.status != Pass {
			setupPassed = false
			result = mergePhase(result, "setup", resultFromState(run.name, sc.path, state, err))
		}
	}

	if setupPassed {
		state, err := r.exec(vm, run.public, context)
		result = mergePhase(result, "test", resultFromState(run.name, sc.path, state, err))
	}

	if cleanup := r.runNamedTeardowns(vm, context, run.name, sc.path); cleanup != nil {
		result = mergePhase(result, "named teardown", *cleanup)
	}

	if fixture, ok := sc.fixtures["test_teardown"]; ok {
		teardownState, teardownErr := r.exec(vm, fixture, context)
		if teardownErr != nil || teardownState.status != Pass {
			result = mergePhase(result, "teardown", resultFromState(run.name, sc.path, teardownState, teardownErr))
		}
	}

	if message, file, line := context.mocks.verify(); message != "" {
		result = mergePhase(result, "mock verification", Result{Name: run.name, File: file, Line: line, Status: Fail, Message: message})
	}

	if context.strict && r.Isolation != "suite" {
		for _, message := range scenarios.StrictFailures() {
			result = mergePhase(result, "strict scenarios", Result{Name: run.name, File: sc.path, Status: Fail, Message: message})
		}
	}

	if err := sc.providers.afterTest(run.name); err != nil {
		result = mergePhase(result, "provider teardown", Result{Name: run.name, File: sc.path, Status: Error, Message: err.Error()})
	}

	result.Duration = time.Since(start)

	return result, nil
}

func hasPublic(publics []backend.Public, name string) bool {
	for _, public := range publics {
		if public.Name == name {
			return true
		}
	}

	return false
}

func (r Runner) exec(vm backend.VM, pub backend.Public, context *executionContext) (*nativeState, error) {
	state := &nativeState{status: Pass}
	context.state = state

	context.publicName = pub.Name
	if err := r.registerModules(vm, context); err != nil {
		return state, err
	}

	_, err := vm.ExecPublic(pub.Index)
	if err != nil {
		applyRuntimeLocation(vm, state, err)
	}

	return state, err
}

func (r Runner) registerModules(vm backend.VM, context *executionContext) error {
	for _, module := range defaultNativeModules() {
		if err := module.Register(vm, context); err != nil {
			return err
		}
	}

	return nil
}

func (r Runner) runNamedTeardowns(vm backend.VM, context *executionContext, name, path string) *Result {
	if len(context.fixtures.teardowns) == 0 {
		return nil
	}

	state := &nativeState{status: Pass}

	context.state = state
	if err := r.registerModules(vm, context); err != nil {
		result := Result{Name: name, File: path, Status: Error, Message: err.Error()}
		return &result
	}

	err := context.fixtures.runTeardowns(vm)
	if err == nil && state.status == Pass {
		return nil
	}

	result := resultFromState(name, path, state, err)

	return &result
}

func (r Runner) selectTests(publics []backend.Public) ([]backend.Public, error) {
	tests := discovery.TestPublics(publics)

	expression, err := parseTagExpression(r.TagExpression)
	if err != nil {
		return nil, err
	}

	if expression != nil {
		tags := testTags(publics)

		filtered := tests[:0]
		for _, test := range tests {
			if expression.matches(tags[test.Name]) {
				filtered = append(filtered, test)
			}
		}

		tests = filtered
	}

	if r.Run != "" {
		re, err := regexp.Compile(r.Run)
		if err != nil {
			return nil, fmt.Errorf("invalid run regex: %w", err)
		}

		filtered := tests[:0]
		for _, test := range tests {
			if re.MatchString(test.Name) {
				filtered = append(filtered, test)
			}
		}

		tests = filtered
	}

	if r.Shuffle {
		seed := r.Seed
		if seed == 0 {
			seed = 1
		}

		rand.New(rand.NewSource(seed)).Shuffle(len(tests), func(i, j int) { tests[i], tests[j] = tests[j], tests[i] })
	}

	return tests, nil
}

type testRun struct {
	public backend.Public
	name   string
}

func (r Runner) testRuns(tests []backend.Public) []testRun {
	repeat := r.Repeat
	if repeat <= 0 {
		repeat = 1
	}

	runs := make([]testRun, 0, len(tests)*repeat)
	for attempt := 1; attempt <= repeat; attempt++ {
		for _, test := range tests {
			name := test.Name
			if repeat > 1 {
				name = fmt.Sprintf("%s [attempt %d/%d]", name, attempt, repeat)
			}

			runs = append(runs, testRun{public: test, name: name})
		}
	}

	return runs
}

func applyRuntimeLocation(vm backend.VM, state *nativeState, err error) {
	var runtimeErr goamx.RuntimeError

	locator, ok := vm.(backend.DebugLocator)
	if !ok || !errors.As(err, &runtimeErr) {
		return
	}

	file, line, function, ok := locator.DebugLocation(runtimeErr.CIP)
	if !ok {
		return
	}

	state.file, state.line = file, line
	if function != "" {
		state.message = function + ": " + err.Error()
	}
}

func findFixtures(publics []backend.Public) map[string]backend.Public {
	out := map[string]backend.Public{}

	for _, pub := range publics {
		switch pub.Name {
		case "test_suite_setup", "test_suite_teardown", "test_setup", "test_teardown":
			out[pub.Name] = pub
		}
	}

	return out
}

func requirePawntestInclude(publics []backend.Public, path string) error {
	for _, pub := range publics {
		if pub.Name == markerPublic {
			return nil
		}
	}

	return fmt.Errorf("%s must include <pawntest>: %w", path, ErrMissingPawntestInclude)
}

func shouldStop(status Status) bool {
	return status == Fail || status == Error || status == XPass
}

func resultFromState(name, path string, state *nativeState, err error) Result {
	result := Result{Name: name, File: path, Status: state.status, Message: state.message, Line: state.line}
	if state.file != "" {
		result.File = state.file
	}

	if err != nil {
		result.Status = Error
		if state.message == "" {
			result.Message = err.Error()
		}
	}

	return result
}

func fixtureMessage(state *nativeState, err error) string {
	if err != nil {
		return err.Error()
	}

	if state.message != "" {
		return state.message
	}

	return "suite setup failed"
}
