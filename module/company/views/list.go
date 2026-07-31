package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

func List(ctx *appctx.AppContext) retui.Element {

	// controller := company.NewController(ctx)

	// companies, err := controller.List()
	// if err != nil {
	// 	return retui.Text(err.Error())
	// }

	return retui.Box(
		retui.Props{},
		retui.NewStyle(),
		// Table
		// Toolbar
		// Search
	)
}
