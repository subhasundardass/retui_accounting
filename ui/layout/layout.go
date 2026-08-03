// internal/ui/layout.go
package layout

import (
	"time"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
)

type LayoutProps struct {
	Ctx     *appctx.AppContext
	Title   string
	Content retui.Element
}

func MasterLayout(ctx *appctx.AppContext, props LayoutProps) retui.Element {
	if ctx == nil {
		return retui.Text("ERROR: ctx is nil in MasterLayout", retui.NewStyle().Foreground(retui.Red))
	}

	mainContent := retui.Box(
		retui.Props{Direction: retui.Column, Width: retui.Grow(1), Gap: 1},
		retui.NewStyle().Border(retui.Border{Color: retui.Gray(1)}).Background(retui.Black),
		props.Content,
	)

	body := retui.Box(
		retui.Props{Direction: retui.Row, Gap: 0, Width: retui.Grow(1), Height: retui.Grow(1)},
		retui.NewStyle(),
		SidebarTree(ctx, retui.Props{}),
		mainContent,
	)

	final := retui.Box(
		retui.Props{Direction: retui.Column, Gap: 0, Width: retui.Grow(1), Height: retui.Grow(1)},
		retui.NewStyle(),
		Header(retui.Props{}),
		body,
		Footer(retui.Props{}),
	)

	//--Config Toast
	components.ConfigureToasts(
		components.WithDefaultPosition(components.ToastTopRight),
		components.WithDefaultDuration(1*time.Second),
	)
	children := append([]retui.Element{final}, components.ToastLayer(nil)...)

	// Toasts painted last, above everything else.
	return retui.Box(
		retui.Props{Direction: retui.Column, Width: retui.Grow(1), Height: retui.Grow(1)},
		retui.NewStyle(),
		children...,
	)
}
