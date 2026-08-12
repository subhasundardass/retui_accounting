package ledger

import (
	"context"

	"github.com/subhasundardass/retui/ent"
)

type Service interface {
	GetLedgers(ctx context.Context) ([]*ent.Ledger, error)
	GetGroups(ctx context.Context) ([]*ent.Ledger_Group, error)
	GetLedgersByGroup(ctx context.Context, groupID int) ([]*ent.Ledger, error)
	GetLedger(ctx context.Context, id int) (*ent.Ledger, error)
	GetCashBankLedgers(ctx context.Context) ([]*ent.Ledger, error)
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
