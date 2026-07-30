package layout

import (
	"fmt"

	"github.com/subhasundardass/retui/retui"
)

func Footer(props retui.Props) retui.Element {

	footer := retui.Box(
		retui.Props{
			Width:   retui.Grow(1),
			Justify: retui.JustifySpaceBetween, // ← Pushes left/right apart
			Align:   retui.AlignCenter,
		},
		retui.NewStyle().Background(retui.Gray(1)),
		retui.Box(
			retui.Props{},
			retui.NewStyle(),
			retui.Text("Ready", retui.Style{}),
		),
		retui.Box(
			retui.Props{},
			retui.NewStyle(),
			retui.Text(fmt.Sprintf("v%s", "0.6.0"), retui.Style{}),
		),
	)

	return footer
}
