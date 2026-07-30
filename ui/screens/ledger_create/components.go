package ledger_create

import (
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

const totalFields = 7 // Name, Code, Group, OpeningBal, IsActive

type Components struct {
	controller *Controller

	state    LedgerState
	setState func(LedgerState)

	win *window.Window
}

type LedgerState struct {
	FocusIndex int

	Name string
	Code string

	GroupID string

	OpeningBal  float64
	IsActive    bool
	Description string

	IsParty bool
	IsCash  bool
	IsBank  bool
}

func NewComponents(controller *Controller) *Components {
	return &Components{
		controller: controller,
	}
}

func (c *Components) RenderWindow() *window.Window {

	c.state, c.setState = retui.UseState(LedgerState{IsActive: true})

	c.win = window.NewWindow().
		SetTitle("Create Ledger").
		SetModal(true).
		Center().
		SetSize(120, 40)

	c.win.SetRenderFn(func() retui.Element {
		c.state, c.setState = retui.UseState(LedgerState{})
		return c.render()
	})

	c.bindKeys()

	return c.win
}

func (c *Components) bindKeys() {

	c.win.OnKeyPress(func(key retui.Key) bool {

		retui.Debugf("Window captured=%q key=%v",
			retui.CapturedFocus(),
			key.Code,
		)

		if retui.CapturedFocus() != "" {
			return false
		}

		switch key.Code {

		// case retui.KeyEscape:

		// 	c.win.Close()
		// 	return true

		case retui.KeyDown:
			retui.Debugf("================ Down%v", totalFields)
			s := c.state
			s.FocusIndex = (s.FocusIndex + 1) % totalFields
			c.setState(s)
			return true

		case retui.KeyUp:
			retui.Debugf("================ Up%v", c.state.FocusIndex)

			s := c.state
			s.FocusIndex = (s.FocusIndex - 1 + totalFields) % totalFields
			c.setState(s)
			return true

		}

		return false
	})
}

func (c *Components) render() retui.Element {

	isFocused := func(index int) bool {
		return c.state.FocusIndex == index
	}

	row := func(label string, input retui.Element) retui.Element {
		return retui.Box(
			retui.Props{Direction: retui.Row, Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(15)},
				retui.NewStyle(),
				retui.Text(label, retui.NewStyle()),
			),

			input,
		)
	}

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Gap:       0,
			Padding:   [4]int{1, 1, 0, 1},
		},
		retui.NewStyle(),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),

			row("Ledger Name", components.TextInput().
				ID("name").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(0)).
				Value(c.state.Name).
				OnChange(func(id, value string) {
					s := c.state
					s.Name = value
					c.setState(s)
				}).
				Render(),
			),

			row("Code", components.TextInput().
				ID("code").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(1)).
				Value(c.state.Code).
				OnChange(func(id, value string) {
					s := c.state
					s.Code = value
					c.setState(s)
				}).
				Render(),
			),
		),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			row("Opening Bal.", components.NumberInput().
				ID("opening_bal").
				Width(40).
				Focused(isFocused(2)).
				Value(c.state.OpeningBal).
				Decimals(2).
				OnChange(func(id string, value float64) {
					s := c.state
					s.OpeningBal = value
					c.setState(s)
				}).
				Render(),
			),

			row("Group",
				components.SelectDropdown().
					ID("ledger_group").
					Focused(isFocused(3)).
					OverlayAbsPos(80, 5).
					OnFilter(func(id, query string) []components.SelectOption {
						return c.controller.LedgerFilterOptions(query)
					}).
					Value(c.state.GroupID).
					OnChange(func(id, value string) {
						s := c.state
						s.GroupID = value
						c.setState(s)
					}).
					Render()),
		),

		//Description
		row("Description", components.TextArea().
			ID("description").
			Width(100).
			Height(2).
			Focused(isFocused(4)).
			Value(c.state.Description).
			OnChange(func(id, value string) {
				s := c.state
				s.Description = value
				c.setState(s)
			}).
			Render(),
		),

		row("Active", components.Checkbox().
			ID("is_active").
			Focused(isFocused(5)).
			Checked(c.state.IsActive).
			OnChange(func(id string, value bool) {
				s := c.state
				s.IsActive = value
				c.setState(s)
			}).
			Render(),
		),

		retui.Box(
			retui.Props{Padding: [4]int{1, 0, 0, 0}},
			retui.NewStyle(),
			components.Button().
				ID("submit").
				Focused(isFocused(6)).
				Label("Submit").
				Render(),
		),
	)
}
