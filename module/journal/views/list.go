package views

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type JournalListComponent struct {
	controller *journal.JournalController
	ctx        *appctx.AppContext
	// form       *FormComponent
}

func NewJournalListComponent(ctx *appctx.AppContext) *JournalListComponent {
	return &JournalListComponent{
		controller: journal.NewController(ctx),
		ctx:        ctx,
		// form:       form,
	}
}

func (c *JournalListComponent) bindKeys() {
	if !retui.IsFocused("journal_list") {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")
		win := JournalCreateWindow(c.ctx)
		win.Show()
	}
}

func (c *JournalListComponent) List(ctx *appctx.AppContext) retui.Element {

	journals, setJournals := retui.UseState([]*ent.Journal{})
	selected, setSelected := retui.UseState(&ent.Journal{})

	retui.UseEffect(func() func() {
		list, err := c.controller.ListWithPagination(0, 40)
		if err != nil {
			retui.Errorf("Error fetching data %s", err.Error())
			return nil
		}

		setJournals(list)
		return nil
	}, []any{journals})

	c.bindKeys()

	return retui.Box(
		retui.Props{Direction: retui.Column},
		retui.NewStyle(),
		c.buildToolbar(selected),
		c.buildTable(journals, setSelected),
	)
}

func (c *JournalListComponent) buildToolbar(selected *ent.Journal) retui.Element {
	title := "Journals"
	if selected != nil {
		title = fmt.Sprintf("Journals  %s", selected.ID)
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

func (c *JournalListComponent) buildTable(
	journals []*ent.Journal,
	setSelected func(*ent.Journal),
) retui.Element {

	rows := make([][]string, len(journals))
	for i, j := range journals {

		rows[i] = []string{
			j.VoucherDate.Format("02/01/2006"),
			string(j.VoucherNo),
			string(*j.ReferenceNo),
			string(j.VoucherType),
			fmt.Sprintf("%.2f", j.TotalDebit),
			fmt.Sprintf("%.2f", j.TotalCredit),
			*j.Narration,
			string(j.JournalStatus),
		}
	}

	tbl := components.Table().
		ID("journal_table").
		Headers([]string{
			"Date",
			"Voucher No",
			"Reference",
			"Type",
			"Debit",
			"Credit",
			"Narration",
			"Status",
		}).
		Alignments([]string{
			"left",
			"left",
			"left",
			"left",
			"right",
			"right",
			"left",
			"center",
		}).
		Focused(true).
		Rows(rows).
		SelectedIndex(0).
		ColumnWidths([]int{
			15,
			15,
			15,
			15,
			20,
			20,
			60,
			10,
		}).
		OnChange(func(i int) {
			if i >= 0 && i < len(journals) {
				setSelected(journals[i])
			}

			if retui.CurrentKey.Code == retui.KeyEnter {
				// fmt.Print(journals[i].ID)
				// c.controller.ShowJournal(journals[i].ID)
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
