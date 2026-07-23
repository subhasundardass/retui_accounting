package journal_line

import (
	"fmt"

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

func (c *Controller) getJounalLineByLedger(ledgerID int) []*ent.Journal_Line {

	journalLine, err := c.repo.JournalLineByLedger(c.ctx.Context, ledgerID)

	if err != nil {
		fmt.Print(err.Error())
		return nil
	}

	return journalLine
}
