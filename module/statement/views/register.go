package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	statement := NewComponent(ctx)

	ui.Register("statement", ui.Screen{
		ID:     "statement",
		Title:  "Statement of Account",
		Render: statement.List,
	})

}
