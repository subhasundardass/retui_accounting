package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/journal"
	"github.com/subhasundardass/retui/ent/journal_line"
	"github.com/subhasundardass/retui/ent/ledger"
	"github.com/subhasundardass/retui/internal/database"
	"github.com/subhasundardass/retui/retui"
)

type JournalRepository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *JournalRepository {
	return &JournalRepository{
		client: client,
	}
}

// List of Journal
func (r *JournalRepository) List(ctx context.Context) ([]*ent.Journal, error) {

	return r.client.Journal.Query().
		Limit(40).All(ctx)
}

// READ - List with pagination
func (r *JournalRepository) ListWithPagination(ctx context.Context, offset, limit int) ([]*ent.Journal, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	journals, err := client.Journal.
		Query().
		Limit(limit).
		Offset(offset).
		All(ctx)

	retui.Infof("Found %d journals", len(journals))
	return journals, err
}

// List of Journal
func (r *JournalRepository) ListByLedger(ctx context.Context, ledgerID int) ([]*ent.Journal, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	return client.Journal.Query().Limit(40).All(ctx)
}

// Get Journal with lines
func (r *JournalRepository) GetJournalWithLine(ctx context.Context, journalID int) (*ent.Journal, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	j, err := client.Journal.
		Query().
		Where(journal.IDEQ(journalID)).
		WithLines(func(q *ent.JournalLineQuery) {
			q.WithLedger().
				Order(ent.Asc(journal_line.FieldLineNo))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// func (r *JournalRepository) ListByDate(ctx context.Context) ([]*ent.Journal, error)      {}
// func (r *JournalRepository) ListByDateRange(ctx context.Context) ([]*ent.Journal, error) {}
// func (r *JournalRepository) GetLines(ctx context.Context) ([]*ent.Journal, error)        {}

// ---Journal Line
type JournalLineRepository struct {
	db *database.DB
}

func NewJournalLineRepository(db *database.DB) *JournalLineRepository {
	return &JournalLineRepository{
		db: db,
	}
}

// --Create with Line
type JournalLineInput struct {
	LedgerCode    string
	Debit         float64
	Credit        float64
	Description   string
	ReferenceType string
	ReferenceID   *int
}

type CreateJournalInput struct {
	Date        time.Time
	VoucherNo   string
	VoucherDate time.Time
	ReferenceNo *string

	Lines []JournalLineInput
}

func (r *JournalRepository) CreateNew(
	ctx context.Context,
	journal *ent.Journal,
	lines []JournalLineInput,
) (*ent.Journal, error) {

	//--Validation

	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// -------------------------
	// Calculate totals
	// -------------------------
	var totalDebit float64
	var totalCredit float64

	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	// -------------------------
	// Create Journal
	// -------------------------
	newJournal, err := tx.Journal.
		Create().
		SetDate(journal.Date).
		SetVoucherType(journal.VoucherType).
		SetVoucherNo(journal.VoucherNo).
		SetVoucherDate(journal.VoucherDate).
		SetReferenceNo(*journal.ReferenceNo).
		SetJournalStatus(journal.JournalStatus).
		SetNarration(*journal.Narration).
		SetTotalDebit(totalDebit).
		SetTotalCredit(totalCredit).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	// -------------------------
	// Create Journal Lines
	// -------------------------
	for i, line := range lines {

		ledger, err := tx.Ledger.Query().
			Where(ledger.CodeEQ(line.LedgerCode)).
			Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("ledger %q not found: %w", line.LedgerCode, err)
		}

		builder := tx.Journal_Line.
			Create().
			SetJournalID(newJournal.ID).
			SetLedgerID(ledger.ID).
			SetDebit(line.Debit).
			SetCredit(line.Credit).
			SetLineNo(i + 1)

		if line.Description != "" {
			builder.SetDescription(line.Description)
		}

		if line.ReferenceType != "" {
			builder.SetReferenceType(line.ReferenceType)
		}

		if line.ReferenceID != nil {
			builder.SetReferenceID(*line.ReferenceID)
		}

		if _, err = builder.Save(ctx); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newJournal, nil
}

// Get Journal line by ledger
func (r *JournalRepository) JournalLineByLedger(ctx context.Context, ledgerID int) ([]*ent.Journal_Line, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	return client.Journal_Line.
		Query().
		Where(journal_line.LedgerIDEQ(ledgerID)).
		WithJournal().
		WithLedger().
		Order(ent.Asc(journal_line.FieldLineNo)).
		Limit(40).
		All(ctx)
}
