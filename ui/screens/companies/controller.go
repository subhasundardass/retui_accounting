package companies

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui/window"
)

type Controller struct {
	window *window.Window
	ctx    *context.AppContext
	repo   *repository.CompanyRepository
}

func NewController(ctx *context.AppContext, repo *repository.CompanyRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
	}
}

// List returns all companies.
func (c *Controller) List() ([]*ent.Company, error) {
	return c.repo.List(c.ctx.Context)
}
