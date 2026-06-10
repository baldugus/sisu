package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchSemestersCommand struct{}

func (cmd *FetchSemestersCommand) Execute(db *database.Database) ([]*types.Semester, error) {
	return db.FetchSemesters()
}
