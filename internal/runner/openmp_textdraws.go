package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testTextDraw struct {
	text                                string
	position, letterSize, textSize      [2]float32
	colour, boxColour, backgroundColour int
	shadow, outline, font, alignment    int
	box, proportional, selectable       bool
	previewModel                        int
	previewRotation                     [4]float32
	previewColours                      [2]int
	visible                             map[int]bool
	playerText                          map[int]string
}

type textDrawSelection struct {
	active bool
	colour int
}

type textDrawState struct {
	next        int
	draws       map[int]*testTextDraw
	playerNext  map[int]int
	playerDraws map[int]map[int]*testTextDraw
	selection   map[int]textDrawSelection
	players     *openMPState
}

func newTextDrawState() *textDrawState {
	return &textDrawState{draws: map[int]*testTextDraw{}, playerNext: map[int]int{}, playerDraws: map[int]map[int]*testTextDraw{}, selection: map[int]textDrawSelection{}}
}

func (state *textDrawState) Clone() scenarioModule {
	clone := newTextDrawState()
	clone.next = state.next
	maps.Copy(clone.playerNext, state.playerNext)
	maps.Copy(clone.selection, state.selection)

	for id, draw := range state.draws {
		clone.draws[id] = cloneTextDraw(draw)
	}

	for playerID, draws := range state.playerDraws {
		clone.playerDraws[playerID] = map[int]*testTextDraw{}
		for id, draw := range draws {
			clone.playerDraws[playerID][id] = cloneTextDraw(draw)
		}
	}

	return clone
}

func cloneTextDraw(draw *testTextDraw) *testTextDraw {
	clone := *draw
	clone.visible = maps.Clone(draw.visible)
	clone.playerText = maps.Clone(draw.playerText)

	return &clone
}

func (state *textDrawState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func newTestTextDraw(x, y float32, text string) *testTextDraw {
	return &testTextDraw{text: text, position: [2]float32{x, y}, previewModel: -1, visible: map[int]bool{}, playerText: map[int]string{}}
}

func (state *textDrawState) draw(params []backend.Cell, player bool) (*testTextDraw, bool) {
	if player {
		if len(params) < 2 {
			return nil, false
		}

		draw, ok := state.playerDraws[int(params[0])][int(params[1])]

		return draw, ok
	}

	if len(params) == 0 {
		return nil, false
	}

	draw, ok := state.draws[int(params[0])]

	return draw, ok
}

func drawOffset(player bool) int {
	if player {
		return 2
	}

	return 1
}

func drawColour(draw *testTextDraw) *int            { return &draw.colour }
func drawBoxColour(draw *testTextDraw) *int         { return &draw.boxColour }
func drawBackgroundColour(draw *testTextDraw) *int  { return &draw.backgroundColour }
func drawShadow(draw *testTextDraw) *int            { return &draw.shadow }
func drawOutline(draw *testTextDraw) *int           { return &draw.outline }
func drawFont(draw *testTextDraw) *int              { return &draw.font }
func drawAlignment(draw *testTextDraw) *int         { return &draw.alignment }
func drawPreviewModel(draw *testTextDraw) *int      { return &draw.previewModel }
func drawBox(draw *testTextDraw) *bool              { return &draw.box }
func drawProportional(draw *testTextDraw) *bool     { return &draw.proportional }
func drawSelectable(draw *testTextDraw) *bool       { return &draw.selectable }
func drawPosition(draw *testTextDraw) *[2]float32   { return &draw.position }
func drawLetterSize(draw *testTextDraw) *[2]float32 { return &draw.letterSize }
func drawTextSize(draw *testTextDraw) *[2]float32   { return &draw.textSize }
