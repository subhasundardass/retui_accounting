package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	// formComp := NewLedgerFormComponent(ctx)

	list_Journal := NewJournalListComponent(ctx)
	create_Journal := NewJournalCreateWindow(ctx)

	ui.Register("journal_list", ui.Screen{
		ID:     "journal_list",
		Title:  "Ledgers",
		Render: list_Journal.List,
	})

	ui.Register("journal_entry", ui.Screen{
		ID:     "journal_entry",
		Title:  "Journal Entry",
		Render: create_Journal.JournalCreateForm,
	})
}
