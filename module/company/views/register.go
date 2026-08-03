package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	formComp := NewFormComponent(ctx)
	listComp := NewComponent(ctx, formComp)

	ui.Register("companies", ui.Screen{
		ID:     "companies",
		Title:  "Companies",
		Render: listComp.List,
	})
	// ui.Register("company_form", ui.Screen{
	// 	ID:     "company_form",
	// 	Title:  "Company Form",
	// 	Render: formComp.Form,
	// })
}
