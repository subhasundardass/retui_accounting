package ledger_create

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui/window"
)

type Controller struct {
	window *window.Window
	ctx    *context.AppContext
	repo   *repository.LedgerRepository
}

func NewController(ctx *context.AppContext, repo *repository.LedgerRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
	}
}
