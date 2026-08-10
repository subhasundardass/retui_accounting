package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/util"
	"github.com/subhasundardass/retui/module/receipt"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/window"
)

//win := receipt.OpenReceiptForm(ctx, func() {
// companiesList.Refresh()
// })

type FormComponent struct {
	controller *receipt.Controller
	state      receipt.FormState

	editing bool
	editID  int

	// onSaved is called after a successful save so the caller (e.g. the
	// companies list) can refresh its data.
	onClose func()
}

func NewFormComponent(ctx *appctx.AppContext, onClose func()) *FormComponent {
	return &FormComponent{
		controller: receipt.NewController(ctx),
		state:      receipt.FormState{},
		onClose:    onClose,
	}
}

// const totalFields = 18

func OpenReceiptForm(ctx *appctx.AppContext, onClose func()) *window.Window {
	c := NewFormComponent(ctx, onClose)
	return c.buildWindow(ctx)
}

func OpenReceiptFormForEdit(ctx *appctx.AppContext, id int, onSaved func()) *window.Window {
	c := NewFormComponent(ctx, onSaved)
	c.editing = true
	c.editID = id
	// TODO: once receipt.Controller exposes Load, hydrate c.state here
	// so the window opens pre-filled instead of blank:
	// state, err := c.controller.Load(id)
	// if err == nil { c.state = state }
	return c.buildWindow(ctx)
}

func (c *FormComponent) buildWindow(ctx *appctx.AppContext) *window.Window {

	win := window.NewWindow().
		SetTitle("Receipt Entry").
		SetSize(120, 10).
		SetModal(true).
		Center()

	win.SetRenderFn(func() retui.Element {

		return c.renderWindow(ctx)
	})

	return win
}

func (c *FormComponent) renderWindow(ctx *appctx.AppContext) retui.Element {

	form := retui.UseForm(c.state)
	v := form.Values()

	return retui.Box(
		retui.Props{
			Gap:       0,
			Height:    retui.Grow(1),
			Width:     retui.Grow(1),
			Direction: retui.Column,
		},
		retui.NewStyle(),
		c.headerComponent(ctx, form, v),
		retui.Box(
			retui.Props{},
			retui.NewStyle().Border(retui.Border{Top: true, Title: &retui.BorderTitle{
				Text:  "Particular",
				Style: retui.NewStyle().Background(retui.Gray(1)),
				Align: retui.AlignStart,
			}}),
		),
	)
}

// --Head
// 1. Receipt Sl No
// 2. Receipt Date
// 3. Narration
// 4. Refference
func (c *FormComponent) headerComponent(ctx *appctx.AppContext, form *retui.Form[receipt.FormState], v receipt.FormState) retui.Element {

	rcpt_no := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(10)},
			retui.NewStyle(),
			retui.Text("Receipt No", retui.NewStyle()),
		),
		components.TextInput().
			ID("rcpt_no").
			Value(util.IntToString(v.SlNo)).
			Width(20).
			Prefix(" : ").
			Focused(v.FocusIndex == 0).
			OnChange(func(id, value string) {
				form.SetField("SlNo", util.StringToInt(value, 0))
			}).
			Render(),
	)

	rcpt_date := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(10)},
			retui.NewStyle(),
			retui.Text("Date", retui.NewStyle()),
		),
		retui.Box(retui.Props{}, retui.NewStyle(),
			components.DateInput().
				ID("date").
				Width(20).
				Prefix(" : ").
				Focused(v.FocusIndex == 1).
				Value(v.Date).
				Format("DD/MM/YYYY").
				OnChange(func(id, value string) {
					if err := form.SetField("Date", value); err != nil {
						retui.Debugf("SetField error: %v", err) // or however you actually log
					}
				}).
				Render(),
		),
	)

	// Reference
	rcpt_ref := retui.Box(
		retui.Props{Gap: 1},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Fixed(10)},
			retui.NewStyle(),
			retui.Text("Reference", retui.NewStyle()),
		),
		retui.Box(retui.Props{}, retui.NewStyle(),
			components.TextInput().
				ID("reference").
				Width(20).
				Prefix(" : ").
				Focused(v.FocusIndex == 2).
				Value(v.Reference).
				OnChange(func(id string, value string) {
					if err := form.SetField("Reference", value); err != nil {
						retui.Debugf("SetField error: %v", err) // or however you actually log
					}
				}).
				Render(),
		),
	)
	rcpt_narration := // Narration
		retui.Box(
			retui.Props{Gap: 1},
			retui.NewStyle(),
			retui.Box(
				retui.Props{Width: retui.Fixed(10)},
				retui.NewStyle(),
				retui.Text("Narration", retui.NewStyle()),
			),
			retui.Box(retui.Props{}, retui.NewStyle(),
				components.TextInput().
					ID("narration").
					Focused(true).
					Width(retui.Grow(1).Value).
					Prefix(" : ").
					Value(v.Narration).
					OnChange(func(id string, value string) {
						if err := form.SetField("Narration", value); err != nil {
							retui.Debugf("SetField error: %v", err) // or however you actually log
						}
					}).
					Render(),
			),
		)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Padding:   [4]int{1, 2, 1, 2},
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Gap: 1,
			},
			retui.NewStyle(),
			rcpt_no,
			rcpt_date,
			rcpt_ref,
		),
		rcpt_narration,
	)
}
