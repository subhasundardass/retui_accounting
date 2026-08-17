package widgets

import "github.com/subhasundardass/retui/retui"

func ShortCutItem(text string, hotkey string) retui.Element {
	return retui.Box(
		retui.Props{
			Gap:     1,
			Width:   retui.Grow(1),
			Padding: [4]int{0, 0, 0, 0},
			Justify: retui.JustifySpaceBetween,
		},
		retui.NewStyle().Foreground(retui.White).Bold(false),
		retui.Text(text, retui.NewStyle().Foreground(retui.SkyBlue)),
		retui.Text(hotkey, retui.NewStyle().Bold(true)),
	)
}
