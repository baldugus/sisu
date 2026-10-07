package commands

import (
	"github.com/go-jet/jet/v2/qrm"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

// DeleteCallCommand removes the last call. Its entries are deleted with it
// (ON DELETE CASCADE), which reverts everything the call did: called students
// go back to the waitlist and promotions back to semester 2.
type DeleteCallCommand struct {
	ID int32
}

func (cmd *DeleteCallCommand) Execute(db *database.Database) error {
	call, err := db.FetchCallByID(cmd.ID)
	if err != nil {
		return err
	}

	if call.Number == 1 {
		return ErrCannotDeleteFirstCall{}
	}

	if call.Status == types.CallStatusDone {
		return ErrCannotDeleteClosedCall{}
	}

	hasCallAfter, err := db.HasCallAfterNumber(call.Number)
	if err != nil {
		return err
	}

	if hasCallAfter {
		return ErrCannotDeleteCallWithLaterCalls{}
	}

	return db.RunInTx(func(tx qrm.DB) error {
		return database.DeleteCall(tx, cmd.ID)
	})
}
