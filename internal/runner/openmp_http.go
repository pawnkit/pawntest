package runner

import (
	"errors"
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

type httpResponse struct {
	code        backend.Cell
	bodyAddress backend.Cell
}

type httpRequest struct {
	index  backend.Cell
	method backend.Cell
	url    string
	data   string
}

type httpResponseKey struct {
	method backend.Cell
	url    string
}

type httpState struct {
	responses map[httpResponseKey][]httpResponse
	requests  []httpRequest
	unmatched []httpRequest
}

const (
	httpGet  backend.Cell = 1
	httpPost backend.Cell = 2
	httpHead backend.Cell = 3
)

func newHTTPState() *httpState {
	return &httpState{responses: map[httpResponseKey][]httpResponse{}}
}

func (state *httpState) Clone() scenarioModule {
	clone := newHTTPState()
	for key, responses := range state.responses {
		clone.responses[key] = append([]httpResponse(nil), responses...)
	}
	clone.requests = append([]httpRequest(nil), state.requests...)
	clone.unmatched = append([]httpRequest(nil), state.unmatched...)

	return clone
}

func (state *httpState) Register(vm backend.VM, context *executionContext) error {
	natives := map[string]backend.NativeFunc{
		"HTTP":                      state.request,
		"__pt_http_response":        state.addResponse,
		"__pt_http_method_response": state.addMethodResponse,
		"__pt_http_requests":        state.assertRequestCount(context.state),
		"__pt_http_request":         state.assertRequest(context.state),
	}

	return registerScenarioNatives(vm, natives, context.mocks, context.allowUnknown)
}

func (state *httpState) assertRequest(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		reader := readNativeParams(ctx, params)
		if err := reader.Require(5, "HTTP request assertion"); err != nil {
			return 0, err
		}
		url, err := reader.String(1)
		if err != nil {
			return 0, err
		}
		data, err := reader.String(2)
		if err != nil {
			return 0, err
		}
		for _, request := range state.requests {
			if request.method == params[0] && request.url == url && request.data == data {
				return 1, nil
			}
		}
		setFailure(result, params, 3, fmt.Sprintf("HTTP request not found: method %d, URL %q, data %q", params[0], url, data), ctx)

		return 0, nil
	}
}

func (state *httpState) addResponse(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	reader := readNativeParams(ctx, params)
	if err := reader.Require(3, "HTTP response"); err != nil {
		return 0, err
	}
	url, err := reader.String(0)
	if err != nil {
		return 0, err
	}
	state.queueResponse(httpResponseKey{url: url}, params[1], params[2])

	return 1, nil
}

func (state *httpState) addMethodResponse(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	reader := readNativeParams(ctx, params)
	if err := reader.Require(4, "HTTP method response"); err != nil {
		return 0, err
	}
	if params[0] < httpGet || params[0] > httpHead {
		return 0, nil
	}
	url, err := reader.String(1)
	if err != nil {
		return 0, err
	}
	state.queueResponse(httpResponseKey{method: params[0], url: url}, params[2], params[3])

	return 1, nil
}

func (state *httpState) queueResponse(key httpResponseKey, code, bodyAddress backend.Cell) {
	state.responses[key] = append(state.responses[key], httpResponse{code: code, bodyAddress: bodyAddress})
}

func (state *httpState) request(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return 0, nil
	}
	if params[1] < httpGet || params[1] > httpHead {
		return 0, nil
	}
	reader := readNativeParams(ctx, params)
	url, err := reader.String(2)
	if err != nil {
		return 0, err
	}
	data, err := reader.String(3)
	if err != nil {
		return 0, err
	}
	callback, err := reader.String(4)
	if err != nil {
		return 0, err
	}
	state.requests = append(state.requests, httpRequest{index: params[0], method: params[1], url: url, data: data})

	key, responses := state.responseQueue(params[1], url)
	if len(responses) == 0 {
		state.unmatched = append(state.unmatched, state.requests[len(state.requests)-1])

		return 0, nil
	}
	response := responses[0]
	state.responses[key] = responses[1:]
	if params[1] == httpHead {
		response.bodyAddress = params[3]
	}
	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support HTTP callbacks")
	}
	if _, err := caller.CallPublic(callback, params[0], response.code, response.bodyAddress); err != nil {
		return 0, fmt.Errorf("HTTP callback: %w", err)
	}

	return 1, nil
}

func (state *httpState) StrictFailures() []string {
	failures := make([]string, 0, len(state.unmatched)+1)

	for _, request := range state.unmatched {
		failures = append(failures, fmt.Sprintf("unconfigured HTTP request: method %d, URL %q", request.method, request.url))
	}

	pending := 0
	for _, responses := range state.responses {
		pending += len(responses)
	}
	if pending != 0 {
		failures = append(failures, fmt.Sprintf("unused HTTP responses: %d", pending))
	}

	return failures
}

func (state *httpState) responseQueue(method backend.Cell, url string) (httpResponseKey, []httpResponse) {
	key := httpResponseKey{method: method, url: url}
	if responses := state.responses[key]; len(responses) != 0 {
		return key, responses
	}

	key.method = 0

	return key, state.responses[key]
}

func (state *httpState) assertRequestCount(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if err := readNativeParams(ctx, params).Require(3, "HTTP request count assertion"); err != nil {
			return 0, err
		}
		actual := len(state.requests)
		if actual != int(params[0]) {
			setFailure(result, params, 1, fmt.Sprintf("HTTP requests: expected %d, got %d", params[0], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
