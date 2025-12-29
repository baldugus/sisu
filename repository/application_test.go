// nolint: varnamelen, nolintlint
// TODO: finish tests.
// TODO: test finding applications with filter of selection 2.
package repository_test

import (
	"testing"

	"changeme/repository"
	"changeme/types"

	"github.com/google/go-cmp/cmp"
)

func TestApplicationRepository_CreateApplication(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		applicationRepository := repository.NewApplicationRepository(d)

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

		selection := &types.Selection{
			ID:   0,
			Name: "foo",
			Kind: types.ApprovedSelection,
			Date: "2023-01-01",
		}
		MustCreateSelection(t, d, selection)

		var id int64 = 1
		rollcall := &types.Rollcall{
			ID:     &id,
			Number: 1,
			Status: "CALLING",
		}
		MustCreateRollcall(t, d)

		applicant := &types.Applicant{}
		MustCreateApplicant(t, d, applicant)

		application := &types.Application{
			ID:                   0,
			Status:               "APPROVED",
			EnrollmentID:         "1",
			Option:               1,
			LanguagesScore:       631.0,
			HumanitiesScore:      632.0,
			NaturalSciencesScore: 633.0,
			MathematicsScore:     634.0,
			EssayScore:           635.0,
			CompositeScore:       636.0,
			Ranking:              2,
			Class:                *class,
			Rollcall:             *rollcall,
			Applicant:            *applicant,
		}

		if err := applicationRepository.CreateApplication(application, selection.ID); err != nil {
			t.Fatal(err)
		}

		result, err := applicationRepository.FindApplicationByID(1)
		if err != nil {
			t.Fatal(err)
		}

		c := repository.NewClassRepository(d)
		rClass, err := c.FindClassByID(result.Class.ID)
		if err != nil {
			t.Fatal(err)
		}
		result.Class = *rClass

		r := repository.NewRollcallRepository(d)
		rRollcall, err := r.FindRollcallByID(result.Rollcall.ID)
		if err != nil {
			t.Fatal(err)
		}
		result.Rollcall = *rRollcall

		if !cmp.Equal(application, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(application, result))
		}
	})
}

func MustCreateApplication(tb testing.TB, db *repository.DB, application *types.Application, selection *types.Selection) {
	tb.Helper()

	if err := repository.NewClassRepository(db).CreateClass(&application.Class); err != nil {
		tb.Fatal(err)
	}

	if _, err := repository.NewRollcallRepository(db).CreateRollcall(); err != nil {
		tb.Fatal(err)
	}

	if err := repository.NewSelectionRepository(db).CreateSelection(selection); err != nil {
		tb.Fatal(err)
	}

	if err := repository.NewApplicationRepository(db).CreateApplication(application, selection.ID); err != nil {
		tb.Fatal(err)
	}
}
