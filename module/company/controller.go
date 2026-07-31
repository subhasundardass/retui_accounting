package company

import (
	appctx "github.com/subhasundardass/retui/internal/context"
)

type Controller struct {
	ctx   *appctx.AppContext
	repo  Repository
	state *State
}

func NewController(ctx *appctx.AppContext) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: *NewRepository(ctx.DB.Client),
		state: &State{
			Errors: make(map[string]string),
		},
	}
}

func (c *Controller) State() *State {
	return c.state
}

func (c *Controller) ResetState() {
	c.state = &State{
		Errors: make(map[string]string),
	}
}
