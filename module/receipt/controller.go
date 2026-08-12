package receipt

import (
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
