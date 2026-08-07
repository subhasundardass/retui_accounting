package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type LedgerFormComponent struct {
	controller *ledger.LedgerController
	win        *window.Window

	state ledger.LedgerState

	editing bool
	editID  int
}

func NewLedgerFormComponent(ctx *appctx.AppContext) *LedgerFormComponent {
	return &LedgerFormComponent{
		controller: ledger.NewController(ctx),
	}
}

func (c *LedgerFormComponent) bindKeys(form *retui.Form[ledger.LedgerState]) {
	c.win.OnKeyPress(func(key retui.Key) bool {
		if retui.CapturedFocus() != "" {
			return false
		}

		v := form.Values()

		switch key.Code {
		case retui.KeyEscape:
			c.win.Close()
			return true

		case retui.KeyDown, retui.KeyTab:
			v.FocusIndex = (v.FocusIndex + 1) % 6
			form.SetValuesSilent(v)
			return true

		case retui.KeyUp, retui.KeyShiftTab:
			v.FocusIndex = (v.FocusIndex - 1 + 6) % 6
			form.SetValuesSilent(v)
			return true
		}

		return false
	})
}

// func (c *LedgerFormComponent) OpenForEdit(id int, ctx *appctx.AppContext) *window.Window {}

func (c *LedgerFormComponent) LedgerCreateForm() *window.Window {

	c.win = window.NewWindow().
		SetTitle("Edit Ledger").
		SetModal(true).
		Center().
		SetSize(150, 60)

	c.win.SetRenderFn(func() retui.Element {
		return c.buildWindow()
	})

	return c.win
}

func (c *LedgerFormComponent) buildWindow() retui.Element {

	form := retui.UseForm(c.state)
	v := form.Values()

	// Wire bindKeys here where form is available
	c.bindKeys(form)

	// Code
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

	//Name
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

	//Group
	group := retui.Box(
		retui.Props{
			Gap: 1,
		},
		retui.NewStyle(),

		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Text("Group :", retui.NewStyle()),
		),
		components.SelectDropdown().
			ID("group").
			Focused(v.FocusIndex == 2).
			OverlayAbsPos(80, 5).
			OnFilter(func(id, query string) []components.SelectOption {
				return c.controller.LedgerGroupFilterOptions(query)
			}).
			Value(v.GroupID).
			OnChange(func(id, value string) {
				if err := form.SetField("GroupID", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	//Description
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
		group,
		description,
	)
}
