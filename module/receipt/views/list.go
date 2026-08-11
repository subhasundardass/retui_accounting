package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/receipt"
	"github.com/subhasundardass/retui/retui"
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

	case retui.KeyF6:
		// win := OpenReceiptForm(c.ctx, func() {})
		// win.Show()
	}
}

func (c *Component) ReceiptBook(ctx *appctx.AppContext) retui.Element {
	c.bindKeys()
	return retui.Box(
		retui.Props{
			Gap: 1,
		},
		retui.NewStyle(),
		c.buildToolbar(),
	)
}

func (c *Component) buildToolbar() retui.Element {

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
		retui.Text("Receipt Book", retui.NewStyle().Bold(true)),
		retui.Text("Create Receipt <F6>", retui.NewStyle().Bold(true).Foreground(retui.Gold)),
	)
}
