package runner

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

type scheduledCallback struct {
	due      int64
	sequence int64
	name     string
}

type scheduler struct {
	now      int64
	sequence int64
	pending  []scheduledCallback
}

func newScheduler() *scheduler { return &scheduler{} }

func registerSchedulerNatives(vm backend.VM, scheduler *scheduler) error {
	natives := map[string]backend.NativeFunc{
		"__pawntest_schedule":     nativeSchedule(scheduler),
		"__pawntest_advance_time": nativeAdvanceTime(scheduler),
		"__pawntest_run_pending":  nativeRunPending(scheduler),
		"__pawntest_now":          nativeNow(scheduler),
	}
	for name, fn := range natives {
		if err := vm.RegisterNative(name, fn); err != nil {
			return err
		}
	}

	return nil
}

func nativeSchedule(scheduler *scheduler) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 2 || params[0] < 0 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		scheduler.sequence++
		scheduler.pending = append(scheduler.pending, scheduledCallback{
			due: scheduler.now + int64(params[0]), sequence: scheduler.sequence, name: name,
		})
		scheduler.sort()

		return backend.Cell(scheduler.sequence), nil
	}
}

func nativeAdvanceTime(scheduler *scheduler) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 1 || params[0] < 0 {
			return 0, nil
		}

		scheduler.now += int64(params[0])
		if err := scheduler.runDue(ctx); err != nil {
			return 0, err
		}

		return backend.Cell(scheduler.now), nil
	}
}

func nativeRunPending(scheduler *scheduler) backend.NativeFunc {
	return func(ctx backend.NativeContext, _ []backend.Cell) (backend.Cell, error) {
		for len(scheduler.pending) > 0 {
			scheduler.now = scheduler.pending[0].due
			if err := scheduler.runDue(ctx); err != nil {
				return 0, err
			}
		}

		return backend.Cell(scheduler.now), nil
	}
}

func nativeNow(scheduler *scheduler) backend.NativeFunc {
	return func(backend.NativeContext, []backend.Cell) (backend.Cell, error) {
		return backend.Cell(scheduler.now), nil
	}
}

func (scheduler *scheduler) runDue(ctx backend.NativeContext) error {
	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return errors.New("runtime does not support scheduled callbacks")
	}

	for len(scheduler.pending) > 0 && scheduler.pending[0].due <= scheduler.now {
		callback := scheduler.pending[0]
		scheduler.pending = scheduler.pending[1:]

		if _, err := caller.CallPublic(callback.name); err != nil {
			return fmt.Errorf("scheduled callback %s: %w", callback.name, err)
		}
	}

	return nil
}

func (scheduler *scheduler) sort() {
	slices.SortStableFunc(scheduler.pending, func(a, b scheduledCallback) int {
		if a.due != b.due {
			return cmp.Compare(a.due, b.due)
		}

		return cmp.Compare(a.sequence, b.sequence)
	})
}
