package company

import (
	"context"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/company"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// List returns all companies.
func (r *Repository) List(ctx context.Context) ([]*ent.Company, error) {
	return r.client.Company.
		Query().
		WithCountryRef().
		WithStateRef().
		All(ctx)
}

// Get returns a company by ID.
func (r *Repository) Get(ctx context.Context, id int) (*ent.Company, error) {
	return r.client.Company.
		Query().
		Where(company.ID(id)).
		WithCountryRef().
		WithStateRef().
		Only(ctx)
}

// Create inserts a new company.
func (r *Repository) Create(ctx context.Context, data FormState) (*ent.Company, error) {
	return r.client.Company.
		Create().
		SetName(data.Name).
		SetCode(data.Code).
		SetLegalName(data.LegalName).
		SetEmail(data.Email).
		SetPhone(data.Phone).
		SetWebsite(data.Website).
		SetTaxID(data.TaxID).
		SetGstin(data.GSTIN).
		SetPan(data.PAN).
		// SetCurrency(data.Currency).
		// SetTimezone(data.Timezone).
		// SetLogo(data.Logo).
		SetAddress(data.Address).
		SetCity(data.City).
		SetCountryRefID(data.Country).
		SetStateRefID(data.State).
		SetPostalCode(data.PostalCode).
		SetActive(data.IsActive).
		Save(ctx)
}

// Update updates an existing company.
func (r *Repository) Update(ctx context.Context, id int, data FormState) (*ent.Company, error) {
	return r.client.Company.
		UpdateOneID(id).
		SetName(data.Name).
		SetCode(data.Code).
		SetLegalName(data.LegalName).
		SetEmail(data.Email).
		SetPhone(data.Phone).
		SetWebsite(data.Website).
		SetTaxID(data.TaxID).
		SetGstin(data.GSTIN).
		SetPan(data.PAN).
		SetAddress(data.Address).
		SetCity(data.City).
		SetCountryRefID(data.Country).
		SetStateRefID(data.State).
		SetPostalCode(data.PostalCode).
		SetActive(data.IsActive).
		Save(ctx)
}

// Delete removes a company.
func (r *Repository) Delete(ctx context.Context, id int) error {
	return r.client.Company.
		DeleteOneID(id).
		Exec(ctx)
}

// ExistsByCode checks if a company code already exists.
func (r *Repository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.client.Company.
		Query().
		Where(company.CodeEQ(code)).
		Exist(ctx)
}

// Count returns the number of companies.
func (r *Repository) Count(ctx context.Context) (int, error) {
	return r.client.Company.
		Query().
		Count(ctx)
}
