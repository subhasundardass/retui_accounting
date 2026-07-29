package createcompany

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui/window"
)

func Window(ctx *context.AppContext) *window.Window {

	repo := repository.NewLedgerRepository(ctx.DB)
	controller := NewController(ctx, repo)
	components := NewComponents(controller)

	return components.RenderWindow()
}
