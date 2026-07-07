package runner

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/pawnkit/pawntest/internal/backend"
)

func registerFuzzNative(vm backend.VM, state *nativeState, seed int64, testName string) error {
	return vm.RegisterNative("__pawntest_fuzz_int", func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 {
			return 0, nil
		}
		callback := readStringParam(ctx, params, 0)
		iterations := int(params[1])
		minimum, maximum := int64(params[2]), int64(params[3])
		if iterations <= 0 || minimum > maximum {
			setFailure(state, params, 4, "FUZZ_INT requires positive iterations and minimum <= maximum", ctx)
			return 0, nil
		}
		width := maximum - minimum + 1
		if width <= 0 {
			setFailure(state, params, 4, "FUZZ_INT range overflows int64", ctx)
			return 0, nil
		}
		caller, ok := ctx.(backend.PublicCaller)
		if !ok {
			return 0, fmt.Errorf("runtime does not support property-test callbacks")
		}
		testSeed := seedForTest(seed, testName)
		random := rand.New(rand.NewSource(testSeed))
		evaluate := func(value int64) (bool, string, error) {
			state.status, state.message, state.failures = Pass, "", nil
			_, err := caller.CallPublic(callback, backend.Cell(value))
			if err != nil {
				return true, err.Error(), nil
			}
			if state.status == Skip {
				return false, "", nil
			}
			return state.status != Pass, state.message, nil
		}

		var failed bool
		var value int64
		var detail string
		for iteration := 0; iteration < iterations; iteration++ {
			if width == 1 {
				value = minimum
			} else {
				value = minimum + random.Int63n(width)
			}
			var err error
			failed, detail, err = evaluate(value)
			if err != nil {
				return 0, err
			}
			if failed || state.status == Skip {
				break
			}
		}
		if state.status == Skip {
			return 0, nil
		}
		if !failed {
			state.status, state.message = Pass, ""
			return 1, nil
		}

		best := value
		var target int64
		if target < minimum {
			target = minimum
		} else if target > maximum {
			target = maximum
		}
		if target != best {
			targetFailed, targetDetail, err := evaluate(target)
			if err != nil {
				return 0, err
			}
			if targetFailed {
				best, detail = target, targetDetail
			} else {
				passing, failing := target, best
				for distance(passing, failing) > 1 {
					candidate := passing + (failing-passing)/2
					candidateFailed, candidateDetail, err := evaluate(candidate)
					if err != nil {
						return 0, err
					}
					if candidateFailed {
						failing, detail = candidate, candidateDetail
					} else {
						passing = candidate
					}
				}
				best = failing
			}
		}
		state.status = Fail
		state.message = fmt.Sprintf("property failed (seed %d, value %d): %s", testSeed, best, detail)
		state.file = readStringParam(ctx, params, 4)
		state.line = int(params[5])
		return 0, nil
	})
}

func distance(left, right int64) uint64 {
	if left < right {
		return uint64(right - left)
	}
	return uint64(left - right)
}

func seedForTest(seed int64, testName string) int64 {
	if seed == 0 {
		seed = 1
	}
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(testName))
	return seed ^ int64(hash.Sum64())
}
