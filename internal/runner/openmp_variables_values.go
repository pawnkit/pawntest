package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *variableState) setInt(player bool) backend.NativeFunc {
	return state.setNumber(player, func(value backend.Cell) testVariable {
		return testVariable{variableType: variableInt, integer: int(value)}
	})
}

func (state *variableState) setFloat(player bool) backend.NativeFunc {
	return state.setNumber(player, func(value backend.Cell) testVariable {
		return testVariable{variableType: variableFloat, floating: cellFloat(value)}
	})
}

func (state *variableState) setNumber(player bool, value func(backend.Cell) testVariable) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, offset, ok := state.scope(params, player)
		if !ok || len(params) < offset+2 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[offset])
		if err != nil {
			return 0, err
		}

		setVariable(scope, name, value(params[offset+1]))

		return 1, nil
	}
}

func (state *variableState) getInt(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		value, ok, err := state.variable(ctx, params, player)
		if err != nil || !ok || value.variableType != variableInt {
			return 0, err
		}

		return backend.Cell(value.integer), nil
	}
}

func (state *variableState) getFloat(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		value, ok, err := state.variable(ctx, params, player)
		if err != nil || !ok || value.variableType != variableFloat {
			return 0, err
		}

		return floatCell(value.floating), nil
	}
}

func (state *variableState) setString(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, offset, ok := state.scope(params, player)
		if !ok || len(params) < offset+2 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[offset])
		if err != nil {
			return 0, err
		}

		value, err := ctx.ReadString(params[offset+1])
		if err != nil {
			return 0, err
		}

		setVariable(scope, name, testVariable{variableType: variableString, text: value})

		return 1, nil
	}
}

func (state *variableState) getString(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, offset, ok := state.scope(params, player)
		if !ok || len(params) < offset+3 {
			return 0, nil
		}

		name, err := ctx.ReadString(params[offset])
		if err != nil {
			return 0, err
		}

		value, exists := scope.values[name]
		if !exists || value.variableType != variableString {
			return 0, nil
		}

		if err := ctx.WriteString(params[offset+1], truncateString(value.text, int(params[offset+2]))); err != nil {
			return 0, err
		}

		return backend.Cell(len(value.text)), nil
	}
}

func (state *variableState) variable(ctx backend.NativeContext, params []backend.Cell, player bool) (testVariable, bool, error) {
	scope, offset, ok := state.scope(params, player)
	if !ok || len(params) <= offset {
		return testVariable{}, false, nil
	}

	name, err := ctx.ReadString(params[offset])
	if err != nil {
		return testVariable{}, false, err
	}

	value, exists := scope.values[name]

	return value, exists, nil
}

func (state *variableState) deleteVariable(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, offset, ok := state.scope(params, player)
		if !ok || len(params) <= offset {
			return 0, nil
		}

		name, err := ctx.ReadString(params[offset])
		if err != nil {
			return 0, err
		}

		if deleteVariable(scope, name) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *variableState) getType(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		value, ok, err := state.variable(ctx, params, player)
		if err != nil || !ok {
			return backend.Cell(variableNone), err
		}

		return backend.Cell(value.variableType), nil
	}
}

func (state *variableState) getUpperIndex(player bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, _, ok := state.scope(params, player)
		if !ok || len(scope.order) == 0 {
			return -1, nil
		}

		return backend.Cell(len(scope.order) - 1), nil
	}
}

func (state *variableState) getNameAtIndex(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		scope, offset, ok := state.scope(params, player)
		if !ok || len(params) < offset+3 || params[offset] < 0 || int(params[offset]) >= len(scope.order) {
			return 0, nil
		}

		name := scope.order[int(params[offset])]
		if err := ctx.WriteString(params[offset+1], truncateString(name, int(params[offset+2]))); err != nil {
			return 0, err
		}

		return backend.Cell(len(name)), nil
	}
}
