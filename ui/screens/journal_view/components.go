package journal_view

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

func (c *Components) RenderScreen() retui.Element {

	journalID := 0
	params := retui.CurrentScreenParams()
	if params != nil {
		if id, ok := params["journalID"].(int); ok {
			journalID = id
		}
	}

	journal, setJournal := retui.UseState[*ent.Journal](nil)

	retui.UseEffect(func() func() {
		j := c.controller.GetJournal(journalID)
		if j != nil {
			setJournal(j)
		}

		return nil
	}, []any{})

	panel := components.Panel().
		Width(retui.Fixed(180)).
		Header(retui.Box(
			retui.Props{
				Direction: retui.Row,
				Padding:   [4]int{0, 1, 0, 1},
				Width:     retui.Grow(1), Height: retui.Grow(1),
				Justify: retui.JustifySpaceBetween,
			},
			retui.NewStyle(),
			retui.Text("Journal", retui.NewStyle().Bold(true)),
		)).
		Children(
			c.headerSection(journal),
		).
		DividerWithText("Journal Entries").
		Children(

			c.lineItemRows(journal)...,
		).
		Render()

	return retui.Box(
		retui.Props{},
		retui.NewStyle(),

		panel,
	)
}

func (c *Components) headerSection(journal *ent.Journal) retui.Element {

	//--Check first
	if journal == nil {
		retui.Debug("journal is nil")
		return retui.Text("Loading...", retui.NewStyle())
	}

	return retui.Box(
		retui.Props{
			Width:   retui.Grow(1),
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Voucher No:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					Width(20).
					ID("vNo").
					Value(journal.VoucherNo).
					Style(retui.NewStyle().Bold(true)).
					Render(),
			),
		),

		// Voucher Date
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Voucher Date:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.DateInput().
					ID("vcDate").
					Width(12).
					Value(journal.VoucherDate.Format("02/01/2006")).
					Format("DD/MM/YYYY").
					Render(),
			),
		),

		// Reference
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Reference:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					ID("vReference").
					Value(*journal.ReferenceNo).
					Width(20).
					Placeholder("Enter Reference").
					Style(retui.NewStyle().Bold(true)).
					Render(),
			),
		),
		// Narration
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Narration:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					ID("vNarration").
					Width(68).
					Value(*journal.Narration).
					Placeholder("Enter Narration").
					Style(retui.NewStyle().Bold(true)).
					Render(),
			),
		),
	)
}

func headerCellLabel(label string, width int) string {
	if len(label) > width {
		label = label[:width]
	}
	return fmt.Sprintf("[ %-*s ]", width, label)
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

func (c *Components) lineItemRows(journal *ent.Journal) []retui.Element {
	rows := make([]retui.Element, 0)

	if journal == nil {
		return rows
	}

	const (
		ledgerColWidth  = 60
		debitColWidth   = 15
		creditColWidth  = 15
		remarksColWidth = 52
	)

	// Header
	rows = append(rows, retui.Box(
		retui.Props{
			Direction: retui.Row,
			Width:     retui.Grow(1), Height: retui.Grow(1),
			Gap:     2,
			Padding: [4]int{0, 1, 1, 1},
		},
		retui.NewStyle(),

		retui.Text(headerCellLabel("Ledger Name", ledgerColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Debit", debitColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Credit", creditColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Remarks", remarksColWidth), retui.NewStyle().Bold(true)),
	))

	rows = append(rows, retui.Box(
		retui.Props{},
		retui.NewStyle(),
		retui.Text("", retui.NewStyle()),
	))

	for _, line := range journal.Edges.Lines {

		ledgerName := ""
		if line.Edges.Ledger != nil {
			ledgerName = line.Edges.Ledger.Name
		}

		rows = append(rows, retui.Box(
			retui.Props{
				Direction: retui.Row,
				Height:    retui.Fit(),
				Gap:       2,
				Padding:   [4]int{0, 1, 0, 1},
			},
			retui.NewStyle(),

			components.TextInput().
				Width(ledgerColWidth+4).
				Value(ledgerName).
				Render(),

			components.NumberInput().
				Width(debitColWidth+4).
				Decimals(2).
				Value(line.Debit).
				Render(),

			components.NumberInput().
				Width(creditColWidth+4).
				Decimals(2).
				Value(line.Credit).
				Render(),

			components.TextInput().
				Width(remarksColWidth+4).
				Value(derefString(line.Description)).
				Render(),
		))
	}

	return rows
}
