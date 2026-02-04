package commands

import "github.com/baldugus/sisu/database"

type CloseCallCommand struct {
	ID int32
}

func (cmd *CloseCallCommand) Execute(db *database.Database) error {
	hasPendingRegistrations, err := db.CallHasPendingRegistrations(cmd.ID)
	if err != nil {
		return err
	}

	if hasPendingRegistrations {
		return ErrCallHasPendingRegistrations{}
	}

	return db.CloseCall(cmd.ID)

}
