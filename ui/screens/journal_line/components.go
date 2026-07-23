package journal_line

import (
	"fmt"
	"strconv"

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

func (c *Components) RenderScreen(ledgerID int) retui.Element {

	journalLine, setJournalLine := retui.UseState([]*ent.Journal_Line{})

	retui.UseEffect(func() func() {
		line := c.controller.getJounalLineByLedger(ledgerID)
		setJournalLine(line)

		return nil
	}, []any{})

	return retui.Box(
		retui.Props{Direction: retui.Column, Gap: 0},
		retui.NewStyle(),

		c.header(journalLine),
		c.buildTable(journalLine),
	)
}

func (c *Components) header(lines []*ent.Journal_Line) retui.Element {
	if len(lines) == 0 {
		return retui.Box(
			retui.Props{Padding: [4]int{0, 1, 0, 1}},
			retui.NewStyle().Bold(true).Border(retui.Border{Top: true, Left: true, Bottom: true, Right: true}),
			retui.Text("No entries", retui.NewStyle().Bold(true)),
		)
	}

	ledger := lines[0].Edges.Ledger
	if ledger == nil {
		return retui.Box(
			retui.Props{Padding: [4]int{0, 1, 0, 1}},
			retui.NewStyle().Bold(true).Border(retui.Border{Top: true, Left: true, Bottom: true, Right: true}),
			retui.Text("Unknown ledger", retui.NewStyle().Bold(true)),
		)
	}

	return retui.Box(
		retui.Props{Padding: [4]int{0, 1, 0, 1}},
		retui.NewStyle().Bold(true).Border(retui.Border{Top: true, Left: true, Bottom: true, Right: true}),
		retui.Text(ledger.Name, retui.NewStyle().Bold(true)),
	)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

func (c *Components) buildTable(lines []*ent.Journal_Line) retui.Element {

	rows := make([][]string, len(lines))
	for i, line := range lines {
		voucherNo := ""
		voucherType := ""
		if line.Edges.Journal != nil {
			voucherNo = line.Edges.Journal.VoucherNo
			voucherType = line.Edges.Journal.VoucherType
		}

		rows[i] = []string{
			line.Edges.Journal.Date.Format("02/01/2006"),
			voucherNo,
			voucherType,
			fmt.Sprintf("%.2f", line.Debit),
			fmt.Sprintf("%.2f", line.Credit),
			derefString(line.Description),
			derefString(line.Edges.Journal.ReferenceNo),
		}
	}

	return components.Table().
		ID("journal_line_table").
		Headers([]string{"Date", "Voucher No", "Type", "Debit", "Credit", "Description", "Reference"}).
		Alignments([]string{"left", "left", "left", "right", "right"}).
		ColumnWidths([]int{10, 15, 10, 15, 15, 30, 15}).
		Focused(true).
		Rows(rows).
		SelectedIndex(0).
		HeaderColor(retui.Cyan).
		OnChange(func(i int) {
			if retui.CurrentKey.Code == retui.KeyEnter {
				jrID := lines[i].JournalID
				retui.SetFocus("journal_view")
				retui.PushScreen("journal_view", retui.ScreenParams{"journalID": jrID})
			}
		}).
		Render()

}
