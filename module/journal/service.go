package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/subhasundardass/retui/ent"
)

type VoucherType string

const (
	VoucherJV VoucherType = "JV" // Journal
	VoucherPV VoucherType = "PV" // Payment
	VoucherRV VoucherType = "RV" // Receipt
	VoucherCV VoucherType = "CV" // Contra
	VoucherSV VoucherType = "SV" // Sales
	VoucherPR VoucherType = "PR" // Purchase
)

type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusApproved  Status = "APPROVED"
	StatusPosted    Status = "POSTED"
	StatusCancelled Status = "CANCELLED"
)

// Mode indicates whether we're creating a new voucher or updating an
// existing one.
type Mode int

const (
	ModeCreate Mode = iota
	ModeUpdate
)

// LineInput is a single debit or credit leg of a voucher. Exactly one
// of Debit/Credit should be non-zero — never both, never neither.
type LineInput struct {
	LedgerID      int
	Debit         float64
	Credit        float64
	Description   string
	ReferenceType string
	ReferenceID   int
}

// VoucherInput is the generic shape every accounting module (receipt,
// payment, sales, purchase, contra, journal) builds and hands to the
// journal Service. The service doesn't know or care which UI produced it.
type VoucherInput struct {
	Type            VoucherType
	VoucherNo       string
	Date            time.Time // "date" field
	VoucherDate     time.Time // "voucher_date" field — usually same as Date, kept separate to match schema
	ReferenceNo     string
	ExternalRef     string
	Status          Status // defaults to StatusDraft if empty
	ApprovedBy      int    // 0 = not set
	FinancialYearID int    // 0 = not set
	Narration       string
	SourceModule    string // e.g. "receipt", "sales"
	SourceType      string
	SourceID        int
	Lines           []LineInput
}

// Service is the entry point every accounting module should depend on
// for posting transactions. It owns validation; callers never need to
// reimplement debit/credit balancing.
type Service interface {
	Save(ctx context.Context, mode Mode, id int, in VoucherInput) (*ent.Journal, error)
	Get(ctx context.Context, id int) (*ent.Journal, error)
	List(ctx context.Context, filter ListFilter) ([]*ent.Journal, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Save(ctx context.Context, mode Mode, id int, in VoucherInput) (*ent.Journal, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}

	switch mode {
	case ModeCreate:
		return s.repo.Create(ctx, in)
	case ModeUpdate:
		if id == 0 {
			return nil, fmt.Errorf("update requires a valid id")
		}
		return s.repo.Update(ctx, id, in)
	default:
		return nil, fmt.Errorf("unknown save mode: %v", mode)
	}
}

func (s *service) Get(ctx context.Context, id int) (*ent.Journal, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, filter ListFilter) ([]*ent.Journal, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
