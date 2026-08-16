package payment

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type PartyLine struct {
	Ledger  int
	Remarks string
	Amount  float64
}

type FormState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode
	//--

	VcNo        string
	Reference   string
	Date        string
	Narration   string
	Amount      float64
	PymtAccount int

	//--
	Lines []PartyLine
}
