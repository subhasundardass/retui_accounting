package ledger_create

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
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

func (c *Controller) LedgerSeedOptions() []components.SelectOption {
	ledgers, err := c.repo.Default(c.ctx.Context, 10)
	if err != nil {
		retui.Debug("LedgerSeedOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, 0, len(ledgers)+1)
	for _, l := range ledgers {
		opts = append(opts, components.SelectOption{Label: l.Name, Value: l.Code})
	}

	return opts
}

func LedgerFilterOptions() {

}

func (c *Controller) LedgerFilterOptions(query string) []components.SelectOption {
	ledgers, err := c.repo.Search(c.ctx.Context, query, 10)
	if err != nil {
		retui.Debug("LedgerFilterOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, len(ledgers))
	for i, l := range ledgers {
		opts[i] = components.SelectOption{Label: l.Name, Value: l.Code}
	}

	return opts
}
