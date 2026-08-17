package statement

import (
	"time"

	"github.com/subhasundardass/retui/module/ledger"
)

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type FormState struct {
	// --- Filters (user input) ---
	LedgerAc   int
	FromDate   string
	ToDate     string
	FocusIndex int // 0=ledger, 1=from date, 2=to date, 3=table

	// --- Query results ---
	Rows           []ledger.Entry
	SelectedIndex  int // selected row within Rows, drives table.SelectedIndex
	OpeningBalance float64
	ClosingBalance float64

	// --- UI lifecycle ---
	Loading bool
	Loaded  bool
	Error   string
}

// Entry is a single statement line item, already carrying its running
// balance so the view layer never has to compute anything.
type Entry struct {
	ID          int64
	Date        time.Time
	VoucherNo   string
	VoucherType string
	Narration   string
	Debit       float64
	Credit      float64
	Balance     float64
}

// Reset clears query results while preserving the user's filter inputs.
// Useful when the user changes a filter after a successful load, so stale
// rows don't linger on screen.
func (s *FormState) ResetResults() {
	s.Rows = nil
	s.SelectedIndex = 0
	s.OpeningBalance = 0
	s.ClosingBalance = 0
	s.Loaded = false
	s.Error = ""
}
