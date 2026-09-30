package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchRegistrationsCommand struct {
	SelectionID *int32
	CourseID    *int32
}

func (cmd *FetchRegistrationsCommand) Execute(db *database.Database) ([]*types.Registration, error) {
	switch {
	case cmd.SelectionID != nil:
		return db.FetchRegistrationsBySelectionID(*cmd.SelectionID)
	case cmd.CourseID != nil:
		return db.FetchRegistrationsByCourseID(*cmd.CourseID)
	default:
		return db.FetchRegistrations()
	}
}
