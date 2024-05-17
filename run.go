package main

import (
	"embed"
	"errors"
	"log/slog"
	"os"

	"changeme/internal/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed repository/migrations
var migrations embed.FS

func main() {
	db, err := sqlx.Connect("sqlite", "test.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		slog.Error("connect to db", slog.Any("err", err))
		os.Exit(1)
	}

	filesystem, err := iofs.New(migrations, "repository/migrations")
	if err != nil {
		slog.Error("new iofs", slog.Any("err", err))
		os.Exit(1)
	}

	var sqliteConfig sqlite.Config

	s, err := sqlite.WithInstance(db.DB, &sqliteConfig)
	if err != nil {
		slog.Error("new migrate sqlite instance", slog.Any("err", err))
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		slog.Error("new migrate instance", slog.Any("err", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migrate up", slog.Any("err", err))
		os.Exit(1)
	}

	query := sql.Query("insert_selection")

	_, err = sqlx.NamedExec(db, query, struct {
		Name        string
		Kind        string
		Institution string
		Course      string
		Date        string
	}{"a", "b", "c", "d", "e"})
	if err != nil {
		slog.Error("test query", slog.Any("err", err))
		os.Exit(1)
	}

	os.Exit(0)
}
