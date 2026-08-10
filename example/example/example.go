package example

import (
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

// Example is the root screen for the RetUI examples.
//
// It demonstrates the current screen-stack API, focus navigation, and the
// flexbox-style layout system.
func Example() retui.Element {
	// The root application has two focus regions: the sidebar and content.
	retui.SetFocusOrder([]string{"sidebar", "content"})

	if retui.CurrentFocus() == "" && !window.IsAnyModalOpen() {
		retui.SetFocus("sidebar")
	}

	if retui.CurrentKey.Code == retui.KeyTab {
		retui.FocusNext()
	}

	currentScreenID := retui.CurrentScreen()

	var content retui.Element
	screen, ok := GetScreen(currentScreenID)
	if !ok {
		content = retui.Text(
			"404 - Page Not Found\n\nScreen: "+currentScreenID,
			retui.NewStyle().Foreground(retui.Red),
		)
	} else {
		content = screen.Render(retui.Props{})
	}

	header := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Height:    retui.Fit(),
			Justify:   retui.JustifySpaceBetween,
		},
		retui.NewStyle().Border(retui.Border{Left: true, Right: true, Bottom: true, Top: true}),
		retui.Text("RetUI Examples", retui.NewStyle().Bold(true)),
		retui.Text("Current screen: "+currentScreenID, retui.NewStyle()),
	)

	mainContent := retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle(),
		content,
	)

	body := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{Width: retui.Grow(1), Height: retui.Grow(1)},
			retui.NewStyle(),
			Sidebar(),
		),
		mainContent,
	)

	return retui.Box(
		retui.Props{
			Direction: retui.Column,
			Width:     retui.Grow(1),
			Height:    retui.Grow(1),
		},
		retui.NewStyle(),
		header,
		body,
	)
}
