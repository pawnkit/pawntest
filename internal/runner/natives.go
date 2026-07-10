package runner

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

type nativeState struct {
	status   Status
	message  string
	file     string
	line     int
	failures []string
}

func registerPawntestNatives(vm backend.VM, state *nativeState, mocks *mockState) error {
	natives := map[string]backend.NativeFunc{
		"__pawntest_assert":            nativeAssert(state),
		"__pawntest_assert_eq":         nativeAssertEQ(state),
		"__pawntest_assert_ne":         nativeAssertNE(state),
		"__pawntest_assert_str_eq":     nativeAssertStrEQ(state),
		"__pawntest_assert_compare":    nativeAssertCompare(state),
		"__pawntest_assert_string":     nativeAssertString(state),
		"__pawntest_assert_float_near": nativeAssertFloatNear(state),
		"__pawntest_assert_array_eq":   nativeAssertArrayEQ(state),
		"__pawntest_assert_between":    nativeAssertBetween(state),
		"__pawntest_assert_has_flag":   nativeAssertHasFlag(state),
		"__pt_assert_array_contains":   nativeAssertArrayContains(state),
		"__pawntest_assert_str_ieq":    nativeAssertStrIEQ(state),
		"__pawntest_assert_eq_message": nativeAssertEQMessage(state),
		"__pawntest_expect_error":      nativeExpectError(state),
		"__pawntest_expect_no_error":   nativeExpectNoError(state),
		"__pawntest_xfail":             nativeXFail(state),
		"__pawntest_fail":              nativeFail(state),
		"__pawntest_skip":              nativeSkip(state),
	}
	for name, fn := range natives {
		if err := vm.RegisterNative(name, fn); err != nil {
			return err
		}
	}

	return nil
}

func nativeAssertBetween(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 || params[0] < params[1] || params[0] > params[2] {
			message := fmt.Sprintf("expected: %s between %d and %d\nactual:   %d", readStringParam(ctx, params, 3), params[1], params[2], params[0])
			setFailure(state, params, 4, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeAssertHasFlag(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 || params[0]&params[1] != params[1] {
			message := fmt.Sprintf("expected: %s has flag %#x\nactual:   %#x", readStringParam(ctx, params, 2), uint32(params[1]), uint32(params[0]))
			setFailure(state, params, 3, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeAssertArrayContains(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 || params[1] < 0 {
			setFailure(state, params, 4, "array membership received invalid parameters", ctx)
			return 0, nil
		}

		for index := range params[1] {
			value, err := ctx.ReadCell(params[0] + index*4)
			if err != nil {
				return 0, err
			}

			if value == params[2] {
				return 1, nil
			}
		}

		message := fmt.Sprintf("expected: %s contains %d\nactual:   value was not present", readStringParam(ctx, params, 3), params[2])
		setFailure(state, params, 4, message, ctx)

		return 0, nil
	}
}

func nativeAssertStrIEQ(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 {
			setFailure(state, params, 4, "case-insensitive string comparison received too few parameters", ctx)
			return 0, nil
		}

		actual, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		expected, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		if !strings.EqualFold(actual, expected) {
			message := fmt.Sprintf("expected: %s = %q (ignoring case)\nactual:   %s = %q", readStringParam(ctx, params, 3), expected, readStringParam(ctx, params, 2), actual)
			setFailure(state, params, 4, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeAssertEQMessage(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 || params[0] != params[1] {
			context := readStringParam(ctx, params, 4)
			message := fmt.Sprintf("%s\nexpected: %s = %d\nactual:   %s = %d", context, readStringParam(ctx, params, 3), params[1], readStringParam(ctx, params, 2), params[0])
			setFailure(state, params, 5, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeExpectNoError(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, nil
		}

		callback := readStringParam(ctx, params, 0)

		caller, ok := ctx.(backend.PublicCaller)
		if !ok {
			return 0, errors.New("runtime does not support no-error callbacks")
		}

		if _, err := caller.CallPublic(callback); err != nil {
			setFailure(state, params, 1, fmt.Sprintf("expected %s not to error: %v", callback, err), ctx)
			return 0, nil
		}

		return 1, nil
	}
}

func nativeXFail(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, nil
		}

		callback := readStringParam(ctx, params, 0)
		reason := readStringParam(ctx, params, 1)

		caller, ok := ctx.(backend.PublicCaller)
		if !ok {
			return 0, errors.New("runtime does not support expected-failure callbacks")
		}

		state.status, state.message, state.failures = Pass, "", nil

		_, err := caller.CallPublic(callback)
		if err != nil || state.status != Pass {
			detail := state.message
			if err != nil {
				detail = err.Error()
			}

			state.status = XFail

			state.message = reason
			if detail != "" {
				state.message += ": " + detail
			}

			return 1, nil
		}

		state.status = XPass
		state.message = "unexpected pass: " + reason
		state.file = readStringParam(ctx, params, 2)
		state.line = int(params[3])

		return 0, nil
	}
}

func nativeAssertCompare(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			setFailure(state, params, 5, "comparison received too few parameters", ctx)
			return 0, nil
		}

		actual, expected, comparison := params[0], params[1], params[2]

		passed := comparison == 1 && actual > expected ||
			comparison == 2 && actual >= expected ||
			comparison == 3 && actual < expected ||
			comparison == 4 && actual <= expected
		if passed {
			return 1, nil
		}

		operator := map[backend.Cell]string{1: ">", 2: ">=", 3: "<", 4: "<="}[comparison]
		actualExpr := readStringParam(ctx, params, 3)
		expectedExpr := readStringParam(ctx, params, 4)
		message := fmt.Sprintf("expected: %s %s %s\nactual:   %d %s %d", actualExpr, operator, expectedExpr, actual, operator, expected)
		setFailure(state, params, 5, message, ctx)

		return 0, nil
	}
}

func nativeAssertString(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			setFailure(state, params, 5, "string comparison received too few parameters", ctx)
			return 0, nil
		}

		actual, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		expected, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		comparison := params[2]

		passed := comparison == 1 && strings.Contains(actual, expected) ||
			comparison == 2 && strings.HasPrefix(actual, expected) ||
			comparison == 3 && strings.HasSuffix(actual, expected)
		if passed {
			return 1, nil
		}

		message := fmt.Sprintf("expected: string condition with %q\nactual:   %q", expected, actual)
		setFailure(state, params, 5, message, ctx)

		return 0, nil
	}
}

func nativeAssertFloatNear(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 {
			setFailure(state, params, 5, "float comparison received too few parameters", ctx)
			return 0, nil
		}

		actual := math.Float32frombits(uint32(params[0]))
		expected := math.Float32frombits(uint32(params[1]))

		tolerance := math.Abs(float64(math.Float32frombits(uint32(params[2]))))
		if math.Abs(float64(actual-expected)) <= tolerance {
			return 1, nil
		}

		message := fmt.Sprintf("expected: %s = %g +/- %g\nactual:   %s = %g", readStringParam(ctx, params, 4), expected, tolerance, readStringParam(ctx, params, 3), actual)
		setFailure(state, params, 5, message, ctx)

		return 0, nil
	}
}

func nativeAssertArrayEQ(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 7 || params[2] < 0 {
			setFailure(state, params, 5, "array comparison received invalid parameters", ctx)
			return 0, nil
		}

		for i := range params[2] {
			actual, err := ctx.ReadCell(params[0] + i*4)
			if err != nil {
				return 0, err
			}

			expected, err := ctx.ReadCell(params[1] + i*4)
			if err != nil {
				return 0, err
			}

			if actual != expected {
				message := fmt.Sprintf("arrays differ at index %d\nexpected: %d\nactual:   %d", i, expected, actual)
				setFailure(state, params, 5, message, ctx)

				return 0, nil
			}
		}

		return 1, nil
	}
}

func nativeExpectError(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, nil
		}

		callback := readStringParam(ctx, params, 0)
		expected := readStringParam(ctx, params, 1)

		caller, ok := ctx.(backend.PublicCaller)
		if !ok {
			return 0, errors.New("runtime does not support expected-error callbacks")
		}

		_, err := caller.CallPublic(callback)
		if err != nil && (expected == "" || strings.Contains(strings.ToLower(err.Error()), strings.ToLower(expected))) {
			return 1, nil
		}

		message := fmt.Sprintf("expected %s to fail with %q", callback, expected)
		if err != nil {
			message += "; got " + err.Error()
		}

		setFailure(state, params, 2, message, ctx)

		return 0, nil
	}
}

func isPawntestNative(name string) bool {
	return strings.HasPrefix(name, "__pawntest_") || strings.HasPrefix(name, "__pt_")
}

func nativeAssert(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 || params[0] == 0 {
			expr := readStringParam(ctx, params, 1)
			setFailure(state, params, 2, "assertion failed: "+expr, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeAssertEQ(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 || params[0] != params[1] {
			actual := readStringParam(ctx, params, 2)
			expected := readStringParam(ctx, params, 3)
			message := fmt.Sprintf("expected: %s\nactual:   %s", describeCell(expected, params[1]), describeCell(actual, params[0]))
			setFailure(state, params, 4, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func describeCell(expression string, value backend.Cell) string {
	valueText := fmt.Sprint(value)
	if strings.TrimSpace(expression) == valueText {
		return valueText
	}

	return fmt.Sprintf("%s (%s)", valueText, expression)
}

func nativeAssertNE(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 || params[0] == params[1] {
			actual := readStringParam(ctx, params, 2)
			expected := readStringParam(ctx, params, 3)
			setFailure(state, params, 4, fmt.Sprintf("expected: %s differs from %s\nactual:   both were %d", actual, expected, params[0]), ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeAssertStrEQ(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 {
			setFailure(state, params, 4, "assertion failed: string comparison received too few parameters", ctx)
			return 0, nil
		}

		actual, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		expected, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		if actual != expected {
			actualExpr := readStringParam(ctx, params, 2)
			expectedExpr := readStringParam(ctx, params, 3)
			message := fmt.Sprintf("expected: %s = %q\nactual:   %s = %q", expectedExpr, expected, actualExpr, actual)
			setFailure(state, params, 4, message, ctx)

			return 0, nil
		}

		return 1, nil
	}
}

func nativeFail(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		message := readStringParam(ctx, params, 0)
		setFailure(state, params, 1, message, ctx)

		return 0, nil
	}
}

func nativeSkip(state *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		state.status = Skip
		state.message = readStringParam(ctx, params, 0)

		state.file = readStringParam(ctx, params, 1)
		if len(params) > 2 {
			state.line = int(params[2])
		}

		return 0, nil
	}
}

func setFailure(state *nativeState, params []backend.Cell, fileParam int, message string, ctx backend.NativeContext) {
	state.status = Fail
	state.failures = append(state.failures, message)
	state.message = strings.Join(state.failures, "\n")

	state.file = readStringParam(ctx, params, fileParam)
	if len(params) > fileParam+1 {
		state.line = int(params[fileParam+1])
	}
}

func readStringParam(ctx backend.NativeContext, params []backend.Cell, index int) string {
	if index < 0 || index >= len(params) {
		return ""
	}

	value, err := ctx.ReadString(params[index])
	if err != nil {
		return fmt.Sprintf("<invalid string: %v>", err)
	}

	return value
}
