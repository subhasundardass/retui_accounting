package createcompany

import (
	"strings"

	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type State struct {
	FocusIndex int

	//--
	Code       string
	Name       string
	LegalName  string
	Email      string
	Phone      string
	Website    string
	Country    string
	State      string
	City       string
	PostalCode string
	Address    string
	TaxID      string
	GSTIN      string
	PAN        string
}

type Components struct {
	controller *Controller

	state    State
	setState func(State)

	win *window.Window

	form CompanyForm
}

type CompanyForm struct {
	Code       string
	Name       string
	LegalName  string
	Email      string
	Phone      string
	Website    string
	Country    string
	State      string
	City       string
	PostalCode string
	Address    string
	TaxID      string
	GSTIN      string
	PAN        string
}

func NewComponents(controller *Controller) *Components {
	return &Components{
		controller: controller,
	}
}

const totalFields = 17

func (c *Components) bindKeys() {

	c.win.OnKeyPress(func(key retui.Key) bool {

		if retui.CapturedFocus() != "" {
			return false
		}

		// retui.Debugf("============", retui.CurrentFocus())

		switch key.Code {
		case retui.KeyDown:
			s := c.state
			s.FocusIndex = (s.FocusIndex + 1) % totalFields
			c.setState(s)
			return true

		case retui.KeyUp:

			s := c.state
			s.FocusIndex = (s.FocusIndex - 1 + totalFields) % totalFields
			c.setState(s)
			return true

		}
		return false
	})
}

func (c *Components) RenderWindow() *window.Window {

	c.state, c.setState = retui.UseState(State{})

	c.win = window.NewWindow().
		SetTitle("Create Company").
		SetModal(true).
		Center().
		SetSize(100, 40)

	c.win.SetRenderFn(func() retui.Element {
		c.state, c.setState = retui.UseState(State{})
		return c.render()
	})

	c.bindKeys()
	return c.win
}

func (c *Components) render() retui.Element {

	isFocused := func(index int) bool {
		return c.state.FocusIndex == index
	}

	row1 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(20),
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(5),
				},
				retui.NewStyle(),
				retui.Text("Code", retui.NewStyle()),
			),
			components.TextInput().
				ID("code").
				Focused(isFocused(0)).
				Value(c.state.Code).
				OnChange(func(id, value string) {
					var b strings.Builder

					for _, r := range strings.ToUpper(value) {
						if r >= 'A' && r <= 'Z' {
							b.WriteRune(r)
						}
					}

					s := c.state
					s.Code = b.String()
					c.setState(s)
				}).
				Render(),
		),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Name", retui.NewStyle()),
			),
			components.TextInput().
				ID("name").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(1)).
				Render(),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Legal Name", retui.NewStyle()),
			),
			components.TextInput().
				ID("legal_name").
				Width(retui.Fixed(40).Value).
				Focused(isFocused(2)).
				Render(),
		),
	)

	row2 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Email", retui.NewStyle()),
			),
			components.TextInput().
				ID("email").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(3)).
				Render(),
		),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Phone", retui.NewStyle()),
			),
			components.TextInput().
				ID("phone").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(4)).
				Render(),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Website", retui.NewStyle()),
			),
			components.TextInput().
				ID("website").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(5)).
				Render(),
		),
	)
	row3 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Country", retui.NewStyle()),
			),
			components.TextInput().
				ID("Country").
				Focused(isFocused(6)).
				Render(),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("State", retui.NewStyle()),
			),
			components.TextInput().
				ID("state").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(7)).
				Render(),
		),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("City", retui.NewStyle()),
			),
			components.TextInput().
				ID("city").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(8)).
				Render(),
		),
	)

	row4 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Postal Code", retui.NewStyle()),
			),
			components.TextInput().
				ID("postal_code").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(9)).
				Render(),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Address", retui.NewStyle()),
			),
			components.TextInput().
				ID("address").
				Width(retui.Fixed(72).Value).
				Focused(isFocused(10)).
				Render(),
		),
	)

	row10 := retui.Box(
		retui.Props{
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(11),
				},
				retui.NewStyle(),
				retui.Text("Reg. No", retui.NewStyle()),
			),
			components.TextInput().
				ID("tax_id").
				Focused(isFocused(11)).
				Render(),
		),

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("GSTIN", retui.NewStyle()),
			),
			components.TextInput().
				ID("gstin").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(12)).
				Render(),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("PAN", retui.NewStyle()),
			),
			components.TextInput().
				ID("pan").
				Width(retui.Fixed(30).Value).
				Focused(isFocused(13)).
				Render(),
		),
	)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Padding:   [4]int{1, 1, 1, 1},
		},
		retui.NewStyle(),

		row1,
		retui.Box(
			retui.Props{
				Margin: [4]int{1, 0, 0, 0},
			},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Contact",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),

			row2,
		),
		row3,
		row4,

		retui.Box(
			retui.Props{
				Margin: [4]int{1, 0, 0, 0},
			},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Tax Information",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),

			row10,
		),

		retui.Box(
			retui.Props{
				Margin: [4]int{1, 0, 0, 0},
			},
			retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(2), Title: &retui.BorderTitle{
				Text:  "Actions",
				Align: retui.AlignStart,
				Style: retui.NewStyle().Foreground(retui.Teal),
			}}),

			retui.Box(
				retui.Props{
					Gap: 1,
				},
				retui.NewStyle(),
				components.Button().
					ID("save").
					Label("Save").
					Focused(isFocused(14)).
					Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
					Render(),

				components.Button().
					ID("cancel").
					Label("Cancel").
					Focused(isFocused(15)).
					Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
					Render(),

				components.Button().
					ID("reset").
					Label("Reset").
					Focused(isFocused(16)).
					Style(retui.NewStyle().Background(retui.Gray(2)).Foreground(retui.BrightWhite)).
					Render(),
			),
		),
	)
}
