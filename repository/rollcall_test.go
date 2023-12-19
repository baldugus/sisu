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

func TestRollcallRepository_CreateRollcall(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		rollcallRepository := repository.NewRollcallRepository(d)
		var id int64 = 1
		rollcall := &types.Rollcall{
			ID:     &id,
			Number: 1,
			Status: "CALLING",
		}

		if _, err := rollcallRepository.CreateRollcall(); err != nil {
			t.Fatal(err)
		}

		result, err := rollcallRepository.FindRollcallByID(&id)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(rollcall, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcall, result))
		}
	})
}

func TestRollcallRepository_FindRollcall(t *testing.T) {
	t.Run("By Number", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		var id int64 = 1
		rollcall := &types.Rollcall{
			ID:     &id,
			Number: 1,
			Status: "CALLING",
		}

		MustCreateRollcall(t, d)

		r := repository.NewRollcallRepository(d)
		result, err := r.FindRollcallByNumber(1)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(rollcall, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcall, result))
		}
	})
}

func TestRollcallRepository_FindRollcalls(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		var firstID int64 = 1
		var secondID int64 = 2
		rollcalls := []*types.Rollcall{
			{
				ID:     &firstID,
				Number: 1,
				Status: "DONE",
			},
			{
				ID:     &secondID,
				Number: 2,
				Status: "CALLING",
			},
		}

		MustCreateRollcall(t, d)
		r := repository.NewRollcallRepository(d)
		status := "DONE"
		_, err := r.UpdateRollcall(1, types.RollcallUpdate{Status: &status})
		if err != nil {
			t.Fatal(err)
		}

		MustCreateRollcall(t, d)

		result, err := r.FindRollcalls(types.RollcallsFilter{Status: nil})
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(rollcalls, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcalls, result))
		}

		status = "CALLING"
		result, err = r.FindRollcalls(types.RollcallsFilter{Status: &status})
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(rollcalls[1:], result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcalls[1:], result))
		}
	})
}

func TestRollcallRepository_UpdateRollcall(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		var firstID int64 = 1
		var secondID int64 = 2
		rollcalls := []*types.Rollcall{
			{
				ID:     &firstID,
				Number: 1,
				Status: "DONE",
			},
			{
				ID:     &secondID,
				Number: 1,
				Status: "CALLING",
			},
		}

		MustCreateRollcall(t, d)

		r := repository.NewRollcallRepository(d)
		status := "DONE"
		result, err := r.UpdateRollcall(1, types.RollcallUpdate{Status: &status})
		if err != nil {
			t.Fatal(err)
		}

		MustCreateRollcall(t, d)

		r = repository.NewRollcallRepository(d)
		status = "CALLING"
		_, err = r.UpdateRollcall(1, types.RollcallUpdate{Status: &status})
		if err != nil && !errors.Is(err, types.ErrActiveRollcalls) {
			t.Fatal(err)
		}

		if cmp.Equal(rollcalls, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcalls, result))
		}

		var id int64 = 1
		findResult, err := r.FindRollcallByID(&id)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(result, findResult) {
			t.Errorf("error after insert:\n%s", cmp.Diff(result, findResult))
		}

		status = "DONE"
		result, err = r.UpdateRollcall(2, types.RollcallUpdate{Status: &status})
		if err != nil {
			t.Fatal(err)
		}

		if cmp.Equal(rollcalls, result) {
			t.Errorf("error after insert:\n%s", cmp.Diff(rollcalls, result))
		}

		var otherID int64 = 2
		findResult, err = r.FindRollcallByID(&otherID)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(result, findResult) {
			t.Errorf("error after insert:\n%s", cmp.Diff(result, findResult))
		}
	})
}

func TestRollcallRepository_DeleteRollcall(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		d := repository.NewDB(db)
		var id int64 = 1
		rollcall := &types.Rollcall{
			ID:     &id,
			Number: 1,
			Status: "CALLING",
		}

		MustCreateRollcall(t, d)

		r := repository.NewRollcallRepository(d)
		if err := r.DeleteRollcall(1); err != nil {
			t.Fatal(err)
		} else if _, err := r.FindRollcallByID(rollcall.ID); !errors.Is(err, sql.ErrNoRows) {
			t.Fatal(err)
		}
	})
}

func MustCreateRollcall(tb testing.TB, db *repository.DB) {
	tb.Helper()

	if _, err := repository.NewRollcallRepository(db).CreateRollcall(); err != nil {
		tb.Fatal(err)
	}
}
