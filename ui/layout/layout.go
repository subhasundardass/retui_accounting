// internal/ui/layout.go
package layout

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
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

	//Main content with border and title
	// borderChars := retui.BorderRounded
	if retui.IsFocused("content") {
		// borderChars = retui.BorderDouble
	}
	mainContent := retui.Box(
		retui.Props{
			Direction: retui.Column,
			Padding:   [4]int{0, 0, 0, 0},
			Width:     retui.Grow(1),
			Gap:       1,
		},
		retui.NewStyle().Border(retui.Border{
			// Top:    true,
			// Right: true,
			// Bottom: true,
			// Left: true,
			// Chars:  borderChars,
			Color: retui.Color{Type: retui.ColorANSI256, R: 40, G: 40, B: 40},
			Title: props.Title,
		}).Background(retui.Black),
		props.Content, // ← Your page content goes here
	)

	//Body: Sidebar + Main content (both should grow)
	body := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Gap:       0,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1), // ← Fill remaining height
		},
		retui.NewStyle(),
		SidebarTree(ctx, retui.Props{}), // Sidebar
		mainContent,                     // Main content
	)

	final := retui.Box(
		retui.Props{
			Direction: retui.Column,
			Gap:       0,
			Padding:   [4]int{0, 0, 0, 0},
			Width:     retui.Grow(1),
			Height:    retui.Grow(1), // ← Fill entire screen
		},
		retui.NewStyle(),
		Header(retui.Props{}), // Fixed height
		body,                  // Takes remaining space
		Footer(retui.Props{}), // Fixed height

	)

	return final
}
