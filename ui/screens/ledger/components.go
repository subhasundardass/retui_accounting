package ledger

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/screens/ledger_create"
)

type Components struct {
	controller *Controller
}

func NewComponents(controller *Controller) *Components {
	return &Components{controller: controller}
}

func (c *Components) RenderScreen() retui.Element {

	groupID := 0
	params := retui.CurrentScreenParams()

	if params != nil {
		if id, ok := params["groupID"].(int); ok {
			groupID = id
		}
	}

	ledgers, setLedger := retui.UseState([]*ent.Ledger{})
	selected, setSelected := retui.UseState(&ent.Ledger{})

	retui.UseEffect(func() func() {
		legs := c.controller.GetLedgers(groupID)
		setLedger(legs)

		return nil
	}, []any{ledgers})

	//--Key
	switch retui.CurrentKey.Code {
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")
		win := ledger_create.Window(c.controller.ctx)
		win.Show()
	}

	rows := make([][]string, len(ledgers))
	for i, rs := range ledgers {
		groupName := ""
		if rs.Edges.Group != nil {
			groupName = rs.Edges.Group.Name
		}

		rows[i] = []string{
			rs.Code,
			rs.Name,
			groupName, // Display group name instead of ID
			rs.Description,
			fmt.Sprintf("%.2f", rs.Balance),
			map[bool]string{true: "Yes", false: "No"}[rs.IsBank],
			map[bool]string{true: "Yes", false: "No"}[rs.IsCash],
			map[bool]string{true: "Yes", false: "No"}[rs.IsParty],
			map[bool]string{true: "Yes", false: "No"}[rs.IsSystem],
		}
	}

	return retui.Box(
		retui.Props{Direction: retui.Column,
			Width: retui.Grow(1), Height: retui.Grow(1),
		},
		retui.NewStyle(),
		c.header(selected),
		c.buildLedgerTable(ledgers, rows, setSelected),
	)
}

func (c *Components) header(selected *ent.Ledger) retui.Element {

	title := "Ledgers"
	if selected != nil {
		title = fmt.Sprintf("Ledgers  %s", selected.Name)
	}

	return retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Justify:   retui.JustifySpaceBetween,
			Align:     retui.AlignCenter,
		},
		retui.NewStyle().Foreground(retui.BrightCyan).
			Border(retui.Border{Bottom: true, Left: true, Right: true, Top: true, Color: retui.Gray(1)}),
		retui.Text(title, retui.NewStyle().Bold(true)),
		retui.Text("Create <F2>", retui.NewStyle().Bold(true).Foreground(retui.Gold)),
	)
}

func (c *Components) buildLedgerTable(
	ledgers []*ent.Ledger,
	rows [][]string,
	setSelected func(*ent.Ledger),

) retui.Element {

	tbl := components.Table().
		ID("ledger_table").
		Headers([]string{"Code", "Name", "Group", "Description", "Balance", "Is Bank", "Is Cash", "Is Party", "Is System"}).
		Alignments([]string{
			"left",   // Code
			"left",   // Name
			"left",   // Group
			"left",   // Description
			"right",  // Balance - RIGHT ALIGNED
			"center", // Bank
			"center", // Cash
			"center", // Party
			"center", // System
		}).
		Focused(true).
		Rows(rows).
		SelectedIndex(0).
		HeaderColor(retui.Cyan).
		ColumnWidths([]int{15, 30, 30, 30, 10}).
		OnChange(func(i int) {
			setSelected(ledgers[i])
			if retui.CurrentKey.Code == retui.KeyEnter {
				c.controller.GetJournalsByLedger(ledgers[i].ID)
			}
		}).
		Render()

	return tbl
}
