package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchCallsCommand struct {
}

func (cmd *FetchCallsCommand) Execute(db *database.Database) ([]*types.Call, error) {
	calls, err := db.FetchCalls()
	if err != nil {
		return nil, err
	}

	return calls, nil
}
