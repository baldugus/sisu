package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
)

type DeleteCallCommand struct {
	ID int32
}

func (cmd *DeleteCallCommand) Execute(db *database.Database) error {
	call, err := db.FetchCallByID(cmd.ID)
	if err != nil {
		return err
	}

	// Cannot delete a closed call
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
		if err := database.RevertRegistrationsToWaitlisted(tx, cmd.ID); err != nil {
			return err
		}

		return database.DeleteCall(tx, cmd.ID)
	})
}
