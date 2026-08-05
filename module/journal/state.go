package journal

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type JournalLine struct {
	LedgerCode string
	Debit      float64
	Credit     float64
	Remarks    string
}

type FormState struct {
	FocusIndex  int
	Errors      map[string]string
	Mode        FormMode
	VcNo        string
	VcDate      string
	VcReference string
	VcNarration string

	Lines []JournalLine
}
