// journal/repository.go
package journal

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/journal"
	"github.com/subhasundardass/retui/ent/journal_line"
)

type Repository interface {
	Create(ctx context.Context, in VoucherInput) (*ent.Journal, error)
	Update(ctx context.Context, id int, in VoucherInput) (*ent.Journal, error)
	Get(ctx context.Context, id int) (*ent.Journal, error)
	List(ctx context.Context, filter ListFilter) ([]*ent.Journal, error)
	Delete(ctx context.Context, id int) error
}

type ListFilter struct {
	Type     VoucherType
	Status   Status
	LedgerID int
	Limit    int
	Offset   int
}

type repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, in VoucherInput) (*ent.Journal, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	j, err := createJournalWithLines(ctx, tx.Client(), in)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return r.Get(ctx, j.ID)
}

func (r *repository) Update(ctx context.Context, id int, in VoucherInput) (*ent.Journal, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	txClient := tx.Client()

	// Replace all lines: delete existing, then insert the new set.
	// Simpler and safer than diffing line-by-line for a form-driven UI.
	if _, err := txClient.Journal_Line.Delete().
		Where(journal_line.JournalIDEQ(id)). // adjust predicate name to match your generated package
		Exec(ctx); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("clear existing lines: %w", err)
	}

	debit, credit := Totals(in)

	update := txClient.Journal.UpdateOneID(id).
		SetDate(in.Date).
		SetVoucherType(string(in.Type)).
		SetVoucherNo(in.VoucherNo).
		SetVoucherDate(in.VoucherDate).
		SetTotalDebit(debit).
		SetTotalCredit(credit)

	applyOptionalHeaderFields(update, in)

	j, err := update.Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("update journal: %w", err)
	}

	if err := createLines(ctx, txClient, j.ID, in.Lines); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return r.Get(ctx, j.ID)
}

func (r *repository) Get(ctx context.Context, id int) (*ent.Journal, error) {
	return r.client.Journal.Query().
		Where(journal.IDEQ(id)). // adjust predicate name to match generated package
		WithLines().
		Only(ctx)
}

func (r *repository) List(ctx context.Context, filter ListFilter) ([]*ent.Journal, error) {
	q := r.client.Journal.Query().
		Order(journal.ByID(sql.OrderDesc())).
		WithLines()

	if filter.Type != "" {
		q = q.Where(journal.VoucherTypeEQ(string(filter.Type)))
	}
	if filter.Status != "" {
		q = q.Where(journal.JournalStatusEQ(journal.JournalStatus(filter.Status)))
	}
	if filter.LedgerID != 0 {
		q = q.Where(journal.HasLinesWith(journal_line.LedgerIDEQ(filter.LedgerID)))
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}

	return q.All(ctx)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	txClient := tx.Client()

	if _, err := txClient.Journal_Line.Delete().
		Where(journal_line.JournalIDEQ(id)).
		Exec(ctx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("delete lines: %w", err)
	}

	if err := txClient.Journal.DeleteOneID(id).Exec(ctx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("delete journal: %w", err)
	}

	return tx.Commit()
}

// ---- helpers ----

func createJournalWithLines(ctx context.Context, txClient *ent.Client, in VoucherInput) (*ent.Journal, error) {
	debit, credit := Totals(in)

	create := txClient.Journal.Create().
		SetDate(in.Date).
		SetVoucherType(string(in.Type)).
		SetVoucherNo(in.VoucherNo).
		SetVoucherDate(in.VoucherDate).
		SetTotalDebit(debit).
		SetTotalCredit(credit)

	applyOptionalHeaderFields(create, in)

	j, err := create.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create journal: %w", err)
	}

	if err := createLines(ctx, txClient, j.ID, in.Lines); err != nil {
		return nil, err
	}

	return j, nil
}

func createLines(ctx context.Context, txClient *ent.Client, journalID int, lines []LineInput) error {
	builders := make([]*ent.JournalLineCreate, 0, len(lines))

	for i, l := range lines {
		b := txClient.Journal_Line.Create().
			SetJournalID(journalID).
			SetLedgerID(l.LedgerID).
			SetDebit(l.Debit).
			SetCredit(l.Credit).
			SetLineNo(i + 1)

		if l.Description != "" {
			b = b.SetDescription(l.Description)
		}
		if l.ReferenceType != "" {
			b = b.SetReferenceType(l.ReferenceType)
		}
		if l.ReferenceID != 0 {
			b = b.SetReferenceID(l.ReferenceID)
		}

		builders = append(builders, b)
	}

	if _, err := txClient.Journal_Line.CreateBulk(builders...).Save(ctx); err != nil {
		return fmt.Errorf("create journal lines: %w", err)
	}

	return nil
}

// applyOptionalHeaderFields sets the header's Nillable/optional fields.
// Works against either a JournalCreate or JournalUpdateOne builder via
// a small local interface, so Create and Update share the same logic.
type headerSetter interface {
	SetNillableReferenceNo(*string) *ent.JournalCreate // placeholder — see note below
}

func applyOptionalHeaderFields(b interface{}, in VoucherInput) {
	switch v := b.(type) {
	case *ent.JournalCreate:
		if in.ReferenceNo != "" {
			v.SetReferenceNo(in.ReferenceNo)
		}
		if in.ExternalRef != "" {
			v.SetExternalRef(in.ExternalRef)
		}
		if in.Status != "" {
			v.SetJournalStatus(journal.JournalStatus(in.Status))
		}
		if in.ApprovedBy != 0 {
			v.SetApprovedBy(in.ApprovedBy)
		}
		if in.FinancialYearID != 0 {
			v.SetFinancialYearID(in.FinancialYearID)
		}
		if in.Narration != "" {
			v.SetNarration(in.Narration)
		}
		if in.SourceModule != "" {
			v.SetSourceModule(in.SourceModule)
		}
		if in.SourceType != "" {
			v.SetSourceType(in.SourceType)
		}
		if in.SourceID != 0 {
			v.SetSourceID(in.SourceID)
		}
	case *ent.JournalUpdateOne:
		if in.ReferenceNo != "" {
			v.SetReferenceNo(in.ReferenceNo)
		}
		if in.ExternalRef != "" {
			v.SetExternalRef(in.ExternalRef)
		}
		if in.Status != "" {
			v.SetJournalStatus(journal.JournalStatus(in.Status))
		}
		if in.ApprovedBy != 0 {
			v.SetApprovedBy(in.ApprovedBy)
		}
		if in.FinancialYearID != 0 {
			v.SetFinancialYearID(in.FinancialYearID)
		}
		if in.Narration != "" {
			v.SetNarration(in.Narration)
		}
		if in.SourceModule != "" {
			v.SetSourceModule(in.SourceModule)
		}
		if in.SourceType != "" {
			v.SetSourceType(in.SourceType)
		}
		if in.SourceID != 0 {
			v.SetSourceID(in.SourceID)
		}
	}
}
