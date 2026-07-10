package runner

import (
	"strings"

	"github.com/pawnkit/pawntest/internal/backend"
)

func (state *databaseState) freeResult(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || state.results[int(params[0])] == nil {
		return 0, nil
	}
	delete(state.results, int(params[0]))

	return 1, nil
}

func (state *databaseState) getRowCount(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	result := state.result(params)
	if result == nil {
		return 0, nil
	}

	return backend.Cell(len(result.rows)), nil
}

func (state *databaseState) selectNextRow(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	result := state.result(params)
	if result == nil || result.row+1 >= len(result.rows) {
		return 0, nil
	}
	result.row++

	return 1, nil
}

func (state *databaseState) getFieldCount(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	result := state.result(params)
	if result == nil {
		return 0, nil
	}

	return backend.Cell(len(result.columns)), nil
}

func (state *databaseState) getFieldName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return 0, nil
	}
	result := state.result(params)
	index := int(params[1])
	if result == nil || index < 0 || index >= len(result.columns) {
		return 0, nil
	}

	return 1, ctx.WriteString(params[2], truncateString(result.columns[index], int(params[3])))
}

func (state *databaseState) getFieldString(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return 0, nil
	}
	value, ok := state.field(params, int(params[1]))
	if !ok {
		return 0, nil
	}

	return 1, ctx.WriteString(params[2], truncateString(value, int(params[3])))
}

func (state *databaseState) getFieldInt(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}
	value, _ := state.field(params, int(params[1]))

	return parseDatabaseInt(value), nil
}

func (state *databaseState) getFieldFloat(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}
	value, _ := state.field(params, int(params[1]))

	return parseDatabaseFloat(value), nil
}

func (state *databaseState) getFieldStringByName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return 0, nil
	}
	name, err := ctx.ReadString(params[1])
	if err != nil {
		return 0, err
	}
	value, ok := state.namedField(params, name)
	if !ok {
		return 0, nil
	}

	return 1, ctx.WriteString(params[2], truncateString(value, int(params[3])))
}

func (state *databaseState) getFieldIntByName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	value, err := state.namedFieldValue(ctx, params)

	return parseDatabaseInt(value), err
}

func (state *databaseState) getFieldFloatByName(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	value, err := state.namedFieldValue(ctx, params)

	return parseDatabaseFloat(value), err
}

func (state *databaseState) getDatabaseHandle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || state.connections[int(params[0])] == nil {
		return 0, nil
	}

	return params[0], nil
}

func (state *databaseState) getResultHandle(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if state.result(params) == nil {
		return 0, nil
	}

	return params[0], nil
}

func (state *databaseState) getConnectionCount(_ backend.NativeContext, _ []backend.Cell) (backend.Cell, error) {
	return backend.Cell(len(state.connections)), nil
}

func (state *databaseState) getResultCount(_ backend.NativeContext, _ []backend.Cell) (backend.Cell, error) {
	return backend.Cell(len(state.results)), nil
}

func (state *databaseState) result(params []backend.Cell) *databaseResult {
	if len(params) == 0 {
		return nil
	}

	return state.results[int(params[0])]
}

func (state *databaseState) field(params []backend.Cell, index int) (string, bool) {
	result := state.result(params)
	if result == nil || result.row >= len(result.rows) || index < 0 || index >= len(result.columns) {
		return "", false
	}

	return result.rows[result.row][index], true
}

func (state *databaseState) namedField(params []backend.Cell, name string) (string, bool) {
	result := state.result(params)
	if result == nil {
		return "", false
	}
	for index, column := range result.columns {
		if strings.EqualFold(column, name) {
			return state.field(params, index)
		}
	}

	return "", false
}

func (state *databaseState) namedFieldValue(ctx backend.NativeContext, params []backend.Cell) (string, error) {
	if len(params) < 2 {
		return "", nil
	}
	name, err := ctx.ReadString(params[1])
	if err != nil {
		return "", err
	}
	value, _ := state.namedField(params, name)

	return value, nil
}
