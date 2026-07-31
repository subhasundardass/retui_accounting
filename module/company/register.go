package company

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/company/views"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {
	ui.Register("companies", ui.RegisterScreen{
		ID:     "companies",
		Title:  "Companies",
		Render: views.List,
	})
	ui.Register("company_form", ui.RegisterScreen{
		ID:     "company_form",
		Title:  "Company Form",
		Render: views.Form,
	})
}
