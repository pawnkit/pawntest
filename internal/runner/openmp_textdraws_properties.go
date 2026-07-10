package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *textDrawState) destroyTextDraw(player bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if _, ok := state.draw(params, player); !ok {
			return 0, nil
		}

		if player {
			delete(state.playerDraws[int(params[0])], int(params[1]))
		} else {
			delete(state.draws, int(params[0]))
		}

		return 1, nil
	}
}

func (state *textDrawState) isValidTextDraw(player bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if _, ok := state.draw(params, player); ok {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *textDrawState) setTextDrawString(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset {
			return 0, nil
		}

		text, err := ctx.ReadString(params[offset])
		if err != nil {
			return 0, err
		}

		draw.text = text

		return 1, nil
	}
}

func (state *textDrawState) getTextDrawString(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+1 {
			return 0, nil
		}

		text := draw.text

		return 1, ctx.WriteString(params[offset], truncateString(text, int(params[offset+1])))
	}
}

func (state *textDrawState) setTextDrawInt(player bool, field func(*testTextDraw) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset {
			return 0, nil
		}

		*field(draw) = int(params[offset])

		return 1, nil
	}
}

func (state *textDrawState) getTextDrawInt(player bool, field func(*testTextDraw) *int) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)
		if !ok {
			return 0, nil
		}

		return backend.Cell(*field(draw)), nil
	}
}

func (state *textDrawState) setTextDrawBool(player bool, field func(*testTextDraw) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset {
			return 0, nil
		}

		*field(draw) = params[offset] != 0

		return 1, nil
	}
}

func (state *textDrawState) getTextDrawBool(player bool, field func(*testTextDraw) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)
		if ok && *field(draw) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *textDrawState) setTextDrawPair(player bool, field func(*testTextDraw) *[2]float32) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+1 {
			return 0, nil
		}

		*field(draw) = [2]float32{cellFloat(params[offset]), cellFloat(params[offset+1])}

		return 1, nil
	}
}

func (state *textDrawState) getTextDrawPair(player bool, field func(*testTextDraw) *[2]float32) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+1 {
			return 0, nil
		}

		for index, value := range *field(draw) {
			if err := ctx.WriteCell(params[offset+index], floatCell(value)); err != nil {
				return 0, err
			}
		}

		return 1, nil
	}
}

func (state *textDrawState) setPreviewRotation(player bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+3 {
			return 0, nil
		}

		draw.previewRotation = [4]float32{cellFloat(params[offset]), cellFloat(params[offset+1]), cellFloat(params[offset+2]), cellFloat(params[offset+3])}

		return 1, nil
	}
}

func (state *textDrawState) getPreviewRotation(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+3 {
			return 0, nil
		}

		for index, value := range draw.previewRotation {
			if err := ctx.WriteCell(params[offset+index], floatCell(value)); err != nil {
				return 0, err
			}
		}

		return 1, nil
	}
}

func (state *textDrawState) setPreviewColours(player bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+1 {
			return 0, nil
		}

		draw.previewColours = [2]int{int(params[offset]), int(params[offset+1])}

		return 1, nil
	}
}

func (state *textDrawState) getPreviewColours(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		draw, ok := state.draw(params, player)

		offset := drawOffset(player)
		if !ok || len(params) <= offset+1 {
			return 0, nil
		}

		return writeVehicleInts(ctx, params[offset:offset+2], draw.previewColours[:])
	}
}
