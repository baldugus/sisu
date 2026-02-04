package database

import (
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var testMigrations embed.FS

// TestDB wraps a Database instance for testing with automatic cleanup.
type TestDB struct {
	*Database
	t *testing.T
}

// NewTestDatabase creates a new test database with migrations applied.
// The database is automatically cleaned up when the test finishes.
func NewTestDatabase(t *testing.T) *TestDB {
	t.Helper()

	// Create a temporary database file
	dbPath := filepath.Join(t.TempDir(), "test.db")

	// Open database connection
	db, err := sqlx.Open(
		"sqlite",
		fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", dbPath),
	)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Run migrations
	if err := runTestMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create Database wrapper
	database := NewDatabase(db, dbPath)

	// Register cleanup
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}
	})

	return &TestDB{
		Database: database,
		t:        t,
	}
}

// runTestMigrations runs all migrations on the test database.
func runTestMigrations(db *sqlx.DB) error {
	filesystem, err := iofs.New(testMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("new iofs: %w", err)
	}

	var sqliteConfig sqlite.Config

	s, err := sqlite.WithInstance(db.DB, &sqliteConfig)
	if err != nil {
		return fmt.Errorf("sqlite with instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		return fmt.Errorf("migrate new with instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
