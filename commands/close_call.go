package commands

import "github.com/baldugus/sisu/database"

type CloseCallCommand struct {
	ID int32
}

func (cmd *CloseCallCommand) Execute(db *database.Database) error {
	hasPending, err := db.CallHasPendingEntries(cmd.ID)
	if err != nil {
		return err
	}

	if hasPending {
		return ErrCallHasPendingRegistrations{}
	}

	return db.CloseCall(cmd.ID)
}
