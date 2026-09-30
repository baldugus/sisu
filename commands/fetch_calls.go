package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type FetchCallsCommand struct{}

func (cmd *FetchCallsCommand) Execute(db *database.Database) ([]*types.CallSummary, error) {
	return db.FetchCallSummaries()
}

// FetchCallEntriesCommand returns the entries of a call with their
// registrations and courses.
type FetchCallEntriesCommand struct {
	CallID int32
}

func (cmd *FetchCallEntriesCommand) Execute(db *database.Database) ([]*types.CallEntryDetail, error) {
	return database.FetchCallEntryDetails(db.DB(), cmd.CallID, nil, nil)
}
