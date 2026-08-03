package company

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
)

type Controller struct {
	ctx  *appctx.AppContext
	repo *Repository
}

func NewController(ctx *appctx.AppContext) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: NewRepository(ctx.DB.Client),
	}
}

func (c *Controller) List() ([]*ent.Company, error) {
	return c.repo.List(c.ctx.Ctx())
}

func (c *Controller) Edit(id int) (*ent.Company, error) {
	comp, err := c.repo.Get(c.ctx.Ctx(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to load company %d: %w", id, err)
	}
	return comp, nil
}

func (c *Controller) Save(mode FormMode, id int, data FormState) (*ent.Company, error) {
	switch mode {
	case ModeCreate:
		if err := ValidateCreate(data); err != nil {
			return nil, err
		}
		return c.repo.Create(c.ctx.Ctx(), data)
	case ModeUpdate:
		return c.repo.Update(c.ctx.Ctx(), id, data)
	default:
		return nil, fmt.Errorf("unknown save mode: %v", mode)
	}
}
