package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

const invalidScenarioID = -1

type testTextLabel struct {
	text                        string
	colour, world               int
	position                    [3]float32
	drawDistance                float32
	lineOfSight                 bool
	parentPlayer, parentVehicle int
}

type textLabelState struct {
	next         int
	labels       map[int]*testTextLabel
	playerNext   map[int]int
	playerLabels map[int]map[int]*testTextLabel
	players      *openMPState
	vehicles     *vehicleState
}

func newTextLabelState() *textLabelState {
	return &textLabelState{labels: map[int]*testTextLabel{}, playerNext: map[int]int{}, playerLabels: map[int]map[int]*testTextLabel{}}
}

func (state *textLabelState) Clone() scenarioModule {
	clone := newTextLabelState()
	clone.next = state.next
	maps.Copy(clone.playerNext, state.playerNext)

	for id, label := range state.labels {
		labelCopy := *label
		clone.labels[id] = &labelCopy
	}

	for playerID, labels := range state.playerLabels {
		clone.playerLabels[playerID] = map[int]*testTextLabel{}

		for id, label := range labels {
			labelCopy := *label
			clone.playerLabels[playerID][id] = &labelCopy
		}
	}

	return clone
}

func (state *textLabelState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()
	state.vehicles = context.scenarios.vehicleState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *textLabelState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_text_label_create":     state.createTextLabel,
		"__pt_text_label_valid":      state.assertValid(result),
		"__pt_text_label_text":       state.assertText(result),
		"Create3DTextLabel":          state.createTextLabel,
		"Delete3DTextLabel":          state.deleteTextLabel,
		"IsValid3DTextLabel":         state.isValidTextLabel,
		"Is3DTextLabelStreamedIn":    state.isTextLabelStreamedIn,
		"Update3DTextLabelText":      state.updateTextLabel,
		"Get3DTextLabelText":         state.getTextLabelText,
		"Get3DTextLabelColor":        state.getTextLabelInt(labelColour),
		"Get3DTextLabelColour":       state.getTextLabelInt(labelColour),
		"Get3DTextLabelPos":          state.getTextLabelPosition,
		"Get3DTextLabelDrawDistance": state.getTextLabelFloat(labelDrawDistance),
		"Set3DTextLabelDrawDistance": state.setTextLabelFloat(labelDrawDistance),
		"Get3DTextLabelLOS":          state.getTextLabelBool(labelLineOfSight),
		"Set3DTextLabelLOS":          state.setTextLabelBool(labelLineOfSight),
		"Get3DTextLabelVirtualWorld": state.getTextLabelInt(labelWorld),
		"Set3DTextLabelVirtualWorld": state.setTextLabelInt(labelWorld),
		"Attach3DTextLabelToPlayer":  state.attachTextLabelToPlayer,
		"Attach3DTextLabelToVehicle": state.attachTextLabelToVehicle,
		"Get3DTextLabelAttachedData": state.getTextLabelAttachedData,

		"__pt_player_text_label_create": state.createPlayerTextLabel,
		"__pt_player_text_label_valid":  state.assertPlayerValid(result),
		"__pt_player_text_label_text":   state.assertPlayerText(result),
		"CreatePlayer3DTextLabel":       state.createPlayerTextLabel,
		"DeletePlayer3DTextLabel":       state.deletePlayerTextLabel,
		"IsValidPlayer3DTextLabel":      state.isValidPlayerTextLabel,
		"UpdatePlayer3DTextLabelText":   state.updatePlayerTextLabel,
		"GetPlayer3DTextLabelText":      state.getPlayerTextLabelText,
		"GetPlayer3DTextLabelColor":     state.getPlayerTextLabelInt(labelColour),
		"GetPlayer3DTextLabelColour":    state.getPlayerTextLabelInt(labelColour),
		"GetPlayer3DTextLabelPos":       state.getPlayerTextLabelPosition,
		"GetPlayer3DTextLabelDrawDist":  state.getPlayerTextLabelFloat(labelDrawDistance),
		"SetPlayer3DTextLabelDrawDist":  state.setPlayerTextLabelFloat(labelDrawDistance),
		"GetPlayer3DTextLabelLOS":       state.getPlayerTextLabelBool(labelLineOfSight),
		"SetPlayer3DTextLabelLOS":       state.setPlayerTextLabelBool(labelLineOfSight),
		"GetPlayer3DTextLabelVirtualW":  state.getPlayerTextLabelInt(labelWorld),
		"GetPlayer3DTextLabelAttached":  state.getPlayerTextLabelAttachedData,
	}
}

func newTestTextLabel(text string, colour int, position [3]float32, drawDistance float32, world int) *testTextLabel {
	return &testTextLabel{text: text, colour: colour, position: position, drawDistance: drawDistance, world: world, parentPlayer: invalidScenarioID, parentVehicle: invalidScenarioID}
}

func (state *textLabelState) createTextLabel(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 7 {
		return -1, nil
	}

	text, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	lineOfSight := len(params) > 7 && params[7] != 0
	id := state.next
	state.next++
	state.labels[id] = newTestTextLabel(text, int(params[1]), cellsToVector(params[2:5]), cellFloat(params[5]), int(params[6]))
	state.labels[id].lineOfSight = lineOfSight

	return backend.Cell(id), nil
}

func (state *textLabelState) label(params []backend.Cell) (*testTextLabel, bool) {
	if len(params) == 0 {
		return nil, false
	}

	label, ok := state.labels[int(params[0])]

	return label, ok
}

func (state *textLabelState) deleteTextLabel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.label(params); !ok {
		return 0, nil
	}

	delete(state.labels, int(params[0]))

	return 1, nil
}

func (state *textLabelState) isValidTextLabel(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.label(params); ok {
		return 1, nil
	}

	return 0, nil
}

func labelColour(label *testTextLabel) *int           { return &label.colour }
func labelWorld(label *testTextLabel) *int            { return &label.world }
func labelDrawDistance(label *testTextLabel) *float32 { return &label.drawDistance }
func labelLineOfSight(label *testTextLabel) *bool     { return &label.lineOfSight }
