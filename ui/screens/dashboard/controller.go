package dashboard

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui/window"
)

type Controller struct {
	window *window.Window
	ctx    *context.AppContext
}

func NewController(ctx *context.AppContext) *Controller {
	return &Controller{
		ctx: ctx,
	}
}
