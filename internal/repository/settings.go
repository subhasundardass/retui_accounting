package repository

import (
	"context"
	"fmt"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/ent/settings"
	"github.com/subhasundardass/retui/internal/database"
)

type Setting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type SettingsRepository struct {
	db *database.DB
}

func NewSettingsRepository(db *database.DB) *SettingsRepository {
	return &SettingsRepository{
		db: db,
	}
}

// List of setting
func (r *SettingsRepository) List(ctx context.Context) ([]*ent.Settings, error) {
	client := r.db.GetClient()
	if client == nil {
		return nil, fmt.Errorf("database client not initialized")
	}
	return client.Settings.Query().All(ctx)
}

// Create a setting
func (r *SettingsRepository) Create(ctx context.Context, s *Setting) (*ent.Settings, error) {
	return r.db.GetClient().Settings.Create().
		SetKey(s.Key).
		SetValue(s.Value).
		Save(ctx)
}

// Get a setting by key
func (r *SettingsRepository) GetByKey(ctx context.Context, key string) (*ent.Settings, error) {
	return r.db.GetClient().Settings.Query().
		Where(settings.Key(key)).
		Only(ctx)
}

// Update a setting within a transaction
func (r *SettingsRepository) UpdateWithTx(ctx context.Context, key, newValue string) error {
	return r.db.WithTx(ctx, func(tx *ent.Tx) error {
		_, err := tx.Settings.Update().
			Where(settings.Key(key)).
			SetValue(newValue).
			Save(ctx)
		return err
	})
}

// Delete a setting
func (r *SettingsRepository) Delete(ctx context.Context, key string) error {
	_, err := r.db.GetClient().Settings.Delete().
		Where(settings.Key(key)).
		Exec(ctx)
	return err
}
