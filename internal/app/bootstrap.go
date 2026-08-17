package app

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/subhasundardass/retui/internal/config"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/database"
	"github.com/subhasundardass/retui/internal/database/seed"
	"github.com/subhasundardass/retui/module/company"
	"github.com/subhasundardass/retui/retui"
)

var (
	globalBootstrap *Bootstrap
	bootstrapMu     sync.RWMutex
)

// Bootstrap holds the application's core dependencies and startup state.
type Bootstrap struct {
	DB     *database.DB
	Config *config.Config
	AppCtx *appctx.AppContext

	Ctx    context.Context
	Cancel context.CancelFunc

	cleanup []func() error

	state *appState

	mu sync.RWMutex
}

// ============================================================================
// INITIALIZATION
// ============================================================================

// NewBootstrap creates and initializes the application.
func NewBootstrap() (*Bootstrap, error) {
	retui.Debug("Initializing Bootstrap...")

	// -------------------------------------------------------------------------
	// Configuration
	// -------------------------------------------------------------------------

	cfg := config.Load()
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	b := &Bootstrap{
		Config:  cfg,
		cleanup: make([]func() error, 0),
		state: &appState{
			config: cfg,
			user: &user{
				name:  "Guest",
				email: "guest@example.com",
				role:  "viewer",
			},
		},
	}

	b.Ctx, b.Cancel = context.WithCancel(context.Background())

	// -------------------------------------------------------------------------
	// Database
	// -------------------------------------------------------------------------

	retui.Debug("Initializing database...")

	db, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if db == nil {
		return nil, fmt.Errorf("database client is nil")
	}

	b.DB = db

	b.RegisterCleanup(func() error {
		retui.Debug("Closing database connection...")

		if b.DB != nil {
			return b.DB.Close()
		}

		return nil
	})

	// -------------------------------------------------------------------------
	// Database schema
	// -------------------------------------------------------------------------

	if err := b.runMigrations(); err != nil {
		retui.Warnf("Migration warning: %v", err)
	}

	// -------------------------------------------------------------------------
	// Seed reference data
	// -------------------------------------------------------------------------

	if err := b.seedDatabase(); err != nil {
		retui.Warnf("Seed warning: %v", err)
	}

	// -------------------------------------------------------------------------
	// Application context
	// -------------------------------------------------------------------------

	b.AppCtx = &appctx.AppContext{}

	b.setContext()

	// -------------------------------------------------------------------------
	// Register application modules
	// -------------------------------------------------------------------------

	b.registerModules()

	// -------------------------------------------------------------------------
	// Initialize company
	// -------------------------------------------------------------------------

	if err := b.initializeCompany(); err != nil {
		return nil, fmt.Errorf("failed to initialize company: %w", err)
	}

	// -------------------------------------------------------------------------
	// Global bootstrap
	// -------------------------------------------------------------------------

	SetBootstrap(b)

	retui.Success("Bootstrap completed successfully")

	return b, nil
}

// ============================================================================
// COMPANY INITIALIZATION
// ============================================================================

// initializeCompany determines the company state for the current session.
//
// The company controller is responsible for:
//   - restoring the previously selected company
//   - detecting whether no companies exist
//   - detecting whether multiple companies exist
//   - setting the active company in AppContext
func (b *Bootstrap) initializeCompany() error {
	if b.AppCtx == nil {
		return fmt.Errorf("app context is nil")
	}

	controller := company.NewController(b.AppCtx)

	result, err := controller.Initialize()
	if err != nil {
		return err
	}

	switch result.State {
	case company.StartupCreateCompany:
		retui.Info("No companies found; company creation required")

	case company.StartupSelectCompany:
		retui.Info("Multiple companies found; company selection required")

	case company.StartupReady:
		if result.Company == nil {
			return fmt.Errorf("startup company state is ready but company is nil")
		}

		retui.Infof(
			"Active company: %s (%d)",
			result.Company.Name,
			result.Company.ID,
		)

	default:
		return fmt.Errorf(
			"unknown company startup state: %d",
			result.State,
		)
	}

	return nil
}

// ============================================================================
// DATABASE FUNCTIONS
// ============================================================================

// InitDB initializes the database connection.
func InitDB(cfg *config.Config) (*database.DB, error) {
	retui.Infof("Connecting to database...")

	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// Config.Load() already resolves SQLitePath to an absolute path.
	// Keep that resolved path instead of resolving it again against cwd.
	dbPath := cfg.SQLitePath

	if dbPath == "" {
		dbPath = filepath.Join(
			config.ExecutableDir(),
			"data",
			"retui.db",
		)
	} else if !filepath.IsAbs(dbPath) {
		absPath, err := filepath.Abs(dbPath)
		if err == nil {
			dbPath = absPath
		}
	}

	retui.Infof("Database path: %s", dbPath)

	db, err := database.New("sqlite", dbPath)
	if err != nil {
		retui.Errorf("Failed to open database: %v", err)

		return nil, fmt.Errorf(
			"failed to initialize database: %w",
			err,
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	retui.Info("Creating/verifying schema...")

	if err := db.Client.Schema.Create(ctx); err != nil {
		retui.Errorf("Failed to create schema: %v", err)

		return nil, fmt.Errorf(
			"failed to create schema: %w",
			err,
		)
	}

	retui.Successf("Database connected: %s", dbPath)

	return db, nil
}

// runMigrations runs database migrations.
func (b *Bootstrap) runMigrations() error {
	retui.Info("Running database migrations...")

	ctx, cancel := context.WithTimeout(
		b.Ctx,
		10*time.Second,
	)
	defer cancel()

	if err := b.DB.Client.Schema.Create(ctx); err != nil {
		return fmt.Errorf(
			"failed to create schema: %w",
			err,
		)
	}

	retui.Success("Migrations completed successfully")

	return nil
}

// seedDatabase populates initial reference data.
//
// Seed functions are expected to be safe to execute repeatedly.
func (b *Bootstrap) seedDatabase() error {
	retui.Info("Seeding database...")

	ctx, cancel := context.WithTimeout(
		b.Ctx,
		10*time.Second,
	)
	defer cancel()

	if err := seed.AccountingSeed(
		ctx,
		b.DB.Client,
	); err != nil {
		return fmt.Errorf(
			"failed to seed accounting data: %w",
			err,
		)
	}

	retui.Success("Database seeded successfully")

	return nil
}

// ============================================================================
// CONTEXT FUNCTIONS
// ============================================================================

// setContext initializes the application context.
func (b *Bootstrap) setContext() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.AppCtx == nil {
		retui.Error("App context not initialized")
		return
	}

	if b.DB == nil {
		retui.Error("DB is nil in bootstrap")
		return
	}

	userName := "Guest"

	if b.state.user != nil &&
		b.state.user.name != "" {
		userName = b.state.user.name
	}

	b.AppCtx.Set(appctx.AppContextValues{
		Context:   b.Ctx,
		AppName:   b.state.config.AppName,
		DarkMode:  b.state.darkMode,
		UserName:  userName,
		DB:        b.DB,
		Config:    b.state.config,
		CompanyID: b.AppCtx.CompanyID(),

		ToggleDark: func() {
			b.mu.Lock()
			b.state.darkMode = !b.state.darkMode
			darkMode := b.state.darkMode
			b.mu.Unlock()

			if b.AppCtx != nil {
				b.AppCtx.ToggleDark()
			}

			retui.Debugf(
				"Dark mode toggled: %v",
				darkMode,
			)
		},
	})

	retui.Debug("Context values set successfully")
}

// UpdateContext updates the application context.
func (b *Bootstrap) UpdateContext() {
	b.setContext()
}

// ============================================================================
// CLEANUP AND SHUTDOWN
// ============================================================================

// RegisterCleanup registers a function to execute during shutdown.
func (b *Bootstrap) RegisterCleanup(fn func() error) {
	if fn == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.cleanup = append(b.cleanup, fn)
}

// Shutdown gracefully shuts down the application.
func (b *Bootstrap) Shutdown() error {
	retui.Debug("Shutting down...")

	// Cancel the application context first so background work can stop.
	if b.Cancel != nil {
		b.Cancel()
	}

	var errs []error

	b.mu.Lock()

	// Execute cleanup in reverse registration order.
	for i := len(b.cleanup) - 1; i >= 0; i-- {
		if err := b.cleanup[i](); err != nil {
			errs = append(errs, err)
		}
	}

	b.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf(
			"errors during shutdown: %v",
			errs,
		)
	}

	retui.Success("Shutdown completed successfully")

	return nil
}

// Close is an alias for Shutdown.
func (b *Bootstrap) Close() error {
	return b.Shutdown()
}

// ============================================================================
// GLOBAL BOOTSTRAP
// ============================================================================

// SetBootstrap stores the global bootstrap instance.
func SetBootstrap(b *Bootstrap) {
	bootstrapMu.Lock()
	defer bootstrapMu.Unlock()

	globalBootstrap = b
}

// GetBootstrap returns the global bootstrap instance.
func GetBootstrap() *Bootstrap {
	bootstrapMu.RLock()
	defer bootstrapMu.RUnlock()

	return globalBootstrap
}
