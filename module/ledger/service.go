package ledger

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/util"
)

type Service interface {
	GetLedgers(ctx context.Context) ([]*ent.Ledger, error)
	GetGroups(ctx context.Context) ([]*ent.Ledger_Group, error)
	GetLedgersByGroup(ctx context.Context, groupID int) ([]*ent.Ledger, error)
	GetLedger(ctx context.Context, id int) (*ent.Ledger, error)
	GetCashBankLedgers(ctx context.Context) ([]*ent.Ledger, error)
	GetStatement(ctx context.Context, ledgerAc int, from, to time.Time) (Result, error)
	GetOpeningBalance(ctx context.Context, ledgerAc int, asOf time.Time) (float64, error)
}

type service struct {
	repo *Repository
}

func NewService(repo *Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetLedgers(ctx context.Context) ([]*ent.Ledger, error) {
	return s.repo.List(ctx)
}

func (s *service) GetGroups(ctx context.Context) ([]*ent.Ledger_Group, error) {
	return s.repo.Groups(ctx)
}

func (s *service) GetLedgersByGroup(
	ctx context.Context,
	groupID int,
) ([]*ent.Ledger, error) {
	return s.repo.ListByGroup(ctx, groupID)
}

func (s *service) GetLedger(
	ctx context.Context,
	id int,
) (*ent.Ledger, error) {
	return s.repo.GetLedger(ctx, id)
}

func (s *service) GetCashBankLedgers(ctx context.Context) ([]*ent.Ledger, error) {
	return s.repo.ListByCodes(ctx, "CASH", "BANK")
}

func (s *service) GetOpeningBalance(ctx context.Context, ledgerAc int, asOf time.Time) (float64, error) {
	if ledgerAc <= 0 {
		return 0, fmt.Errorf("statement: ledger account is required")
	}
	return s.repo.OpeningBalance(ctx, ledgerAc, asOf)
}

// ======
func (s *service) GetStatement(ctx context.Context, ledgerAc int, from, to time.Time) (Result, error) {
	if ledgerAc <= 0 {
		return Result{}, fmt.Errorf("statement: ledger account is required")
	}
	if to.Before(from) {
		return Result{}, fmt.Errorf("statement: to date cannot be before from date")
	}

	opening, err := s.repo.OpeningBalance(ctx, ledgerAc, from)
	if err != nil {
		return Result{}, err
	}

	lines, err := s.repo.Entries(ctx, ledgerAc, from, to)
	if err != nil {
		return Result{}, err
	}

	sort.SliceStable(lines, func(i, j int) bool {
		vi, vj := lines[i].Edges.Journal, lines[j].Edges.Journal
		if vi == nil || vj == nil {
			return false
		}
		return vi.VoucherDate.Before(vj.VoucherDate)
	})

	entries := make([]Entry, 0, len(lines))
	running := opening
	for _, l := range lines {
		j := l.Edges.Journal
		if j == nil {
			continue // line loaded without its journal edge — skip rather than panic
		}
		running += l.Debit - l.Credit
		entries = append(entries, Entry{
			ID:          int64(l.ID),
			Date:        j.VoucherDate,
			VoucherNo:   j.VoucherNo,
			VoucherType: j.VoucherType,
			Narration:   util.Deref(j.Narration),
			Debit:       l.Debit,
			Credit:      l.Credit,
			Balance:     running,
		})
	}

	return Result{
		Rows:           entries,
		OpeningBalance: opening,
		ClosingBalance: running,
	}, nil
}
