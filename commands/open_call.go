package commands

import "github.com/baldugus/sisu/database"

// OpenCallCommand reopens a closed call. Undo happens in reverse order, so only
// the last call can be reopened.
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

	return db.OpenCall(cmd.ID)
}
