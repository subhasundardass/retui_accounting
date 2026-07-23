package repository

import (
	"context"
	"fmt"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/ledger_group"
	"github.com/subhasundardass/retui/internal/database"
)

type LedgerGroupRepository struct {
	db *database.DB
}

func NewLedgerGroupRepository(db *database.DB) *LedgerGroupRepository {
	return &LedgerGroupRepository{
		db: db,
	}
}

// CREATE - Create a new Ledger Group
func (r *LedgerGroupRepository) Create(ctx context.Context, input ent.Ledger_Group) (*ent.Ledger_Group, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	// Create new ledger group
	newGroup, err := client.Ledger_Group.
		Create().
		SetName(input.Name).
		SetDescription(input.Description).
		// Add other fields as needed
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create ledger group: %w", err)
	}

	return newGroup, nil
}

// List of Ledger Group
func (r *LedgerGroupRepository) List(ctx context.Context) ([]*ent.Ledger_Group, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	return client.Ledger_Group.Query().Limit(40).All(ctx)
}

// READ - List with pagination
// func (r *LedgerGroupRepository) ListWithPagination(ctx context.Context, limit, offset int) ([]*ent.Ledger_Group, error) {
// 	client := r.db.GetClient()
// 	if client == nil {
// 		return nil, fmt.Errorf("database client not initialized")
// 	}
// 	return client.Ledger_Group.Query().
// 		Limit(limit).
// 		Offset(offset).
// 		All(ctx)
// }

// READ - Get a single Ledger Group by ID
func (r *LedgerGroupRepository) GetByID(ctx context.Context, id int) (*ent.Ledger_Group, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	group, err := client.Ledger_Group.
		Query().
		Where(ledger_group.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get ledger group by id %d: %w", id, err)
	}

	return group, nil
}

// UPDATE - Update an existing Ledger Group
func (r *LedgerGroupRepository) Update(ctx context.Context, id int, input ent.Ledger_Group) (*ent.Ledger_Group, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	// Check if record exists
	exists, err := client.Ledger_Group.Query().
		Where(ledger_group.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("ledger group with id %d not found", id)
	}

	// Update the record
	updatedGroup, err := client.Ledger_Group.
		UpdateOneID(id).
		SetName(input.Name).
		SetDescription(input.Description).
		// Add other fields as needed
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update ledger group: %w", err)
	}

	return updatedGroup, nil
}

// UPDATE - Partial update (only update specified fields)
func (r *LedgerGroupRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*ent.Ledger_Group, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}

	// Check if record exists
	exists, err := client.Ledger_Group.Query().
		Where(ledger_group.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("ledger group with id %d not found", id)
	}

	// Build update query
	update := client.Ledger_Group.UpdateOneID(id)

	if name, ok := updates["name"].(string); ok {
		update.SetName(name)
	}
	if description, ok := updates["description"].(string); ok {
		update.SetDescription(description)
	}
	// Add other fields as needed

	updatedGroup, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update ledger group: %w", err)
	}

	return updatedGroup, nil
}

// DELETE - Delete a Ledger Group by ID
func (r *LedgerGroupRepository) Delete(ctx context.Context, id int) error {
	client := r.db.GetClient()
	if client == nil {
		return fmt.Errorf("database client not initialized")
	}

	// Check if record exists
	exists, err := client.Ledger_Group.Query().
		Where(ledger_group.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("ledger group with id %d not found", id)
	}

	// Delete the record
	err = client.Ledger_Group.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete ledger group: %w", err)
	}

	return nil
}
