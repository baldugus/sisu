package commands

import "github.com/baldugus/sisu/database"

// OpenCallCommand reopens a closed call. Undo happens in reverse order, so only
// the last call can be reopened, and only while no semester was closed after it.
type OpenCallCommand struct {
	ID int32
}

func (cmd *OpenCallCommand) Execute(db *database.Database) error {
	callNumber, err := db.GetCallNumber(cmd.ID)
	if err != nil {
		return err
	}

	hasCallAfter, err := db.HasCallAfterNumber(callNumber)
	if err != nil {
		return err
	}

	if hasCallAfter {
		return ErrCannotReopenCallWithLaterCalls{}
	}

	if err := requireNoSemesterClosedSince(db, callNumber); err != nil {
		return err
	}

	return db.OpenCall(cmd.ID)
}

// requireNoSemesterClosedSince fails when a semester was closed after the given
// call, since that closure must be undone first.
func requireNoSemesterClosedSince(db *database.Database, callNumber int32) error {
	semesters, err := db.FetchSemesters()
	if err != nil {
		return err
	}

	for _, s := range semesters {
		if s.ClosedAfterCall != nil && *s.ClosedAfterCall >= callNumber {
			return ErrSemesterClosedAfterCall{}
		}
	}

	return nil
}
