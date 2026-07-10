package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testMenu struct {
	title        string
	columns      int
	position     [2]float32
	widths       [2]float32
	headers      [2]string
	items        [2][]string
	disabled     bool
	disabledRows map[int]bool
	visible      map[int]bool
}

type menuState struct {
	next    int
	menus   map[int]*testMenu
	players *openMPState
}

func newMenuState() *menuState {
	return &menuState{menus: map[int]*testMenu{}}
}

func (state *menuState) Clone() scenarioModule {
	clone := newMenuState()

	clone.next = state.next
	for id, menu := range state.menus {
		menuCopy := *menu
		menuCopy.items[0] = append([]string(nil), menu.items[0]...)
		menuCopy.items[1] = append([]string(nil), menu.items[1]...)
		menuCopy.disabledRows = maps.Clone(menu.disabledRows)
		menuCopy.visible = maps.Clone(menu.visible)
		clone.menus[id] = &menuCopy
	}

	return clone
}

func (state *menuState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *menuState) natives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_menu_select":    state.selectMenuRow,
		"__pt_menu_exit":      state.exitMenu,
		"__pt_menu_valid":     state.assertValid(result),
		"__pt_menu_visible":   state.assertVisible(result),
		"__pt_menu_items":     state.assertItems(result),
		"CreateMenu":          state.createMenu,
		"DestroyMenu":         state.destroyMenu,
		"AddMenuItem":         state.addMenuItem,
		"SetMenuColumnHeader": state.setMenuColumnHeader,
		"ShowMenuForPlayer":   state.showMenuForPlayer,
		"HideMenuForPlayer":   state.hideMenuForPlayer,
		"IsValidMenu":         state.isValidMenu,
		"DisableMenu":         state.disableMenu,
		"DisableMenuRow":      state.disableMenuRow,
		"GetPlayerMenu":       state.getPlayerMenu,
		"IsMenuDisabled":      state.isMenuDisabled,
		"IsMenuRowDisabled":   state.isMenuRowDisabled,
		"GetMenuColumns":      state.getMenuColumns,
		"GetMenuItems":        state.getMenuItems,
		"GetMenuPos":          state.getMenuPosition,
		"GetMenuColumnWidth":  state.getMenuWidths,
		"GetMenuColumnHeader": state.getMenuColumnHeader,
		"GetMenuItem":         state.getMenuItem,
	}
}

func (state *menuState) createMenu(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 5 {
		return -1, nil
	}

	title, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	secondWidth := float32(0)
	if len(params) > 5 {
		secondWidth = cellFloat(params[5])
	}

	id := state.next
	state.next++
	state.menus[id] = &testMenu{
		title: title, columns: int(params[1]), position: [2]float32{cellFloat(params[2]), cellFloat(params[3])},
		widths: [2]float32{cellFloat(params[4]), secondWidth}, disabledRows: map[int]bool{}, visible: map[int]bool{},
	}

	return backend.Cell(id), nil
}

func (state *menuState) menu(params []backend.Cell) (*testMenu, bool) {
	if len(params) == 0 {
		return nil, false
	}

	menu, ok := state.menus[int(params[0])]

	return menu, ok
}

func (state *menuState) destroyMenu(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.menu(params); !ok {
		return 0, nil
	}

	delete(state.menus, int(params[0]))

	return 1, nil
}

func (state *menuState) isValidMenu(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if _, ok := state.menu(params); ok {
		return 1, nil
	}

	return 0, nil
}
