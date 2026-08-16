package layout

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

func Header(props retui.Props) retui.Element {

	appctx := context.Use() // ← one line, done

	header := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Justify:   retui.JustifySpaceBetween,
			Align:     retui.AlignCenter,
		},
		retui.NewStyle().Background(retui.Gray(1)),
		retui.Text(appctx.AppName(), retui.NewStyle().Foreground(retui.White).Bold(true)),
		retui.Text("v0.0.15", retui.NewStyle()),
	)

	return header
}
