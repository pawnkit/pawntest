package runner

import (
	"fmt"
	"math"
	"strconv"

	"github.com/pawnkit/pawntest/internal/backend"
)

var floatNativeNames = map[string]struct{}{
	"float":       {},
	"strfloat":    {},
	"floatmul":    {},
	"floatdiv":    {},
	"floatadd":    {},
	"floatsub":    {},
	"floatfract":  {},
	"floatround":  {},
	"floatcmp":    {},
	"floatsqroot": {},
	"floatpower":  {},
	"floatlog":    {},
	"floatsin":    {},
	"floatcos":    {},
	"floattan":    {},
	"floatabs":    {},
	"floatint":    {},
}

func isFloatNative(name string) bool {
	_, ok := floatNativeNames[name]
	return ok
}

func registerFloatNatives(vm backend.VM) error {
	for name := range floatNativeNames {
		if err := vm.RegisterNative(name, func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
			return callFloatNative(name, ctx, params)
		}); err != nil {
			return err
		}
	}

	return nil
}

func callFloatNative(name string, ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	floatParam := func(index int) float64 {
		return float64(math.Float32frombits(uint32(params[index])))
	}
	floatCell := func(value float64) backend.Cell {
		return backend.Cell(int32(math.Float32bits(float32(value))))
	}
	require := func(count int) error {
		if len(params) < count {
			return fmt.Errorf("native %s expects at least %d arguments", name, count)
		}

		return nil
	}

	switch name {
	case "float":
		if err := require(1); err != nil {
			return 0, err
		}

		return floatCell(float64(params[0])), nil
	case "strfloat":
		if err := require(1); err != nil {
			return 0, err
		}

		value, err := ctx.ReadString(params[0])
		if err != nil {
			return 0, err
		}

		parsed, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0, nil
		}

		return floatCell(parsed), nil
	case "floatmul", "floatdiv", "floatadd", "floatsub", "floatcmp", "floatpower":
		if err := require(2); err != nil {
			return 0, err
		}

		a, b := floatParam(0), floatParam(1)

		switch name {
		case "floatmul":
			return floatCell(a * b), nil
		case "floatdiv":
			return floatCell(a / b), nil
		case "floatadd":
			return floatCell(a + b), nil
		case "floatsub":
			return floatCell(a - b), nil
		case "floatpower":
			return floatCell(math.Pow(a, b)), nil
		default:
			if a < b {
				return -1, nil
			}

			if a > b {
				return 1, nil
			}

			return 0, nil
		}
	case "floatlog":
		if err := require(1); err != nil {
			return 0, err
		}

		base := 10.0
		if len(params) > 1 {
			base = floatParam(1)
		}

		return floatCell(math.Log(floatParam(0)) / math.Log(base)), nil
	case "floatfract", "floatsqroot", "floatsin", "floatcos", "floattan", "floatabs":
		if err := require(1); err != nil {
			return 0, err
		}

		value := floatParam(0)

		switch name {
		case "floatfract":
			_, value = math.Modf(value)
		case "floatsqroot":
			value = math.Sqrt(value)
		case "floatsin", "floatcos", "floattan":
			mode := backend.Cell(0)
			if len(params) > 1 {
				mode = params[1]
			}

			if mode == 1 {
				value *= math.Pi / 180
			}

			if mode == 2 {
				value *= math.Pi / 200
			}

			if name == "floatsin" {
				value = math.Sin(value)
			}

			if name == "floatcos" {
				value = math.Cos(value)
			}

			if name == "floattan" {
				value = math.Tan(value)
			}
		case "floatabs":
			value = math.Abs(value)
		}

		return floatCell(value), nil
	case "floatround":
		if err := require(1); err != nil {
			return 0, err
		}

		mode := backend.Cell(0)
		if len(params) > 1 {
			mode = params[1]
		}

		value := floatParam(0)

		switch mode {
		case 1:
			value = math.Floor(value)
		case 2:
			value = math.Ceil(value)
		case 3:
			value = math.Trunc(value)
		default:
			value = math.Floor(value + 0.5)
		}

		return backend.Cell(int32(value)), nil
	case "floatint":
		if err := require(1); err != nil {
			return 0, err
		}

		return backend.Cell(int32(math.Trunc(floatParam(0)))), nil
	}

	return 0, fmt.Errorf("unsupported float native %s", name)
}
