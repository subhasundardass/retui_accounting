package views

import (
	"fmt"
	"strings"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/reports"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type CashBookComponent struct {
	controller *reports.Controller
	ctx        *appctx.AppContext
}

func NewCashBookComponent(ctx *appctx.AppContext) *CashBookComponent {
	return &CashBookComponent{
		controller: reports.NewController(ctx),
		ctx:        ctx,
	}
}

func (c *CashBookComponent) bindKeys() {
	key := retui.CurrentKey
	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()

	}
}

func (c *CashBookComponent) Book(ctx *appctx.AppContext) retui.Element {

	journals, setJournal := retui.UseState([]*ent.Journal{})

	retui.UseEffect(func() func() {

		getJournal := func() {
			journals, err := c.controller.GetCashBookEntries()
			if err != nil {
				retui.Debug("Failed to load journals:", err)
				return
			}

			setJournal(journals)
		}

		getJournal()
		return nil
	}, []any{})

	c.bindKeys()
	return retui.Box(
		retui.Props{
			Direction: retui.Column,
		},
		retui.NewStyle(),
		c.buildToolbarCashbook(),
		c.buildTableCashBook(journals),
	)
}

func (c *CashBookComponent) buildToolbarCashbook() retui.Element {

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
		retui.Text("Cash Book :", retui.NewStyle().Bold(true)),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Gap:   1,
					Width: retui.Fixed(10),
				},
				retui.NewStyle(),
				retui.Text("Narration", retui.NewStyle()),
			),
			retui.Box(
				retui.Props{
					Gap: 1,
				},
				retui.NewStyle(),
				retui.Text("Text", retui.NewStyle()),
			),
		),
	)
}

func (c *CashBookComponent) buildTableCashBook(journals []*ent.Journal) retui.Element {
	rows := make([][]string, len(journals))

	for i, j := range journals {
		var partyName string
		var byName string

		for _, line := range j.Edges.Lines {
			ledger, err := line.Edges.LedgerOrErr()
			if err != nil {
				continue
			}

			name := strings.TrimSpace(ledger.Name)
			code := strings.ToUpper(strings.TrimSpace(ledger.Code))

			if code == "CASH" || code == "BANK" ||
				strings.EqualFold(name, "Cash In Hand") ||
				strings.EqualFold(name, "Bank") {
				byName = name
			} else {
				partyName = name
			}
		}

		rows[i] = []string{
			j.VoucherDate.Format("02/01/2006"),
			j.VoucherNo,
			util.Deref(j.ReferenceNo),
			j.VoucherType,
			partyName,
			byName,
			fmt.Sprintf("%.2f", j.TotalDebit),
			util.Deref(j.Narration),
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
			"Party",
			"By",
			"Amount",
			"Narration",
			"Status",
		}).
		Alignments([]string{
			"left",
			"left",
			"left",
			"left",
			"left",
			"left",
			"right",
			"left",
			"center",
		}).
		Focused(true).
		Rows(rows).
		SelectedIndex(0).
		ColumnWidths([]int{
			15,
			10,
			15,
			10,
			30,
			30,
			15,
			30,
			10,
		}).
		OnChange(func(i int) {
			if i < 0 || i >= len(journals) {
				return
			}

			// setSelected(journals[i])

			if retui.CurrentKey.Code == retui.KeyEnter {
				// c.controller.ShowJournal(journals[i].ID)
			}
		}).
		Render()

	return retui.Box(
		retui.Props{
			Height: retui.Fixed(33),
		},
		retui.NewStyle(),
		tbl,
	)
}
