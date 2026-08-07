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
	formComp   *GroupFormComponent
}

func NewLedgerGroupComponent(ctx *appctx.AppContext, form *GroupFormComponent) *LedgerGroupComponent {
	return &LedgerGroupComponent{
		ctx:        ctx,
		controller: ledger.NewController(ctx),
		formComp:   form,
	}
}

func (c *LedgerGroupComponent) bindKeys() {
	if !retui.IsFocused("ledger_group") {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyAltC:
		retui.Debugf("Alt+C Pressed.......")
		// win :=
		if c.formComp == nil {
			c.formComp = NewGroupFormComponent(c.ctx)
			c.formComp.SetOnSaved(func() {
				// refresh list after save
			})
		}
		win := c.formComp.GroupCreateForm(c.ctx)
		win.Show()
		retui.CurrentKey.Consumed = true
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
		title = fmt.Sprintf("Groups   %s", selected.Name)
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
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Text("Create <Alt+C>", retui.NewStyle().Bold(true).Foreground(retui.Gold)),
		),
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
		HeaderColor(retui.Cyan).
		ColumnWidths([]int{20, 30, 20, 80}).
		ShowBorders(true).
		SelectedIndex(0).
		Focused(true).
		OnChange(func(i int) {

			if i < 0 || i >= len(rows) {
				return
			}

			setSelected(groups[i])

			// Edit:
			if retui.CurrentKey.Code == retui.KeyEnter {
				// if err := c.formComp.LoadForEdit(companies[i].ID); err != nil {
				// 	retui.Debugf("failed to load company for edit: %v", err)
				// 	return
				// }

				// win := c.formComp.GroupEditForm(c.ctx)
				// if win != nil {
				// 	win.Show()
				// }

				// c.controller.EditGroup(groups[i].ID)
				c.formComp.OpenForEdit(groups[i].ID, c.ctx)
			}

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
