package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchSelectionCommand struct {
	Kind types.SelectionKind
}

func (cmd *FetchSelectionCommand) Execute(db *database.Database) (*types.Selection, error) {
	selection, err := db.FetchSelection(cmd.Kind)
	if err != nil {
		return nil, err
	}

	return selection, nil
}
