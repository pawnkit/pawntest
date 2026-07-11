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

type httpState struct {
	responses map[string][]httpResponse
	requests  []httpRequest
}

const (
	httpGet  backend.Cell = 1
	httpPost backend.Cell = 2
	httpHead backend.Cell = 3
)

func newHTTPState() *httpState {
	return &httpState{responses: map[string][]httpResponse{}}
}

func (state *httpState) Clone() scenarioModule {
	clone := newHTTPState()
	for url, responses := range state.responses {
		clone.responses[url] = append([]httpResponse(nil), responses...)
	}
	clone.requests = append([]httpRequest(nil), state.requests...)

	return clone
}

func (state *httpState) Register(vm backend.VM, context *executionContext) error {
	natives := map[string]backend.NativeFunc{
		"HTTP":               state.request,
		"__pt_http_response": state.addResponse,
		"__pt_http_requests": state.assertRequestCount(context.state),
		"__pt_http_request":  state.assertRequest(context.state),
	}

	return registerScenarioNatives(vm, natives, context.mocks, context.allowUnknown)
}

func (state *httpState) assertRequest(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 5 {
			return 0, errors.New("HTTP request assertion expects 5 arguments")
		}
		url, err := ctx.ReadString(params[1])
		if err != nil {
			return 0, err
		}
		data, err := ctx.ReadString(params[2])
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
	if len(params) < 3 {
		return 0, errors.New("HTTP response expects 3 arguments")
	}
	url, err := ctx.ReadString(params[0])
	if err != nil {
		return 0, err
	}
	state.responses[url] = append(state.responses[url], httpResponse{code: params[1], bodyAddress: params[2]})

	return 1, nil
}

func (state *httpState) request(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return 0, nil
	}
	if params[1] < httpGet || params[1] > httpHead {
		return 0, nil
	}
	url, err := ctx.ReadString(params[2])
	if err != nil {
		return 0, err
	}
	data, err := ctx.ReadString(params[3])
	if err != nil {
		return 0, err
	}
	callback, err := ctx.ReadString(params[4])
	if err != nil {
		return 0, err
	}
	state.requests = append(state.requests, httpRequest{index: params[0], method: params[1], url: url, data: data})

	responses := state.responses[url]
	if len(responses) == 0 {
		return 0, nil
	}
	response := responses[0]
	state.responses[url] = responses[1:]
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

func (state *httpState) assertRequestCount(result *nativeState) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if len(params) < 3 {
			return 0, errors.New("HTTP request count assertion expects 3 arguments")
		}
		actual := len(state.requests)
		if actual != int(params[0]) {
			setFailure(result, params, 1, fmt.Sprintf("HTTP requests: expected %d, got %d", params[0], actual), ctx)
			return 0, nil
		}

		return 1, nil
	}
}
