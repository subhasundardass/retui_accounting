package journal_entry

import (
	"fmt"

	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

type Components struct {
	controller *Controller
}

// JournalLine represents a single debit/credit line in the journal voucher.
type JournalLine struct {
	LedgerCode string
	Debit      float64
	Credit     float64
	Remarks    string
}

func NewComponents(controller *Controller) *Components {
	return &Components{controller: controller}
}

// newEmptyJournalLines returns n blank journal lines, used to seed initial state.
func newEmptyJournalLines(n int) []JournalLine {
	lines := make([]JournalLine, n)
	return lines
}

// updateLine returns a copy of lines with the field(s) at idx mutated.
func updateLine(lines []JournalLine, idx int, mutate func(*JournalLine)) []JournalLine {
	updated := make([]JournalLine, len(lines))
	copy(updated, lines)
	l := updated[idx]
	mutate(&l)
	updated[idx] = l
	return updated
}

// removeLine returns a copy of lines with the line at idx removed.
func removeLine(lines []JournalLine, idx int) []JournalLine {
	updated := make([]JournalLine, 0, len(lines)-1)
	updated = append(updated, lines[:idx]...)
	updated = append(updated, lines[idx+1:]...)
	return updated
}

// appendLine returns a copy of lines with a new blank line appended.
func appendLine(lines []JournalLine) []JournalLine {
	updated := make([]JournalLine, len(lines), len(lines)+1)
	copy(updated, lines)
	return append(updated, JournalLine{})
}

func (c *Components) RenderScreen() retui.Element {

	// State management for header fields
	vNo, setVNo := retui.UseState("")
	vcDate, setVcDate := retui.UseState("")
	vReference, setVReference := retui.UseState("")
	vNarration, setVNarration := retui.UseState("")

	// State management for journal line items (starts with 5 blank rows)
	lines, setLines := retui.UseState(newEmptyJournalLines(2))

	// 4 header fields, then 5 fields per journal line (ledger, debit, credit, remarks, delete)
	totalFields := 4 + len(lines)*5
	focusIndex, setFocusIndex := retui.UseState(0)

	if retui.IsFocused("journal_entry") {
		switch retui.CurrentKey.Code {
		case retui.KeyEscape:
			retui.SetFocus("sidebar")
		case retui.KeyDown:
			setFocusIndex((focusIndex + 1) % totalFields)
		case retui.KeyUp:
			setFocusIndex((focusIndex - 1 + totalFields) % totalFields)

		case retui.KeyF10:
			// Save functionality
			for _, line := range lines {
				retui.Infof("%+v\n", line)
			}
			entry := JournalEntry{
				VoucherNo:   vNo,
				VoucherDate: vcDate,
				Reference:   vReference,
				Narration:   vNarration,
				Lines:       lines,
			}

			jrnl, err := c.controller.SaveJournal(entry)
			if err != nil {
				retui.Error(err)
				window.AlertError("Error!", retui.Text(err.Error(), retui.NewStyle()))
			} else {
				window.AlertError(
					"Success",
					retui.Text(
						fmt.Sprintf("Journal %s saved successfully.", jrnl.VoucherNo),
						retui.NewStyle(),
					),
				)

				// Reset header
				setVNo("")
				setVcDate("")
				setVReference("")
				setVNarration("")

				// Reset lines
				setLines(newEmptyJournalLines(2))

				// Reset focus
				setFocusIndex(0)
			}

		case retui.KeyF4:
			// Add a new blank journal line, move focus to its first field
			newIndex := 4 + len(lines)*5
			setLines(appendLine(lines))
			setFocusIndex(newIndex)
		case retui.KeyEnter:
			// If focus is on a line's [Delete] cell, remove that line
			if focusIndex >= 4 {
				relative := focusIndex - 4
				row := relative / 5
				col := relative % 5
				if col == 4 && row < len(lines) {
					// Calculate new length before removal
					newLen := len(lines) - 1
					setLines(removeLine(lines, row))

					// Calculate new total fields after removal
					newTotal := 4 + newLen*5

					// Adjust focus if needed
					newFocus := focusIndex
					if focusIndex >= newTotal {
						newFocus = newTotal - 1
					}
					if newFocus < 0 {
						newFocus = 0
					}
					setFocusIndex(newFocus)
				}
			}

		}
	}

	isFocused := func(idx int) bool { return focusIndex == idx }

	//======

	panel := components.Panel().
		FixedWidth(180). // Use FixedWidth
		Header(retui.Box(
			retui.Props{
				Direction: retui.Row,
				Padding:   [4]int{0, 1, 0, 1},
				Width:     retui.Grow(1),
				Height:    retui.Fit(),
				Justify:   retui.JustifySpaceBetween,
			},
			retui.NewStyle(),
			retui.Text("Journal Voucher", retui.NewStyle().Bold(true)),
			retui.Text("F10: Save   F4: Add Line   Enter: Delete Line", retui.NewStyle().Bold(true)),
		)).
		Children(
			c.headerSection(
				vNo, setVNo,
				vcDate, setVcDate,
				vReference, setVReference,
				vNarration, setVNarration,
				isFocused),
		).
		DividerWithText("Journal Entries").
		Children(
			c.lineItemRows(lines, setLines, isFocused)...,
		).
		DividerWithText("Total").
		// ContentGap(2).
		Children(
			c.footer(lines),
		).
		Render()

	return retui.Box(
		retui.Props{},
		retui.NewStyle(),

		panel,
	)
}

func (c *Components) headerSection(
	vNo string, setVNo func(string),
	vcDate string, setVcDate func(string),
	vReference string, setVReference func(string),
	vNarration string, setVNarration func(string),
	isFocused func(int) bool,
) retui.Element {
	return retui.Box(
		retui.Props{
			// Direction: retui.Column,
			Width:   retui.Grow(1),
			Gap:     1,
			Padding: [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),

		// Voucher No
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Voucher No:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					Width(20).
					ID("vNo").
					Focused(isFocused(0)).
					Value(vNo).
					Placeholder("Enter Voucher No").
					Style(retui.NewStyle().Bold(true).Background(retui.Hex("#ffffff"))).
					OnChange(func(id string, value string) { setVNo(value) }).
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
					Focused(isFocused(1)).
					Value(vcDate).
					Format("DD/MM/YYYY").
					OnChange(func(id, value string) {
						setVcDate(value)
					}).
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
					Focused(isFocused(2)).
					Value(vReference).
					Width(20).
					Placeholder("Enter Reference").
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) { setVReference(value) }).
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
					Focused(isFocused(3)).
					Width(68).
					Value(vNarration).
					Placeholder("Enter Narration").
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) { setVNarration(value) }).
					Render(),
			),
		),
	)
}

// headerCellLabel pads label to width and wraps it in brackets so the
// header row visually lines up with the "[ value ]" border TextInput
// draws for each data cell beneath it.
func headerCellLabel(label string, width int) string {
	if len(label) > width {
		label = label[:width]
	}
	return fmt.Sprintf("[ %-*s ]", width, label)
}

// lineItemRows renders the journal-entry table rows:
func (c *Components) lineItemRows(
	lines []JournalLine,
	setLines func([]JournalLine),
	isFocused func(int) bool,
) []retui.Element {
	rows := make([]retui.Element, 0, len(lines)+1)

	// Column widths must match the TextInput widths used in the data rows
	// below (ledgerField/debitField/creditField/remarksField), so the
	// bracket borders line up.
	const (
		ledgerIDColWidth = 0
		ledgerColWidth   = 60
		debitColWidth    = 15
		creditColWidth   = 15
		remarksColWidth  = 52
		actionColWidth   = 13
	)

	// Column header row — manually bracketed to match the "[ value ]"
	// border TextInput renders for each data cell below, since Box/Style
	// don't have a confirmed Border() API to draw a real border here.
	rows = append(rows, retui.Box(
		retui.Props{Direction: retui.Row, Height: retui.Fit(), Gap: 2, Padding: [4]int{0, 1, 0, 1}},
		retui.NewStyle(),
		retui.Text(headerCellLabel("Ledger Name", ledgerColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Debit", debitColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Credit", creditColWidth), retui.NewStyle().Bold(true)),
		retui.Text(headerCellLabel("Remarks", remarksColWidth), retui.NewStyle().Bold(true)),
		retui.Text("[ Action ]", retui.NewStyle().Bold(true)),
	))

	for i := range lines {
		i := i // capture loop var for closures
		line := lines[i]
		baseIdx := 4 + i*5

		ledgerField := components.SelectDropdown().
			ID(fmt.Sprintf("lineLedger_%d", i)).
			Width(ledgerColWidth+4).
			Style(retui.NewStyle().Bold(true).Background(retui.Red)).
			OverlayAbsPos(40, 5).
			Height(10).
			Options(c.controller.LedgerSeedOptions(line.LedgerCode)). // ← MUST be non-empty
			OnFilter(func(id, query string) []components.SelectOption {
				return c.controller.LedgerFilterOptions(query)
			}).
			Value(line.LedgerCode).
			Focused(isFocused(baseIdx + 0)).
			OnChange(func(id, value string) {
				setLines(updateLine(lines, i, func(l *JournalLine) {
					l.LedgerCode = value
				}))
			}).
			Render()

		debitField := components.NumberInput().
			ID(fmt.Sprintf("lineDebit_%d", i)).
			Width(debitColWidth + 4).
			Focused(isFocused(baseIdx + 1)).
			Decimals(2).
			Value(line.Debit).
			// SelectAllOnFocus(true).
			Style(retui.NewStyle().Bold(true)).
			OnChange(func(id string, value float64) {
				setLines(updateLine(lines, i, func(l *JournalLine) { l.Debit = value }))
			}).
			Render()

		creditField := components.NumberInput().
			ID(fmt.Sprintf("lineCredit_%d", i)).
			Width(creditColWidth + 4).
			Focused(isFocused(baseIdx + 2)).
			Value(line.Credit).
			Decimals(2).
			Style(retui.NewStyle().Bold(true)).
			OnChange(func(id string, value float64) {
				setLines(updateLine(lines, i, func(l *JournalLine) { l.Credit = value }))
			}).
			Render()

		remarksField := components.TextInput().
			ID(fmt.Sprintf("lineRemarks_%d", i)).
			Width(remarksColWidth + 4).
			Focused(isFocused(baseIdx + 3)).
			Value(line.Remarks).
			// Placeholder("Remarks").
			Style(retui.NewStyle().Bold(true)).
			OnChange(func(id string, value string) {
				setLines(updateLine(lines, i, func(l *JournalLine) { l.Remarks = value }))
			}).
			Render()

		deleteButton := components.Button().
			ID("submit").
			Focused(isFocused(baseIdx + 4)).
			Width(actionColWidth).
			Label("Delete").
			Style(retui.NewStyle()).
			ActiveStyle(retui.NewStyle().Foreground(retui.White).Background(retui.Red).Bold(true)).
			Render()

		rows = append(rows, retui.Box(
			retui.Props{Direction: retui.Row, Height: retui.Fit(), Gap: 2, Padding: [4]int{0, 1, 0, 1}},
			retui.NewStyle(),
			retui.Box(retui.Props{}, retui.NewStyle(), ledgerField),
			retui.Box(retui.Props{}, retui.NewStyle(), debitField),
			retui.Box(retui.Props{}, retui.NewStyle(), creditField),
			retui.Box(retui.Props{}, retui.NewStyle(), remarksField),
			retui.Box(retui.Props{}, retui.NewStyle(), deleteButton),
		))

	}

	return rows
}

// Add this method to calculate totals
func (c *Components) calculateTotals(lines []JournalLine) (totalDebit, totalCredit float64) {
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	return
}

func (c *Components) footer(lines []JournalLine) retui.Element {
	totalDebit, totalCredit := c.calculateTotals(lines)
	balanced := totalDebit == totalCredit // Calculate balance properly

	return retui.Box(
		retui.Props{
			Width:   retui.Grow(1),
			Justify: retui.JustifySpaceBetween,
			Padding: [4]int{0, 1, 0, 1},
		}, retui.NewStyle(),

		retui.Box(
			retui.Props{}, retui.NewStyle(),

			retui.Box(
				retui.Props{
					Direction: retui.Row,
					Height:    retui.Fit(),
					Gap:       2,
				},
				retui.NewStyle(),

				// Ledger column - show "TOTAL"
				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text(
						fmt.Sprintf("%s :", "TOTAL :"),
						retui.NewStyle().Bold(true).Foreground(retui.Blue),
					),
				),

				// Debit column
				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text(
						fmt.Sprintf("%.2f", totalDebit),
						retui.NewStyle().Bold(true).Foreground(retui.Green),
					),
				),

				// Credit column
				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text(
						fmt.Sprintf("%.2f", totalCredit),
						retui.NewStyle().Bold(true).Foreground(retui.Green),
					),
				),
			),
		),

		retui.Text(
			fmt.Sprintf("Balanced :%t", balanced), // Use %t for boolean
			func() retui.Style {
				style := retui.NewStyle().Bold(true)
				if balanced {
					return style.Foreground(retui.Green)
				}
				return style.Foreground(retui.Red)
			}(),
		),
	)
}
