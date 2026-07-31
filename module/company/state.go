package company

type State struct {
	Loading bool
	Saving  bool

	CompanyID int

	Code    string
	Name    string
	Email   string
	Phone   string
	Address string

	Errors map[string]string
}
