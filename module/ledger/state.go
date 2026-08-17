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
	// Form state
	FocusIndex int
	Errors     map[string]string
	Mode       FormMode

	// Core fields
	Code        string
	Name        string
	Alias       string
	GroupID     int
	Description string

	// Accounting fields
	// OpeningBalance     float64
	// OpeningBalanceDate string // Using string for form input, can convert to time.Time
	// Balance float64

	// Party classification
	PartyType string // CUSTOMER, SUPPLIER, BOTH, INTERNAL

	// Address fields
	AddressLine1 string
	AddressLine2 string
	City         string
	StateID      int
	CountryID    int
	Pincode      string

	// Contact fields
	Phone         string
	Mobile        string
	Email         string
	ContactPerson string

	// Tax registration
	GSTRegistrationType string // REGULAR, COMPOSITION, UNREGISTERED, CONSUMER, SEZ, OVERSEAS
	GSTIN               string
	PAN                 string

	// Bank details
	BankName      string
	BankAccountNo string
	BankIFSC      string
	BankBranch    string

	// Status fields
	// IsSystem bool
	// IsParty  bool
	// IsBank   bool
	// IsCash   bool
	IsActive bool
}
