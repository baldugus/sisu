package repository_test

// TODO: tests?

import (
	"testing"

	"changeme/repository"
	"changeme/types"
)

func MustCreateApplicant(tb testing.TB, db *repository.DB, applicant *types.Applicant) {
	tb.Helper()

	if err := repository.NewApplicantRepository(db).CreateApplicant(applicant); err != nil {
		tb.Fatal(err)
	}
}
