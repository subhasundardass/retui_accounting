package views

import (
	"fmt"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/widgets"
)

type JournalCreateComponent struct {
	controller *journal.Controller
	ctx        *appctx.AppContext
	// form       *FormComponent
}

const (
	ledgerWidth  = 30
	debitWidth   = 10
	creditWidth  = 10
	remarksWidth = 50
)

func NewJournalCreateWindow(ctx *appctx.AppContext) *JournalCreateComponent {
	return &JournalCreateComponent{
		controller: journal.NewController(ctx),
		ctx:        ctx,
	}
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
	totalFields := 4 + len(v.Lines)*4 // confirm: 4 JournalLine fields shown here, is there a 5th (delete action)?

	moveFocus := func(delta int) {
		if totalFields == 0 {
			return
		}

		v.FocusIndex = (v.FocusIndex + delta + totalFields) % totalFields

		// Focus movement shouldn't mark the form dirty.
		form.SetValuesSilent(v)

	}

	switch key.Code {
	case retui.KeyDown, retui.KeyTab:
		moveFocus(+1)
	case retui.KeyUp, retui.KeyShiftTab:
		moveFocus(-1)
	case retui.KeyEscape:
		retui.PopScreen()
	case retui.KeyF2:
		retui.Debugf("F2 Pressed.......")

	case retui.KeyF4:
		v := form.Values()
		v.Lines = append(v.Lines, journal.JournalLine{})

		// Move focus to the Ledger field of the new row
		v.FocusIndex = 4 + (len(v.Lines)-1)*4
		form.SetValues(v)

	case retui.KeyF10:
		// Save functionality
		for _, line := range v.Lines {
			retui.Infof("%+v\n", line)
		}
		entry := journal.FormState{
			VcNo:        v.VcNo,
			VcDate:      v.VcDate,
			VcReference: v.VcReference,
			VcNarration: v.VcNarration,
			Lines:       v.Lines,
		}

		jrnl, err := c.controller.SaveJournal(entry)
		if err != nil {
			components.ShowError(err.Error())
			return
		}

		components.ShowSuccess(fmt.Sprintf("Journal %s saved.", jrnl.VoucherNo))
		//--Reset
		form.Reset()
		v := form.Values() // read AFTER reset
		v.Lines = []journal.JournalLine{{}, {}}
		form.SetValues(v) // push back ✗

	default:
		return
	}

	retui.CurrentKey.Consumed = true
}

func (c *JournalCreateComponent) JournalCreateForm(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(journal.FormState{
		Lines: []journal.JournalLine{{}, {}},
	})

	panel := components.Panel().
		// Width(retui.Percent(80)).
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
			c.footerSecton(form),
		).
		Render()

	return retui.Box(
		retui.Props{
			Gap: 1,
		},
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

		// Voucher No
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
						if err := form.SetField("VcNo", value); err != nil {
							retui.Debugf("SetField error: %v", err) // or however you actually log
						}
					}).
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
					Focused(v.FocusIndex == 1).
					Value(v.VcDate).
					Format("DD/MM/YYYY").
					OnChange(func(id, value string) {
						if err := form.SetField("VcDate", value); err != nil {
							retui.Debugf("SetField error: %v", err) // or however you actually log
						}
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
					ID("vcReference").
					Focused(v.FocusIndex == 2).
					Value(v.VcReference).
					Width(20).
					Placeholder("Enter Reference").
					OnChange(func(id string, value string) {
						if err := form.SetField("VcReference", value); err != nil {
							retui.Debugf("SetField error: %v", err) // or however you actually log
						}
					}).
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
					ID("vcNarration").
					Focused(v.FocusIndex == 3).
					Width(65).
					Value(v.VcNarration).
					Placeholder("Narration").
					Style(retui.NewStyle().Bold(true)).
					OnChange(func(id string, value string) {
						if err := form.SetField("VcNarration", value); err != nil {
							retui.Debugf("SetField error: %v", err) // or however you actually log
						}
					}).
					Render(),
			),
		),
	)
}

func (c *JournalCreateComponent) footerSecton(form *retui.Form[journal.FormState]) retui.Element {
	totalDebit, totalCredit := c.calculateTotals(form.Values().Lines)
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

func (c *JournalCreateComponent) lineItemRows(
	form *retui.Form[journal.FormState],
) []retui.Element {

	v := form.Values()

	rows := []retui.Element{
		c.lineHeader(),
	}

	for i := range v.Lines {
		rows = append(rows, c.lineRow(form, i))
	}

	return rows
}

func headerCell(label string, width int) retui.Element {

	width = retui.CurrentScreenWidth * width / 100

	return retui.Box(
		retui.Props{
			Width: retui.Fixed(width),
		},
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

func (c *JournalCreateComponent) lineRow(
	form *retui.Form[journal.FormState],
	index int,
) retui.Element {

	v := form.Values()

	line := v.Lines[index]
	base := 4 + index*4

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
		// c.deleteButton(form, index, base),
	)
}

func (c *JournalCreateComponent) ledgerField(
	form *retui.Form[journal.FormState],
	index, focus int,
	line journal.JournalLine,
	width int,
) retui.Element {

	wid := retui.CurrentScreenWidth * width / 100

	return widgets.LedgerComponent(
		c.ctx,
		fmt.Sprintf("ledger_%d", index),
		line.LedgerCode,
		wid,
		form.Values().FocusIndex == focus,
		func(id, value string) {
			c.updateLine(form, index, func(line *journal.JournalLine) {
				line.LedgerCode = value
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

// func (c *JournalCreateComponent) deleteButton(
// 	form *retui.Form[journal.FormState],
// 	index, focus int,
// ) retui.Element {

// 	return components.Button().
// 		ID(fmt.Sprintf("delete_%d", index)).
// 		Label("Delete").
// 		Width(13).
// 		Focused(form.Values().FocusIndex == focus+4).
// 		Style(retui.NewStyle()).
// 		ActiveStyle(
// 			retui.NewStyle().
// 				Foreground(retui.White).
// 				Background(retui.Red).
// 				Bold(true),
// 		).
// 		Render()
// }
