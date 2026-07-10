package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *textDrawState) natives(result *nativeState) map[string]backend.NativeFunc {
	natives := map[string]backend.NativeFunc{
		"__pt_textdraw_create": state.createTextDraw(false), "__pt_textdraw_valid": state.assertValid(result, false),
		"__pt_textdraw_text": state.assertText(result, false), "__pt_textdraw_visible": state.assertVisible(result, false),
		"TextDrawCreate": state.createTextDraw(false), "TextDrawDestroy": state.destroyTextDraw(false),
		"IsValidTextDraw": state.isValidTextDraw(false), "TextDrawSetString": state.setTextDrawString(false),
		"TextDrawGetString": state.getTextDrawString(false), "TextDrawSetPos": state.setTextDrawPair(false, drawPosition),
		"TextDrawGetPos": state.getTextDrawPair(false, drawPosition), "TextDrawLetterSize": state.setTextDrawPair(false, drawLetterSize),
		"TextDrawGetLetterSize": state.getTextDrawPair(false, drawLetterSize), "TextDrawTextSize": state.setTextDrawPair(false, drawTextSize),
		"TextDrawGetTextSize": state.getTextDrawPair(false, drawTextSize), "TextDrawAlignment": state.setTextDrawInt(false, drawAlignment),
		"TextDrawGetAlignment": state.getTextDrawInt(false, drawAlignment), "TextDrawColor": state.setTextDrawInt(false, drawColour),
		"TextDrawColour": state.setTextDrawInt(false, drawColour), "TextDrawGetColor": state.getTextDrawInt(false, drawColour),
		"TextDrawGetColour": state.getTextDrawInt(false, drawColour), "TextDrawUseBox": state.setTextDrawBool(false, drawBox),
		"TextDrawIsBox": state.getTextDrawBool(false, drawBox), "TextDrawBoxColor": state.setTextDrawInt(false, drawBoxColour),
		"TextDrawBoxColour": state.setTextDrawInt(false, drawBoxColour), "TextDrawGetBoxColor": state.getTextDrawInt(false, drawBoxColour),
		"TextDrawGetBoxColour": state.getTextDrawInt(false, drawBoxColour), "TextDrawSetShadow": state.setTextDrawInt(false, drawShadow),
		"TextDrawGetShadow": state.getTextDrawInt(false, drawShadow), "TextDrawSetOutline": state.setTextDrawInt(false, drawOutline),
		"TextDrawGetOutline": state.getTextDrawInt(false, drawOutline), "TextDrawBackgroundColor": state.setTextDrawInt(false, drawBackgroundColour),
		"TextDrawBackgroundColour": state.setTextDrawInt(false, drawBackgroundColour), "TextDrawGetBackgroundColor": state.getTextDrawInt(false, drawBackgroundColour),
		"TextDrawGetBackgroundColour": state.getTextDrawInt(false, drawBackgroundColour), "TextDrawFont": state.setTextDrawInt(false, drawFont),
		"TextDrawGetFont": state.getTextDrawInt(false, drawFont), "TextDrawSetProportional": state.setTextDrawBool(false, drawProportional),
		"TextDrawIsProportional": state.getTextDrawBool(false, drawProportional), "TextDrawSetSelectable": state.setTextDrawBool(false, drawSelectable),
		"TextDrawIsSelectable": state.getTextDrawBool(false, drawSelectable), "TextDrawSetPreviewModel": state.setTextDrawInt(false, drawPreviewModel),
		"TextDrawGetPreviewModel": state.getTextDrawInt(false, drawPreviewModel), "TextDrawSetPreviewRot": state.setPreviewRotation(false),
		"TextDrawGetPreviewRot": state.getPreviewRotation(false), "TextDrawSetPreviewVehCol": state.setPreviewColours(false),
		"TextDrawGetPreviewVehCol": state.getPreviewColours(false), "TextDrawShowForPlayer": state.showGlobalTextDraw,
		"TextDrawHideForPlayer": state.hideGlobalTextDraw, "TextDrawShowForAll": state.showGlobalTextDrawForAll,
		"TextDrawHideForAll": state.hideGlobalTextDrawForAll, "IsTextDrawVisibleForPlayer": state.isGlobalTextDrawVisible,
		"TextDrawSetStringForPlayer": state.setGlobalTextDrawStringForPlayer, "SelectTextDraw": state.selectTextDraw,
		"CancelSelectTextDraw": state.cancelSelectTextDraw,

		"__pt_player_textdraw_create": state.createTextDraw(true), "__pt_player_textdraw_valid": state.assertValid(result, true),
		"__pt_player_textdraw_text": state.assertText(result, true), "__pt_player_textdraw_visible": state.assertVisible(result, true),
		"CreatePlayerTextDraw": state.createTextDraw(true), "PlayerTextDrawDestroy": state.destroyTextDraw(true),
		"IsValidPlayerTextDraw": state.isValidTextDraw(true), "PlayerTextDrawSetString": state.setTextDrawString(true),
		"PlayerTextDrawGetString": state.getTextDrawString(true), "PlayerTextDrawSetPos": state.setTextDrawPair(true, drawPosition),
		"PlayerTextDrawGetPos": state.getTextDrawPair(true, drawPosition), "PlayerTextDrawLetterSize": state.setTextDrawPair(true, drawLetterSize),
		"PlayerTextDrawGetLetterSize": state.getTextDrawPair(true, drawLetterSize), "PlayerTextDrawTextSize": state.setTextDrawPair(true, drawTextSize),
		"PlayerTextDrawGetTextSize": state.getTextDrawPair(true, drawTextSize), "PlayerTextDrawAlignment": state.setTextDrawInt(true, drawAlignment),
		"PlayerTextDrawGetAlignment": state.getTextDrawInt(true, drawAlignment), "PlayerTextDrawColor": state.setTextDrawInt(true, drawColour),
		"PlayerTextDrawColour": state.setTextDrawInt(true, drawColour), "PlayerTextDrawGetColor": state.getTextDrawInt(true, drawColour),
		"PlayerTextDrawGetColour": state.getTextDrawInt(true, drawColour), "PlayerTextDrawUseBox": state.setTextDrawBool(true, drawBox),
		"PlayerTextDrawIsBox": state.getTextDrawBool(true, drawBox), "PlayerTextDrawBoxColor": state.setTextDrawInt(true, drawBoxColour),
		"PlayerTextDrawBoxColour": state.setTextDrawInt(true, drawBoxColour), "PlayerTextDrawGetBoxColor": state.getTextDrawInt(true, drawBoxColour),
		"PlayerTextDrawGetBoxColour": state.getTextDrawInt(true, drawBoxColour), "PlayerTextDrawSetShadow": state.setTextDrawInt(true, drawShadow),
		"PlayerTextDrawGetShadow": state.getTextDrawInt(true, drawShadow), "PlayerTextDrawSetOutline": state.setTextDrawInt(true, drawOutline),
		"PlayerTextDrawGetOutline": state.getTextDrawInt(true, drawOutline), "PlayerTextDrawBackgroundColor": state.setTextDrawInt(true, drawBackgroundColour),
		"PlayerTextDrawBackgroundColour": state.setTextDrawInt(true, drawBackgroundColour), "PlayerTextDrawGetBackgroundCol": state.getTextDrawInt(true, drawBackgroundColour),
		"PlayerTextDrawFont": state.setTextDrawInt(true, drawFont), "PlayerTextDrawGetFont": state.getTextDrawInt(true, drawFont),
		"PlayerTextDrawSetProportional": state.setTextDrawBool(true, drawProportional), "PlayerTextDrawIsProportional": state.getTextDrawBool(true, drawProportional),
		"PlayerTextDrawSetSelectable": state.setTextDrawBool(true, drawSelectable), "PlayerTextDrawIsSelectable": state.getTextDrawBool(true, drawSelectable),
		"PlayerTextDrawSetPreviewModel": state.setTextDrawInt(true, drawPreviewModel), "PlayerTextDrawGetPreviewModel": state.getTextDrawInt(true, drawPreviewModel),
		"PlayerTextDrawSetPreviewRot": state.setPreviewRotation(true), "PlayerTextDrawGetPreviewRot": state.getPreviewRotation(true),
		"PlayerTextDrawSetPreviewVehCol": state.setPreviewColours(true), "PlayerTextDrawGetPreviewVehCol": state.getPreviewColours(true),
		"PlayerTextDrawShow": state.showPlayerTextDraw, "PlayerTextDrawHide": state.hidePlayerTextDraw,
		"IsPlayerTextDrawVisible": state.isPlayerTextDrawVisible,
	}

	return natives
}

func (state *textDrawState) createTextDraw(player bool) backend.NativeFunc {
	return func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		offset := 0

		if player {
			if len(params) < 4 || !state.validPlayer(params[0]) {
				return -1, nil
			}

			offset = 1
		} else if len(params) < 3 {
			return -1, nil
		}

		text, err := ctx.ReadString(params[offset+2])
		if err != nil {
			return -1, err
		}

		draw := newTestTextDraw(cellFloat(params[offset]), cellFloat(params[offset+1]), text)
		if player {
			playerID := int(params[0])
			id := state.playerNext[playerID]

			state.playerNext[playerID]++
			if state.playerDraws[playerID] == nil {
				state.playerDraws[playerID] = map[int]*testTextDraw{}
			}

			state.playerDraws[playerID][id] = draw

			return backend.Cell(id), nil
		}

		id := state.next
		state.next++
		state.draws[id] = draw

		return backend.Cell(id), nil
	}
}

func (state *textDrawState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	p, ok := state.players.player(id)

	return ok && p.connected
}
