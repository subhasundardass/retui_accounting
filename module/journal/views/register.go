package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	// formComp := NewLedgerFormComponent(ctx)
	listJournalComp := NewJournalListComponent(ctx)

	ui.Register("journal_list", ui.Screen{
		ID:     "journal_list",
		Title:  "Ledgers",
		Render: listJournalComp.List,
	})

	// ui.Register("company_form", ui.Screen{
	// 	ID:     "company_form",
	// 	Title:  "Company Form",
	// 	Render: formComp.Form,
	// })
}
