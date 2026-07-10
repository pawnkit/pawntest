package runner

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

func registerMockControlNatives(vm backend.VM, mocks *mockState) error {
	natives := map[string]backend.NativeFunc{
		"__pawntest_mock_return":        nativeMockReturn(mocks),
		"__pawntest_mock_return_once":   nativeMockReturnOnce(mocks),
		"__pawntest_mock_output":        nativeMockOutput(mocks),
		"__pawntest_mock_output_string": nativeMockOutputString(mocks),
		"__pawntest_mock_callback":      nativeMockCallback(mocks),
		"__pawntest_mock_reset":         nativeMockReset(mocks),
		"__pawntest_expect_calls":       nativeExpectCalls(mocks),
		"__pawntest_expect_arg":         nativeExpectArg(mocks),
		"__pawntest_expect_string_arg":  nativeExpectStringArg(mocks),
		"__pawntest_expect_order":       nativeExpectOrder(mocks),
	}
	for name, fn := range natives {
		if err := vm.RegisterNative(name, fn); err != nil {
			return err
		}
	}

	return nil
}

func registerUnknownNativeMocks(vm backend.VM, mocks *mockState, allowUnmocked bool) error {
	natives, err := vm.Natives()
	if err != nil {
		return err
	}

	for _, name := range natives {
		if isPawntestNative(name) || isFloatNative(name) {
			continue
		}

		if err := vm.RegisterNative(name, mockUnknownNative(name, mocks, allowUnmocked)); err != nil {
			return err
		}
	}

	return nil
}

func mockUnknownNative(name string, mocks *mockState, allowUnmocked bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		configured := mocks.configured(name)

		value, mocked := mocks.returnValue(name)
		if !mocked && !configured && !allowUnmocked {
			return 0, fmt.Errorf(
				"native %s is not mocked; call MOCK_RETURN(%s, value) before invoking it or use --allow-unknown-natives",
				name,
				name,
			)
		}

		mocks.recordCall(name, ctx, params)

		if err := mocks.applyOutputs(name, ctx, params); err != nil {
			return 0, err
		}

		if err := mocks.invokeCallback(name, ctx, params); err != nil {
			return 0, err
		}

		return value, nil
	}
}

type mockState struct {
	returns       map[string]backend.Cell
	returnQueues  map[string][]backend.Cell
	calls         map[string]int
	args          map[string][][]backend.Cell
	outputs       map[string]map[int]backend.Cell
	stringOutputs map[string]map[int]mockStringOutput
	callbacks     map[string]string
	expectations  []*mockExpectation
	expectedOrder []mockLocation
	actualOrder   []string
}

type mockExpectation struct {
	kind                 string
	name                 string
	minimum, maximum     int
	callIndex, argIndex  int
	value                backend.Cell
	stringValue          string
	actualString         string
	actualStringCaptured bool
	mockLocation
}

type mockLocation struct {
	name string
	file string
	line int
}

type mockStringOutput struct {
	value     string
	maxLength int
}

func truncateString(value string, maxLength int) string {
	if maxLength <= 1 {
		return ""
	}

	runes := []rune(value)
	if len(runes) >= maxLength {
		runes = runes[:maxLength-1]
	}

	return string(runes)
}

func newMockState() *mockState {
	return &mockState{
		returns:       map[string]backend.Cell{},
		returnQueues:  map[string][]backend.Cell{},
		calls:         map[string]int{},
		args:          map[string][][]backend.Cell{},
		outputs:       map[string]map[int]backend.Cell{},
		stringOutputs: map[string]map[int]mockStringOutput{},
		callbacks:     map[string]string{},
	}
}

func (mocks *mockState) returnValue(name string) (backend.Cell, bool) {
	if queue := mocks.returnQueues[name]; len(queue) > 0 {
		value := queue[0]
		mocks.returnQueues[name] = queue[1:]

		return value, true
	}

	value, ok := mocks.returns[name]

	return value, ok
}

func (mocks *mockState) configured(name string) bool {
	_, hasReturn := mocks.returns[name]

	return hasReturn || len(mocks.returnQueues[name]) > 0 ||
		len(mocks.outputs[name]) > 0 || len(mocks.stringOutputs[name]) > 0 || mocks.callbacks[name] != ""
}

func (mocks *mockState) recordCall(name string, ctx backend.NativeContext, params []backend.Cell) {
	callIndex := len(mocks.args[name])
	mocks.calls[name]++
	mocks.args[name] = append(mocks.args[name], slices.Clone(params))

	mocks.actualOrder = append(mocks.actualOrder, name)
	for _, expectation := range mocks.expectations {
		if expectation.kind != "string_arg" || expectation.name != name || expectation.callIndex != callIndex {
			continue
		}

		if expectation.argIndex < 0 || expectation.argIndex >= len(params) {
			continue
		}

		value, err := ctx.ReadString(params[expectation.argIndex])
		if err == nil {
			expectation.actualString = value
			expectation.actualStringCaptured = true
		}
	}
}

func (mocks *mockState) verify() (string, string, int) {
	for _, expectation := range mocks.expectations {
		message := ""

		switch expectation.kind {
		case "calls":
			message = mocks.verifyCallCount(expectation)
		case "arg":
			message = mocks.verifyCellArgument(expectation)
		case "string_arg":
			message = mocks.verifyStringArgument(expectation)
		}

		if message != "" {
			return message, expectation.file, expectation.line
		}
	}

	return mocks.verifyCallOrder()
}

func (mocks *mockState) verifyCallCount(expectation *mockExpectation) string {
	count := mocks.calls[expectation.name]
	if count >= expectation.minimum && count <= expectation.maximum {
		return ""
	}

	return fmt.Sprintf(
		"expected %s to be called %d..%d times, got %d\nrecorded calls: %s",
		expectation.name,
		expectation.minimum,
		expectation.maximum,
		count,
		mocks.describeCalls(expectation.name),
	)
}

func (mocks *mockState) verifyCellArgument(expectation *mockExpectation) string {
	calls := mocks.args[expectation.name]
	if expectation.callIndex < 0 || expectation.callIndex >= len(calls) ||
		expectation.argIndex < 0 || expectation.argIndex >= len(calls[expectation.callIndex]) {
		return fmt.Sprintf(
			"expected argument %d of %s call %d, but the call or argument was missing\nrecorded calls: %s",
			expectation.argIndex,
			expectation.name,
			expectation.callIndex,
			mocks.describeCalls(expectation.name),
		)
	}

	call := calls[expectation.callIndex]

	actual := call[expectation.argIndex]
	if actual == expectation.value {
		return ""
	}

	return fmt.Sprintf(
		"expected argument %d of %s call %d to be %d, got %d\nrecorded call: %s",
		expectation.argIndex,
		expectation.name,
		expectation.callIndex,
		expectation.value,
		actual,
		describeCall(call),
	)
}

func (mocks *mockState) verifyStringArgument(expectation *mockExpectation) string {
	if expectation.actualStringCaptured && expectation.actualString == expectation.stringValue {
		return ""
	}

	return fmt.Sprintf(
		"expected string argument %d of %s call %d to be %q, got %q\nrecorded calls: %s",
		expectation.argIndex,
		expectation.name,
		expectation.callIndex,
		expectation.stringValue,
		expectation.actualString,
		mocks.describeCalls(expectation.name),
	)
}

func (mocks *mockState) verifyCallOrder() (string, string, int) {
	for index, expected := range mocks.expectedOrder {
		if index >= len(mocks.actualOrder) || mocks.actualOrder[index] != expected.name {
			actual := "<missing>"
			if index < len(mocks.actualOrder) {
				actual = mocks.actualOrder[index]
			}

			return fmt.Sprintf("expected native call %d to be %s, got %s", index, expected.name, actual), expected.file, expected.line
		}
	}

	return "", "", 0
}

func (mocks *mockState) describeCalls(name string) string {
	calls := mocks.args[name]
	if len(calls) == 0 {
		return "none"
	}

	parts := make([]string, 0, len(calls))
	for index, call := range calls {
		parts = append(parts, fmt.Sprintf("#%d %s", index, describeCall(call)))
	}

	return strings.Join(parts, "; ")
}

func describeCall(args []backend.Cell) string {
	parts := make([]string, len(args))
	for index, arg := range args {
		parts[index] = fmt.Sprint(arg)
	}

	return "(" + strings.Join(parts, ", ") + ")"
}

func (mocks *mockState) applyOutputs(name string, ctx backend.NativeContext, params []backend.Cell) error {
	for index, output := range mocks.outputs[name] {
		if index >= 0 && index < len(params) {
			if err := ctx.WriteCell(params[index], output); err != nil {
				return err
			}
		}
	}

	for index, output := range mocks.stringOutputs[name] {
		if index >= 0 && index < len(params) {
			value := truncateString(output.value, output.maxLength)
			if err := ctx.WriteString(params[index], value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (mocks *mockState) invokeCallback(name string, ctx backend.NativeContext, params []backend.Cell) error {
	callback := mocks.callbacks[name]
	if callback == "" {
		return nil
	}

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return errors.New("runtime does not support mock callbacks")
	}

	if _, err := caller.CallPublic(callback, params...); err != nil {
		return fmt.Errorf("mock callback %s: %w", callback, err)
	}

	return nil
}

func nativeMockReturn(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 2 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		mocks.returns[name] = params[1]

		return 1, nil
	}
}

func nativeMockReturnOnce(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 2 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		mocks.returnQueues[name] = append(mocks.returnQueues[name], params[1])

		return 1, nil
	}
}

func nativeMockOutput(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		if mocks.outputs[name] == nil {
			mocks.outputs[name] = map[int]backend.Cell{}
		}

		mocks.outputs[name][int(params[1])] = params[2]

		return 1, nil
	}
}

func nativeMockOutputString(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 4 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		value, err := ctx.ReadString(params[3])
		if err != nil {
			return 0, err
		}

		if mocks.stringOutputs[name] == nil {
			mocks.stringOutputs[name] = map[int]mockStringOutput{}
		}

		mocks.stringOutputs[name][int(params[1])] = mockStringOutput{value: value, maxLength: int(params[2])}

		return 1, nil
	}
}

func nativeMockCallback(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 2 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		callback, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		mocks.callbacks[name] = callback

		return 1, nil
	}
}

func nativeExpectCalls(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, nil
		}

		expectation, err := expectationLocation(ctx, params, 0, 3, 4)
		if err != nil {
			return 0, err
		}

		expectation.kind = "calls"
		expectation.minimum = int(params[1])
		expectation.maximum = int(params[2])
		mocks.expectations = append(mocks.expectations, expectation)

		return 1, nil
	}
}

func nativeExpectArg(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 {
			return 0, nil
		}

		expectation, err := expectationLocation(ctx, params, 0, 4, 5)
		if err != nil {
			return 0, err
		}

		expectation.kind = "arg"
		expectation.callIndex = int(params[1])
		expectation.argIndex = int(params[2])
		expectation.value = params[3]
		mocks.expectations = append(mocks.expectations, expectation)

		return 1, nil
	}
}

func nativeExpectStringArg(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 6 {
			return 0, nil
		}

		expectation, err := expectationLocation(ctx, params, 0, 4, 5)
		if err != nil {
			return 0, err
		}

		expectation.kind = "string_arg"
		expectation.callIndex = int(params[1])
		expectation.argIndex = int(params[2])

		expectation.stringValue, err = ctx.ReadString(params[3])
		if err != nil {
			return 0, err
		}

		mocks.expectations = append(mocks.expectations, expectation)

		return 1, nil
	}
}

func nativeExpectOrder(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		file, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}

		mocks.expectedOrder = append(mocks.expectedOrder, mockLocation{name: name, file: file, line: int(params[2])})

		return 1, nil
	}
}

func expectationLocation(ctx backend.NativeContext, params []backend.Cell, nameIndex, fileIndex, lineIndex int) (*mockExpectation, error) {
	name, err := ctx.ReadString(params[nameIndex])
	if err != nil {
		return nil, err
	}

	file, err := ctx.ReadString(params[fileIndex])
	if err != nil {
		return nil, err
	}

	return &mockExpectation{name: name, mockLocation: mockLocation{file: file, line: int(params[lineIndex])}}, nil
}

func nativeMockReset(mocks *mockState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 1 {
			mocks.calls = map[string]int{}
			mocks.returns = map[string]backend.Cell{}
			mocks.returnQueues = map[string][]backend.Cell{}
			mocks.args = map[string][][]backend.Cell{}
			mocks.outputs = map[string]map[int]backend.Cell{}
			mocks.stringOutputs = map[string]map[int]mockStringOutput{}
			mocks.callbacks = map[string]string{}

			return 1, nil
		}

		name, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		delete(mocks.calls, name)
		delete(mocks.returns, name)
		delete(mocks.returnQueues, name)
		delete(mocks.args, name)
		delete(mocks.outputs, name)
		delete(mocks.stringOutputs, name)
		delete(mocks.callbacks, name)

		return 1, nil
	}
}
