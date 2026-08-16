package ledger

import "time"

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

// Result is the fully computed statement handed back to the view layer.
type Result struct {
	Rows           []Entry
	OpeningBalance float64
	ClosingBalance float64
}
