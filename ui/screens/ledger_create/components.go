package ledger_create

import (
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

const totalFields = 5 // Name, Code, Group, OpeningBal, IsActive

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

	GroupID   int
	GroupName string

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

		switch key.Code {

		case retui.KeyEscape:
			c.win.Close()
			return true

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

		case retui.KeyEnter:
			c.submit()
			return true
		}

		return false
	})
}

func (c *Components) submit() {
	if c.state.Name == "" || c.state.Code == "" || c.state.GroupID == 0 {
		return // nothing selected/filled yet, ignore
	}

	// openingBal, _ := strconv.ParseFloat(c.state.OpeningBal, 64)

	// c.controller.CreateLedger(ent.Ledger{
	// 	Name:       c.state.Name,
	// 	Code:       c.state.Code,
	// 	GroupID:    c.state.GroupID,
	// 	OpeningBal: openingBal,
	// 	IsActive:   c.state.IsActive,
	// })

	c.win.Close()
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
		},
		retui.NewStyle(),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),

			row("Ledger Name", components.TextInput().
				ID("name").
				Width(40).
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
				Width(40).
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

			row("Active", components.Checkbox().
				ID("is_active").
				Focused(isFocused(3)).
				Checked(c.state.IsActive).
				OnChange(func(id string, value bool) {
					s := c.state
					s.IsActive = value
					c.setState(s)
				}).
				Render(),
			),
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

		// row("Group", components.Select().
		// 	ID("group").
		// 	Width(40).
		// 	Focused(isFocused(2)).
		// 	Options(c.groupOptions()).
		// 	Value(strconv.Itoa(c.state.GroupID)).
		// 	OnSelect(func(id, groupID, groupName string) {
		// 		gid, err := strconv.Atoi(groupID)
		// 		if err != nil {
		// 			return
		// 		}
		// 		s := c.state
		// 		s.GroupID = gid
		// 		s.GroupName = groupName
		// 		c.setState(s)
		// 	}).
		// 	Render(),
		// ),

	)
}

// func (c *Components) groupOptions() []components.SelectOption {
// 	groups := c.controller.ListGroups()
// 	opts := make([]components.SelectOption, len(groups))
// 	for i, g := range groups {
// 		opts[i] = components.SelectOption{ID: strconv.Itoa(g.ID), Label: g.Name}
// 	}
// 	return opts
// }
