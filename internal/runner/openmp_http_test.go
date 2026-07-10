package runner

import (
	"slices"
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestHTTPScenarioReturnsConfiguredResponse(t *testing.T) {
	base := &mockVM{natives: map[string]backend.NativeFunc{}, strings: map[backend.Cell]string{
		100: "api.example.test/player", 200: `{"name":"Alice"}`, 300: "", 400: "OnHTTPResponse",
	}}
	vm := &dialogMockVM{mockVM: base}
	registry := newScenarioRegistry()
	context := &executionContext{state: &nativeState{status: Pass}, mocks: newMockState(), scenarios: registry}
	if err := registry.Register(vm, context); err != nil {
		t.Fatal(err)
	}

	callScenarioNative(t, base, "__pt_http_response", 100, 200, 200)
	result, err := base.natives["HTTP"](vm, []backend.Cell{7, 1, 100, 300, 400})
	if err != nil {
		t.Fatal(err)
	}
	if result != 1 || vm.callback != "OnHTTPResponse" {
		t.Fatalf("result=%d callback=%q", result, vm.callback)
	}
	if want := []backend.Cell{7, 200, 200}; !slices.Equal(vm.args, want) {
		t.Fatalf("callback args = %v, want %v", vm.args, want)
	}

	state := httpScenarioState(t, registry)
	if len(state.requests) != 1 || state.requests[0].url != "api.example.test/player" || state.requests[0].method != 1 {
		t.Fatalf("unexpected requests: %#v", state.requests)
	}
	if matched := callScenarioNative(t, base, "__pt_http_request", 1, 100, 300, 500, 1); matched != 1 {
		t.Fatal("HTTP request assertion did not match")
	}
	unmatched, err := base.natives["HTTP"](vm, []backend.Cell{8, 1, 100, 300, 400})
	if err != nil {
		t.Fatal(err)
	}
	if unmatched != 0 {
		t.Fatalf("unmatched HTTP = %d, want 0", unmatched)
	}
}

func TestHTTPScenarioCloneIsolatesResponses(t *testing.T) {
	state := newHTTPState()
	state.responses["example.test"] = []httpResponse{{code: 200, bodyAddress: 100}}
	state.requests = []httpRequest{{url: "example.test"}}

	clone, ok := state.Clone().(*httpState)
	if !ok {
		t.Fatal("cloned scenario was not HTTP state")
	}
	clone.responses["example.test"][0].code = 404
	clone.requests[0].url = "changed.test"
	if state.responses["example.test"][0].code != 200 || state.requests[0].url != "example.test" {
		t.Fatal("HTTP clone shared state")
	}
}

func httpScenarioState(t *testing.T, registry *scenarioRegistry) *httpState {
	t.Helper()

	state, ok := registry.modules[16].(*httpState)
	if !ok {
		t.Fatal("scenario registry did not contain HTTP state")
	}

	return state
}
