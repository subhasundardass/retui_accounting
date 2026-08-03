package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

// type JournalCreateComponent struct {
// 	controller *journal.JournalController
// 	ctx        *appctx.AppContext
// 	// form       *FormComponent
// }

// func NewJournalCreateWindow(ctx *appctx.AppContext) *JournalCreateComponent {
// 	return &JournalCreateComponent{
// 		controller: journal.NewController(ctx),
// 		ctx:        ctx,
// 	}
// }

func JournalCreateWindow(ctx *context.AppContext) *window.Window {

	// controller := journal.NewController(ctx)

	win := window.NewWindow().
		SetTitle("Create Company").
		SetModal(true).
		Center().
		SetSize(100, 40)

	win.SetRenderFn(func() retui.Element {

		return renderWindow()
	})

	return win
}

func renderWindow() retui.Element {

	// c.bindKeys()
	return retui.Box(
		retui.Props{
			Gap: 1,
		},
		retui.NewStyle(),
		retui.Text("Text,", retui.NewStyle()),
	)
}
