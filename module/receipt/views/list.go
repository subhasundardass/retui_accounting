package views

import (
	"fmt"
	"strings"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/receipt"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

// Component holds the "companies" screen's controller and cached render
// state. One instance is created at registration time and its bound
// methods (e.g. List) are passed as ui.Screen.Render — this avoids
// recreating the controller and re-querying the DB on every frame.
type Component struct {
	controller *receipt.Controller
	ctx        *appctx.AppContext
}

func NewComponent(ctx *appctx.AppContext) *Component {
	return &Component{
		controller: receipt.NewController(ctx),
		ctx:        ctx,
	}
}

func (c *Component) bindKeys() {
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

func (c *Component) ReceiptBook(ctx *appctx.AppContext) retui.Element {

	journals, setJournals := retui.UseState([]*ent.Journal{})
	selected, setSelected := retui.UseState(&ent.Journal{})

	retui.UseEffect(func() func() {
		list, err := c.controller.List(0, 40)
		if err != nil {
			retui.Errorf("Error fetching data %s", err.Error())
			return nil
		}

		setJournals(list)
		return nil
	}, []any{})

	c.bindKeys()
	return retui.Box(
		retui.Props{
			Direction: retui.Column,
		},
		retui.NewStyle(),
		c.buildToolbar(selected),
		c.buildTable(journals, setSelected),
	)
}

func (c *Component) buildToolbar(selected *ent.Journal) retui.Element {
	title := "Receipt Book  "
	if selected != nil {
		title = fmt.Sprintf("Receipt Book  %d", selected.ID)
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

func (c *Component) buildTable(
	journals []*ent.Journal,
	setSelected func(*ent.Journal),
) retui.Element {
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

			setSelected(journals[i])

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
