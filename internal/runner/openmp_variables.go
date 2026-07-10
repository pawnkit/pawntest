package runner

import (
	"maps"
	"slices"

	"github.com/pawnkit/pawntest/internal/backend"
)

const (
	variableNone = iota
	variableInt
	variableString
	variableFloat
)

type testVariable struct {
	variableType int
	integer      int
	text         string
	floating     float32
}

type variableScope struct {
	values map[string]testVariable
	order  []string
}

type variableState struct {
	server  variableScope
	players map[int]*variableScope
	model   *openMPState
}

func newVariableState() *variableState {
	return &variableState{server: newVariableScope(), players: map[int]*variableScope{}}
}

func newVariableScope() variableScope {
	return variableScope{values: map[string]testVariable{}}
}

func (state *variableState) Clone() scenarioModule {
	clone := newVariableState()

	clone.server = cloneVariableScope(&state.server)
	for playerID, scope := range state.players {
		scopeCopy := cloneVariableScope(scope)
		clone.players[playerID] = &scopeCopy
	}

	return clone
}

func cloneVariableScope(scope *variableScope) variableScope {
	return variableScope{values: maps.Clone(scope.values), order: append([]string(nil), scope.order...)}
}

func (state *variableState) Register(vm backend.VM, context *executionContext) error {
	state.model = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *variableState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_svar_int": state.assertInt(result, false), "__pt_svar_float": state.assertFloat(result, false),
		"__pt_svar_string": state.assertString(result, false), "__pt_pvar_int": state.assertInt(result, true),
		"__pt_pvar_float": state.assertFloat(result, true), "__pt_pvar_string": state.assertString(result, true),
		"SetSVarInt": state.setInt(false), "GetSVarInt": state.getInt(false),
		"SetSVarString": state.setString(false), "GetSVarString": state.getString(false),
		"SetSVarFloat": state.setFloat(false), "GetSVarFloat": state.getFloat(false),
		"DeleteSVar": state.deleteVariable(false), "GetSVarsUpperIndex": state.getUpperIndex(false),
		"GetSVarNameAtIndex": state.getNameAtIndex(false), "GetSVarType": state.getType(false),
		"SetPVarInt": state.setInt(true), "GetPVarInt": state.getInt(true),
		"SetPVarString": state.setString(true), "GetPVarString": state.getString(true),
		"SetPVarFloat": state.setFloat(true), "GetPVarFloat": state.getFloat(true),
		"DeletePVar": state.deleteVariable(true), "GetPVarsUpperIndex": state.getUpperIndex(true),
		"GetPVarNameAtIndex": state.getNameAtIndex(true), "GetPVarType": state.getType(true),
	}
}

func (state *variableState) scope(params []backend.Cell, player bool) (*variableScope, int, bool) {
	if !player {
		return &state.server, 0, true
	}

	if len(params) == 0 || state.model == nil {
		return nil, 0, false
	}

	playerID := int(params[0])

	p, ok := state.model.players[playerID]
	if !ok || !p.connected {
		return nil, 0, false
	}

	if state.players[playerID] == nil {
		scope := newVariableScope()
		state.players[playerID] = &scope
	}

	return state.players[playerID], 1, true
}

func setVariable(scope *variableScope, name string, value testVariable) {
	if _, exists := scope.values[name]; !exists {
		scope.order = append(scope.order, name)
	}

	scope.values[name] = value
}

func deleteVariable(scope *variableScope, name string) bool {
	if _, exists := scope.values[name]; !exists {
		return false
	}

	delete(scope.values, name)

	index := slices.Index(scope.order, name)
	if index >= 0 {
		scope.order = slices.Delete(scope.order, index, index+1)
	}

	return true
}
