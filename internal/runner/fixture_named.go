package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

type namedFixtureState struct {
	teardowns []string
}

func registerNamedFixtureNative(vm backend.VM, state *nativeState, fixtures *namedFixtureState) error {
	return vm.RegisterNative("__pawntest_use_fixture", func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, nil
		}
		setup := readStringParam(ctx, params, 0)
		teardown := readStringParam(ctx, params, 1)
		caller, ok := ctx.(backend.PublicCaller)
		if !ok {
			return 0, fmt.Errorf("runtime does not support named fixtures")
		}
		if _, err := caller.CallPublic(setup); err != nil {
			setFailure(state, params, 2, fmt.Sprintf("fixture %s failed: %v", setup, err), ctx)
			return 0, nil
		}
		if state.status != Pass {
			return 0, nil
		}
		fixtures.teardowns = append(fixtures.teardowns, teardown)
		return 1, nil
	})
}

func (fixtures *namedFixtureState) runTeardowns(vm backend.VM) error {
	caller, ok := vm.(backend.PublicCaller)
	if !ok && len(fixtures.teardowns) > 0 {
		return fmt.Errorf("runtime does not support named fixture teardown")
	}
	var failures []error
	for len(fixtures.teardowns) > 0 {
		index := len(fixtures.teardowns) - 1
		name := fixtures.teardowns[index]
		fixtures.teardowns = fixtures.teardowns[:index]
		if _, err := caller.CallPublic(name); err != nil {
			failures = append(failures, fmt.Errorf("fixture %s failed: %w", name, err))
		}
	}
	return errors.Join(failures...)
}
