package repository

import (
	"context"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/company"
	"github.com/subhasundardass/retui/internal/database"
)

type CompanyRepository struct {
	db *database.DB
}

func NewCompanyRepository(db *database.DB) *CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

// List returns all companies ordered by name.
func (r *CompanyRepository) List(ctx context.Context) ([]*ent.Company, error) {
	return r.db.Client.Company.
		Query().
		Order(ent.Asc(company.FieldName)).
		All(ctx)
}

// Get returns a company by ID.
func (r *CompanyRepository) Get(ctx context.Context, id int) (*ent.Company, error) {
	return r.db.Client.Company.
		Get(ctx, id)
}

// GetByCode returns a company by its code.
func (r *CompanyRepository) GetByCode(ctx context.Context, code string) (*ent.Company, error) {
	return r.db.Client.Company.
		Query().
		Where(company.CodeEQ(code)).
		Only(ctx)
}

// Create inserts a new company.
func (r *CompanyRepository) Create(ctx context.Context, req *ent.Company) (*ent.Company, error) {
	return r.db.Client.Company.
		Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetLegalName(req.LegalName).
		SetEmail(req.Email).
		SetPhone(req.Phone).
		SetWebsite(req.Website).
		SetTaxID(req.TaxID).
		SetGstin(req.Gstin).
		SetPan(req.Pan).
		SetCurrency(req.Currency).
		SetTimezone(req.Timezone).
		SetLogo(req.Logo).
		SetAddress(req.Address).
		SetCity(req.City).
		SetCountryRefID(req.Edges.StateRef.CountryID).
		SetStateRef(req.Edges.StateRef).
		SetPostalCode(req.PostalCode).
		SetActive(req.Active).
		Save(ctx)
}

// Update updates an existing company.
func (r *CompanyRepository) Update(ctx context.Context, req *ent.Company) (*ent.Company, error) {
	return r.db.Client.Company.
		UpdateOneID(req.ID).
		SetName(req.Name).
		SetCode(req.Code).
		SetLegalName(req.LegalName).
		SetEmail(req.Email).
		SetPhone(req.Phone).
		SetWebsite(req.Website).
		SetTaxID(req.TaxID).
		SetGstin(req.Gstin).
		SetPan(req.Pan).
		SetCurrency(req.Currency).
		SetTimezone(req.Timezone).
		SetLogo(req.Logo).
		SetAddress(req.Address).
		SetCity(req.City).
		SetCountryRef(req.Edges.CountryRef).
		SetStateRef(req.Edges.StateRef).
		SetPostalCode(req.PostalCode).
		SetActive(req.Active).
		Save(ctx)
}

// Delete removes a company.
func (r *CompanyRepository) Delete(ctx context.Context, id int) error {
	return r.db.Client.Company.
		DeleteOneID(id).
		Exec(ctx)
}

// ExistsByCode returns true if a company with the given code exists.
func (r *CompanyRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.db.Client.Company.
		Query().
		Where(company.CodeEQ(code)).
		Exist(ctx)
}

// Active returns all active companies.
func (r *CompanyRepository) Active(ctx context.Context) ([]*ent.Company, error) {
	return r.db.Client.Company.
		Query().
		Where(company.ActiveEQ(true)).
		Order(ent.Asc(company.FieldName)).
		All(ctx)
}

// Count returns the total number of companies.
func (r *CompanyRepository) Count(ctx context.Context) (int, error) {
	return r.db.Client.Company.
		Query().
		Count(ctx)
}
