package journal

import (
	"context"
	"fmt"
	"strconv"
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
		Order(ent.Desc(journal.FieldCreateTime)).
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

// Create New
func (r *JournalRepository) CreateNew(
	ctx context.Context,
	journal *ent.Journal,
	lines []JournalLine,
) (*ent.Journal, error) {

	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// Use committed flag — cleaner than relying on err shadowing
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var totalDebit, totalCredit float64
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	newJournal, err := tx.Journal.
		Create().
		SetDate(journal.Date).
		SetVoucherType(journal.VoucherType).
		SetVoucherNo(journal.VoucherNo).
		SetVoucherDate(journal.VoucherDate).
		SetJournalStatus(journal.JournalStatus).
		SetTotalDebit(totalDebit).
		SetTotalCredit(totalCredit).
		Save(ctx)

	// Guard nil pointer dereference
	if journal.ReferenceNo != nil {
		newJournal, err = tx.Journal.UpdateOne(newJournal).
			SetReferenceNo(*journal.ReferenceNo).
			Save(ctx)
	}
	if journal.Narration != nil {
		newJournal, err = tx.Journal.UpdateOne(newJournal).
			SetNarration(*journal.Narration).
			Save(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Declare ledger and lineErr outside loop — no shadowing
	var ledgerResult *ent.Ledger
	for i, line := range lines {

		ledgerID, err := strconv.Atoi(line.LedgerCode)
		if err != nil {
			return nil, fmt.Errorf("invalid ledger code %q: %w", line.LedgerCode, err)
		}

		ledgerResult, err = tx.Ledger.Query(). // = not :=
							Where(ledger.IDEQ(ledgerID)).
							Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("ledger %q not found: %w", line.LedgerCode, err)
		}

		builder := tx.Journal_Line.
			Create().
			SetJournalID(newJournal.ID).
			SetLedgerID(ledgerResult.ID).
			SetDebit(line.Debit).
			SetCredit(line.Credit).
			SetLineNo(i + 1)

		if line.Remarks != "" {
			builder.SetDescription(line.Remarks)
		}

		if _, err = builder.Save(ctx); err != nil { // = not :=
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	committed = true // rollback skipped from here on
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
