package dashboard

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/dashboard/views"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {
	ui.Register("dashboard", ui.Screen{
		ID:     "dashboard",
		Title:  "Dashboard",
		Render: views.Dashboard,
	})

}
