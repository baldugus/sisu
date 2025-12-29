// nolint: varnamelen,nolintlint
package repository_test

import (
	"database/sql"
	"embed"
	"errors"
	"testing"

	"changeme/repository"
	"changeme/types"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func TestSelectionRepository_CreateSelection(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		selectionRepository := repository.NewSelectionRepository(d)
		selection := &types.Selection{
			ID:   0,
			Name: "foo",
			Kind: types.ApprovedSelection,
			Date: "2023-01-01",
		}

		if err := selectionRepository.CreateSelection(selection); err != nil {
			t.Fatal(err)
		}

		result, err := selectionRepository.FindSelectionByID(1)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(selection, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(selection, result))
		}
	})
}

func TestSelectionRepository_DeleteSelection(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		selection := &types.Selection{
			ID:   0,
			Name: "foo",
			Kind: types.ApprovedSelection,
			Date: "2023",
		}

		MustCreateSelection(t, d, selection)

		s := repository.NewSelectionRepository(d)
		if err := s.DeleteSelection(1); err != nil {
			t.Fatal(err)
		} else if _, err := s.FindSelectionByID(selection.ID); !errors.Is(err, sql.ErrNoRows) {
			t.Fatal(err)
		}
	})
}

func TestSelectionRepository_FindSelection(t *testing.T) {
	t.Run("By Kind", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		selection := &types.Selection{
			ID:   0,
			Name: "arquivo",
			Kind: types.ApprovedSelection,
			Date: "2023",
		}

		MustCreateSelection(t, d, selection)

		s := repository.NewSelectionRepository(d)
		result, err := s.FindSelectionByKind(types.ApprovedSelection)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(selection, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(selection, result))
		}
	})
}

func TestSelectionRepository_FindSelections(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		selections := []*types.Selection{
			{
				ID:   0,
				Name: "arquivo",
				Kind: types.ApprovedSelection,
				Date: "2023",
			},
			{
				ID:   0,
				Name: "arquivo2",
				Kind: types.InterestedSelection,
				Date: "2023",
			},
		}

		for _, s := range selections {
			MustCreateSelection(t, d, s)
		}

		s := repository.NewSelectionRepository(d)
		result, err := s.FindSelections()
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(selections, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(selections, result))
		}
	})
}

//go:embed migrations
var migrations embed.FS

func MustOpenDB(tb testing.TB) *sqlx.DB {
	tb.Helper()

	db, err := sqlx.Open(
		"sqlite",
		":memory:?_pragma=foreign_keys(1)",
	)
	if err != nil {
		tb.Fatalf("sqlx open: %v", err)
	}

	filesystem, err := iofs.New(migrations, "migrations")
	if err != nil {
		tb.Fatalf("iofs new: %v", err)
	}

	s, err := sqlite.WithInstance(db.DB, &sqlite.Config{}) //nolint: exhaustruct
	if err != nil {
		tb.Fatalf("sqlite with instance: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", filesystem, "sqlite", s)
	if err != nil {
		tb.Fatalf("migrate new with instance: %v", err)
	}

	if err := m.Up(); err != nil {
		tb.Fatalf("migrate up: %v", err)
	}

	return db
}

func MustCloseDB(tb testing.TB, db *sqlx.DB) {
	tb.Helper()

	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}

func MustCreateSelection(tb testing.TB, db *repository.DB, selection *types.Selection) {
	tb.Helper()

	if err := repository.NewSelectionRepository(db).CreateSelection(selection); err != nil {
		tb.Fatal(err)
	}
}
