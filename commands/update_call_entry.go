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
//
//   - initial / waitlist entries: pending, enrolled or absent.
//   - promotion entries: pending, enrolled (accepted) or declined.
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

		if !outcomeAllowed(entry.Kind, cmd.Outcome) {
			return ErrInvalidStatusTransition{}
		}

		if entry.Outcome == cmd.Outcome {
			return nil
		}

		// A promotion request only makes sense for a student who holds the seat.
		if cmd.Outcome == types.CallEntryOutcomeAbsent && entry.WantsPromotion {
			if err := database.UpdateCallEntryWantsPromotion(tx, cmd.CallID, cmd.RegistrationID, false); err != nil {
				return err
			}
		}

		return database.UpdateCallEntryOutcome(tx, cmd.CallID, cmd.RegistrationID, cmd.Outcome)
	})
}

func outcomeAllowed(kind types.CallEntryKind, outcome types.CallEntryOutcome) bool {
	switch outcome {
	case types.CallEntryOutcomePending, types.CallEntryOutcomeEnrolled:
		return true
	case types.CallEntryOutcomeAbsent:
		return kind != types.CallEntryKindPromotion
	case types.CallEntryOutcomeDeclined:
		return kind == types.CallEntryKindPromotion
	default:
		return false
	}
}

// SetWantsPromotionCommand records that a semester-2 student asked to move to
// semester 1 in a later call (or withdraws the request).
type SetWantsPromotionCommand struct {
	CallID         int32
	RegistrationID int32
	WantsPromotion bool
}

func (cmd *SetWantsPromotionCommand) Execute(db *database.Database) error {
	return db.RunInTx(func(tx qrm.DB) error {
		entry, err := fetchEditableEntry(tx, cmd.CallID, cmd.RegistrationID)
		if err != nil {
			return err
		}

		if entry.WantsPromotion == cmd.WantsPromotion {
			return nil
		}

		if cmd.WantsPromotion &&
			(entry.Semester != 2 || entry.Kind == types.CallEntryKindPromotion || entry.Outcome == types.CallEntryOutcomeAbsent) {
			return ErrPromotionNotAllowed{}
		}

		return database.UpdateCallEntryWantsPromotion(tx, cmd.CallID, cmd.RegistrationID, cmd.WantsPromotion)
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
