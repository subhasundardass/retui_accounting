package repository

import "github.com/subhasundardass/retui/internal/database"

type CountryStateRepository struct {
	db *database.DB
}

func NewCountryStateRepository(db *database.DB) *CountryStateRepository {
	return &CountryStateRepository{
		db: db,
	}
}
