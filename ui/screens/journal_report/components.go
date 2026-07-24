package journal_report

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
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

	journals, setJournals := retui.UseState([]*ent.Journal{})
	selected, setSelected := retui.UseState((*ent.Journal)(nil))

	retui.UseEffect(func() func() {
		jrnls, err := c.controller.GetJournalsPaginated(0, 40)
		if err != nil {
			retui.Error(err)
			return nil
		}

		setJournals(jrnls)

		if len(jrnls) > 0 {
			setSelected(jrnls[0])
		}

		return nil
	}, []any{})

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

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Gap:       0,
		},
		retui.NewStyle(),

		c.header(selected),

		c.buildTable(
			journals,
			rows,
			setSelected,
		),
	)
}

func (c *Components) header(selected *ent.Journal) retui.Element {

	title := "Journal"
	if selected != nil {
		title = fmt.Sprintf("Journal  %s", selected.Date)
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
	)
}

func (c *Components) buildTable(
	journals []*ent.Journal,
	rows [][]string,
	setSelected func(*ent.Journal),
) retui.Element {

	return components.Table().
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
		// Width(110).
		// Height(33).
		SelectedIndex(0).
		HeaderColor(retui.Cyan).
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
				c.controller.ShowJournal(journals[i].ID)
			}

		}).
		Render()

}
