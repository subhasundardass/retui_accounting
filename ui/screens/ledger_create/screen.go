package ledger_create

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui/window"
)

func Window(ctx *context.AppContext) *window.Window {

	// if retui.CurrentKey.Code == retui.KeyEscape {
	// 	retui.PopScreen()
	// 	retui.FocusPrev()
	// }

	repo := repository.NewLedgerRepository(ctx.DB)
	controller := NewController(ctx, repo)
	components := NewComponents(controller)

	return components.RenderWindow()
}
