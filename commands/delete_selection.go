package commands

import (
	"fmt"
	"slices"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
)

type DeleteSelectionCommand struct {
	Kind types.SelectionKind
}

func (cmd *DeleteSelectionCommand) Execute(db *database.Database) error {
	existingKinds, err := db.FetchSelectionKinds()
	if err != nil {
		return fmt.Errorf("fetch selection kinds: %w", err)
	}

	if !slices.Contains(existingKinds, cmd.Kind) {
		return ErrSelectionNotFound{}
	}

	if cmd.Kind == types.SelectionKindApproved {
		if slices.Contains(existingKinds, types.SelectionKindWaitlist) {
			return ErrCannotDeleteApprovedWithWaitlist{}
		}
	}

	hasCallsAfterFirst, err := db.HasCallAfterNumber(1)
	if err != nil {
		return fmt.Errorf("check calls after first: %w", err)
	}

	if hasCallsAfterFirst {
		return ErrCannotDeleteWithMultipleCalls{}
	}

	hasOpenCall, err := db.HasOpenCall()
	if err != nil {
		return fmt.Errorf("check open call: %w", err)
	}

	if !hasOpenCall {
		return ErrCannotDeleteWithClosedFirstCall{}
	}

	hasModified, err := db.SelectionHasModifiedRegistrations(cmd.Kind)
	if err != nil {
		return fmt.Errorf("check modified registrations: %w", err)
	}

	if hasModified {
		return ErrSelectionHasModifiedRegistrations{}
	}

	err = db.RunInTx(func(tx qrm.DB) error {
		if err := database.DeleteAllRegistrations(tx); err != nil {
			return fmt.Errorf("delete registrations: %w", err)
		}
		if err := database.DeleteAllCandidates(tx); err != nil {
			return fmt.Errorf("delete candidates: %w", err)
		}
		if err := database.DeleteAllCourses(tx); err != nil {
			return fmt.Errorf("delete courses: %w", err)
		}
		if err := database.DeleteAllQuotas(tx); err != nil {
			return fmt.Errorf("delete quotas: %w", err)
		}
		if err := database.DeleteAllCalls(tx); err != nil {
			return fmt.Errorf("delete calls: %w", err)
		}
		if err := database.DeleteSelection(tx, cmd.Kind); err != nil {
			return fmt.Errorf("delete selection: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete selection transaction: %w", err)
	}

	return nil
}
