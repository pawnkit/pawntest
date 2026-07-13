package runner

import (
	"fmt"

	"github.com/pawnkit/pawntest/internal/backend"
)

type nativeParams struct {
	ctx    backend.NativeContext
	values []backend.Cell
}

func readNativeParams(ctx backend.NativeContext, values []backend.Cell) nativeParams {
	return nativeParams{ctx: ctx, values: values}
}

func (params nativeParams) Require(count int, name string) error {
	if len(params.values) < count {
		return fmt.Errorf("%s expects %d arguments, got %d", name, count, len(params.values))
	}

	return nil
}

func (params nativeParams) Cell(index int) backend.Cell {
	return params.values[index]
}

func (params nativeParams) Int(index int) int {
	return int(params.values[index])
}

func (params nativeParams) String(index int) (string, error) {
	return params.ctx.ReadString(params.values[index])
}
