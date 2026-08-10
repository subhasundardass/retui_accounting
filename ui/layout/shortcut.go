package layout

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

func ShortcutPanel(ctx *context.AppContext, props retui.Props) retui.Element {

	return retui.Box(
		retui.Props{
			Width:  retui.Percent(10),
			Height: retui.Grow(1),
		},
		retui.NewStyle().Background(retui.Hex("#00875f")),
	)
}
