package ledger_group

import (
	context "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

func Screen(ctx *context.AppContext, props retui.Props) retui.Element {

	if retui.CurrentKey.Code == retui.KeyEscape && !window.IsAnyModalOpen() {
		retui.PopScreen()
		retui.PopFocus()
		// return retui.Box(retui.Props{}, retui.NewStyle())
	}

	repo := repository.NewLedgerGroupRepository(ctx.DB)
	controller := NewController(ctx, repo)
	components := NewComponents(controller)

	return components.RenderScreen()
}
