package ui

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"

	"github.com/subhasundardass/retui/ui/screens/dashboard"
	"github.com/subhasundardass/retui/ui/screens/journal_entry"
	"github.com/subhasundardass/retui/ui/screens/journal_line"
	"github.com/subhasundardass/retui/ui/screens/journal_report"
	"github.com/subhasundardass/retui/ui/screens/journal_view"
	"github.com/subhasundardass/retui/ui/screens/ledger"
	"github.com/subhasundardass/retui/ui/screens/ledger_group"
)

// ─── Screen Definition ────────────────────────────────────────────────────

type Screen struct {
	ID     string
	Title  string
	Render func(ctx *appctx.AppContext, props retui.Props) retui.Element
}

// ─── Helper Functions ─────────────────────────────────────────────────────

func GetScreen(id string) (Screen, bool) {
	screen, ok := Routes[id]
	return screen, ok
}

// Registry holds all available screens
var Routes = map[string]Screen{

	// Add new screens here...
	"dashboard": {
		ID:     "dashboard",
		Title:  "Dashboard",
		Render: dashboard.Screen,
	},
	"ledger_group": {
		ID:     "ledger_group",
		Title:  "Groups",
		Render: ledger_group.Screen,
	},
	"ledger_list": {
		ID:     "ledger_list",
		Title:  "Ledgers",
		Render: ledger.Screen,
	},
	"journal_entry": {
		ID:     "journnal_entry",
		Title:  "Journal Entry",
		Render: journal_entry.Screen,
	},
	"journal_report": {
		ID:     "journal_report",
		Title:  "Journals",
		Render: journal_report.Screen,
	},
	"journal_view": {
		ID:     "journal_view",
		Title:  "Journal View",
		Render: journal_view.Screen,
	},
	"journal_line": {
		ID:     "journal_line",
		Title:  "Journal line",
		Render: journal_line.Screen,
	},
}
