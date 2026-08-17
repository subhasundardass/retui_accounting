package views

import (
	"strconv"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/statement"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/ui/widgets"
)

type Component struct {
	controller *statement.Controller
	ctx        *appctx.AppContext
	state      statement.FormState
}

func NewComponent(ctx *appctx.AppContext) *Component {
	return &Component{
		controller: statement.NewController(ctx),
		ctx:        ctx,
	}
}

const focusCount = 4

func (c *Component) bindKeys(form *retui.Form[statement.FormState]) {
	key := retui.CurrentKey
	v := form.Values()

	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	const tableFocus = focusCount - 1

	switch key.Code {
	case retui.KeyEscape:
		if v.FocusIndex == tableFocus {
			// Leave the table back to the filters, keep results loaded.
			// v.FocusIndex = 0
			// form.SetValuesSilent(v)
			form.Reset()
			return
		} else {
			if v.FocusIndex == 0 {
				retui.PopScreen()
			}
		}
		// form.Reset()

	case retui.KeyDown, retui.KeyUp:
		// While the table is focused, Up/Down belongs to its own row
		// navigation (components.Table wires OnChange itself). Don't
		// steal the key here.
		if v.FocusIndex == tableFocus {
			return
		}
		if key.Code == retui.KeyDown {
			v.FocusIndex = (v.FocusIndex + 1) % focusCount
		} else {
			v.FocusIndex = (v.FocusIndex - 1 + focusCount) % focusCount
		}
		form.SetValuesSilent(v)

	case retui.KeyTab:
		v.FocusIndex = (v.FocusIndex + 1) % focusCount
		form.SetValuesSilent(v)

	case retui.KeyShiftTab:
		v.FocusIndex = (v.FocusIndex - 1 + focusCount) % focusCount
		form.SetValuesSilent(v)

	case retui.KeyEnter:
		if v.FocusIndex == 2 { // "to date" field
			c.loadStatement(form)
		}
		// Enter on the table row is handled inside buildTable's OnChange,
		// where we still have the concrete Entry, not just an index.
	}
}

func (c *Component) List(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(c.state)
	c.bindKeys(form)

	return retui.Box(
		retui.Props{
			Justify:   retui.JustifySpaceBetween,
			Direction: retui.Column,
		},
		retui.NewStyle(),
		c.buildToolbar(ctx, form),
		c.buildTable(ctx, form),
	)
}

func (c *Component) buildToolbar(ctx *appctx.AppContext, form *retui.Form[statement.FormState]) retui.Element {

	v := form.Values()

	ledger := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(10),
			},
			retui.NewStyle(),
			retui.Text("Ledger : ", retui.NewStyle()),
		),
		widgets.LedgerComponent(
			ctx, "ledgerAccount", strconv.Itoa(v.LedgerAc), 40, v.FocusIndex == 0,
			func(id, value string) {
				if err := form.SetField("LedgerAc", util.StringToInt(value, 0)); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			},
		),
	)

	from_date := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(10),
			},
			retui.NewStyle(),
			retui.Text("From Date : ", retui.NewStyle()),
		),
		components.DateInput().
			Prefix(" : ").
			ID("date").
			Focused(v.FocusIndex == 1).
			Value(v.FromDate).
			Format("DD/MM/YYYY").
			OnChange(func(id, value string) {
				if err := form.SetField("FromDate", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

	to_date := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap:   1,
				Width: retui.Fixed(10),
			},
			retui.NewStyle(),
			retui.Text("To Date : ", retui.NewStyle()),
		),
		components.DateInput().
			Prefix(" : ").
			ID("date").
			Focused(v.FocusIndex == 2).
			Value(v.ToDate).
			Format("DD/MM/YYYY").
			OnChange(func(id, value string) {
				if err := form.SetField("ToDate", value); err != nil {
					retui.Debugf("SetField error: %v", err)
				}
			}).
			Render(),
	)

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

		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Text("Statement of Account", retui.NewStyle()),
		),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			retui.Box(
				retui.Props{
					Gap: 1,
				},
				retui.NewStyle(),
				ledger,
				from_date,
				to_date,
			),
		),
	)
}

func (c *Component) buildTable(ctx *appctx.AppContext, form *retui.Form[statement.FormState]) retui.Element {
	v := form.Values()

	if v.Loading {
		return retui.Box(
			retui.Props{Direction: retui.Row, Justify: retui.JustifyCenter, Padding: [4]int{2, 0, 2, 0}},
			retui.NewStyle(),
			retui.Text("Loading…", retui.NewStyle().Foreground(retui.Gray(2))),
		)
	}

	if v.Error != "" {
		return retui.Box(
			retui.Props{Direction: retui.Row, Justify: retui.JustifyCenter, Padding: [4]int{2, 0, 2, 0}},
			retui.NewStyle(),
			retui.Text("Error: "+v.Error, retui.NewStyle().Foreground(retui.Red)),
		)
	}

	rows := make([][]string, len(v.Rows))
	for i, e := range v.Rows {
		rows[i] = []string{
			e.Date.Format("02/01/2006"),
			e.VoucherNo,
			e.VoucherType,
			e.Narration,
			formatAmount(e.Debit),
			formatAmount(e.Credit),
			formatAmount(e.Balance),
		}
	}

	tbl := components.Table().
		ID("statement_table").
		Headers([]string{
			"Date", "Voucher No", "Type", "Narration", "Debit", "Credit", "Balance",
		}).
		Alignments([]string{
			"left", "left", "left", "left", "right", "right", "right",
		}).
		ColumnWidths([]int{12, 12, 10, 40, 14, 14, 14}).
		Focused(v.FocusIndex == focusCount-1).
		Rows(rows).
		SelectedIndex(v.SelectedIndex).
		OnChange(func(i int) {
			if i < 0 || i >= len(v.Rows) {
				return
			}
			if err := form.SetField("SelectedIndex", i); err != nil {
				retui.Debugf("SetField error: %v", err)
			}
			if retui.CurrentKey.Code == retui.KeyEnter {
				// c.controller.ShowVoucher(v.Rows[i].ID)
			}
		}).
		Render()

	summary := retui.Box(
		retui.Props{Direction: retui.Row, Justify: retui.JustifySpaceBetween, Padding: [4]int{0, 1, 0, 1}},
		retui.NewStyle().Border(retui.Border{Top: true, Color: retui.Gray(1)}),
		retui.Text("Opening Balance: "+formatAmount(v.OpeningBalance), retui.NewStyle()),
		retui.Text("Closing Balance: "+formatAmount(v.ClosingBalance), retui.NewStyle().Bold(true)),
	)

	if len(v.Rows) == 0 {
		return retui.Box(
			retui.Props{Direction: retui.Row, Justify: retui.JustifyCenter, Padding: [4]int{2, 0, 2, 0}},
			retui.NewStyle(),
			retui.Text("Select a ledger and date range, then press Enter to load the statement",
				retui.NewStyle().Foreground(retui.Gray(2))),
		)
	}

	return retui.Box(
		retui.Props{Direction: retui.Column, Height: retui.Fixed(33)},
		retui.NewStyle(),
		tbl,
		summary,
	)
}

func (c *Component) loadStatement(form *retui.Form[statement.FormState]) {
	v := form.Values()

	v.Loading = true
	v.Error = ""
	form.SetValues(v)

	result, err := c.controller.Load(v.LedgerAc, v.FromDate, v.ToDate)

	v = form.Values() // re-read in case fields changed mid-query
	v.Loading = false
	if err != nil {
		v.Error = err.Error()
		v.Loaded = false
		v.Rows = nil
		v.SelectedIndex = 0
		v.OpeningBalance = 0
		v.ClosingBalance = 0
	} else {
		v.Error = ""
		v.Loaded = true
		v.Rows = result.Rows
		v.SelectedIndex = 0
		v.OpeningBalance = result.OpeningBalance
		v.ClosingBalance = result.ClosingBalance
		v.FocusIndex = focusCount - 1 // jump straight to the table
	}
	form.SetValues(v)
}

func formatAmount(v float64) string {
	if v == 0 {
		return "-"
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}
