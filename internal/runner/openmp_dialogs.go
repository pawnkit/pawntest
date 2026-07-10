package runner

import (
	"errors"
	"fmt"
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testDialog struct {
	id, style        int
	title, body      string
	button1, button2 string
	visible          bool
}

type dialogState struct {
	dialogs map[int]testDialog
	players *openMPState
}

func newDialogState() *dialogState {
	return &dialogState{dialogs: map[int]testDialog{}}
}

func (state *dialogState) Clone() scenarioModule {
	clone := newDialogState()
	maps.Copy(clone.dialogs, state.dialogs)

	return clone
}

func (state *dialogState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *dialogState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_dialog_respond": state.respondDialog,
		"__pt_dialog_visible": state.assertVisible(result),
		"__pt_dialog_title":   state.assertString(result, "title", func(dialog testDialog) string { return dialog.title }),
		"__pt_dialog_body":    state.assertString(result, "body", func(dialog testDialog) string { return dialog.body }),
		"ShowPlayerDialog":    state.showPlayerDialog,
		"HidePlayerDialog":    state.hidePlayerDialog,
		"GetPlayerDialogData": state.getPlayerDialogData,
		"GetPlayerDialogID":   state.getPlayerDialogID,
		"GetPlayerDialog":     state.getPlayerDialogID,
	}
}

func (state *dialogState) showPlayerDialog(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	if params[1] < 0 {
		delete(state.dialogs, int(params[0]))

		return 1, nil
	}

	strings := make([]string, 4)

	for index, param := range params[3:7] {
		value, err := ctx.ReadString(param)
		if err != nil {
			return 0, err
		}

		strings[index] = value
	}

	state.dialogs[int(params[0])] = testDialog{
		id: int(params[1]), style: int(params[2]), title: strings[0], body: strings[1],
		button1: strings[2], button2: strings[3], visible: true,
	}

	return 1, nil
}

func (state *dialogState) hidePlayerDialog(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 || !state.validPlayer(params[0]) {
		return 0, nil
	}

	dialog := state.dialogs[int(params[0])]
	dialog.visible = false
	state.dialogs[int(params[0])] = dialog

	return 1, nil
}

func (state *dialogState) getPlayerDialogID(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return -1, nil
	}

	dialog, ok := state.dialogs[int(params[0])]
	if !ok || !dialog.visible {
		return -1, nil
	}

	return backend.Cell(dialog.id), nil
}

func (state *dialogState) getPlayerDialogData(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 10 {
		return 0, nil
	}

	dialog, ok := state.dialogs[int(params[0])]
	if !ok || !dialog.visible {
		return 0, nil
	}

	if err := ctx.WriteCell(params[1], backend.Cell(dialog.style)); err != nil {
		return 0, err
	}

	values := []string{dialog.title, dialog.body, dialog.button1, dialog.button2}
	for index, value := range values {
		address := params[2+index*2]

		size := int(params[3+index*2])
		if err := ctx.WriteString(address, truncateString(value, size)); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (state *dialogState) respondDialog(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 4 {
		return 0, errors.New("dialog response expects 4 arguments")
	}

	dialog, ok := state.dialogs[int(params[0])]
	if !ok || !dialog.visible {
		return 0, nil
	}

	caller, ok := ctx.(backend.PublicCaller)
	if !ok {
		return 0, errors.New("runtime does not support dialog callbacks")
	}

	dialog.visible = false
	state.dialogs[int(params[0])] = dialog

	response, err := caller.CallPublic("OnDialogResponse", params[0], backend.Cell(dialog.id), params[1], params[2], params[3])
	if err != nil {
		return 0, fmt.Errorf("dialog response callback: %w", err)
	}

	return response, nil
}

func (state *dialogState) validPlayer(id backend.Cell) bool {
	if state.players == nil {
		return false
	}

	player, ok := state.players.player(id)

	return ok && player.connected
}
