package ledger

import (
	"context"

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
		Limit(40).All(ctx)
}

func (r *Repository) ListByGroup(ctx context.Context, groupID int) ([]*ent.Ledger, error) {

	return r.client.Ledger.Query().
		Where(ledger.GroupIDEQ(groupID)).
		WithGroup().
		Limit(40).All(ctx)
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

// // Default returns the first `limit` ledgers, used to seed the select
// // dropdown's initial Options before the user has typed anything (see
// // OnFilter's "empty query" gap in select.go — buildSelectElement never
// // calls OnFilter when filterText == "", so something has to populate
// // Options up front).
// func (r *LedgerRepository) Default(ctx context.Context, limit int) ([]*ent.Ledger, error) {
// 	client := r.client
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	return client.Ledger.Query().
// 		WithGroup().
// 		Order(ent.Asc(ledger.FieldName)).
// 		Limit(limit).
// 		All(ctx)
// }

// // Search returns ledgers whose name contains query (case-insensitive),
// // up to limit results. Intended to be called directly from the select
// // dropdown's OnFilter on every keystroke — see the wiring example.
// func (r *LedgerRepository) Search(ctx context.Context, query string, limit int) ([]*ent.Ledger, error) {
// 	client := r.client
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if query == "" {
// 		return r.Default(ctx, limit)
// 	}
// 	return client.Ledger.Query().
// 		Where(ledger.NameContainsFold(query)).
// 		WithGroup().
// 		Order(ent.Asc(ledger.FieldName)).
// 		Limit(limit).
// 		All(ctx)
// }

// // Create
// func (r *LedgerRepository) Insert(ctx context.Context, data *ent.Ledger) (*ent.Ledger, error) {
// 	client := r.client
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}

// 	// Validate required fields
// 	if data.Code == "" {
// 		return nil, fmt.Errorf("Code is required")
// 	}

// 	// Validate required fields
// 	if data.Name == "" {
// 		return nil, fmt.Errorf("name is required")
// 	}

// 	create := client.Ledger.Create().
// 		SetName(data.Name).
// 		SetDescription(data.Description)

// 	ledger, error := create.Save(ctx)
// 	if error != nil {
// 		retui.Errorf("Error Saving %s", error.Error())
// 	}
// 	return ledger, error
// }

// // Update updates an existing ledger
// func (r *LedgerRepository) Update(ctx context.Context, id int, data *ent.Ledger) (*ent.Ledger, error) {
// 	client := r.client
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}

// 	// Check if ledger exists
// 	exists, err := r.Exists(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !exists {
// 		return nil, fmt.Errorf("ledger with ID %d not found", id)
// 	}

// 	update := client.Ledger.UpdateOneID(id)

// 	// Only update fields that are provided
// 	if data.Name != "" {
// 		update = update.SetName(data.Name)
// 	}
// 	if data.Code != "" {
// 		update = update.SetCode(data.Code)
// 	}

// 	if data.Description != "" {
// 		update = update.SetDescription(data.Description)
// 	}

// 	return update.Save(ctx)
// }

// // Delete deletes a ledger by ID
// func (r *LedgerRepository) Delete(ctx context.Context, id int) error {
// 	client := r.client
// 	if client == nil {
// 		return fmt.Errorf("database client not initialized")
// 	}

// 	// Check if ledger exists
// 	exists, err := r.Exists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("ledger with ID %d not found", id)
// 	}

// 	return client.Ledger.DeleteOneID(id).Exec(ctx)
// }

// // Exists checks if a ledger exists
// func (r *LedgerRepository) Exists(ctx context.Context, id int) (bool, error) {
// 	client := r.client
// 	if client == nil {
// 		return false, fmt.Errorf("database client not initialized")
// 	}

// 	return client.Ledger.Query().
// 		Where(ledger.ID(id)).
// 		Exist(ctx)
// }

// // Get By Code
// func (r *LedgerRepository) GetByCode(ctx context.Context, code string) (*ent.Ledger, error) {
// 	client := r.client
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}
// 	return client.Ledger.Query().
// 		Where(ledger.CodeEQ(code)).
// 		WithGroup().
// 		Only(ctx)
// }
