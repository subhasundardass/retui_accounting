package app

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
	"github.com/subhasundardass/retui/ui"
	"github.com/subhasundardass/retui/ui/layout"
)

// Root is the main application entry point rendered by retui on every frame.
func Root(ctx *appctx.AppContext, props retui.Props) retui.Element {

	// retui.SetFocusOrder([]string{"sidebar", "content"})

	if retui.CurrentKey.Code == retui.KeyTab {
		retui.FocusNext()
	}

	// ── Focus ────────────────────────────────────────────────────────────────
	if retui.CurrentFocus() == "" && !window.IsAnyModalOpen() {
		retui.SetFocus("sidebar")
	}

	// ── Get screen  ─────────────────────────────────
	currentID := retui.CurrentScreen()

	var content retui.Element
	screen, ok := ui.GetScreen(currentID)

	if !ok {
		screen, ok = ui.GetScreen("dashboard")

		if !ok {
			content = retui.Text("404 - No screen registered", retui.NewStyle().Foreground(retui.Red))
		} else {
			content = screen.Render(ctx)
		}
		retui.SetFocus("sidebar")
	} else {
		content = screen.Render(ctx)
	}

	title := "Unknown"
	if s, ok := ui.GetScreen(currentID); ok {
		title = s.Title
	}

	mainLayout := layout.MasterLayout(ctx, layout.LayoutProps{
		Ctx:     ctx,
		Title:   title,
		Content: content,
	})

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle(),
		mainLayout,
	)
}
