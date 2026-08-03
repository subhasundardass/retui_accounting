package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	// formComp := NewLedgerFormComponent(ctx)
	listLedgerComp := NewLedgerComponent(ctx)
	listLedgerGroupComp := NewLedgerGroupComponent(ctx)

	ui.Register("ledger_list", ui.Screen{
		ID:     "ledger_list",
		Title:  "Ledgers",
		Render: listLedgerComp.List,
	})
	ui.Register("ledger_group", ui.Screen{
		ID:     "ledger_group",
		Title:  "Groups",
		Render: listLedgerGroupComp.List,
	})
	// ui.Register("company_form", ui.Screen{
	// 	ID:     "company_form",
	// 	Title:  "Company Form",
	// 	Render: formComp.Form,
	// })
}
