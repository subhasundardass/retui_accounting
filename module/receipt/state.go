package receipt

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type FormState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode
	//--

	SlNo      int
	Date      string
	Reference string
	Narration string
}
