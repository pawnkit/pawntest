package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

func assertScenarioCount(result *nativeState, scope string, value func() int) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New(scope + " count assertion expects 3 arguments")
		}

		actual := value()
		if actual != int(params[0]) {
			setFailure(result, params, 1, fmt.Sprintf("%s: expected %d, got %d", scope, params[0], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
