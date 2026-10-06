package commands

import (
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

// SetCallEntryOutcomeCommand records what happened to a registration in a call.
// Only entries of the open call can change.
type SetCallEntryOutcomeCommand struct {
	CallID         int32
	RegistrationID int32
	Outcome        types.CallEntryOutcome
}

func (cmd *SetCallEntryOutcomeCommand) Execute(db *database.Database) error {
	return db.RunInTx(func(tx qrm.DB) error {
		entry, err := fetchEditableEntry(tx, cmd.CallID, cmd.RegistrationID)
		if err != nil {
			return err
		}

		if !cmd.Outcome.IsValid() {
			return ErrInvalidStatusTransition{}
		}

		if entry.Outcome == cmd.Outcome {
			return nil
		}

		return database.UpdateCallEntryOutcome(tx, cmd.CallID, cmd.RegistrationID, cmd.Outcome)
	})
}

// fetchEditableEntry returns the entry if it exists and its call is open.
func fetchEditableEntry(tx qrm.DB, callID, registrationID int32) (*types.CallEntry, error) {
	call, err := database.FetchCallByID(tx, callID)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCallNotFound{}
		}

		return nil, err
	}

	if call.Status != types.CallStatusCalling {
		return nil, ErrCallNotOpen{}
	}

	entry, err := database.FetchCallEntry(tx, callID, registrationID)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRegistrationNotFound{}
		}

		return nil, err
	}

	return entry, nil
}
