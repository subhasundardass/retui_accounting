package receipt

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type PartyLine struct {
	Ledger  int
	Remarks string
	Amount  float32
}

type FormState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode
	//--

	SlNo        int
	Reference   string
	Date        string
	Narration   string
	Amount      float32
	RcptAccount int

	//--
	Lines []PartyLine
}
