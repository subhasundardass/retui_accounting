package ledger

import (
	"context"
	"fmt"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/ledger"
	"github.com/subhasundardass/retui/ent/ledger_group"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// List of Ledger
func (r *Repository) List(ctx context.Context) ([]*ent.Ledger, error) {

	return r.client.Ledger.Query().
		WithGroup().
		Order(ent.Asc(ledger_group.FieldName)).
		All(ctx)
}

func (r *Repository) GetLedger(ctx context.Context, id int) (*ent.Ledger, error) {
	return r.client.Ledger.
		Query().
		Where(ledger.ID(id)).
		WithGroup().
		Only(ctx)
}

func (r *Repository) ListByCodes(
	ctx context.Context,
	codes ...string,
) ([]*ent.Ledger, error) {
	return r.client.Ledger.
		Query().
		Where(ledger.HasGroupWith(ledger_group.CodeIn(codes...))).
		All(ctx)
}

// Default returns the first `limit` ledgers, used to seed the select
// dropdown's initial Options before the user has typed anything (see
// OnFilter's "empty query" gap in select.go — buildSelectElement never
// calls OnFilter when filterText == "", so something has to populate
// Options up front).
func (r *Repository) Default(ctx context.Context, limit int) ([]*ent.Ledger, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	if limit <= 0 {
		limit = 10
	}
	return client.Ledger.Query().
		WithGroup().
		Order(ent.Asc(ledger.FieldName)).
		Limit(limit).
		All(ctx)
}

func (r *Repository) DefaultGroup(ctx context.Context, limit int) ([]*ent.Ledger_Group, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	if limit <= 0 {
		limit = 10
	}
	return client.Ledger_Group.Query().
		Order(ent.Asc(ledger_group.FieldName)).
		Limit(limit).
		All(ctx)
}

// Search returns ledgers whose name contains query (case-insensitive),
// up to limit results. Intended to be called directly from the select
// dropdown's OnFilter on every keystroke — see the wiring example.
func (r *Repository) Search(ctx context.Context, query string, limit int) ([]*ent.Ledger, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	if limit <= 0 {
		limit = 10
	}
	if query == "" {
		return r.Default(ctx, limit)
	}
	return client.Ledger.Query().
		Where(ledger.NameContainsFold(query)).
		WithGroup().
		Order(ent.Asc(ledger.FieldName)).
		Limit(limit).
		All(ctx)
}

func (r *Repository) SearchGroup(ctx context.Context, query string, limit int) ([]*ent.Ledger_Group, error) {
	client := r.client
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	if limit <= 0 {
		limit = 10
	}
	if query == "" {
		return r.DefaultGroup(ctx, limit)
	}
	return client.Ledger_Group.Query().
		Where(ledger_group.NameContainsFold(query)).
		Order(ent.Asc(ledger_group.FieldName)).
		Limit(limit).
		All(ctx)
}

func (r *Repository) ListByGroup(ctx context.Context, groupID int) ([]*ent.Ledger, error) {

	return r.client.Ledger.Query().
		Where(ledger.GroupIDEQ(groupID)).
		WithGroup().
		All(ctx)
}

func (r *Repository) LedgerCreate(
	ctx context.Context,
	in LedgerState,
) (*ent.Ledger, error) {
	create := r.client.Ledger.
		Create().
		SetCode(in.Code).
		SetName(in.Name).
		SetAlias(in.Alias).
		SetDescription(in.Description).
		SetAddressLine1(in.AddressLine1).
		SetAddressLine2(in.AddressLine2).
		SetCity(in.City).
		SetPincode(in.Pincode).
		SetPhone(in.Phone).
		SetMobile(in.Mobile).
		SetEmail(in.Email).
		SetContactPerson(in.ContactPerson).
		SetGstRegistrationType(
			ledger.GstRegistrationType(in.GSTRegistrationType),
		).
		SetGstin(in.GSTIN).
		SetPan(in.PAN).
		SetBankName(in.BankName).
		SetBankAccountNo(in.BankAccountNo).
		SetBankIfsc(in.BankIFSC).
		SetBankBranch(in.BankBranch).
		SetIsActive(in.IsActive)

	// Set group if provided.
	if in.GroupID > 0 {
		create.SetGroupID(in.GroupID)
	}

	// Set state if provided.
	if in.StateID > 0 {
		create.SetStateID(in.StateID)
	}

	// Set country if provided.
	if in.CountryID > 0 {
		create.SetCountryID(in.CountryID)
	}

	return create.Save(ctx)
}

func (r *Repository) LedgerUpdate(ctx context.Context, id int, in LedgerState) (*ent.Ledger, error) {
	update := r.client.Ledger.
		UpdateOneID(id).
		SetCode(in.Code).
		SetName(in.Name).
		SetAlias(in.Alias).
		SetGroupID(in.GroupID).
		SetDescription(in.Description).
		// SetPartyType(in.PartyType).
		SetAddressLine1(in.AddressLine1).
		SetAddressLine2(in.AddressLine2).
		SetCity(in.City).
		SetPincode(in.Pincode).
		SetPhone(in.Phone).
		SetMobile(in.Mobile).
		SetEmail(in.Email).
		SetContactPerson(in.ContactPerson).
		SetGstRegistrationType(ledger.GstRegistrationType(in.GSTRegistrationType)).
		SetGstin(in.GSTIN).
		SetPan(in.PAN).
		SetBankName(in.BankName).
		SetBankAccountNo(in.BankAccountNo).
		SetBankIfsc(in.BankIFSC).
		SetBankBranch(in.BankBranch).
		SetIsActive(in.IsActive)

	// Set state and country if provided
	if in.StateID > 0 {
		update.SetStateID(in.StateID)
	}
	if in.CountryID > 0 {
		update.SetCountryID(in.CountryID)
	}

	return update.Save(ctx)
}

// -- Groups
func (r *Repository) Groups(ctx context.Context) ([]*ent.Ledger_Group, error) {

	return r.client.Ledger_Group.Query().
		Limit(40).All(ctx)
}

func (r *Repository) GetGroup(ctx context.Context, id int) (*ent.Ledger_Group, error) {
	return r.client.Ledger_Group.
		Query().
		Where(ledger_group.ID(id)).
		Only(ctx)
}

func (r *Repository) GroupCreate(ctx context.Context, in LedgerGroupState) (*ent.Ledger_Group, error) {
	return r.client.Ledger_Group.
		Create().
		SetCode(in.Code).
		SetName(in.Name).
		SetNature(ledger_group.Nature(in.Nature)).
		SetIsSystem(in.IsSystem).
		SetDescription(in.Description).
		Save(ctx)
}

func (r *Repository) GroupUpdate(ctx context.Context, id int, in LedgerGroupState) (*ent.Ledger_Group, error) {
	return r.client.Ledger_Group.
		UpdateOneID(id).
		SetCode(in.Code).
		SetName(in.Name).
		SetNature(ledger_group.Nature(in.Nature)).
		SetIsSystem(in.IsSystem).
		SetDescription(in.Description).
		Save(ctx)
}
