package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchRegistrationCommand struct {
	RegistrationID int32
}

func (cmd *FetchRegistrationCommand) Execute(db *database.Database) (*types.RegistrationDetail, error) {
	return db.FetchRegistrationByID(cmd.RegistrationID)
}
