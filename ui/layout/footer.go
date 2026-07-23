package layout

import (
	"fmt"

	"github.com/subhasundardass/retui/retui"
)

func Footer(props retui.Props) retui.Element {

	footer := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Justify:   retui.JustifySpaceBetween, // ← Pushes left/right apart
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{},
			retui.NewStyle(),
			retui.Text("Ready", retui.Style{}),
		),
		retui.Box(
			retui.Props{},
			retui.NewStyle(),
			retui.Text(fmt.Sprintf("v%s", "1.0.1"), retui.Style{}),
		),
	)

	return footer
}
