package views

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type LedgerComponent struct {
	ctx        *appctx.AppContext
	controller *ledger.LedgerController
}

func NewLedgerComponent(ctx *appctx.AppContext) *LedgerComponent {
	return &LedgerComponent{
		controller: ledger.NewController(ctx),
	}
}

func (c *LedgerComponent) bindKeys() {
	if !retui.IsFocused("ledger_list") {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")

	}
}

func (c *LedgerComponent) List(ctx *appctx.AppContext) retui.Element {

	groupID := 0
	params := retui.CurrentScreenParams()

	if params != nil {
		if id, ok := params["groupID"].(int); ok {
			groupID = id
		}
	}

	ledgers, setLedgers := retui.UseState([]*ent.Ledger{})
	selected, setSelected := retui.UseState(&ent.Ledger{})

	retui.UseEffect(func() func() {
		list, err := c.controller.List(groupID)
		if err != nil {
			retui.Errorf("Error fetching data %s", err.Error())
			return nil
		}

		setLedgers(list)
		return nil
	}, []any{ledgers})

	c.bindKeys()

	//-
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
		retui.Props{Direction: retui.Column},
		retui.NewStyle(),
		c.buildToolbar(selected),
		c.buildTable(ledgers, rows, setSelected),
	)
}

func (c *LedgerComponent) buildToolbar(selected *ent.Ledger) retui.Element {
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

func (c *LedgerComponent) buildTable(
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
		ColumnWidths([]int{15, 30, 30, 30, 10}).
		OnChange(func(i int) {
			setSelected(ledgers[i])
			if retui.CurrentKey.Code == retui.KeyEnter {
				// c.controller.GetJournalsByLedger(ledgers[i].ID)
			}
		}).
		Render()

	return retui.Box(
		retui.Props{
			// Direction: retui.Column,
			Height: retui.Fixed(33),
		},
		retui.NewStyle(),
		tbl,
	)
}
