package ledger

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type LedgerGroupState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode
	//--

	Code        string
	Name        string
	Nature      string
	IsSystem    bool
	Description string
}

type LedgerState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode
	//--

	Code        string
	Name        string
	Description string
	GroupID     string
}
