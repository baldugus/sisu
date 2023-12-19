// nolint: varnamelen,nolintlint
package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"changeme/repository"
	"changeme/types"

	"github.com/google/go-cmp/cmp"
)

func TestClassRepository_CreateClass(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		classRepository := repository.NewClassRepository(d)
		class := &types.Class{
			ID: 0,
			Period: types.Period{
				ID:   0,
				Name: "Manhã",
			},
			Quota: types.Quota{
				ID:   0,
				Name: "Ampla Concorrência",
			},
			Seats:        30,
			MinimumScore: 640.0,
		}

		if err := classRepository.CreateClass(class); err != nil {
			t.Fatal(err)
		}

		result, err := classRepository.FindClassByID(1)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(class, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(class, result))
		}
	})
}

func TestClassRepository_DeleteClass(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		class := &types.Class{
			ID: 0,
			Period: types.Period{
				ID:   0,
				Name: "Manhã",
			},
			Quota: types.Quota{
				ID:   0,
				Name: "Ampla Concorrência",
			},
			Seats:        30,
			MinimumScore: 640.0,
		}

		MustCreateClass(t, d, class)

		c := repository.NewClassRepository(d)
		if err := c.DeleteClass(1); err != nil {
			t.Fatal(err)
		} else if _, err := c.FindClassByID(class.ID); !errors.Is(err, sql.ErrNoRows) {
			t.Fatal(err)
		}
	})
}

func TestClassRepository_FindClasses(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		classes := []*types.Class{
			{
				ID: 0,
				Period: types.Period{
					ID:   0,
					Name: "Manhã",
				},
				Quota: types.Quota{
					ID:   0,
					Name: "Ampla Concorrência",
				},
				Seats:        30,
				MinimumScore: 640.0,
			},
			{
				ID: 0,
				Period: types.Period{
					ID:   0,
					Name: "Noite",
				},
				Quota: types.Quota{
					ID:   0,
					Name: "Baixa Concorrência",
				},
				Seats:        31,
				MinimumScore: 641.0,
			},

			{
				ID: 0,
				Period: types.Period{
					ID:   0,
					Name: "Manhã",
				},
				Quota: types.Quota{
					ID:   0,
					Name: "Baixa Concorrência",
				},
				Seats:        31,
				MinimumScore: 641.0,
			},

			{
				ID: 0,
				Period: types.Period{
					ID:   0,
					Name: "Noite",
				},
				Quota: types.Quota{
					ID:   0,
					Name: "Ampla Concorrência",
				},
				Seats:        31,
				MinimumScore: 641.0,
			},
		}

		for _, c := range classes {
			MustCreateClass(t, d, c)
		}

		c := repository.NewClassRepository(d)
		result, err := c.FindClasses()
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(classes, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(classes, result))
		}
	})
}

func MustCreateClass(tb testing.TB, db *repository.DB, class *types.Class) {
	tb.Helper()

	if err := repository.NewClassRepository(db).CreateClass(class); err != nil {
		tb.Fatal(err)
	}
}
