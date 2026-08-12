package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	formGroupComp := NewGroupFormComponent(ctx)
	formLedgerComp := NewLedgerFormComponent(ctx)

	listLedgerComp := NewLedgerComponent(ctx, formLedgerComp)
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
	ui.Register("ledger_edit", ui.Screen{
		ID:     "ledger_edit",
		Title:  "Ledger Edit",
		Render: formLedgerComp.LedgerEditForm,
	})
	ui.Register("ledger_create", ui.Screen{
		ID:     "ledger_create",
		Title:  "Ledger Create",
		Render: formLedgerComp.LedgerCreateForm,
	})

}
