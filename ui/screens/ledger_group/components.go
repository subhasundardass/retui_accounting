package ledger_group

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type Components struct {
	controller          *Controller
	groupEditComponents *GroupEditComponents
}

func NewComponents(controller *Controller) *Components {
	return &Components{
		controller:          controller,
		groupEditComponents: NewGroupEditComponents(controller),
	}
}

func (c *Components) RenderScreen() retui.Element {

	//--init
	groups := c.controller.LoadGroups()
	selectedIndex, setSelectedIndex := retui.UseState(0)
	selectedGroup, setSelectedGroup := retui.UseState(ent.Ledger_Group{})

	rows := make([][]string, len(groups))
	for i, rs := range groups {
		rows[i] = []string{rs.Code, rs.Name, string(rs.Nature), rs.Description}
	}

	header := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			// Width:     retui.Grow(1),
			Justify: retui.JustifySpaceBetween,
			Align:   retui.AlignCenter,
		},
		retui.NewStyle().Foreground(retui.BrightCyan).
			Border(retui.Border{Bottom: true, Left: true, Right: true, Top: true, Color: retui.Gray(1)}),
		retui.Text("Ledger Groups", retui.NewStyle()),
		retui.Text("F2: Edit   F4: Create New", retui.NewStyle().Bold(true)),
	)

	switch retui.CurrentKey.Code {
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......%v", selectedGroup)
		retui.Debugf("F2 Pressed.......%v", selectedGroup.Code)
		editWindow := c.groupEditComponents.WindowEditGroup(selectedGroup)
		editWindow.Show()
		retui.SetFocus(editWindow.ID)

	}

	tbl := components.Table().
		Headers([]string{"Code", "Name", "Nature", "Description"}).
		Rows(rows).
		// Width(210).
		// Height(33).
		HeaderColor(retui.Cyan).
		ColumnWidths([]int{20, 30, 20, 80}).
		ShowBorders(true).
		SelectedIndex(selectedIndex).
		Focused(true).
		OnChange(func(i int) {
			setSelectedIndex(i)
			setSelectedGroup(*groups[i])

			if retui.CurrentKey.Code == retui.KeyEnter {
				selectedGroup := groups[i]
				c.controller.ShowLedgers(selectedGroup.ID)
			}
		})

	return retui.Box(
		retui.Props{Direction: retui.Column, Gap: 0, Width: retui.Grow(1), Height: retui.Grow(1)},
		retui.NewStyle(),
		header,
		tbl.Render(),
	)
}
