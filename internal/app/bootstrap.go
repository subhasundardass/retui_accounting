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
	"github.com/subhasundardass/retui/retui"
)

var (
	globalBootstrap *Bootstrap
	bootstrapMu     sync.RWMutex
)

// Bootstrap holds the application's core dependencies
type Bootstrap struct {
	DB      *database.DB
	Config  *config.Config
	AppCtx  *appctx.AppContext
	Ctx     context.Context
	Cancel  context.CancelFunc
	cleanup []func() error
	state   *appState
	mu      sync.RWMutex
}

// ============================================================================
// INITIALIZATION
// ============================================================================

// NewBootstrap creates a new Bootstrap instance
func NewBootstrap() (*Bootstrap, error) {
	retui.Debug("Initializing Bootstrap...")

	cfg := config.Load()
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	b := &Bootstrap{
		Config:  cfg,
		cleanup: make([]func() error, 0),
		state: &appState{
			darkMode: cfg.Theme == "dark",
			// Previously this rebuilt a fresh Config with only AppName
			// copied over, silently dropping Version/Theme/Debug/
			// DateFormat/Timezone/etc. even though Load() had already
			// read them. Reuse the fully loaded config as-is.
			config: cfg,
			user: &user{
				name:  "Guest",
				email: "guest@example.com",
				role:  "viewer",
			},
		},
	}

	b.Ctx, b.Cancel = context.WithCancel(context.Background())

	retui.Debug("Initializing database...")
	client, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("database client is nil")
	}
	b.DB = client
	b.AppCtx = &appctx.AppContext{}

	b.RegisterCleanup(func() error {
		retui.Debug("Closing database connection...")
		if b.DB != nil {
			return b.DB.Close()
		}
		return nil
	})

	if err := b.runMigrations(); err != nil {
		retui.Warnf("Migration warning: %v", err)
	}

	if err := b.seedDatabase(); err != nil {
		retui.Warnf("Seed warning: %v", err)
	}

	b.setContext()
	SetBootstrap(b)
	b.registerModules()

	retui.Success("Bootstrap completed successfully")
	return b, nil
}

// ============================================================================
// DATABASE FUNCTIONS
// ============================================================================

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*database.DB, error) {
	retui.Infof("Connecting to database...")

	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// config.Load() already resolves SQLitePath to an absolute path
	// anchored at the config file's directory (or the executable's own
	// directory when running on defaults) — see config.resolveDataPath.
	// That anchoring is what makes `build/retui` write its database to
	// `build/data/retui.db` regardless of the caller's working
	// directory. Re-deriving it here against cwd, as before, would
	// undo that guarantee for anyone who cd's elsewhere before running
	// the binary, so we just trust it — with a last-resort fallback in
	// case InitDB is ever called directly with a hand-built Config.
	dbPath := cfg.SQLitePath
	if dbPath == "" {
		dbPath = filepath.Join(config.ExecutableDir(), "data", "retui.db")
	} else if !filepath.IsAbs(dbPath) {
		absPath, err := filepath.Abs(dbPath)
		if err == nil {
			dbPath = absPath
		}
	}
	retui.Infof("Database path: %s", dbPath)

	// Create DB using your driver abstraction
	db, err := database.New("sqlite", dbPath)
	if err != nil {
		retui.Errorf("Failed to open database: %v", err)
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run schema migration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	retui.Info("Creating/verifying schema...")
	if err := db.Client.Schema.Create(ctx); err != nil {
		retui.Errorf("Failed to create schema: %v", err)
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	retui.Successf("Database connected: %s", dbPath)
	return db, nil
}

// runMigrations runs database migrations
func (b *Bootstrap) runMigrations() error {
	retui.Info("Running database migrations...")

	ctx, cancel := context.WithTimeout(b.Ctx, 10*time.Second)
	defer cancel()

	if err := b.DB.Client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	retui.Success("Migrations completed successfully")
	return nil
}

// seedDatabase populates initial reference data. Safe to call on every
// startup — seed functions are expected to no-op if data already exists.
func (b *Bootstrap) seedDatabase() error {
	retui.Info("Seeding database...")

	ctx, cancel := context.WithTimeout(b.Ctx, 10*time.Second)
	defer cancel()

	//--Accounting
	if err := seed.AccountingSeed(ctx, b.DB.Client); err != nil {
		return fmt.Errorf("failed to seed accounting data: %w", err)
	}

	retui.Success("Database seeded successfully")
	return nil
}

// ============================================================================
// CONTEXT FUNCTIONS
// ============================================================================

// setContext sets the application context values
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
	if b.state.user != nil && b.state.user.name != "" {
		userName = b.state.user.name
	}

	b.AppCtx.Set(appctx.AppContextValues{
		// CurrentPage: retui.CurrentScreen(),
		AppName:  b.state.config.AppName,
		DarkMode: b.state.darkMode,
		UserName: userName,
		DB:       b.DB,
		Context:  b.Ctx,
		Config:   b.state.config,
		// GetStack:   retui.ScreenStackSnapshot,
		// GetCurrent: retui.CurrentScreen,

		ToggleDark: func() {
			b.mu.Lock()
			b.state.darkMode = !b.state.darkMode
			b.mu.Unlock()
			if b.AppCtx != nil {
				b.AppCtx.ToggleDark()
			}
			retui.Debugf("Dark mode toggled: %v", b.state.darkMode)
		},
	})

	retui.Debug("Context values set successfully")
}

// UpdateContext updates the context with current state
func (b *Bootstrap) UpdateContext() {
	b.setContext()
}

// ============================================================================
// CLEANUP AND SHUTDOWN
// ============================================================================

// RegisterCleanup adds a cleanup function to be called on shutdown
func (b *Bootstrap) RegisterCleanup(fn func() error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cleanup = append(b.cleanup, fn)
}

// Shutdown gracefully shuts down the application
func (b *Bootstrap) Shutdown() error {
	retui.Debug("Shutting down...")

	// NOTE: this used to `return b.DB.Close()` here, which meant
	// Cancel() and every registered cleanup func (including the one
	// that closes b.DB itself, registered in NewBootstrap) were
	// skipped entirely whenever a DB was present. DB.Close() is
	// idempotent-safe to call twice, but we no longer need to since
	// closing it is already one of the registered cleanup functions.

	// Cancel context so any goroutines watching b.Ctx unwind first.
	if b.Cancel != nil {
		b.Cancel()
	}

	// Run cleanup functions in reverse order (most-recently-registered
	// first), matching typical defer-stack teardown semantics.
	var errs []error
	b.mu.Lock()
	for i := len(b.cleanup) - 1; i >= 0; i-- {
		if err := b.cleanup[i](); err != nil {
			errs = append(errs, err)
		}
	}
	b.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	retui.Success("Shutdown completed successfully ✅")
	return nil
}

// Close is an alias for Shutdown for compatibility
func (b *Bootstrap) Close() error {
	return b.Shutdown()
}

// SetBootstrap stores the bootstrap instance globally
func SetBootstrap(b *Bootstrap) {
	bootstrapMu.Lock()
	defer bootstrapMu.Unlock()
	globalBootstrap = b
}

// GetBootstrap returns the global bootstrap instance
func GetBootstrap() *Bootstrap {
	bootstrapMu.RLock()
	defer bootstrapMu.RUnlock()
	return globalBootstrap
}
