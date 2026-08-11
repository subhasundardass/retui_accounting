package views

import (
	"fmt"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/receipt"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/widgets"
)

const lw = 15
const totalFields = 6   // SlNo, Reference, Date, Amount, RcptAccount, Narration
const fieldsPerLine = 3 // Ledger, Remarks, Amount

type FormComponent struct {
	controller *receipt.Controller
	state      receipt.FormState

	editing bool
	editID  int

	// onSaved is called after a successful save so the caller (e.g. the
	// companies list) can refresh its data.
	onClose func()
}

func NewFormComponent(ctx *appctx.AppContext) *FormComponent {
	return &FormComponent{
		controller: receipt.NewController(ctx),

		state: receipt.FormState{
			Lines: []receipt.PartyLine{{}},
		},
	}
}

func bindKeys(form *retui.Form[receipt.FormState]) {
	key := retui.CurrentKey
	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	v := form.Values()
	total := totalFields + len(v.Lines)*fieldsPerLine

	moveFocus := func(delta int) {
		if total == 0 {
			return
		}
		v.FocusIndex = (v.FocusIndex + delta + total) % total
		form.SetValuesSilent(v)
	}

	switch key.Code {
	case retui.KeyDown, retui.KeyTab:
		moveFocus(+1)
	case retui.KeyUp, retui.KeyShiftTab:
		moveFocus(-1)
	case retui.KeyEscape:
		retui.PopScreen()

	case retui.KeyF10:
		// Save functionality
		// v := form.Values()

		// entry := receipt.FormState{
		// 	SlNo:        v.SlNo,
		// 	Reference:   v.Reference,
		// 	Date:        v.Date,
		// 	Amount:      v.Amount,
		// 	RcptAccount: v.RcptAccount,
		// 	Narration:   v.Narration,
		// 	Lines:       v.Lines,
		// }

		// jrnl, err := form.Controller.SaveJournal(entry) // NOTE: form has no Controller field in the
		// // original snippet; wire this to
		// // FormComponent.controller instead — see
		// // note below the code block.
		// if err != nil {
		// 	components.ShowError(err.Error())
		// 	return
		// }

		// components.ShowSuccess(fmt.Sprintf("Journal %s saved.", jrnl.VoucherNo))

		// --Reset
		form.Reset()
		nv := form.Values() // read AFTER reset
		nv.Lines = []receipt.PartyLine{{}, {}}
		nv.FocusIndex = 0
		form.SetValues(nv)

	default:
		return
	}

	retui.CurrentKey.Consumed = true
}

func (c *FormComponent) Receipt(ctx *appctx.AppContext) retui.Element {

	return retui.Box(
		retui.Props{
			Gap:   1,
			Width: retui.Grow(1),
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Width: retui.Grow(5),
			},
			retui.NewStyle(),
			c.new(ctx),
		),
		retui.Box(
			retui.Props{
				Width: retui.Grow(1),
			},
			retui.NewStyle(),
			components.Panel().
				Header(retui.Text(" Shortcut", retui.NewStyle().Bold(true))).
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Company
						widgets.ShortCutItem("Select Company", "F1"),
						widgets.ShortCutItem("Change Data", "F2"),
						widgets.ShortCutItem("Contra Voucher", "F4"),
						// Data
						widgets.ShortCutItem("Change Data", "F5"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Vouchers
						widgets.ShortCutItem("F4-Contra Voucher", "F3"),
						widgets.ShortCutItem("F5-Payment Voucher", "F3"),
						widgets.ShortCutItem("F6-Receipt Voucher", "F3"),
						widgets.ShortCutItem("F7-Journal Voucher", "F3"),

						// Sales / Purchase
						widgets.ShortCutItem("F8-Sale Voucher", "F4"),
						widgets.ShortCutItem("F9-Purchase Voucher", "F4"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Features / Configuration
						widgets.ShortCutItem("F11-Company Features", "F5"),
						widgets.ShortCutItem("F12-Configuration", "F5"),
					),
				).
				Divider().
				Children(
					retui.Box(
						retui.Props{
							Gap:       0,
							Direction: retui.Column,
							Width:     retui.Grow(1),
							Padding:   [4]int{0, 1, 0, 1},
						},
						retui.NewStyle(),

						// Features / Configuration
						widgets.ShortCutItem("F11-Company Features", "F5"),
					),
				).
				Render(),
		),
	)
}

func (c *FormComponent) new(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(c.state) // single call, once per render
	bindKeys(form)                 // single call too — see note below

	return components.Panel().
		Header(retui.Box(
			retui.Props{
				Direction: retui.Row,
				Padding:   [4]int{0, 1, 0, 1},
				Width:     retui.Grow(1),
				Height:    retui.Fit(),
				Justify:   retui.JustifySpaceBetween,
			},
			retui.NewStyle(),
			retui.Text("Receipt Entry", retui.NewStyle().Bold(true)),
			retui.Text("F10: Save   F5: Reset", retui.NewStyle().Bold(true)),
		)).
		Children(
			retui.Box(
				retui.Props{
					Gap:     1,
					Padding: [4]int{0, 1, 0, 1},
				},
				retui.NewStyle(),
				c.buildHead(form),
			),
		).
		DividerWithText("Particulars").
		Children(
			retui.Box(
				retui.Props{
					Gap:     2,
					Padding: [4]int{0, 1, 0, 1},
					Width:   retui.Grow(1),
					Margin:  [4]int{0, 0, 0, 0},
				},
				retui.NewStyle(),
				retui.Box(
					retui.Props{
						Gap:   4,
						Width: retui.Grow(3),
					},
					retui.NewStyle(),
					retui.Text("Party", retui.NewStyle()),
				),
				retui.Box(
					retui.Props{
						Gap:   4,
						Width: retui.Grow(3),
					},
					retui.NewStyle(),
					retui.Text("Remarks", retui.NewStyle()),
				),
				retui.Box(
					retui.Props{
						Gap:     1,
						Width:   retui.Grow(1),
						Justify: retui.JustifyEnd,
					},
					retui.NewStyle(),
					retui.Text("Amount", retui.NewStyle()),
				),
			),
		).
		Divider().
		Children(
			retui.Box(
				retui.Props{
					Gap:     1,
					Padding: [4]int{0, 1, 0, 1},
					Width:   retui.Grow(1),
				},
				retui.NewStyle(),
				c.buildPartyRow(ctx, form),
			),
		).
		Render()
}

// ---- Header Section
func (c *FormComponent) buildHead(form *retui.Form[receipt.FormState]) retui.Element {

	// form := retui.UseForm(c.state)
	v := form.Values()

	bindKeys(form)

	slNo := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Receipt No", retui.NewStyle()),
		),
		components.TextInput().
			ID("code").
			Value(util.IntToString(v.SlNo)).
			Prefix(" : ").
			Focused(v.FocusIndex == 0).
			OnChange(func(id, value string) {
				if err := form.SetField("SlNo", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)
	ref := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Reference", retui.NewStyle()),
		),
		components.TextInput().
			Prefix(" : ").
			ID("ref").
			Value(v.Reference).
			Focused(v.FocusIndex == 1).
			OnChange(func(id, value string) {
				if err := form.SetField("Reference", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)
	//  Date
	date := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Date", retui.NewStyle()),
		),
		components.DateInput().
			Prefix(" : ").
			ID("date").
			Focused(v.FocusIndex == 2).
			Value(v.Date).
			Format("DD/MM/YYYY").
			OnChange(func(id, value string) {
				if err := form.SetField("Date", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	// Amount
	amount := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Amount", retui.NewStyle()),
		),
		components.NumberInput().
			ID("amount").
			Prefix(" : ").
			Decimals(2).
			Focused(v.FocusIndex == 3).
			Value(float64(v.Amount)).
			OnChange(func(id string, value float64) {
				if err := form.SetField("Amount", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)
	// Receipt Account
	rcptAccount := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Receipt By", retui.NewStyle()),
		),
		components.NumberInput().
			ID("rcptAccount").
			Prefix(" : ").
			Decimals(2).
			Focused(v.FocusIndex == 4).
			Value(float64(v.RcptAccount)).
			OnChange(func(id string, value float64) {
				if err := form.SetField("RcptAccount", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	// Narration
	narration := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(lw),
			},
			retui.NewStyle(),
			retui.Text("Narration", retui.NewStyle()),
		),
		components.TextInput().
			Prefix(" : ").
			ID("narration").
			Value(v.Narration).
			Focused(v.FocusIndex == 5).
			OnChange(func(id, value string) {
				if err := form.SetField("Narration", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap: 2,
			},
			retui.NewStyle(),
			slNo,
			ref,
			retui.Spacer(),
			date,
		),
		retui.Box(
			retui.Props{
				Gap: 2,
			},
			retui.NewStyle(),
			amount,
			rcptAccount,
		),
		retui.Box(
			retui.Props{
				Gap: 2,
			},
			retui.NewStyle(),
			narration,
		),
	)
}

// -- Party Section
func (c *FormComponent) buildPartyRow(ctx *appctx.AppContext, form *retui.Form[receipt.FormState]) retui.Element {

	bindKeys(form)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Height:    retui.Fixed(20),
		},
		retui.NewStyle(),

		//--Append Rows
		c.lineItemRows(ctx, form)...,
	)
}

func (c *FormComponent) lineItemRows(ctx *appctx.AppContext, form *retui.Form[receipt.FormState]) []retui.Element {

	v := form.Values()
	rows := []retui.Element{}

	for i := range v.Lines {
		rows = append(rows, c.lineRow(ctx, form, i))
	}

	return rows
}

// - row
func (c *FormComponent) lineRow(ctx *appctx.AppContext, form *retui.Form[receipt.FormState], i int) retui.Element {

	v := form.Values()
	bindKeys(form)

	base := totalFields + i*fieldsPerLine // focus index of this row's first field

	setLine := func(mutate func(*receipt.PartyLine)) {
		nv := form.Values()
		line := nv.Lines[i]
		mutate(&line)
		nv.Lines[i] = line
		form.SetValues(nv)
	}

	partyLedger := retui.Box(
		retui.Props{},
		retui.NewStyle(),

		widgets.LedgerComponent(
			ctx,
			fmt.Sprintf("partyLedger-%d", i),
			util.IntToString(v.Lines[i].Ledger),
			70,
			v.FocusIndex == base,
			func(id, value string) {
				setLine(func(l *receipt.PartyLine) { l.Ledger = util.StringToInt(value, 0) })
			},
		),
	)
	partyRemarks := retui.Box(
		retui.Props{},
		retui.NewStyle(),
		components.TextInput().
			ID(fmt.Sprintf("partyRemarks-%d", i)).
			Value(v.Lines[i].Remarks).
			Focused(v.FocusIndex == base+1).
			Width(65).
			OnChange(func(id, value string) {
				setLine(func(l *receipt.PartyLine) { l.Remarks = value })
			}).
			Render(),
	)

	partyAmount := retui.Box(
		retui.Props{},
		retui.NewStyle(),
		components.NumberInput().
			ID(fmt.Sprintf("partyAmount-%d", i)).
			Decimals(2).
			Value(float64(v.Lines[i].Amount)).
			Focused(v.FocusIndex == base+2).
			Width(30).
			OnChange(func(id string, value float64) {
				setLine(func(l *receipt.PartyLine) { l.Amount = float32(value) })
			}).
			OnKeyPress(func(s string, key retui.Key) bool {
				if key.Code == retui.KeyEnter {
					v := form.Values()
					v.Lines = append(v.Lines, receipt.PartyLine{})
					v.FocusIndex = totalFields + (len(v.Lines)-1)*fieldsPerLine
					form.SetValues(v)
					return true // Consume the event
				}
				return false // Allow other key events to pass through
			}).
			Render(),
	)

	return retui.Box(
		retui.Props{},
		retui.NewStyle(),
		partyLedger,
		partyRemarks,
		partyAmount,
	)
}
