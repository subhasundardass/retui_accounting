package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	formGroupComp := NewGroupFormComponent(ctx)

	listLedgerComp := NewLedgerComponent(ctx)
	listLedgerGroupComp := NewLedgerGroupComponent(ctx, formGroupComp)

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

}
