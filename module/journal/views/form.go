package views

import (
	"fmt"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/widgets"
)

type JournalCreateComponent struct {
	controller *journal.Controller
	ctx        *appctx.AppContext

	// editing/editID mirror the same pattern used in receipt.FormComponent:
	// form state itself never carries "am I editing" — that's UI-session
	// state, not voucher data.
	editing bool
	editID  int
}

const (
	ledgerWidth  = 30
	debitWidth   = 10
	creditWidth  = 10
	remarksWidth = 50

	journalFieldsPerLine = 4 // Ledger, Debit, Credit, Remarks
	journalHeaderFields  = 4 // VcNo, VcDate, VcReference, VcNarration
)

func NewJournalCreateWindow(ctx *appctx.AppContext) *JournalCreateComponent {
	return &JournalCreateComponent{
		controller: journal.NewController(ctx),
		ctx:        ctx,
	}
}

// LoadForEdit switches the form into update mode for an existing
// journal. Not wired to a caller yet — hook this up from wherever a
// journal list/detail view lets the user choose "edit".
func (c *JournalCreateComponent) LoadForEdit(id int, state journal.FormState) {
	c.editing = true
	c.editID = id
	// NOTE: pushing `state` into the live retui.UseForm hook depends on
	// how that hook reseeds on re-render — confirm it does before relying
	// on this to populate the form with existing values.
}

func (c *JournalCreateComponent) bindKeys(form *retui.Form[journal.FormState]) {
	key := retui.CurrentKey
	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	v := form.Values()
	totalFields := journalHeaderFields + len(v.Lines)*journalFieldsPerLine

	moveFocus := func(delta int) {
		if totalFields == 0 {
			return
		}
		v.FocusIndex = (v.FocusIndex + delta + totalFields) % totalFields
		form.SetValuesSilent(v) // focus movement shouldn't mark the form dirty
	}

	switch key.Code {
	case retui.KeyDown, retui.KeyTab:
		moveFocus(+1)
	case retui.KeyUp, retui.KeyShiftTab:
		moveFocus(-1)
	case retui.KeyEscape:
		retui.PopScreen()

	case retui.KeyF4:
		c.addLine(form)

	case retui.KeyF10:
		c.save(form)

	default:
		return
	}

	retui.CurrentKey.Consumed = true
}

// addLine appends a blank line and moves focus to its first field.
// Shared by the F4 shortcut and "Enter" in the last remarks field.
func (c *JournalCreateComponent) addLine(form *retui.Form[journal.FormState]) {
	v := form.Values()
	v.Lines = append(v.Lines, journal.JournalLine{})
	v.FocusIndex = journalHeaderFields + (len(v.Lines)-1)*journalFieldsPerLine
	form.SetValues(v)
}

// save persists the current form values and resets the form on success.
func (c *JournalCreateComponent) save(form *retui.Form[journal.FormState]) {
	v := form.Values()

	mode := journal.ModeCreate
	id := 0
	if c.editing {
		mode = journal.ModeUpdate
		id = c.editID
	}

	jrnl, err := c.controller.Save(mode, id, v)
	if err != nil {
		components.ShowError(err.Error())
		return
	}

	components.ShowSuccess(fmt.Sprintf("Journal %s saved.", jrnl.VoucherNo))

	// Reset back to a fresh two-line form. If retui.Form.Reset() already
	// reseeds from the initial UseForm(...) value (which also starts with
	// two blank lines), this SetValues call is redundant — but it's cheap
	// insurance against Reset() zeroing Lines to nil instead.
	form.Reset()
	fresh := form.Values()
	fresh.Lines = []journal.JournalLine{{}, {}}
	form.SetValues(fresh)
}

func (c *JournalCreateComponent) JournalCreateForm(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(journal.FormState{
		Lines: []journal.JournalLine{{}, {}},
	})

	panel := components.Panel().
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
			c.headerSection(form),
		).
		DividerWithText("Journal Entries").
		Children(
			c.lineItemRows(form)...,
		).
		DividerWithText("Total").
		Children(
			c.footerSection(form),
		).
		Render()

	return retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		panel,
	)
}

func (c *JournalCreateComponent) calculateTotals(lines []journal.JournalLine) (totalDebit, totalCredit float64) {
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	return
}

func (c *JournalCreateComponent) headerSection(form *retui.Form[journal.FormState]) retui.Element {
	v := form.Values()
	c.bindKeys(form)

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
					ID("vcNo").
					Focused(v.FocusIndex == 0).
					Width(15).
					Value(v.VcNo).
					OnChange(func(id string, value string) {
						c.setField(form, "VcNo", value)
					}).
					Render(),
			),
		),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Voucher Date:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.DateInput().
					ID("vcDate").
					Width(12).
					Focused(v.FocusIndex == 1).
					Value(v.VcDate).
					Format("DD/MM/YYYY").
					OnChange(func(id, value string) {
						c.setField(form, "VcDate", value)
					}).
					Render(),
			),
		),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Reference:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					ID("vcReference").
					Focused(v.FocusIndex == 2).
					Value(v.VcReference).
					Width(20).
					Placeholder("Enter Reference").
					OnChange(func(id string, value string) {
						c.setField(form, "VcReference", value)
					}).
					Render(),
			),
		),

		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(retui.Props{Width: retui.Fit()}, retui.NewStyle(), retui.Text("Narration:", retui.NewStyle())),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					ID("vcNarration").
					Focused(v.FocusIndex == 3).
					Width(65).
					Value(v.VcNarration).
					Placeholder("Narration").
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) {
						c.setField(form, "VcNarration", value)
					}).
					Render(),
			),
		),
	)
}

// setField wraps form.SetField with consistent error logging so every
// OnChange callback isn't repeating the same three lines.
func (c *JournalCreateComponent) setField(form *retui.Form[journal.FormState], field string, value any) {
	if err := form.SetField(field, value); err != nil {
		retui.Debugf("SetField(%s) error: %v", field, err)
	}
}

func (c *JournalCreateComponent) footerSection(form *retui.Form[journal.FormState]) retui.Element {
	totalDebit, totalCredit := c.calculateTotals(form.Values().Lines)
	balanced := journal.Balanced(totalDebit, totalCredit) // was: totalDebit == totalCredit (unsafe float compare)

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

				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text("TOTAL :", retui.NewStyle().Bold(true).Foreground(retui.Blue)),
				),
				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text(fmt.Sprintf("%.2f", totalDebit), retui.NewStyle().Bold(true).Foreground(retui.Green)),
				),
				retui.Box(
					retui.Props{Width: retui.Fit()},
					retui.NewStyle(),
					retui.Text(fmt.Sprintf("%.2f", totalCredit), retui.NewStyle().Bold(true).Foreground(retui.Green)),
				),
			),
		),

		retui.Text(
			fmt.Sprintf("Balanced :%t", balanced),
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

func (c *JournalCreateComponent) lineItemRows(form *retui.Form[journal.FormState]) []retui.Element {
	v := form.Values()

	rows := []retui.Element{c.lineHeader()}
	for i := range v.Lines {
		rows = append(rows, c.lineRow(form, i))
	}
	return rows
}

func headerCell(label string, width int) retui.Element {
	width = retui.CurrentScreenWidth * width / 100
	return retui.Box(
		retui.Props{Width: retui.Fixed(width)},
		retui.NewStyle(),
		retui.Text(label, retui.NewStyle().Bold(true)),
	)
}

func (c *JournalCreateComponent) lineHeader() retui.Element {
	return retui.Box(
		retui.Props{
			Direction: retui.Row,
			Gap:       2,
			Padding:   [4]int{0, 1, 0, 1},
		},
		retui.NewStyle(),
		headerCell("Ledger", ledgerWidth),
		headerCell("Debit", debitWidth),
		headerCell("Credit", creditWidth),
		headerCell("Remarks", remarksWidth),
	)
}

func (c *JournalCreateComponent) lineRow(form *retui.Form[journal.FormState], index int) retui.Element {
	v := form.Values()
	line := v.Lines[index]
	base := journalHeaderFields + index*journalFieldsPerLine

	return retui.Box(
		retui.Props{
			Direction: retui.Row,
			Gap:       2,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
		},
		retui.NewStyle(),
		c.ledgerField(form, index, base, line, ledgerWidth),
		c.debitField(form, index, base, line, debitWidth),
		c.creditField(form, index, base, line, creditWidth),
		c.remarksField(form, index, base, line, remarksWidth),
	)
}

func (c *JournalCreateComponent) ledgerField(
	form *retui.Form[journal.FormState],
	index, focus int,
	line journal.JournalLine,
	width int,
) retui.Element {
	wid := retui.CurrentScreenWidth * width / 100

	// LedgerComponent works with the string form of a ledger ID (same
	// pattern as receipt.PartyLine) — never a "code" like "CASH"/"BANK".
	return widgets.LedgerComponent(
		c.ctx,
		fmt.Sprintf("ledger_%d", index),
		util.IntToString(line.LedgerID),
		wid,
		form.Values().FocusIndex == focus,
		func(id, value string) {
			c.updateLine(form, index, func(l *journal.JournalLine) {
				l.LedgerID = util.StringToInt(value, 0)
			})
		},
	)
}

func (c *JournalCreateComponent) debitField(
	form *retui.Form[journal.FormState],
	index, focus int,
	line journal.JournalLine,
	width int,
) retui.Element {
	wid := retui.CurrentScreenWidth * width / 100
	return components.NumberInput().
		ID(fmt.Sprintf("debit_%d", index)).
		Width(wid).
		Decimals(2).
		Focused(form.Values().FocusIndex == focus+1).
		Value(line.Debit).
		OnChange(func(id string, value float64) {
			c.updateLine(form, index, func(l *journal.JournalLine) {
				l.Debit = value
			})
		}).
		Render()
}

func (c *JournalCreateComponent) creditField(
	form *retui.Form[journal.FormState],
	index, focus int,
	line journal.JournalLine,
	width int,
) retui.Element {
	wid := retui.CurrentScreenWidth * width / 100
	return components.NumberInput().
		ID(fmt.Sprintf("credit_%d", index)).
		Width(wid).
		Decimals(2).
		Focused(form.Values().FocusIndex == focus+2).
		Value(line.Credit).
		OnChange(func(id string, value float64) {
			c.updateLine(form, index, func(l *journal.JournalLine) {
				l.Credit = value
			})
		}).
		Render()
}

func (c *JournalCreateComponent) remarksField(
	form *retui.Form[journal.FormState],
	index, focus int,
	line journal.JournalLine,
	width int,
) retui.Element {
	wid := retui.CurrentScreenWidth * width / 100
	return components.TextInput().
		ID(fmt.Sprintf("remarks_%d", index)).
		Width(wid).
		Focused(form.Values().FocusIndex == focus+3).
		Value(line.Remarks).
		OnChange(func(id, value string) {
			c.updateLine(form, index, func(l *journal.JournalLine) {
				l.Remarks = value
			})
		}).
		OnKeyPress(func(s string, k retui.Key) bool {
			if k.Code == retui.KeyEnter {
				c.addLine(form)
				return true
			}
			return false
		}).
		Render()
}

func (c *JournalCreateComponent) updateLine(
	form *retui.Form[journal.FormState],
	index int,
	update func(*journal.JournalLine),
) {
	v := form.Values()
	if index < 0 || index >= len(v.Lines) {
		return
	}
	update(&v.Lines[index])
	form.SetValues(v)
}
