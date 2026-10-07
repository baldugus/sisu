package commands

import (
	"github.com/go-jet/jet/v2/qrm"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchSemestersCommand struct{}

func (cmd *FetchSemestersCommand) Execute(db *database.Database) ([]*types.Semester, error) {
	return db.FetchSemesters()
}

// CloseSemesterCommand stops a semester from receiving anyone in later calls.
// The closure is recorded against the last call, so it can be undone while no
// call was created after it.
type CloseSemesterCommand struct {
	Number int32
}

func (cmd *CloseSemesterCommand) Execute(db *database.Database) error {
	semester, err := fetchSemester(db, cmd.Number)
	if err != nil {
		return err
	}

	if semester.Status == types.SemesterStatusClosed {
		return ErrSemesterAlreadyClosed{}
	}

	hasOpenCall, err := db.HasOpenCall()
	if err != nil {
		return err
	}

	if hasOpenCall {
		return ErrCannotCloseSemesterWithOpenCall{}
	}

	lastCallNumber, err := db.GetLastCallNumber()
	if err != nil {
		return err
	}

	if lastCallNumber == 0 {
		return ErrNoCalls{}
	}

	return db.RunInTx(func(tx qrm.DB) error {
		return database.SetSemesterClosedAfterCall(tx, cmd.Number, &lastCallNumber)
	})
}

// ReopenSemesterCommand undoes CloseSemesterCommand while no call was created
// after the closure.
type ReopenSemesterCommand struct {
	Number int32
}

func (cmd *ReopenSemesterCommand) Execute(db *database.Database) error {
	semester, err := fetchSemester(db, cmd.Number)
	if err != nil {
		return err
	}

	if semester.ClosedAfterCall == nil {
		return ErrSemesterNotClosed{}
	}

	hasCallAfter, err := db.HasCallAfterNumber(*semester.ClosedAfterCall)
	if err != nil {
		return err
	}

	if hasCallAfter {
		return ErrCannotReopenSemesterWithLaterCalls{}
	}

	return db.RunInTx(func(tx qrm.DB) error {
		return database.SetSemesterClosedAfterCall(tx, cmd.Number, nil)
	})
}

func fetchSemester(db *database.Database, number int32) (*types.Semester, error) {
	if number != 1 && number != 2 {
		return nil, ErrInvalidSemester{}
	}

	return database.FetchSemester(db.DB(), number)
}
