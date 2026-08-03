package views

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type LedgerGroupComponent struct {
	ctx        *appctx.AppContext
	controller *ledger.LedgerController
}

func NewLedgerGroupComponent(ctx *appctx.AppContext) *LedgerGroupComponent {
	return &LedgerGroupComponent{
		controller: ledger.NewController(ctx),
	}
}

func (c *LedgerGroupComponent) bindKeys() {
	if !retui.IsFocused("ledger_group") {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")
	}
}

func (c *LedgerGroupComponent) List(ctx *appctx.AppContext) retui.Element {

	groups, setGroups := retui.UseState([]*ent.Ledger_Group{})
	selected, setSelected := retui.UseState(&ent.Ledger_Group{})

	retui.UseEffect(func() func() {
		list, err := c.controller.Groups()
		if err != nil {
			retui.Errorf("Error fetching Groups %s", err.Error())
			return nil
		}

		setGroups(list)
		return nil
	}, []any{groups})

	rows := make([][]string, len(groups))
	for i, rs := range groups {
		rows[i] = []string{rs.Code, rs.Name, string(rs.Nature), rs.Description}
	}

	c.bindKeys()

	return retui.Box(
		retui.Props{Direction: retui.Column},
		retui.NewStyle(),
		c.buildToolbar(selected),
		c.buildGroupTable(groups, rows, setSelected),
	)
}

func (c *LedgerGroupComponent) buildToolbar(selected *ent.Ledger_Group) retui.Element {
	title := "Groups"
	if selected != nil {
		title = fmt.Sprintf("Groups :  %s", selected.Name)
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

func (c *LedgerGroupComponent) buildGroupTable(
	groups []*ent.Ledger_Group,
	rows [][]string,
	setSelected func(*ent.Ledger_Group),
) retui.Element {
	tbl := components.Table().
		Headers([]string{"Code", "Name", "Nature", "Description"}).
		Rows(rows).
		// Width(210).
		// Height(33).
		HeaderColor(retui.Cyan).
		ColumnWidths([]int{20, 30, 20, 80}).
		ShowBorders(true).
		SelectedIndex(0).
		Focused(true).
		OnChange(func(i int) {
			setSelected(groups[i])

		}).Render()

	return retui.Box(
		retui.Props{
			// Direction: retui.Column,
			Height: retui.Fixed(33),
		},
		retui.NewStyle(),
		tbl,
	)

}
