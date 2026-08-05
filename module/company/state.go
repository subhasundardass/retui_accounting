package company

type FormMode int

const (
	ModeCreate FormMode = iota
	ModeUpdate
)

type FormState struct {
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode

	Code       string
	Name       string
	LegalName  string
	Email      string
	Phone      string
	Website    string
	Country    int
	State      int
	City       string
	PostalCode string
	Address    string
	TaxID      string
	GSTIN      string
	PAN        string
	IsActive   bool
}
