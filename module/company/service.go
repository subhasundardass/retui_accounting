package company

import (
	"context"

	"github.com/subhasundardass/retui/ent"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// List returns all companies.
func (s *Service) List(ctx context.Context) ([]*ent.Company, error) {
	return s.repo.List(ctx)
}

// Get returns a company by ID.
func (s *Service) Get(ctx context.Context, id int) (*ent.Company, error) {
	return s.repo.Get(ctx, id)
}
