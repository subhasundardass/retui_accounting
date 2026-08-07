package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type GroupFormComponent struct {
	controller *ledger.LedgerController
	win        *window.Window

	state    ledger.LedgerGroupState
	setState func(ledger.LedgerGroupState)

	editing bool
	editID  int

	// onSaved is called after a successful save so the caller (e.g. the
	// group list) can refresh its data.
	onSaved func()
}

var natureOptions = []components.SelectOption{
	{Label: "Asset", Value: "ASSET"},
	{Label: "Liability", Value: "LIABILITY"},
	{Label: "Income", Value: "INCOME"},
	{Label: "Expense", Value: "EXPENSE"},
	{Label: "Equity", Value: "EQUITY"},
}

func NewGroupFormComponent(ctx *appctx.AppContext) *GroupFormComponent {
	return &GroupFormComponent{
		controller: ledger.NewController(ctx),
		state:      ledger.LedgerGroupState{Mode: ledger.ModeCreate},
	}
}

// SetOnSaved registers a callback fired after a successful save.
func (c *GroupFormComponent) SetOnSaved(fn func()) {
	c.onSaved = fn
}

// func (c *GroupFormComponent) save() {
// 	mode := ledger.ModeCreate
// 	id := 0
// 	if c.editing {
// 		mode = ledger.ModeUpdate
// 		id = c.editID
// 	}

// 	_, err := c.controller.CreateOrUpdate(mode, id, c.state)
// 	if err != nil {
// 		retui.Debugf("Save failed: %v", err)
// 		components.ShowError("Save failed: " + err.Error())
// 		// TODO: surface this error in the UI (e.g. an Errors/status field on state)
// 		return
// 	}

// 	if c.onSaved != nil {
// 		c.onSaved()
// 	}
// 	c.win.Close()
// }

// NOTE: bindKeys is currently called from inside buildWindow, which is the
// SetRenderFn callback and therefore runs on every render. If
// window.OnKeyPress APPENDS handlers rather than replacing the previous one,
// this will register a new closure every render and key presses will fire
// multiple times. Verify window.OnKeyPress's replace-vs-append behavior; if
// it appends, move this binding to run once (e.g. at window creation) instead
// of on every render.
func (c *GroupFormComponent) bindKeys(form *retui.Form[ledger.LedgerGroupState]) {
	c.win.OnKeyPress(func(key retui.Key) bool {
		if retui.CapturedFocus() != "" {
			return false
		}

		v := form.Values()

		switch key.Code {
		case retui.KeyDown, retui.KeyTab:
			v.FocusIndex = (v.FocusIndex + 1) % 6
			form.SetValuesSilent(v)
			return true

		case retui.KeyUp, retui.KeyShiftTab:
			v.FocusIndex = (v.FocusIndex - 1 + 6) % 6
			form.SetValuesSilent(v)
			return true

		case retui.KeyEscape:
			c.win.Close()
			return true

		case retui.KeyF10:
			// err := c.controller.SaveGroup(form.Values())
			// if err != nil {
			// 	components.ShowError(err.Error())
			// 	return true
			// }
			// components.ShowSuccess("Group saved successfully.")
			// if c.onSaved != nil {
			// 	c.onSaved()
			// }
			c.win.Close()
			return true
		}
		return false
	})
}

// GroupCreateForm opens (or recreates) the window in create mode. It always
// resets state and rebuilds the window so it can't return a stale window
// left over from a previous edit session.
func (c *GroupFormComponent) GroupCreateForm(ctx *appctx.AppContext) *window.Window {
	c.editing = false
	c.editID = 0
	c.state = ledger.LedgerGroupState{Mode: ledger.ModeCreate}

	c.win = window.NewWindow().
		SetTitle("Create Group").
		SetModal(true).
		Center().
		SetSize(80, 40)

	c.win.SetRenderFn(func() retui.Element {
		return c.buildWindow() // bindKeys wired inside buildWindow
	})

	return c.win
}

// GroupEditForm builds the window used for editing. Call OpenForEdit first
// to populate c.state with the record being edited.
func (c *GroupFormComponent) GroupEditForm(ctx *appctx.AppContext) *window.Window {
	c.win = window.NewWindow().
		SetTitle("Edit Group").
		SetModal(true).
		Center().
		SetSize(80, 40)

	c.win.SetRenderFn(func() retui.Element {
		return c.buildWindow() // bindKeys wired inside buildWindow
	})

	return c.win
}

// OpenForEdit fetches an existing group, populates form state, and opens
// the edit window.
func (c *GroupFormComponent) OpenForEdit(id int, ctx *appctx.AppContext) *window.Window {
	state, err := c.controller.GetGroup(id)
	if err != nil {
		components.ShowError(err.Error())
		return nil
	}

	c.editing = true
	c.editID = id

	c.state = ledger.LedgerGroupState{
		Mode:        ledger.ModeUpdate,
		Code:        state.Code,
		Name:        state.Name,
		Nature:      state.Nature,
		IsSystem:    state.IsSystem,
		Description: state.Description,
	}

	win := c.GroupEditForm(ctx)
	win.Show()
	return win
}

func (c *GroupFormComponent) buildWindow() retui.Element {
	// c.state is always the correct seed value here: GroupCreateForm resets
	// it to a zero-value ModeCreate state, and OpenForEdit populates it with
	// the fetched record before GroupEditForm/buildWindow runs.
	form := retui.UseForm(c.state)
	v := form.Values()

	// Wire bindKeys here where form is available
	c.bindKeys(form)

	code := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Text("Code :", retui.NewStyle()),
		),
		components.TextInput().
			ID("code").
			Value(v.Code).
			Focused(v.FocusIndex == 0).
			OnChange(func(id, value string) {
				if err := form.SetField("Code", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	name := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Text("Name :", retui.NewStyle()),
		),
		components.TextInput().
			ID("name").
			Value(v.Name).
			Focused(v.FocusIndex == 1).
			OnChange(func(id, value string) {
				if err := form.SetField("Name", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	nature := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Text("Nature :", retui.NewStyle()),
		),
		components.SelectDropdown().
			ID("nature").
			Value(v.Nature).
			Focused(v.FocusIndex == 2).
			Options(natureOptions).
			// OnFilter(func(id, query string) []components.SelectOption {
			// 	return FilterOptions(natureOptions, query)
			// }).
			OnChange(func(id, value string) {
				if err := form.SetField("Nature", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			OverlayAbsPos(40, 10).
			Render(),
	)

	description := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Text("Desscription :", retui.NewStyle()),
		),
		components.TextArea().
			ID("description").
			Value(v.Description).
			Height(2).
			Focused(v.FocusIndex == 3).
			OnChange(func(id, value string) {
				if err := form.SetField("Description", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	return retui.Box(
		retui.Props{Direction: retui.Column, Padding: [4]int{1, 2, 1, 2}},
		retui.NewStyle(),
		code,
		name,
		nature,
		description,

		retui.Box(
			retui.Props{Gap: 1, Margin: [4]int{1}},
			retui.NewStyle(),
			components.Button().
				ID("save").
				Label("Save").
				Focused(v.FocusIndex == 4).
				Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
				OnKeyPress(func(id string, key retui.Key) bool {
					if key.Code == retui.KeyEnter {
						// c.save()
						return true
					}
					return false
				}).
				Render(),

			components.Button().
				ID("reset").
				Label("Reset").
				Hidden(c.editing).
				Focused(v.FocusIndex == 5).
				Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
				OnKeyPress(func(id string, key retui.Key) bool {
					if key.Code == retui.KeyEnter {
						form.Reset()
						return true
					}
					return false
				}).
				Render(),
		),
	)
}
