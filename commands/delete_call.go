package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/go-jet/jet/v2/qrm"
)

type DeleteCallCommand struct {
	ID int32
}

func (cmd *DeleteCallCommand) Execute(db *database.Database) error {
	callNumber, err := db.GetCallNumber(cmd.ID)
	if err != nil {
		return err
	}

	hasCallAfter, err := db.HasCallAfterNumber(callNumber)
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
