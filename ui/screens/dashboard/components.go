package dashboard

import (
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type Components struct {
	controller *Controller
}

func NewComponents(controller *Controller) *Components {
	return &Components{controller: controller}
}

func (c *Components) RenderScreen() retui.Element {
	return retui.Box(
		retui.Props{Direction: retui.Column, Gap: 0},
		retui.NewStyle(),

		components.Panel().
			Header(
				retui.Box(
					retui.Props{
						Width:   retui.Percent(100),
						Padding: [4]int{0, 1, 0, 1}},
					retui.NewStyle().Bold(true).Foreground(retui.Cyan),
					retui.Text("Dashboard", retui.NewStyle()),
				),
			).
			Children(
				retui.Box(
					retui.Props{
						Gap:     1,
						Padding: [4]int{0, 1, 0, 1},
						Width:   retui.Grow(1),
						Height:  retui.Grow(1),
					},
					retui.NewStyle(),

					c.Component1(),
					c.Component1(),
					c.Component1(),
					c.Component1(),
					c.Component1(),
					c.Component1(),
				),
			).
			Render(),

		retui.Box(
			retui.Props{},
			retui.NewStyle(),
			components.Panel().
				Header(
					retui.Box(
						retui.Props{
							Width:   retui.Fixed(100),
							Padding: [4]int{0, 1, 0, 1}},
						retui.NewStyle().Bold(true).Foreground(retui.Cyan),
						retui.Text("Recent Purchase", retui.NewStyle()),
					),
				).
				Children(
					retui.Box(
						retui.Props{
							Gap:     1,
							Padding: [4]int{0, 1, 0, 1},
							// Width:   retui.Grow(1),
							Height: retui.Percent(20),
						},
						retui.NewStyle(),
						retui.Text("...", retui.NewStyle()),
					),
				).
				Render(),
			//Right
			components.Panel().
				FixedHeight(10).
				Header(
					retui.Box(
						retui.Props{
							Padding: [4]int{0, 1, 0, 1}},
						retui.NewStyle().Bold(true).Foreground(retui.Cyan),
						retui.Text("Recent Sales", retui.NewStyle()),
					),
				).
				Children(
					retui.Box(
						retui.Props{
							Gap:     1,
							Padding: [4]int{0, 1, 0, 1},
							Width:   retui.Grow(1),
							// Height: retui.Fixed(10),
						},
						retui.NewStyle(),
						retui.Box(
							retui.Props{
								Direction: retui.Column,
								// Height:    retui.Fixed(10),
							},
							retui.NewStyle(),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
							retui.Text("sample", retui.NewStyle()),
						),
					),
				).
				Render(),
		),
	)
}

func (c *Components) Component1() retui.Element {
	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Fixed(3),
			Padding:   [4]int{1, 1, 1, 1}},
		retui.NewStyle().Background(retui.Gray(1)),
		retui.Text("sadasd", retui.NewStyle()),
		retui.Text("10,350.00", retui.NewStyle().Bold(true).Foreground(retui.Turquoise)),
	)
}
