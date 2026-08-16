package views

import (
	"fmt"
	"strconv"

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
	// onSave func()
}

func NewFormComponent(ctx *appctx.AppContext) *FormComponent {
	return &FormComponent{
		controller: receipt.NewController(ctx),
		state: receipt.FormState{
			Lines: []receipt.PartyLine{{}},
		},
	}
}

func (c *FormComponent) bindKeys(form *retui.Form[receipt.FormState]) {
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
		c.save(form)

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
				Width: retui.Grow(1),
			},
			retui.NewStyle(),
			c.new(ctx),
		),
	)
}

func (c *FormComponent) LoadForEdit(id int, state receipt.FormState) {
	c.editing = true
	c.editID = id
	c.state = state
	// also need to push `state` into the form hook, depending on how UseForm reseeds
}

func (c *FormComponent) new(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(c.state) // single call, once per render
	c.bindKeys(form)               // single call too — see note below

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
				c.buildHead(ctx, form),
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
func (c *FormComponent) buildHead(ctx *appctx.AppContext, form *retui.Form[receipt.FormState]) retui.Element {

	// form := retui.UseForm(c.state)
	v := form.Values()

	c.bindKeys(form)

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
			ID("vcNo").
			Value(v.VcNo).
			Prefix(" : ").
			Focused(v.FocusIndex == 0).
			OnChange(func(id, value string) {
				if err := form.SetField("VcNo", value); err != nil {
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
	// Receipt Account (Cash/Bank)
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
		widgets.CashBankDropdown(
			ctx, "rcptAccount", strconv.Itoa(v.RcptAccount), 40, v.FocusIndex == 4,
			func(id, value string) {
				if err := form.SetField("RcptAccount", util.StringToInt(value, 0)); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			},
		),
		// components.SelectDropdown().
		// 	ID("rcptAccount").
		// 	Options(c.controller.CashBankOptions()).
		// 	Prefix(" : ").
		// 	Value(strconv.Itoa(v.RcptAccount)).
		// 	Focused(v.FocusIndex == 4).
		// 	OnChange(func(s1, value string) {
		// 		if err := form.SetField("RcptAccount", util.StringToInt(value, 0)); err != nil {
		// 			retui.Debugf("SetField error: %v", err)
		// 		}
		// 	}).
		// 	Render(),
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

	c.bindKeys(form)

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
	c.bindKeys(form)

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
				setLine(func(l *receipt.PartyLine) { l.Amount = float64(value) })
			}).
			OnKeyPress(func(s string, key retui.Key) bool {
				if key.Code == retui.KeyEnter {
					v := form.Values()
					if i == len(v.Lines)-1 { // only grow when on the last row
						v.Lines = append(v.Lines, receipt.PartyLine{})
						v.FocusIndex = totalFields + (len(v.Lines)-1)*fieldsPerLine
						form.SetValues(v)
					} else {
						v.FocusIndex = totalFields + (i+1)*fieldsPerLine // move to next existing row
						form.SetValues(v)
					}
					return true
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

// --Save
func (c *FormComponent) save(form *retui.Form[receipt.FormState]) {

	state := form.Values()

	if err := receipt.ValidateForm(state); err != nil {
		components.ShowError(err.Error())
		return // validation failed — don't reset, keep the user's input
	}

	mode := receipt.ModeCreate
	id := 0
	if c.editing {
		mode = receipt.ModeUpdate
		id = c.editID
	}

	_, err := c.controller.Save(mode, id, state)
	if err != nil {
		retui.Debugf("Save failed: %v", err)
		components.ShowError("Save failed: " + err.Error())
		return // save failed — don't reset, keep the user's input
	}

	components.ShowSuccess("Receipt saved")

	// Only reached on success. Clear edit-mode state too, or the next
	// unrelated "new entry" save would incorrectly PATCH the old record.
	c.editing = false
	c.editID = 0
	c.state = receipt.FormState{Lines: []receipt.PartyLine{{}}}

	form.Reset()

	// Reset() may only clear fields tracked via SetField — force the slice
	// and focus back explicitly so stale party lines can't linger.
	fresh := form.Values()
	fresh.Lines = []receipt.PartyLine{{}}
	fresh.FocusIndex = 0
	form.SetValues(fresh)
}
