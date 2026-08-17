package layout

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/company"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/ui/widgets"
)

func Header(ctx *appctx.AppContext, props retui.Props) retui.Element {

	if ctx == nil {
		return retui.Text(
			"Loading...",
			retui.NewStyle(),
		)
	}

	companyController := company.NewController(ctx)

	companySwitcher := widgets.CompanySwitcher(
		ctx,
		companyController,
		"company",
		func(id, value string) {
			// Optional application logic.
		},
	)

	header := retui.Box(
		retui.Props{
			Direction: retui.Row,
			Padding:   [4]int{0, 1, 0, 1},
			Width:     retui.Grow(1),
			Justify:   retui.JustifySpaceBetween,
			Align:     retui.AlignCenter,
		},
		retui.NewStyle().Background(retui.Gray(1)),

		retui.Text(
			ctx.AppName(),
			retui.NewStyle().
				Foreground(retui.White).
				Bold(true),
		),

		companySwitcher,
	)

	return header
}
