package journal_view

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui/window"
)

type Controller struct {
	window *window.Window
	ctx    *context.AppContext
	repo   *repository.JournalRepository
}

func NewController(ctx *context.AppContext, repo *repository.JournalRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
	}
}

//==Get Journal

func (c *Controller) GetJournal(id int) *ent.Journal {

	jr, err := c.repo.GetJournalWithLine(c.ctx.Context, id)

	if err != nil {
		return nil
	}

	// retui.Debugf("============%v", jr)

	return jr
}
