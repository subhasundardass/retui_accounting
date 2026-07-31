package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

func Form(ctx *appctx.AppContext) retui.Element {

	// controller := company.NewController(ctx)
	// state := controller.State()

	return retui.Box(
		retui.Props{},
		retui.NewStyle(),

		// Inputs

		// Save
		// Cancel
	)
}
