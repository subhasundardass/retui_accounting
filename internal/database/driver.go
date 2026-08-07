package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/retui"
)

type PostgresDriver struct{}
type SQLiteDriver struct{}

// Driver interface
type Driver interface {
	Open(url string) (*ent.Client, error)
}

// GetDriver returns backend driver
func GetDriver(name string) (Driver, error) {
	switch name {
	case "postgres", "postgresql":
		return PostgresDriver{}, nil
	case "sqlite":
		return SQLiteDriver{}, nil
	default:
		return nil, errors.New("unsupported database driver: " + name)
	}
}

// Postgres
func (d PostgresDriver) Open(url string) (*ent.Client, error) {

	drv, err := entsql.Open(dialect.Postgres, url)
	if err != nil {
		return nil, err
	}

	return ent.NewClient(ent.Driver(drv)), nil
}

// SQLite
func (d SQLiteDriver) Open(url string) (*ent.Client, error) {

	retui.Infof("SQLiteDriver.Open called with url: %s", url)

	if err := ensureDir(url); err != nil {
		retui.Errorf("ensureDir failed: %v", err)
		return nil, err
	}
	retui.Info("ensureDir succeeded ✓")

	// Add WAL mode + busy timeout + foreign keys
	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",
		url,
	)
	retui.Infof("Opening with DSN: %s", dsn)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		retui.Errorf("sql.Open failed: %v", err)
		return nil, err
	}

	// CRITICAL for SQLite — limit to one connection
	// SQLite is not safe for concurrent writes from multiple connections
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		retui.Errorf("db.Ping failed: %v", err)
		return nil, err
	}
	retui.Info("db.Ping succeeded ✓")

	drv := entsql.OpenDB(dialect.SQLite, db)
	retui.Info("entsql.OpenDB succeeded ✓")

	return ent.NewClient(ent.Driver(drv)), nil
}

func ensureDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0755)
}
