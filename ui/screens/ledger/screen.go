package ledger

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
)

func Screen(ctx *context.AppContext, props retui.Props) retui.Element {

	if retui.CurrentKey.Code == retui.KeyEscape {
		retui.PopScreen()
		retui.PopFocus()
		return retui.Box(retui.Props{}, retui.NewStyle())
	}

	repo := repository.NewLedgerRepository(ctx.DB)
	journalRepo := repository.NewJournalRepository(ctx.DB)
	controller := NewController(ctx, repo, journalRepo)
	components := NewComponents(controller)

	return components.RenderScreen()
}
